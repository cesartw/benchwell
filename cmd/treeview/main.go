package main

import (
	"context"
	"log"
	"os"
	"sync"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/logger"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"bitbucket.org/goreorto/sqlhero/ui/controls"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Create and initialize the window
func setupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(600, 300)
	return win
}

func main() {
	conf, err := config.New("")
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(os.Stdout)
	eng := sqlengine.New(conf, log)
	defer eng.Dispose()

	ctx, err := eng.Connect(context.TODO(), conf.Connection[0].GetDSN())
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := eng.Databases(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx, err = eng.UseDatabase(ctx, dbs[3])
	if err != nil {
		log.Fatal(err)
	}

	tables, err := eng.Tables(ctx)
	if err != nil {
		log.Fatal(err)
	}

	gtk.Init(nil)

	win := setupWindow("Go Feature Timeline")

	tableList, err := NewUITableList(tables)
	if err != nil {
		log.Fatal(err)
	}

	resultGrid, err := controls.NewUIResult(nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	tableList.Connect("table-activated", func(_ *gtk.ListBox) {
		table, ok := tableList.CurrentTable()
		if !ok {
			return
		}
		log.Debugf("===table %s", table)

		def, data, err := eng.FetchTable(ctx, table, 0, 10)
		if err != nil {
			log.Fatal(err)
		}

		err = resultGrid.UpdateData(def, data)
		if err != nil {
			log.Fatal(err)
		}
	})

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal(err)
	}
	box.PackStart(tableList, false, true, 0)
	box.PackStart(resultGrid, true, true, 0)

	win.Add(box)
	win.ShowAll()
	gtk.Main()
}

type MVar struct {
	value interface{}
	sync.RWMutex
}

func (mv *MVar) Set(v interface{}) {
	mv.Lock()
	defer mv.Unlock()

	mv.value = v
}

func (mv *MVar) Get() interface{} {
	mv.RLock()
	defer mv.RUnlock()

	return mv.value
}

type UITableList struct {
	currentTable MVar
	names        []string
	*gtk.ListBox
}

type uiTableListRow struct {
	data interface{}
	*gtk.ListBoxRow
}

var sg, _ = glib.SignalNew("table-activated")

func NewUITableList(names []string) (*UITableList, error) {
	list := &UITableList{
		names: names,
	}

	var err error
	list.ListBox, err = gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}
	list.SetProperty("activate-on-single-click", false)
	list.Connect("row-activated", list.rowActivated)

	for _, name := range names {
		label, err := gtk.LabelNew(name)
		if err != nil {
			return nil, err
		}

		row, err := gtk.ListBoxRowNew()
		if err != nil {
			return nil, err
		}

		row.SetHAlign(gtk.ALIGN_START)
		row.Add(label)
		list.Add(row)
	}

	return list, nil
}

func (u *UITableList) rowActivated(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	u.currentTable.Set(u.names[row.GetIndex()])
	u.ListBox.Emit("table-activated", u.names[row.GetIndex()])
}

func (u *UITableList) CurrentTable() (string, bool) {
	v := u.currentTable.Get()
	if v == nil {
		return "", false
	}
	table, ok := u.currentTable.Get().(string)

	return table, ok
}
