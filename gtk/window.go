package gtk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
)

// keeps track of tabs on all windows
var tabs = map[string][]*ToolTab{}

// hold the tab being moved from one window to another.
// it really doesn't know whether it was just removed or if it's being moved to another window
var transit = struct {
	tab *ToolTab
	src *gtk.Notebook
}{}

type windowCtrl interface {
	OnNewDatabaseTab()
	OnNewHTTPTab()
	OnCloseTab(string)
	OnSaveEnv(*config.Env) error
	OnDeleteEnv(*config.Env) error
}

type Window struct {
	*gtk.ApplicationWindow
	id          string
	nb          *gtk.Notebook
	box         *gtk.Box // holds nb and statusbar
	statusBar   *gtk.Statusbar
	statusBarID uint

	envCb    *gtk.ComboBox
	envModel *gtk.ListStore
	btnEnv   *gtk.Button

	Menu struct {
		NewConnection  *glib.SimpleAction
		NewDatabaseTab *glib.SimpleAction
		NewHTTPTab     *glib.SimpleAction
		Close          *glib.SimpleAction
	}
	ctrl windowCtrl

	tabIndex int
}

func (w Window) Init(app *gtk.Application, ctrl windowCtrl) (*Window, error) {
	defer config.LogStart("Window.Init", nil)()

	var err error

	w.id = uuid.New().String()
	w.ApplicationWindow, err = gtk.ApplicationWindowNew(app)
	w.SetTitle("BenchWell")
	w.SetSizeRequest(1432, 867)
	w.ctrl = ctrl
	w.ApplicationWindow.Window.Connect("focus-in-event", func() {
		config.ActiveWindow = w.ApplicationWindow
	})
	config.ActiveWindow = w.ApplicationWindow

	w.nb, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	w.nb.Show()
	w.nb.SetProperty("scrollable", true)
	w.nb.SetName("MainNotebook")
	w.nb.SetGroupName("MainWindow")
	w.nb.PopupEnable()

	w.nb.Connect("switch-page", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		w.tabIndex = i
	})
	w.nb.Connect("page-removed", w.onPageRemoved)
	w.nb.Connect("page-added", w.onPageAdded)
	w.nb.Connect("page-reordered", w.onTabReorder)

	switch config.GUI.TabPosition.String() {
	case "bottom":
		w.nb.SetProperty("tab-pos", gtk.POS_BOTTOM)
	default:
		w.nb.SetProperty("tab-pos", gtk.POS_TOP)
	}

	w.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	w.box.Show()
	w.box.SetVExpand(true)
	w.box.SetHExpand(true)

	w.statusBar, err = gtk.StatusbarNew()
	if err != nil {
		return nil, err
	}
	w.statusBar.Show()

	w.box.PackStart(w.nb, true, true, 0)
	w.box.PackEnd(w.statusBar, false, false, 0)

	w.statusBarID = w.statusBar.GetContextId("main")

	w.ApplicationWindow.Add(w.box)
	header, err := w.headerMenu()
	if err != nil {
		return nil, err
	}
	header.Show()
	w.SetTitlebar(header)

	w.Show()
	// TODO: when we get a systray
	//w.HideOnDelete()

	// add main tab
	w.Menu.NewDatabaseTab.Connect("activate", ctrl.OnNewDatabaseTab)
	w.Menu.NewHTTPTab.Connect("activate", ctrl.OnNewHTTPTab)
	w.Menu.Close.Connect("activate", func() {
		tab := w.CurrentTab()
		if tab == nil {
			return
		}
		ctrl.OnCloseTab(tab.id)
	})

	tabs[w.id] = []*ToolTab{}

	return &w, nil
}

func (w *Window) AddToolTab(tab *ToolTab) error {
	defer config.LogStart("Window.AddToolTab", nil)()

	w.nb.AppendPage(tab.Content(), tab.Header())

	tab.SetTitle(tab.tabCtrl.Title())
	w.nb.ChildSetProperty(tab.Content(), "tab-fill", false)
	w.nb.SetMenuLabelText(tab.Content(), tab.Title())
	w.nb.SetTabReorderable(tab.Content(), true)
	w.nb.SetCurrentPage(w.nb.PageNum(tab.Content()))
	tabs[w.id] = append(tabs[w.id], tab)
	w.nb.SetTabDetachable(tab.Content(), true)

	return nil
}

func (w *Window) onPageRemoved(_ *gtk.Notebook, _ *gtk.Widget, i int) {
	defer config.LogStart("Window.onPageRemoved", nil)()
	transit.tab = tabs[w.id][i]
	transit.src = w.nb
	tabs[w.id] = append(tabs[w.id][:i], tabs[w.id][i+1:]...)
}

func (w *Window) onPageAdded(_ *gtk.Notebook, c *gtk.Widget, i int) {
	defer config.LogStart("Window.onPageAdded", nil)()
	if transit.tab == nil {
		return
	}
	if transit.tab.w.id == w.id {
		transit.tab = nil
		transit.src = nil
		return
	}

	if i > len(tabs[w.id]) {
		tabs[w.id] = append(tabs[w.id], transit.tab)
	} else {
		rest := make([]*ToolTab, len(tabs[w.id][i:]))
		copy(rest, tabs[w.id][i:])
		tabs[w.id] = append(tabs[w.id][:i], transit.tab)
		tabs[w.id] = append(tabs[w.id], rest...)
	}

	transit.tab.SetWindowCtrl(w.ctrl)
	transit.tab.w = w
	transit.tab = nil
	transit.src = nil
}

func (w *Window) RemovePage(id string) {
	defer config.LogStart("Window.RemovePage", nil)()
	for i, tab := range tabs[w.id] {
		if tab.id != id {
			continue
		}

		w.nb.RemovePage(i)
		break
	}
}

func (w *Window) CurrentPage() int {
	defer config.LogStart("Window.CurrentPage", nil)()

	return w.nb.GetCurrentPage()
}

func (w *Window) CurrentTab() *ToolTab {
	defer config.LogStart("Window.CurrentTab", nil)()

	if len(tabs[w.id]) == 0 {
		return nil
	}
	return tabs[w.id][w.CurrentPage()]
}

func (w *Window) TabByID(id string) *ToolTab {
	for _, tab := range tabs[w.id] {
		if tab.id == id {
			return tab
		}
	}
	return nil
}

func (w *Window) Remove(wd gtk.IWidget) {
	defer config.LogStart("Window.Remove", nil)()

	w.nb.Remove(wd)
}

func (w Window) PushStatus(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05")}, args...)
	w.statusBar.Push(w.statusBarID, fmt.Sprintf("[%s] "+format, args...))
}

func (w *Window) OnPageRemoved(f interface{}) {
	defer config.LogStart("Window.OnPageRemoved", nil)()

	w.nb.Connect("page-removed", f)
}

func (w *Window) PageCount() int {
	defer config.LogStart("Window.PageCount", nil)()

	return w.nb.GetNPages()
}

func (w *Window) headerMenu() (*gtk.HeaderBar, error) {
	defer config.LogStart("Window.headerMenu", nil)()

	w.Menu.NewDatabaseTab = glib.SimpleActionNew("new.db", nil)
	w.Menu.NewHTTPTab = glib.SimpleActionNew("new.http", nil)
	w.Menu.Close = glib.SimpleActionNew("close", nil)
	w.AddAction(w.Menu.NewDatabaseTab)
	w.AddAction(w.Menu.NewHTTPTab)
	w.AddAction(w.Menu.Close)

	header, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	header.SetShowCloseButton(true)
	header.SetTitle("BenchWell")
	header.SetSubtitle(config.Version)

	// Create a new window menu button
	windowBtnMenu, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	windowBtnMenu.Show()
	addImg, err := BWImageNewFromFile("add-tab", "orange", ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}
	addImg.Show()

	windowBtnMenu.SetImage(addImg)

	// Set up the menu model for the button
	windowMenu := glib.MenuNew()
	if windowMenu == nil {
		return nil, errors.New("nil menu")
	}
	windowBtnMenu.SetMenuModel(&windowMenu.MenuModel)

	windowMenu.Append("Window", "app.new")
	windowMenu.Append("Database", "win.new.db")
	windowMenu.Append("HTTP", "win.new.http")

	// Create a new app menu button
	appBtnMenu, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	appBtnMenu.Show()

	prefImg, err := BWImageNewFromFile("config", "orange", ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}
	prefImg.Show()

	appBtnMenu.SetImage(prefImg)

	// Set up the menu model for the button
	appMenu := glib.MenuNew()
	if appMenu == nil {
		return nil, errors.New("nil menu")
	}
	appBtnMenu.SetMenuModel(&appMenu.MenuModel)

	appMenu.Append("Preferences", "app.preferences")
	appMenu.Append("Dark toggle", "app.darkmode")

	env, err := w.envselector()
	if err != nil {
		return nil, err
	}
	env.Show()

	// add the menu button to the header
	header.PackStart(windowBtnMenu)
	header.PackEnd(appBtnMenu)
	header.PackEnd(env)

	// Assemble the window
	return header, nil
}

func (w *Window) Go(job func(context.Context) func()) func() {
	defer config.LogStart("Window.Go", nil)()

	cancel := make(chan struct{}, 0)
	done := make(chan struct{}, 0)
	onDone := func() {}
	ctxCancel := func() {}

	go func() {
		go func() {
			var ctx context.Context
			ctx, ctxCancel = context.WithCancel(context.Background())
			onDone = job(ctx)
			close(done)
		}()

		select {
		case <-done:
			if onDone != nil {
				_, err := glib.IdleAdd(onDone)
				if err != nil {
					log.Fatal("IdleAdd() failed:", err)
				}
			}
		case <-cancel:
		}
	}()

	return func() {
		ctxCancel()
		close(cancel)
	}
}

func (w *Window) onTabReorder(_ *gtk.Notebook, _ *gtk.Widget, landing int) {
	defer config.LogStart("Window.onTabReorder", nil)()

	// https://play.golang.com/p/YMfQouxHuvr
	movingTab := tabs[w.id][w.tabIndex]

	tabs[w.id] = append(tabs[w.id][:w.tabIndex], tabs[w.id][w.tabIndex+1:]...)

	lh := make([]*ToolTab, len(tabs[w.id][:landing]))
	rh := make([]*ToolTab, len(tabs[w.id][landing:]))
	copy(lh, tabs[w.id][:landing])
	copy(rh, tabs[w.id][landing:])

	tabs[w.id] = append(lh, movingTab)
	tabs[w.id] = append(tabs[w.id], rh...)
}

func (w *Window) envselector() (*gtk.Grid, error) {
	var err error

	w.envModel, err = gtk.ListStoreNew(glib.TYPE_INT64, glib.TYPE_STRING)
	if err != nil {
		return nil, err
	}
	w.envCb, err = gtk.ComboBoxNewWithModel(w.envModel)
	if err != nil {
		return nil, err
	}
	w.envCb.Show()
	w.envCb.SetProperty("id-column", 0)
	w.envCb.SetIDColumn(0)
	w.envCb.SetEntryTextColumn(1)
	area, err := w.envCb.GetProperty("cell-area")
	fmt.Println("=====", area.(*gtk.CellArea), err)

	for _, env := range config.Environments {
		fmt.Println("=====", env.Name)
		iter := w.envModel.Append()
		//w.envModel.SetValue(iter, 0, env.ID)
		w.envModel.SetValue(iter, 1, env.Name)
		//w.envCb.Append(fmt.Sprintf("%d", env.ID), env.Name)
	}

	w.btnEnv, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	w.btnEnv.Show()
	prefImg, err := BWImageNewFromFile("config", "orange", ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}
	w.btnEnv.SetImage(prefImg)

	popover, err := gtk.PopoverNew(w.btnEnv)
	if err != nil {
		return nil, err
	}

	w.btnEnv.Connect("clicked", func() {
		popover.Show()
	})

	envpop, err := EnvironmentsPopover{}.Init(w, w.ctrl)
	if err != nil {
		return nil, err
	}
	envpop.Show()
	popover.Add(envpop)

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	grid.Attach(w.envCb, 0, 0, 4, 1)
	grid.Attach(w.btnEnv, 4, 0, 1, 1)
	BWAddClass(grid, "linked")

	return grid, nil
}

type environmenctrl interface {
	OnSaveEnv(*config.Env) error
	OnDeleteEnv(*config.Env) error
}

type EnvironmentsPopover struct {
	*gtk.Box
	w      *Window
	btnAdd *gtk.Button
	envs   []*EnvironmentPanel

	stack *gtk.Stack

	environmenctrl
}

func (e EnvironmentsPopover) Init(w *Window, ctrl environmenctrl) (*EnvironmentsPopover, error) {
	var err error
	e.environmenctrl = ctrl
	e.w = w

	e.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	header, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	header.Show()

	e.btnAdd, err = BWButtonNewFromIconName("add", "white", ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}
	e.btnAdd.Show()
	BWAddClass(e.btnAdd, "suggested-action")
	e.btnAdd.Connect("clicked", func() {
		env := &config.Env{Name: "new"}
		err := ctrl.OnSaveEnv(env)
		if err != nil {
			e.w.PushStatus("saving environment: %s", err.Error())
			return
		}

		envpanel, err := EnvironmentPanel{}.Init(env, &e)
		if err != nil {
			e.w.PushStatus("error creating EnvironmentPanel: %s", err.Error())
			return
		}
		envpanel.Show()
		//w.envCb.Append(fmt.Sprintf("%d", env.ID), env.Name)
		//w.envCb.SetActiveID(fmt.Sprintf("%d", env.ID))

		config.Environments = append(config.Environments, env)
		e.stack.AddTitled(envpanel, fmt.Sprintf("%d", len(config.Environments)), env.Name)
		e.stack.SetVisibleChildName(env.Name)
	})
	header.PackEnd(e.btnAdd)

	paned, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	paned.Show()

	adapterStack, err := gtk.StackSwitcherNew()
	if err != nil {
		return nil, err
	}
	adapterStack.SetOrientation(gtk.ORIENTATION_VERTICAL)
	adapterStack.Show()
	adapterStack.SetVExpand(true)
	adapterStack.SetHExpand(true)

	e.stack, err = gtk.StackNew()
	if err != nil {
		return nil, err
	}
	e.stack.Show()
	e.stack.SetHomogeneous(true)

	for i, env := range config.Environments {
		envpanel, err := EnvironmentPanel{}.Init(env, &e)
		if err != nil {
			return nil, err
		}
		envpanel.Show()

		e.stack.AddTitled(envpanel, fmt.Sprintf("%d", i), env.Name)
		e.envs = append(e.envs, envpanel)
	}

	if len(config.Environments) > 0 {
		e.stack.SetVisibleChildName("0")
	}
	e.stack.SetVExpand(true)
	e.stack.SetHExpand(true)

	adapterStack.SetStack(e.stack)

	paned.Pack1(adapterStack, false, true)
	paned.Pack2(e.stack, true, false)

	e.PackStart(header, true, true, 0)
	e.PackStart(paned, false, false, 0)

	return &e, nil
}

func (e *EnvironmentsPopover) OnSaveEnv(env *config.Env) error {
	err := e.environmenctrl.OnSaveEnv(env)
	if err != nil {
		e.w.PushStatus("Fail to save env: %s", err.Error())
		return err
	}

	for _, panel := range e.envs {
		if panel.env != env {
			continue
		}

		e.stack.ChildSetProperty(panel, "title", env.Name)
		//e.w.envCb.Remov
		break
	}

	return nil
}

func (e *EnvironmentsPopover) OnDeleteEnv(env *config.Env) error {
	err := e.environmenctrl.OnDeleteEnv(env)
	if err != nil {
		e.w.PushStatus("Fail to delete env: %s", err.Error())
		return err
	}

	for _, panel := range e.envs {
		if panel.env != env {
			continue
		}

		e.stack.Remove(panel)

		//e.w.envCb.Remove(i)
		break
	}

	return nil
}

type EnvironmentPanel struct {
	*gtk.Box
	env       *config.Env
	name      *gtk.Entry
	btnSave   *gtk.Button
	btnRemove *gtk.Button
}

func (e EnvironmentPanel) Init(env *config.Env, ctrl environmenctrl) (*EnvironmentPanel, error) {
	var err error

	e.env = env
	e.name, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	e.name.Show()
	e.name.SetText(env.Name)
	e.name.SetPlaceholderText("Name")

	e.btnSave, err = gtk.ButtonNewWithLabel("Save")
	if err != nil {
		return nil, err
	}
	e.btnSave.Show()
	e.btnSave.Connect("clicked", func() {
		env.Name, _ = e.name.GetText()
		err := ctrl.OnSaveEnv(e.env)
		if err != nil {
			return
		}
	})

	e.btnRemove, err = gtk.ButtonNewWithLabel("Remove")
	if err != nil {
		return nil, err
	}
	e.btnRemove.Show()
	BWAddClass(e.btnRemove, "destructive-action")
	e.btnRemove.Connect("clicked", func() {
		err := ctrl.OnDeleteEnv(e.env)
		if err != nil {
			return
		}
	})

	btnBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	btnBox.Show()
	btnBox.PackEnd(e.btnRemove, false, false, 5)
	btnBox.PackEnd(e.btnSave, false, false, 5)

	e.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	vbox.Show()

	vv, err := KeyValues{}.Init(&env.Variables, func() {})
	if err != nil {
		return nil, err
	}
	vv.Show()

	vbox.PackStart(e.name, false, false, 5)
	vbox.PackStart(vv, true, true, 5)
	vbox.PackEnd(btnBox, false, false, 5)
	e.PackStart(vbox, true, true, 5)

	return &e, nil
}
