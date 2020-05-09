package gtk

import (
	"fmt"
	"regexp"

	"bitbucket.org/goreorto/sqlaid/assets"
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func (f *App) NewConnectionScreen() (*ConnectionScreen, error) {
	cs := &ConnectionScreen{}
	return cs, cs.init()
}

type ConnectionScreen struct {
	*gtk.Paned
	dbCombo     *gtk.ComboBox
	tableFilter *gtk.SearchEntry
	tableList   *List
	tabber      *gtk.Notebook

	databaseNames []string
	dbStore       *gtk.ListStore

	activeDatabase MVar

	tableMenu    *gtk.Menu
	editMenu     *gtk.MenuItem
	schemaMenu   *gtk.MenuItem
	truncateMenu *gtk.MenuItem
	deleteMenu   *gtk.MenuItem
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

	c.Paned.Pack1(sideBar, false, true)

	c.tableFilter, err = gtk.SearchEntryNew()
	if err != nil {
		return err
	}
	c.tableFilter.Connect("search-changed", c.onSearch)
	c.tableFilter.SetPlaceholderText("Filter table: .*")

	// TODO: figure out how to focus on accelerator
	//k, mod := gtk.AcceleratorParse("<Control>f")
	//c.tableFilter.AddAccelerator("activate", nil, k, mod, gtk.ACCEL_VISIBLE)

	tableListSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return err
	}

	c.dbCombo, err = gtk.ComboBoxNewWithModelAndEntry(c.dbStore)
	if err != nil {
		return err
	}
	c.dbCombo.SetEntryTextColumn(0)

	c.tableList, err = NewList(ListOptions{
		SelectOnRightClick: true,
		IconFunc: func(name fmt.Stringer) *gdk.Pixbuf {
			def, ok := name.(driver.TableDef)
			if ok {
				if def.Type == driver.TableTypeRegular {
					return assets.Table
				} else {
					return assets.TableCustom
				}
			}
			return assets.Table
		},
	})
	if err != nil {
		return err
	}
	c.tableList.SetHExpand(true)
	c.tableList.SetVExpand(true)
	c.tableList.OnButtonPress(c.onTableListButtonPress)

	tableListSW.Add(c.tableList)

	sideBar.PackStart(c.dbCombo, false, true, 0)
	sideBar.PackStart(c.tableFilter, false, true, 4)
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
	// TODO: move to config
	//c.tabber.SetProperty("tab-pos", gtk.POS_BOTTOM)
	c.tabber.SetProperty("scrollable", true)
	c.tabber.SetProperty("enable-popup", true)

	mainSection.Add(c.tabber)
	mainSection.SetVExpand(true)
	mainSection.SetHExpand(true)

	c.Paned.Pack2(mainSection, true, false)

	err = c.initTableMenu()
	if err != nil {
		return err
	}

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

func (c *ConnectionScreen) Close() bool {
	if c.tabber.GetCurrentPage() == -1 {
		return false
	}

	c.tabber.RemovePage(c.tabber.GetCurrentPage())
	return true
}

func (c *ConnectionScreen) SetDatabases(dbs []string) {
	c.databaseNames = dbs

	for _, name := range dbs {
		c.dbStore.SetValue(c.dbStore.Append(), 0, name)
	}
}

func (c *ConnectionScreen) SetTables(tables driver.TableDefs) {
	c.tableList.UpdateItems(tables.ToStringer())
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

func (c *ConnectionScreen) ActiveTable() (driver.TableDef, bool) {
	i, ok := c.tableList.SelectedItem()
	if !ok {
		return driver.TableDef{}, false
	}

	return i.(driver.TableDef), true
}

func (c *ConnectionScreen) ShowTableSchemaModal(tableName, schema string) {
	modal, err := gtk.DialogNewWithButtons(fmt.Sprintf("Table %s", tableName), nil,
		gtk.DIALOG_DESTROY_WITH_PARENT|gtk.DIALOG_MODAL,
		[]interface{}{"Ok", gtk.RESPONSE_ACCEPT},
	)
	if err != nil {
		return
	}

	modal.SetDefaultSize(400, 400)
	content, err := modal.GetContentArea()
	if err != nil {
		return
	}

	textView, err := gtk.TextViewNew()
	if err != nil {
		return
	}
	textView.SetVExpand(true)
	textView.SetHExpand(true)

	schema, err = ChromaHighlight(schema)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	buff, err := textView.GetBuffer()
	if err != nil {
		return
	}

	buff.InsertMarkup(buff.GetStartIter(), schema)

	textView.Show()
	content.Add(textView)

	modal.Run()
}

func (c *ConnectionScreen) SelectedTable() (string, bool) {
	i, ok := c.tableList.SelectedItem()
	if !ok {
		return "", ok
	}
	return i.String(), true
}

func (c *ConnectionScreen) OnEditMenu(fn interface{}) {
	c.editMenu.Connect("activate", fn)
}

func (c *ConnectionScreen) OnSchemaMenu(fn interface{}) {
	c.schemaMenu.Connect("activate", fn)
}

func (c *ConnectionScreen) OnTruncateMenu(fn interface{}) {
	c.truncateMenu.Connect("activate", fn)
}

func (c *ConnectionScreen) OnDeleteMenu(fn interface{}) {
	c.deleteMenu.Connect("activate", fn)
}

func (c *ConnectionScreen) onTableListButtonPress(_ *gtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	c.tableMenu.ShowAll()
	c.tableMenu.PopupAtPointer(e)
}

func (c *ConnectionScreen) onDatabaseSelected() {
	iter, err := c.dbCombo.GetActiveIter()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	v, err := c.dbStore.GetValue(iter, 0)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	dbName, err := v.GetString()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.activeDatabase.Set(dbName)
}
func (c *ConnectionScreen) initTableMenu() error {
	var err error
	c.tableMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.editMenu, err = menuItemWithImage("Edit", "gtk-edit")
	if err != nil {
		return err
	}

	c.schemaMenu, err = menuItemWithImage("Schema", "gtk-info")
	if err != nil {
		return err
	}

	c.truncateMenu, err = menuItemWithImage("Truncate", "gtk-clear")
	if err != nil {
		return err
	}

	c.deleteMenu, err = menuItemWithImage("Delete", "gtk-delete")
	if err != nil {
		return err
	}

	c.tableMenu.Add(c.editMenu)
	c.tableMenu.Add(c.schemaMenu)
	c.tableMenu.Add(c.truncateMenu)

	return nil
}

func (c *ConnectionScreen) onSearch(e *gtk.SearchEntry) {
	buff, err := e.GetBuffer()
	if err != nil {
		return
	}

	txt, err := buff.GetText()
	if err != nil {
		return
	}

	rg, err := regexp.Compile(txt)
	if err != nil {
		rg = regexp.MustCompile(fmt.Sprintf(".*%s.*", regexp.QuoteMeta(txt)))
	}
	c.tableList.SetFilterRegex(rg)
}
