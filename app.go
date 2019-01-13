package sqlhero

import (
	"context"
	"time"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/connect"
	"bitbucket.org/goreorto/sqlhero/server"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type keybindable interface {
	Keybinds() map[tcell.Key]tview.Primitive
}

// App ...
type App struct {
	*tview.Application
	layout        *layout
	currentScreen keybindable
	eng           *sqlengine.Engine
	conf          *config.Config
}

// New ...
func New(conf *config.Config, eng *sqlengine.Engine) *App {
	app := &App{
		Application: tview.NewApplication(),
		eng:         eng,
		conf:        conf,
		layout:      newLayout(),
	}

	screenConnect := app.newConnectScreen()
	app.currentScreen = screenConnect

	app.layout.SetScreen(screenConnect)
	app.SetRoot(app.layout, true)
	app.SetFocus(app.layout.screen)

	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if p, ok := app.currentScreen.Keybinds()[e.Key()]; ok {
			app.SetFocus(p)
			return nil
		}

		return e
	})

	return app
}

func (app *App) newServerScreen() *server.Screen {
	screenServer := server.New(app.Application)

	screenServer.OnSelectDatabase = func(db string) {
		ctx, err := app.eng.UseDatabase(screenServer.Context(), db)
		if err != nil {
			app.layout.SetStatus("Error selecting db %s", err.Error())
			return
		}

		screenServer.SetContext(ctx)

		tables, err := app.eng.Tables(ctx)
		if err != nil {
			app.layout.SetStatus("Error fetching tables %s", err.Error())
			return
		}

		app.layout.SetStatus("Using `%s`", db)
		screenServer.SetTables(tables)
		app.SetFocus(screenServer.TableList())
	}

	screenServer.OnSelectTable = func(tableName string) {
		def, rows, err := app.eng.FetchTable(screenServer.Context(), tableName, 0, 20)
		if err != nil {
			app.layout.SetStatus("Error fetching table %s", err.Error())
			return
		}

		app.layout.SetStatus("Count: %d Offset: %d Limit: %d", len(rows), 0, 20)
		screenServer.SetData(tableName, def, rows)
		app.SetFocus(screenServer.RecordTable())
	}

	screenServer.OnSaveRecord = func(tableName string, def []driver.ColDef, values, oldValues []*string) bool {
		_, err := app.eng.UpdateRecord(screenServer.Context(), tableName, def, values, oldValues)
		if err != nil {
			app.layout.SetStatus("Error saving record %s", err.Error())
			return false
		}

		app.layout.SetStatus("Saved")
		app.QueueUpdateDrawAfter(func() {
			app.layout.SetStatus("")
		}, 5*time.Second)
		return true
	}

	screenServer.OnInsertRecord = func(tableName string, def []driver.ColDef, values []*string) []*string {
		data, err := app.eng.InsertRecord(screenServer.Context(), tableName, def, values)
		if err != nil {
			app.layout.SetStatus("Error inserting record %s", err.Error())
			return nil
		}

		app.layout.SetStatus("Inserted")
		app.QueueUpdateDrawAfter(func() {
			app.layout.SetStatus("")
		}, 5*time.Second)

		return data
	}

	screenServer.OnDeleteRecord = func(tableName string, def []driver.ColDef, row, oldRow []*string) bool {
		return true
	}

	screenServer.OnReload = func(tableName string) {
		screenServer.OnSelectTable(tableName)
	}

	return screenServer
}

func (app *App) newConnectScreen() *connect.Screen {
	screen := connect.New(app.Application, app.conf)

	screen.OnConnect = func(c config.Connection) {
		app.layout.SetStatus("Connecting...")

		ctx, err := app.eng.Connect(context.Background(), c.DSN())
		if err != nil {
			app.layout.SetStatus("Error connecting %s", err.Error())
			return
		}

		screenServer := app.newServerScreen()
		screenServer.SetContext(ctx)

		dbs, err := app.eng.Databases(ctx)
		if err != nil {
			app.layout.SetStatus("Error fetching database %s", err.Error())
			return
		}

		app.currentScreen = screenServer
		screenServer.SetDatabases(dbs)

		app.layout.SetStatus("Connected")
		app.QueueUpdateDrawAfter(func() {
			app.layout.SetStatus("")
		}, 5*time.Second)
		app.layout.SetScreen(screenServer)
		app.SetFocus(screenServer)
	}

	return screen
}

// QueueUpdateDrawAfter ...
func (app *App) QueueUpdateDrawAfter(f func(), d time.Duration) {
	t := time.NewTimer(d)
	go func() {
		<-t.C
		app.QueueUpdateDraw(f)
	}()
}
