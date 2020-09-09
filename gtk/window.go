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
}

type Window struct {
	*gtk.ApplicationWindow
	id          string
	nb          *gtk.Notebook
	box         *gtk.Box // holds nb and statusbar
	statusBar   *gtk.Statusbar
	statusBarID uint

	Menu struct {
		NewConnection  *glib.SimpleAction
		NewDatabaseTab *glib.SimpleAction
		NewHTTPTab     *glib.SimpleAction
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
	w.nb.SetProperty("scrollable", true)
	w.nb.SetName("MainNotebook")
	w.nb.SetGroupName("MainWindow")
	w.nb.PopupEnable()

	w.nb.Connect("switch-page", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		w.tabIndex = i
	})
	w.nb.Connect("page-removed", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		transit.tab = tabs[w.id][i]
		transit.src = w.nb
		tabs[w.id] = append(tabs[w.id][:i], tabs[w.id][i+1:]...)
	})
	w.nb.Connect("page-added", func(_ *gtk.Notebook, c *gtk.Widget, i int) {
		if transit.tab == nil {
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

		transit.tab.SetWindowCtrl(ctrl)
		transit.tab = nil
		transit.src = nil
	})

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
	w.box.SetVExpand(true)
	w.box.SetHExpand(true)

	w.statusBar, err = gtk.StatusbarNew()
	if err != nil {
		return nil, err
	}

	w.box.PackStart(w.nb, true, true, 0)
	w.box.PackEnd(w.statusBar, false, false, 0)

	w.statusBarID = w.statusBar.GetContextId("main")

	w.ApplicationWindow.Add(w.box)

	header, err := w.headerMenu()
	if err != nil {
		return nil, err
	}
	w.SetTitlebar(header)

	w.ShowAll()
	// TODO: when we get a systray
	//w.HideOnDelete()

	// add main tab
	w.Menu.NewDatabaseTab.Connect("activate", ctrl.OnNewDatabaseTab)
	w.Menu.NewHTTPTab.Connect("activate", ctrl.OnNewHTTPTab)

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

func (w *Window) RemovePage(id string) {
	defer config.LogStart("Window.RemoveCurrentPage", nil)()
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
	w.AddAction(w.Menu.NewDatabaseTab)
	w.AddAction(w.Menu.NewHTTPTab)

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
	addImg, err := BWImageNewFromFile("add-tab", "orange", ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}

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
	prefImg, err := BWImageNewFromFile("config", "orange", ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}

	appBtnMenu.SetImage(prefImg)

	// Set up the menu model for the button
	appMenu := glib.MenuNew()
	if appMenu == nil {
		return nil, errors.New("nil menu")
	}
	appBtnMenu.SetMenuModel(&appMenu.MenuModel)

	appMenu.Append("Preferences", "app.preferences")
	appMenu.Append("Dark toggle", "app.darkmode")

	// add the menu button to the header
	header.PackStart(windowBtnMenu)
	header.PackEnd(appBtnMenu)

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
