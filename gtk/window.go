package gtk

import (
	"errors"
	"fmt"
	"time"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Window struct {
	*gtk.ApplicationWindow
	nb          *gtk.Notebook
	box         *gtk.Box // holds nb and statusbar
	statusBar   *gtk.Statusbar
	statusBarID uint

	Menu struct {
		NewConnection *glib.SimpleAction
		NewTab        *glib.SimpleAction
		NewSubTab     *glib.SimpleAction
		LoadFile      *glib.SimpleAction
		SaveQuery     *glib.SimpleAction
		CloseTab      *glib.SimpleAction
	}
}

func (w Window) Init(app *gtk.Application, ctrl interface {
	OnNewTab()
	OnNewSubTab()
	OnCloseTab()
	OnFileSelected(string)
	OnSaveQuery(string, string)
}) (*Window, error) {
	var err error
	w.ApplicationWindow, err = gtk.ApplicationWindowNew(app)
	w.SetTitle("SQLaid")
	w.SetSizeRequest(1024, 768)

	w.nb, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	w.nb.SetName("MainNotebook")

	switch config.Env.GUI.ConnectionTabPosition {
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
	w.Menu.NewTab.Connect("activate", ctrl.OnNewTab)
	// action menu for sub tabs
	w.Menu.NewSubTab.Connect("activate", ctrl.OnNewSubTab)
	w.Menu.CloseTab.Connect("activate", ctrl.OnCloseTab)
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
			config.Env.Log.Error("open file dialog", err)
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
		config.Env.Log.Error("save file dialog", err)
		return
	}
	defer openfileDialog.Destroy()

	response := openfileDialog.Run()
	if response == gtk.RESPONSE_CANCEL {
		return
	}

	f(query, openfileDialog.GetFilename())
}

func (w *Window) AddTab(label *gtk.Label, wd gtk.IWidget, removed func()) error {
	header, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
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

	w.nb.AppendPage(wd, header)
	w.nb.SetTabReorderable(wd, true)
	w.nb.SetCurrentPage(w.nb.PageNum(wd))

	btn.Connect("clicked", func() {
		index := w.nb.PageNum(wd)
		w.nb.RemovePage(index)
		removed()
	})

	return nil
}

func (w *Window) RemoveCurrentPage() {
	w.nb.RemovePage(w.CurrentPage())
}

func (w *Window) CurrentPage() int {
	return w.nb.GetCurrentPage()
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
	header, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	header.SetShowCloseButton(true)
	header.SetTitle("SQLAID")
	header.SetSubtitle(config.Env.Version)

	// Create a new menu button
	mbtn, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	img, _ := gtk.ImageNewFromIconName("gtk-add", gtk.ICON_SIZE_MENU)
	mbtn.SetImage(img)

	// Set up the menu model for the button
	menu := glib.MenuNew()
	if menu == nil {
		return nil, errors.New("nil menu")
	}

	w.Menu.NewTab = glib.SimpleActionNew("new", nil)
	w.Menu.NewSubTab = glib.SimpleActionNew("tabnew", nil)
	w.Menu.LoadFile = glib.SimpleActionNew("file.load", nil)
	w.Menu.CloseTab = glib.SimpleActionNew("close", nil)
	w.AddAction(w.Menu.NewTab)
	w.AddAction(w.Menu.NewSubTab)
	w.AddAction(w.Menu.LoadFile)
	w.AddAction(w.Menu.CloseTab)

	menu.Append("New window", "app.new")
	menu.Append("New connection", "win.new")
	menu.Append("New tab", "win.tabnew")
	//menu.Append("- Table Tab", "win.close")
	menu.Append("Open File", "win.file.load")
	menu.Append("Preferences", "app.preferences")
	menu.Append("Dark toggle", "app.darkmode")

	mbtn.SetMenuModel(&menu.MenuModel)

	// add the menu button to the header
	header.PackStart(mbtn)

	// Assemble the window
	return header, nil
}
