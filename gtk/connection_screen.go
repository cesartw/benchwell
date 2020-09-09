package gtk

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/sqlengine"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
)

type connectionScreenCtrl interface {
	OnDatabaseSelected()
	OnTableSelected()
	OnSchemaMenu()
	OnNewTabMenu()
	OnRefreshMenu()
	OnEditTable()
	OnTruncateTable()
	OnDeleteTable()
	OnCopySelect()
	OnCopyLog()
	OnUpdateRecord([]driver.ColDef, []interface{}) error
	OnCreateRecord([]driver.ColDef, []interface{}) ([]interface{}, error)
	OnExecQuery(string)
	OnTextChange(string, int) //query, cursor position
	OnRefresh()
	OnDelete()
	OnCreate()
	OnCopyInsert([]driver.ColDef, []interface{})
	OnFileSelected(string)
	OnSaveQuery(string, string)
	OnSaveFav(string, string)
	ParseValue(driver.ColDef, string) (interface{}, error)

	SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error)
	String() string
	OnApplyConditions()
}

type ConnectionScreen struct {
	*gtk.Paned
	*ResultView
	ctrl connectionScreenCtrl

	w           *Window
	hPaned      *gtk.Paned
	dbCombo     *gtk.ComboBoxText
	tableFilter *gtk.SearchEntry
	tableList   *List
	logview     *List

	databaseNames []string

	activeDatabase MVar

	tablesMenu struct {
		tableMenu      *gtk.Menu
		editMenu       *gtk.MenuItem
		newTabMenu     *gtk.MenuItem
		schemaMenu     *gtk.MenuItem
		truncateMenu   *gtk.MenuItem
		deleteMenu     *gtk.MenuItem
		refreshMenu    *gtk.MenuItem
		copySelectMenu *gtk.MenuItem
	}

	logMenu struct {
		logMenu   *gtk.Menu
		clearMenu *gtk.MenuItem
		copyMenu  *gtk.MenuItem
	}
}

func (c ConnectionScreen) Init(
	w *Window,
	ctrl connectionScreenCtrl,
) (*ConnectionScreen, error) {
	defer config.LogStart("ConnectionScreen.Init", nil)()

	var err error
	c.w = w
	c.ctrl = ctrl

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	c.Paned.SetWideHandle(true)
	c.Paned.Show()

	c.hPaned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	c.hPaned.SetWideHandle(true)
	c.hPaned.SetHExpand(true)
	c.hPaned.SetVExpand(true)
	c.hPaned.Show()

	// Sidebar

	sideBar, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	sideBar.Show()

	c.hPaned.Pack1(sideBar, false, true)

	c.tableFilter, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	c.tableFilter.Show()
	c.tableFilter.Connect("search-changed", c.onSearch)
	c.tableFilter.SetPlaceholderText("Filter table: .*")

	// TODO: figure out how to focus on accelerator
	//k, mod := gtk.AcceleratorParse("<Control>f")
	//c.tableFilter.AddAccelerator("activate", nil, k, mod, gtk.ACCEL_VISIBLE)

	tableListSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	tableListSW.Show()

	c.dbCombo, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	c.dbCombo.SetIDColumn(0)
	c.dbCombo.Show()

	c.tableList, err = List{}.Init(c.w, &ListOptions{
		SelectOnRightClick: true,
		IconFunc: func(name fmt.Stringer) (string, int) {
			def, ok := name.(driver.TableDef)
			if ok {
				if def.Type == driver.TableTypeRegular {
					return "table", ICON_SIZE_BUTTON
				} else {
					return "table-v", ICON_SIZE_BUTTON
				}
			}
			return "table", ICON_SIZE_BUTTON
		},
	}, ctrl)
	if err != nil {
		return nil, err
	}
	c.tableList.SetHExpand(true)
	c.tableList.SetVExpand(true)
	c.tableList.OnButtonPress(c.onTableListButtonPress)
	c.tableList.Show()

	tableListSW.Add(c.tableList)

	sideBar.PackStart(c.dbCombo, false, true, 0)
	sideBar.PackStart(c.tableFilter, false, true, 4)
	sideBar.PackStart(tableListSW, true, true, 0)

	// main section

	mainSection, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	mainSection.Show()

	c.ResultView, err = ResultView{}.Init(
		c.w,
		ctrl,
	)
	if err != nil {
		return nil, err
	}
	c.ResultView.Show()

	mainSection.Add(c.ResultView)
	mainSection.SetVExpand(true)
	mainSection.SetHExpand(true)

	c.hPaned.Pack2(mainSection, true, false)

	err = c.initTableMenu()
	if err != nil {
		return nil, err
	}

	err = c.initLogMenu()
	if err != nil {
		return nil, err
	}

	// signals

	logSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	logSW.Show()

	c.logview, err = List{}.Init(c.w, &ListOptions{SelectOnRightClick: true}, ctrl)
	if err != nil {
		return nil, err
	}
	c.logview.SetName("logger")
	c.logview.SetSizeRequest(-1, 30)
	c.logview.OnButtonPress(c.onLogViewButtonPress)

	logSW.Add(c.logview)
	c.Paned.Pack1(c.hPaned, false, false)
	c.Paned.Pack2(logSW, false, true)

	c.dbCombo.Connect("changed", c.onDatabaseSelected)
	c.dbCombo.Connect("changed", ctrl.OnDatabaseSelected)
	c.tableList.Connect("row-activated", ctrl.OnTableSelected)
	c.tablesMenu.schemaMenu.Connect("activate", ctrl.OnSchemaMenu)
	c.tablesMenu.refreshMenu.Connect("activate", ctrl.OnRefreshMenu)
	c.tablesMenu.newTabMenu.Connect("activate", ctrl.OnNewTabMenu)
	c.tablesMenu.editMenu.Connect("activate", ctrl.OnEditTable)
	c.tablesMenu.truncateMenu.Connect("activate", ctrl.OnTruncateTable)
	c.tablesMenu.deleteMenu.Connect("activate", ctrl.OnDeleteTable)
	c.tablesMenu.copySelectMenu.Connect("activate", ctrl.OnCopySelect)

	c.logMenu.clearMenu.Connect("activate", c.onClearLog)
	//c.logMenu.copyMenu.Connect("activate", ctrl.OnCopyLog)
	c.Show()

	return &c, nil
}

func (c *ConnectionScreen) Log(s string) {
	defer config.LogStart("ConnectionScreen.Log", nil)()

	c.logview.PrependItem(Stringer(s))
}

func (c *ConnectionScreen) SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error) {
	defer config.LogStart("ConnectionScreen.SetTableDef", nil)()

	return c.ctrl.SetTableDef(ctx, tableDef)
}

func (c *ConnectionScreen) SetQuery(query string) {
	defer config.LogStart("ConnectionScreen.SetQuery", nil)()

	c.ResultView.SetQuery(query)
}

func (c *ConnectionScreen) SetDatabases(dbs []string) {
	defer config.LogStart("ConnectionScreen.SetDatabases", nil)()

	c.databaseNames = dbs

	for _, name := range dbs {
		c.dbCombo.Append(name, name)
		//c.dbStore.SetValue(c.dbStore.Append(), 0, name)
	}
}

func (c *ConnectionScreen) SetTables(tables driver.TableDefs) {
	defer config.LogStart("ConnectionScreen.SetTables", nil)()

	c.tableList.UpdateItems(tables.ToStringer())
}

func (c *ConnectionScreen) SetActiveDatabase(dbName string) {
	defer config.LogStart("ConnectionScreen.SetActiveDatabase", nil)()

	for i, db := range c.databaseNames {
		if db == dbName {
			c.activeDatabase.Set(dbName)
			c.dbCombo.SetActive(i)
			return
		}
	}
}

func (c *ConnectionScreen) ActiveDatabase() (string, bool) {
	defer config.LogStart("ConnectionScreen.ActiveDatabase", nil)()

	if c.activeDatabase.Get() == nil {
		return "", false
	}
	return c.activeDatabase.Get().(string), true
}

func (c *ConnectionScreen) ActiveTable() (driver.TableDef, bool) {
	defer config.LogStart("ConnectionScreen.ActiveTable", nil)()

	i, ok := c.tableList.SelectedItem()
	if !ok {
		return driver.TableDef{}, false
	}

	return i.(driver.TableDef), true
}

func (c *ConnectionScreen) SetActiveTable(t driver.TableDef) {
	defer config.LogStart("ConnectionScreen.SetActiveTable", nil)()

	c.tableList.selectedItem.Set(t)
}

func (c *ConnectionScreen) ShowTableSchemaModal(tableName, schema string) {
	defer config.LogStart("ConnectionScreen.ShowTableSchemaModal", nil)()

	modal, err := gtk.DialogNewWithButtons(fmt.Sprintf("Table %s", tableName), c.w,
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

	sourceView, err := SourceView{}.Init(c.w, SourceViewOptions{
		Highlight: true,
		Undoable:  true,
		Language:  "sql",
	}, c.ctrl)
	if err != nil {
		return
	}
	sourceView.SetVExpand(true)
	sourceView.SetHExpand(true)

	buff, err := sourceView.GetBuffer()
	if err != nil {
		return
	}

	buff.InsertMarkup(buff.GetStartIter(), schema)

	sourceView.Show()
	content.Add(sourceView)

	modal.Run()
	modal.Destroy()
}

func (c *ConnectionScreen) SelectedTable() (string, bool) {
	defer config.LogStart("ConnectionScreen.SelectedTable", nil)()

	i, ok := c.tableList.SelectedItem()
	if !ok {
		return "", ok
	}
	return i.String(), true
}

func (c *ConnectionScreen) onTableListButtonPress(_ *gtk.ListBox, e *gdk.Event) bool {
	defer config.LogStart("ConnectionScreen.onTableListButtonPress", nil)()

	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return false
	}

	c.tablesMenu.tableMenu.Show()
	c.tablesMenu.tableMenu.PopupAtPointer(e)
	return true
}

func (c *ConnectionScreen) onLogViewButtonPress(_ *gtk.ListBox, e *gdk.Event) bool {
	defer config.LogStart("ConnectionScreen.onLogViewButtonPress", nil)()

	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return false
	}

	c.logMenu.logMenu.Show()
	c.logMenu.logMenu.PopupAtPointer(e)
	return true
}

func (c *ConnectionScreen) onDatabaseSelected() {
	defer config.LogStart("ConnectionScreen.onDatabaseSelected", nil)()

	dbName := c.dbCombo.GetActiveText()
	c.activeDatabase.Set(dbName)
}

func (c *ConnectionScreen) initTableMenu() error {
	defer config.LogStart("ConnectionScreen.initTableMenu", nil)()

	var err error
	c.tablesMenu.tableMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.tablesMenu.editMenu, err = BWMenuItemWithImage("Edit", "edit-table")
	if err != nil {
		return err
	}

	c.tablesMenu.newTabMenu, err = BWMenuItemWithImage("New tab", "add-tab")
	if err != nil {
		return err
	}

	c.tablesMenu.schemaMenu, err = BWMenuItemWithImage("Schema", "config")
	if err != nil {
		return err
	}

	c.tablesMenu.truncateMenu, err = BWMenuItemWithImage("Truncate", "truncate")
	if err != nil {
		return err
	}

	c.tablesMenu.deleteMenu, err = BWMenuItemWithImage("Delete", "delete-table")
	if err != nil {
		return err
	}

	c.tablesMenu.refreshMenu, err = BWMenuItemWithImage("Refresh", "refresh")
	if err != nil {
		return err
	}

	c.tablesMenu.copySelectMenu, err = BWMenuItemWithImage("Copy SELECT", "copy")
	if err != nil {
		return err
	}

	cowboy, err := BWMenuItemWithImage("Cowboy", "cowboy")
	if err != nil {
		return err
	}
	c.tablesMenu.tableMenu.Add(c.tablesMenu.newTabMenu)
	c.tablesMenu.tableMenu.Add(c.tablesMenu.copySelectMenu)
	c.tablesMenu.tableMenu.Add(c.tablesMenu.schemaMenu)
	c.tablesMenu.tableMenu.Add(c.tablesMenu.editMenu)
	c.tablesMenu.tableMenu.Add(c.tablesMenu.refreshMenu)
	c.tablesMenu.tableMenu.Add(cowboy)

	cowboyMenu, err := gtk.MenuNew()
	if err != nil {
		return err
	}

	cowboyMenu.Add(c.tablesMenu.truncateMenu)
	cowboyMenu.Add(c.tablesMenu.deleteMenu)
	cowboy.SetSubmenu(cowboyMenu)
	cowboy.Show()

	return nil
}

func (c *ConnectionScreen) initLogMenu() error {
	defer config.LogStart("ConnectionScreen.initLogMenu", nil)()

	var err error
	c.logMenu.logMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.logMenu.clearMenu, err = BWMenuItemWithImage("Clear", "truncate")
	if err != nil {
		return err
	}

	//c.logMenu.copyMenu, err = menuItemWithImage("Copy", "gtk-copy")
	//if err != nil {
	//return err
	//}

	c.logMenu.logMenu.Add(c.logMenu.clearMenu)
	//c.logMenu.logMenu.Add(c.logMenu.copyMenu)

	return nil
}

func (c *ConnectionScreen) onSearch(e *gtk.SearchEntry) {
	defer config.LogStart("ConnectionScreen.onSearch", nil)()

	buff, err := e.GetBuffer()
	if err != nil {
		return
	}

	txt, err := buff.GetText()
	if err != nil {
		return
	}

	rg, err := regexp.Compile(strings.ToLower(txt))
	if err != nil {
		rg = regexp.MustCompile(fmt.Sprintf(".*%s.*", regexp.QuoteMeta(txt)))
	}
	c.tableList.SetFilterRegex(rg)
}

func (c *ConnectionScreen) onClearLog() {
	defer config.LogStart("ConnectionScreen.onClearLog", nil)()

	c.logview.Clear()
}
