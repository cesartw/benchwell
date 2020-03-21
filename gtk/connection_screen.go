package gtk

import (
	"log"

	"bitbucket.org/goreorto/sqlhero/gtk/controls"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func newConnectionScreen(ctx sqlengine.Context) (*ConnectionScreen, error) {
	cs := &ConnectionScreen{ctx: ctx}
	return cs, cs.init()
}

type ConnectionScreen struct {
	ctx sqlengine.Context
	*gtk.Paned
	tableList *controls.List
	result    *controls.Result
	dbCombo   *gtk.ComboBox

	databaseNames []string
	dbStore       *gtk.ListStore

	activeDatabase controls.MVar
}

func (c *ConnectionScreen) init() error {
	var err error

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return err
	}

	c.Paned.SetHExpand(true)
	c.Paned.SetVExpand(true)

	frame1, err := gtk.FrameNew("")
	if err != nil {
		return err
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}

	frame1.Add(box)

	frame2, err := gtk.FrameNew("")
	if err != nil {
		return err
	}

	frame1.SetShadowType(gtk.SHADOW_IN)
	frame1.SetSizeRequest(300, -1)
	frame2.SetShadowType(gtk.SHADOW_IN)
	frame2.SetSizeRequest(50, -1)

	c.dbStore, err = gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		return err
	}

	c.dbCombo, err = gtk.ComboBoxNewWithModelAndEntry(c.dbStore)
	if err != nil {
		return err
	}
	c.dbCombo.SetEntryTextColumn(0)

	c.dbCombo.Connect("changed", c.onDatabaseSelected)

	sw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return err
	}
	sw2, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return err
	}

	c.tableList, err = controls.NewList(nil)
	if err != nil {
		return err
	}
	sw.Add(c.tableList)

	c.result, err = controls.NewResult(nil, nil)
	if err != nil {
		return err
	}
	sw2.Add(c.result)
	frame2.Add(sw2)

	c.tableList.Connect("row-activated", c.onTableActivated)
	c.tableList.Connect("row-selected", c.onTableSelected)

	c.tableList.SetHExpand(true)
	c.tableList.SetVExpand(true)

	box.PackStart(c.dbCombo, false, true, 0)
	box.PackStart(sw, true, true, 0)

	c.Paned.Pack1(frame1, false, true)
	c.Paned.Pack2(frame2, true, false)

	c.Paned.ShowAll()

	return nil
}

func (c *ConnectionScreen) onTableSelected() {
	//index, ok := c.tablenList.SelectedItemIndex()
	//if !ok {
	//return
	//}
}

func (c *ConnectionScreen) onTableActivated() {
	//index, ok := c.tablenList.ActiveItemIndex()
	//if !ok {
	//return
	//}
}

func (c *ConnectionScreen) SetDatabases(dbs []string) {
	c.databaseNames = dbs

	for _, name := range dbs {
		c.dbStore.SetValue(c.dbStore.Append(), 0, name)
	}
}

func (c *ConnectionScreen) SetTables(tables []string) {
	c.tableList.UpdateItems(tables)
}

func (c *ConnectionScreen) SetTableData(cols []driver.ColDef, data [][]interface{}) error {
	return c.result.UpdateData(cols, data)
}

func (c *ConnectionScreen) onDatabaseSelected() {
	iter, err := c.dbCombo.GetActiveIter()
	if err != nil {
		log.Fatal(err)
		return
	}

	v, err := c.dbStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	dbName, err := v.GetString()
	if err != nil {
		log.Fatal(err)
		return
	}

	c.activeDatabase.Set(dbName)
}

func (c *ConnectionScreen) OnDatabaseSelected(f interface{}) {
	c.dbCombo.Connect("changed", f)
}

func (c *ConnectionScreen) OnTableSelected(f interface{}) {
	c.tableList.Connect("row-activated", f)
}

func (c *ConnectionScreen) ActiveDatabase() string {
	return c.activeDatabase.Get().(string)
}

func (c *ConnectionScreen) ActiveTable() (string, bool) {
	return c.tableList.SelectedItem()
}

func (c *ConnectionScreen) Dispose() {
}
