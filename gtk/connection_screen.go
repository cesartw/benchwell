package gtk

import (
	"log"

	"bitbucket.org/goreorto/sqlhero/gtk/controls"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func (f *Factory) NewConnectionScreen() (*ConnectionScreen, error) {
	cs := &ConnectionScreen{}
	return cs, cs.init()
}

type ConnectionScreen struct {
	*gtk.Paned
	tableList *controls.List
	result    *controls.Result
	dbCombo   *gtk.ComboBox
	tabber    *gtk.Notebook

	databaseNames []string
	dbStore       *gtk.ListStore

	activeDatabase controls.MVar
}

func (c *ConnectionScreen) init() error {
	var err error

	c.dbStore, err = gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		return err
	}

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return err
	}

	c.Paned.SetHExpand(true)
	c.Paned.SetVExpand(true)

	// Sidebar

	sideBar, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}

	sideBar.SetSizeRequest(300, -1)

	c.Paned.Pack1(sideBar, false, true)

	tableListSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return err
	}

	c.dbCombo, err = gtk.ComboBoxNewWithModelAndEntry(c.dbStore)
	if err != nil {
		return err
	}
	c.dbCombo.SetEntryTextColumn(0)

	c.tableList, err = controls.NewList(controls.ListOptions{})
	if err != nil {
		return err
	}

	c.tableList.SetHExpand(true)
	c.tableList.SetVExpand(true)
	tableListSW.Add(c.tableList)

	sideBar.PackStart(c.dbCombo, false, true, 0)
	sideBar.PackStart(tableListSW, true, true, 0)

	// main section

	mainSection, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}

	c.tabber, err = gtk.NotebookNew()
	if err != nil {
		return err
	}
	c.tabber.SetVExpand(true)
	c.tabber.SetHExpand(true)

	mainSection.Add(c.tabber)
	mainSection.SetVExpand(true)
	mainSection.SetHExpand(true)

	c.Paned.Pack2(mainSection, true, false)

	// signals

	c.dbCombo.Connect("changed", c.onDatabaseSelected)

	c.Paned.ShowAll()

	return nil
}

func (c *ConnectionScreen) AddTab(title string, content gtk.IWidget, switchNow bool) error {
	header, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return err
	}

	label, err := gtk.LabelNew(title)
	if err != nil {
		return err
	}

	image, err := gtk.ImageNewFromIconName("window-close", gtk.ICON_SIZE_MENU)
	if err != nil {
		return err
	}

	btn, err := gtk.ButtonNew()
	if err != nil {
		return err
	}
	btn.SetImage(image)
	btn.SetRelief(gtk.RELIEF_NONE)
	c.tabber.SetProperty("scrollable", true)
	c.tabber.SetProperty("enable-popup", true)

	header.PackStart(label, true, true, 0)
	header.PackEnd(btn, false, false, 0)
	header.ShowAll()

	c.tabber.AppendPage(content, header)
	c.tabber.SetTabReorderable(content, true)

	btn.Connect("clicked", func() {
		index := c.tabber.PageNum(content)
		if index == -1 {
			return
		}
		c.tabber.RemovePage(index)
	})

	if switchNow {
		c.tabber.SetCurrentPage(c.tabber.GetNPages() - 1)
	}

	return nil
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

func (c *ConnectionScreen) SetActiveDatabase(dbName string) {
	for i, db := range c.databaseNames {
		if db == dbName {
			c.activeDatabase.Set(dbName)
			c.dbCombo.SetActive(i)
			return
		}
	}
}

func (c *ConnectionScreen) ActiveDatabase() (string, bool) {
	if c.activeDatabase.Get() == nil {
		return "", false
	}
	return c.activeDatabase.Get().(string), true
}

func (c *ConnectionScreen) ActiveTable() (string, bool) {
	return c.tableList.SelectedItem()
}

func (c *ConnectionScreen) Dispose() {
}
