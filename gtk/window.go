package gtk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
)

type windowCtrl interface {
	OnNewTab()
	OnNewSubTab()
	OnCloseTab()
	OnFileSelected(string)
	OnSaveQuery(string, string)
	Config() *config.Config
}

type Window struct {
	*gtk.ApplicationWindow
	nb          *gtk.Notebook
	box         *gtk.Box // holds nb and statusbar
	statusBar   *gtk.Statusbar
	statusBarID uint

	Menu struct {
		NewConnection *glib.SimpleAction
		NewToolTab    *glib.SimpleAction
		NewSubToolTab *glib.SimpleAction
		LoadFile      *glib.SimpleAction
		SaveQuery     *glib.SimpleAction
		CloseToolTab  *glib.SimpleAction
	}
	ctrl windowCtrl

	tabs     []*ToolTab
	tabIndex int
}

func (w Window) Init(app *gtk.Application, ctrl windowCtrl) (*Window, error) {
	var err error
	w.ApplicationWindow, err = gtk.ApplicationWindowNew(app)
	w.SetTitle("BenchWell")
	w.SetSizeRequest(1024, 768)
	w.ctrl = ctrl

	w.nb, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	w.nb.SetName("MainNotebook")

	w.nb.Connect("switch-page", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		w.tabIndex = i
	})
	w.nb.Connect("page-removed", func(_ *gtk.Notebook, _ *gtk.Widget, i int) {
		w.tabs = append(w.tabs[:i], w.tabs[i+1:]...)
	})

	w.nb.Connect("page-reordered", w.onTabReorder)

	switch w.ctrl.Config().GUI.ConnectionTabPosition.String() {
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
	w.Menu.NewToolTab.Connect("activate", ctrl.OnNewTab)
	// action menu for sub nb
	w.Menu.NewSubToolTab.Connect("activate", ctrl.OnNewSubTab)
	w.Menu.CloseToolTab.Connect("activate", ctrl.OnCloseTab)
	w.Menu.LoadFile.Connect("activate", w.OnOpenFile(ctrl.OnFileSelected))
	//w.Menu.SaveQuery.Connect("activate", w.OnSaveQuery(ctrl.OnSaveQuery))

	return &w, nil
}

func (w *Window) OnOpenFile(f func(string)) func() {
	return func() {
		openfileDialog, err := gtk.FileChooserDialogNewWith2Buttons("Select file", w, gtk.FILE_CHOOSER_ACTION_OPEN,
			"Open", gtk.RESPONSE_OK,
			"Cancel", gtk.RESPONSE_CANCEL,
		)
		if err != nil {
			w.ctrl.Config().Error("open file dialog", err)
			return
		}
		defer openfileDialog.Destroy()

		response := openfileDialog.Run()
		if response == gtk.RESPONSE_OK && openfileDialog.GetFilename() != "" {
			f(openfileDialog.GetFilename())
		}
	}
}

func (w *Window) OnSaveQuery(query string, f func(string, string)) {
	openfileDialog, err := gtk.FileChooserDialogNewWith2Buttons("Save file", w, gtk.FILE_CHOOSER_ACTION_SAVE,
		"Save", gtk.RESPONSE_OK,
		"Cancel", gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		w.ctrl.Config().Error("save file dialog", err)
		return
	}
	defer openfileDialog.Destroy()

	response := openfileDialog.Run()
	if response == gtk.RESPONSE_CANCEL {
		return
	}

	f(query, openfileDialog.GetFilename())
}

func (w *Window) AddToolTab(tab *ToolTab) error {
	w.nb.AppendPage(tab.Content(), tab.Label())
	w.nb.SetTabReorderable(tab.Content(), true)
	w.nb.SetCurrentPage(w.nb.PageNum(tab.Content()))
	w.tabs = append(w.tabs, tab)

	// TODO: fix
	//btn.Connect("clicked", func() {
	//index := w.nb.PageNum(wd)
	//w.nb.RemovePage(index)
	//tab.Removed()
	//})

	return nil
}

func (w *Window) RemoveCurrentPage() {
	w.nb.RemovePage(w.CurrentPage())
}

func (w *Window) CurrentPage() int {
	return w.nb.GetCurrentPage()
}

func (w *Window) CurrentTab() *ToolTab {
	if len(w.tabs) == 0 {
		return nil
	}
	return w.tabs[w.CurrentPage()]
}

func (w *Window) Remove(wd gtk.IWidget) {
	w.nb.Remove(wd)
}

func (w Window) PushStatus(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05")}, args...)
	w.statusBar.Push(w.statusBarID, fmt.Sprintf("[%s] "+format, args...))
}

func (w *Window) OnPageRemoved(f interface{}) {
	w.nb.Connect("page-removed", f)
}

func (w *Window) PageCount() int {
	return w.nb.GetNPages()
}

func (w *Window) headerMenu() (*gtk.HeaderBar, error) {
	w.Menu.NewToolTab = glib.SimpleActionNew("new", nil)
	w.Menu.NewSubToolTab = glib.SimpleActionNew("tabnew", nil)
	w.Menu.LoadFile = glib.SimpleActionNew("file.load", nil)
	w.Menu.CloseToolTab = glib.SimpleActionNew("close", nil)
	w.AddAction(w.Menu.NewToolTab)
	w.AddAction(w.Menu.NewSubToolTab)
	w.AddAction(w.Menu.LoadFile)
	w.AddAction(w.Menu.CloseToolTab)

	header, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	header.SetShowCloseButton(true)
	header.SetTitle("BenchWell")
	header.SetSubtitle(w.ctrl.Config().Version)

	// Create a new window menu button
	windowBtnMenu, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	addImg, err := BWImageNewFromFile("add-tab", ICON_SIZE_MENU)
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

	windowMenu.Append("New window", "app.new")
	windowMenu.Append("New connection", "win.new")
	windowMenu.Append("New tab", "win.tabnew")
	windowMenu.Append("Open File", "win.file.load")
	//menu.Append("- ToolTable ToolTab", "win.close")

	// Create a new app menu button
	appBtnMenu, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	prefImg, _ := BWImageNewFromFile("config", ICON_SIZE_MENU)
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
	// https://play.golang.com/p/YMfQouxHuvr
	movingTab := w.tabs[w.tabIndex]

	w.tabs = append(w.tabs[:w.tabIndex], w.tabs[w.tabIndex+1:]...)

	lh := make([]*ToolTab, len(w.tabs[:landing]))
	rh := make([]*ToolTab, len(w.tabs[landing:]))
	copy(lh, w.tabs[:landing])
	copy(rh, w.tabs[landing:])

	w.tabs = append(lh, movingTab)
	w.tabs = append(w.tabs, rh...)
}
