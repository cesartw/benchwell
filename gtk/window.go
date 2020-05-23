package gtk

import (
	"errors"
	"fmt"

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
		CloseTab      *glib.SimpleAction
	}
}

func (w Window) Init(app *gtk.Application) (*Window, error) {
	var err error
	w.ApplicationWindow, err = gtk.ApplicationWindowNew(app)
	w.SetTitle("SQLaid")
	w.SetSizeRequest(1024, 768)

	w.nb, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
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

	return &w, nil
}

func (w *Window) OnTabClick(f interface{}) {
	w.nb.Connect("button-press-event", f)
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
		w.RemoveCurrentPage()
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
	w.statusBar.Push(w.statusBarID, fmt.Sprintf(format, args...))
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

	// Set up the menu model for the button
	menu := glib.MenuNew()
	if menu == nil {
		return nil, errors.New("nil menu")
	}

	w.Menu.NewTab = glib.SimpleActionNew("new", nil)
	w.Menu.NewSubTab = glib.SimpleActionNew("tabnew", nil)
	w.Menu.CloseTab = glib.SimpleActionNew("close", nil)
	w.AddAction(w.Menu.NewTab)
	w.AddAction(w.Menu.NewSubTab)
	w.AddAction(w.Menu.CloseTab)

	menu.Append("Open window", "app.new")
	menu.Append("+ Connection Tab", "win.new")
	menu.Append("+ Table Tab", "win.tabnew")
	menu.Append("- Table Tab", "win.close")
	menu.Append("Preferences", "app.preferences")

	mbtn.SetMenuModel(&menu.MenuModel)

	// add the menu button to the header
	header.PackStart(mbtn)

	// Assemble the window
	return header, nil
}
