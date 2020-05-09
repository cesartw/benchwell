package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type ConnectionCtrl struct {
	*TabCtrl

	ctx  sqlengine.Context
	scr  *gtk.ConnectionScreen
	conn *config.Connection

	tabs []*TableCtrl
}

func (c ConnectionCtrl) Init(ctx sqlengine.Context, p *TabCtrl, conn *config.Connection) (*ConnectionCtrl, error) {
	c.TabCtrl = p
	c.ctx = ctx
	c.conn = conn

	dbNames, err := c.engine.Databases(c.ctx)
	if err != nil {
		return nil, err
	}

	c.scr, err = c.app.NewConnectionScreen()
	if err != nil {
		return nil, err
	}

	c.scr.SetDatabases(dbNames)

	c.scr.ShowAll()

	c.scr.OnDatabaseSelected(c.onDatabaseSelected)
	c.scr.OnTableSelected(c.onTableSelected)

	if c.conn.Database != "" {
		c.scr.SetActiveDatabase(c.conn.Database)
		c.onDatabaseSelected()
	}

	err = c.AddTab()
	if err != nil {
		return nil, err
	}

	c.scr.OnSchemaMenu(c.onSchemaMenu)

	return &c, nil
}

func (c *ConnectionCtrl) Close() bool {
	return c.scr.Close()
}

func (c *ConnectionCtrl) AddTab() error {
	tab, err := TableCtrl{}.init(c.ctx, c, driver.TableDef{})
	if err != nil {
		return err
	}

	return c.scr.AddTab("New", tab.Screen().(ggtk.IWidget), true)
}

func (c *ConnectionCtrl) onDatabaseSelected() {
	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.window.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}
	c.ctx, err = c.engine.UseDatabase(c.ctx, dbName)

	tables, err := c.engine.Tables(c.ctx)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.scr.SetTables(tables)
}

func (c *ConnectionCtrl) onTableSelected() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}
	tab, err := TableCtrl{}.init(c.ctx, c, tableDef)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.scr.AddTab(tableDef.Name, tab.Screen().(ggtk.IWidget), true)
	tab.OnConnect()
}

func (c *ConnectionCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectionCtrl) onSchemaMenu() {
	tableName, ok := c.scr.SelectedTable()
	if !ok {
		return
	}

	schema, err := c.engine.GetCreateTable(c.ctx, tableName)
	if err != nil {
		config.Env.Log.Error(err, "getting table schema")
	}

	c.scr.ShowTableSchemaModal(tableName, schema)
}
