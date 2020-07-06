package gtk

import (
	"fmt"
	"regexp"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/assets"
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type tab struct {
	*ConnectionTab
	ctrl interface {
		OnTabRemove()
		SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error)
		SetQuery(ctx *sqlengine.Context, query string) (bool, error)
		String() string
	}
}
type connectionScreenCtrl interface {
	OnDatabaseSelected()
	OnTableSelected()
	OnSchemaMenu()
	OnRefreshMenu()
	OnNewTabMenu()
	OnEditTable()
	OnTruncateTable()
	OnDeleteTable()
	OnCopySelect()
	OnCopyLog()
	Config() *config.Config
}

type ConnectionScreen struct {
	*gtk.Paned
	ctrl connectionScreenCtrl

	w           *Window
	hPaned      *gtk.Paned
	dbCombo     *gtk.ComboBoxText
	tableFilter *gtk.SearchEntry
	tableList   *List
	tabber      *gtk.Notebook
	tabs        []tab
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

	// tab switching
	tabIndex int
}

func (c ConnectionScreen) Init(
	w *Window,
	ctrl connectionScreenCtrl,
) (*ConnectionScreen, error) {
	var err error
	c.w = w
	c.ctrl = ctrl

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	c.Paned.SetWideHandle(true)

	c.hPaned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	c.hPaned.SetWideHandle(true)
	c.hPaned.SetHExpand(true)
	c.hPaned.SetVExpand(true)

	// Sidebar

	sideBar, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	c.hPaned.Pack1(sideBar, false, true)

	c.tableFilter, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	c.tableFilter.Connect("search-changed", c.onSearch)
	c.tableFilter.SetPlaceholderText("Filter table: .*")

	// TODO: figure out how to focus on accelerator
	//k, mod := gtk.AcceleratorParse("<Control>f")
	//c.tableFilter.AddAccelerator("activate", nil, k, mod, gtk.ACCEL_VISIBLE)

	tableListSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	c.dbCombo, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	c.dbCombo.SetIDColumn(0)

	c.tableList, err = List{}.Init(c.w, &ListOptions{
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
	}, ctrl)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	c.tabber, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}

	switch c.ctrl.Config().GUI.TableTabPosition.String() {
	case "bottom":
		c.tabber.SetProperty("tab-pos", gtk.POS_BOTTOM)
	default:
		c.tabber.SetProperty("tab-pos", gtk.POS_TOP)
	}

	c.tabber.SetVExpand(true)
	c.tabber.SetHExpand(true)
	// TODO: move to config
	//c.tabber.SetProperty("tab-pos", gtk.POS_BOTTOM)
	c.tabber.SetProperty("scrollable", true)
	c.tabber.SetProperty("enable-popup", false)
	c.tabber.Connect("switch-page", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		c.tabIndex = i
		if i >= len(c.tabs) {
			return
		}

		c.dbCombo.SetActiveID(c.tabs[i].database)
	})

	c.tabber.Connect("page-removed", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		c.tabs = append(c.tabs[:i], c.tabs[i+1:]...)
	})

	c.tabber.Connect("page-reordered", c.onTabReorder)

	mainSection.Add(c.tabber)
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

	c.logview, err = List{}.Init(c.w, &ListOptions{SelectOnRightClick: true}, ctrl)
	if err != nil {
		return nil, err
	}
	c.logview.SetName("logger")
	c.logview.SetSizeRequest(-1, 30)
	c.logview.OnButtonPress(c.onLogViewButtonPress)
	//c.logview.SetEditable(false)
	//c.logview.SetPixelsAboveLines(5)
	//c.logview.SetPixelsBelowLines(5)

	logSW.Add(c.logview)
	c.Paned.Pack1(c.hPaned, false, false)
	c.Paned.Pack2(logSW, false, true)

	c.Paned.ShowAll()
	c.hPaned.ShowAll()

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

	return &c, nil
}

func (c *ConnectionScreen) onTabReorder(_ *gtk.Notebook, _ *gtk.Widget, landing int) {
	// https://play.golang.com/p/YMfQouxHuvr
	movingTab := c.tabs[c.tabIndex]

	c.tabs = append(c.tabs[:c.tabIndex], c.tabs[c.tabIndex+1:]...)

	lh := make([]tab, len(c.tabs[:landing]))
	rh := make([]tab, len(c.tabs[landing:]))
	copy(lh, c.tabs[:landing])
	copy(rh, c.tabs[landing:])

	c.tabs = append(lh, movingTab)
	c.tabs = append(c.tabs, rh...)
}

func (c *ConnectionScreen) Log(s string) {
	c.logview.PrependItem(Stringer(s))
}

// TODO: This sucks
func (c *ConnectionScreen) CtrlMod() bool {
	return c.tableList.CtrlMod()
}

func (c *ConnectionScreen) CurrentTabIndex() int {
	return c.tabber.GetCurrentPage()
}

func (c *ConnectionScreen) SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error) {
	if c.tabber.GetCurrentPage() == -1 {
		return false, nil
	}

	return c.tabs[c.tabber.GetCurrentPage()].ctrl.SetTableDef(ctx, tableDef)
}

func (c *ConnectionScreen) SetQuery(ctx *sqlengine.Context, query string) (bool, error) {
	if c.tabber.GetCurrentPage() == -1 {
		return false, nil
	}
	return c.tabs[c.tabber.GetCurrentPage()].ctrl.SetQuery(ctx, query)
}

func (c *ConnectionScreen) AddTab(
	ctab *ConnectionTab,
	ctrl interface {
		OnTabRemove()
		SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error)
		SetQuery(ctx *sqlengine.Context, query string) (bool, error)
		String() string
	},
	switchNow bool,
) error {
	c.tabber.AppendPage(ctab.content, ctab.header)
	c.tabber.SetTabReorderable(ctab.content, true)

	ctab.btn.Connect("clicked", func() {
		index := c.tabber.PageNum(ctab.content)
		if index == -1 {
			return
		}

		c.tabber.RemovePage(index)
		ctrl.OnTabRemove()
	})

	if switchNow {
		c.tabber.SetCurrentPage(c.tabber.GetNPages() - 1)
	}

	c.tabs = append(c.tabs, tab{ConnectionTab: ctab, ctrl: ctrl})

	return nil
}

func (c *ConnectionScreen) Close() bool {
	if c.tabber.GetCurrentPage() == -1 {
		return false
	}

	c.tabs[c.tabber.GetCurrentPage()].btn.Emit("clicked")
	return true
}

func (c *ConnectionScreen) SetDatabases(dbs []string) {
	c.databaseNames = dbs

	for _, name := range dbs {
		c.dbCombo.Append(name, name)
		//c.dbStore.SetValue(c.dbStore.Append(), 0, name)
	}
}

func (c *ConnectionScreen) SetTables(tables driver.TableDefs) {
	c.tableList.UpdateItems(tables.ToStringer())
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

	schema, err = ChromaHighlight(c.ctrl.Config().EditorTheme(), schema)
	if err != nil {
		c.ctrl.Config().Error(err)
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
	modal.Destroy()
}

func (c *ConnectionScreen) SelectedTable() (string, bool) {
	i, ok := c.tableList.SelectedItem()
	if !ok {
		return "", ok
	}
	return i.String(), true
}

func (c *ConnectionScreen) onTableListButtonPress(_ *gtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	c.tablesMenu.tableMenu.ShowAll()
	c.tablesMenu.tableMenu.PopupAtPointer(e)
}

func (c *ConnectionScreen) onLogViewButtonPress(_ *gtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	c.logMenu.logMenu.ShowAll()
	c.logMenu.logMenu.PopupAtPointer(e)
}

func (c *ConnectionScreen) onDatabaseSelected() {
	dbName := c.dbCombo.GetActiveText()
	c.activeDatabase.Set(dbName)
}

func (c *ConnectionScreen) initTableMenu() error {
	var err error
	c.tablesMenu.tableMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.tablesMenu.editMenu, err = menuItemWithImage("Edit", "gtk-edit")
	if err != nil {
		return err
	}

	c.tablesMenu.newTabMenu, err = menuItemWithImage("New tab", "gtk-new")
	if err != nil {
		return err
	}

	c.tablesMenu.schemaMenu, err = menuItemWithImage("Schema", "gtk-info")
	if err != nil {
		return err
	}

	c.tablesMenu.truncateMenu, err = menuItemWithImage("Truncate", "gtk-clear")
	if err != nil {
		return err
	}

	c.tablesMenu.deleteMenu, err = menuItemWithImage("Delete", "gtk-delete")
	if err != nil {
		return err
	}

	c.tablesMenu.refreshMenu, err = menuItemWithImage("Refresh", "gtk-refresh")
	if err != nil {
		return err
	}

	c.tablesMenu.copySelectMenu, err = menuItemWithImage("Copy SELECT", "gtk-copy")
	if err != nil {
		return err
	}

	cowboy, err := menuItemWithImage("Cowboy", "gtk-delete")
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
	cowboy.ShowAll()

	return nil
}

func (c *ConnectionScreen) initLogMenu() error {
	var err error
	c.logMenu.logMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.logMenu.clearMenu, err = menuItemWithImage("Clear", "gtk-clear")
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

func (c *ConnectionScreen) onClearLog() {
	c.logview.Clear()
}

type ConnectionTab struct {
	label   *gtk.Label
	btn     *gtk.Button
	content *gtk.Box
	header  *gtk.Box

	//index    int
	database string
}

type ConnectionTabOpts struct {
	Database string
	Title    string
	Content  gtk.IWidget
}

func (c ConnectionTab) Init(opts ConnectionTabOpts) (*ConnectionTab, error) {
	var err error

	c.database = opts.Database
	c.content, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	c.content.PackStart(opts.Content, true, true, 0)
	c.content.SetVExpand(true)
	c.content.SetHExpand(true)
	c.content.Show()

	c.header, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	c.label, err = gtk.LabelNew(opts.Title)
	if err != nil {
		return nil, err
	}

	image, err := gtk.ImageNewFromIconName("window-close", gtk.ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}

	c.btn, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}

	c.btn.SetImage(image)
	c.btn.SetRelief(gtk.RELIEF_NONE)

	c.header.PackStart(c.label, true, true, 0)
	c.header.PackEnd(c.btn, false, false, 0)
	c.header.ShowAll()

	return &c, nil
}

func (c *ConnectionTab) SetTitle(title string) {
	c.label.SetText(title)
}
