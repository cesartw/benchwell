package gtk

import (
	"errors"
	"log"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type App struct {
	*gtk.Application
	mainWindow *Window
	Menu       struct {
		Application struct {
			New         *glib.SimpleAction
			Open        *glib.SimpleAction
			Save        *glib.SimpleAction
			Preferences *glib.SimpleAction
		}
	}
}

func New(appid string) (*App, error) {
	var err error
	f := &App{}

	f.Application, err = gtk.ApplicationNew(appid, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}

	f.Connect("startup", func() {
		f.Menu.Application.New = glib.SimpleActionNew("new_conn", nil)
		f.Menu.Application.Open = glib.SimpleActionNew("open", nil)
		f.Menu.Application.Save = glib.SimpleActionNew("save", nil)
		f.Menu.Application.Preferences = glib.SimpleActionNew("prefernces", nil)

		menu := glib.MenuNew()
		menu.Append("New Connection", "app.new_conn")
		menu.Append("Open", "app.open")
		menu.Append("Save", "app.save")
		menu.Append("Preferences", "app.preferences")

		f.AddAction(f.Menu.Application.New)
		f.AddAction(f.Menu.Application.Open)
		f.AddAction(f.Menu.Application.Save)
		f.AddAction(f.Menu.Application.Preferences)

		f.SetAppMenu(&menu.MenuModel)
	})

	f.Connect("activate", func() {
		f.mainWindow, err = f.newMainScreen()
		if err != nil {
			log.Fatal(err)
		}

		header, err := f.headerMenu()
		if err != nil {
			log.Fatal(err)
		}
		f.mainWindow.SetTitlebar(header)

		f.AddWindow(f.mainWindow)
		f.mainWindow.ShowAll()
	})

	return f, nil
}

func (f *App) NewConnectScreen() (*ConnectScreen, error) {
	return newConnectScreen()
}

func (f *App) AddTab(label *gtk.Label, w gtk.IWidget) {
	f.mainWindow.AddTab(label, w)
}

func (f *App) OnPageRemoved(fn interface{}) {
	f.mainWindow.OnPageRemoved(fn)
}

func (f *App) OnTabClick(fn interface{}) {
	f.mainWindow.OnTabClick(fn)
}

func (f *App) PageCount() int {
	return f.mainWindow.PageCount()
}

func (f *App) Remove(w gtk.IWidget) {
	f.mainWindow.Remove(w)
}

func (f *App) Show() {
	f.mainWindow.Show()
}

func (f *App) PushStatus(status string, args ...interface{}) {
	f.mainWindow.PushStatus(status, args...)
}

func (f *App) newMainScreen() (*Window, error) {
	w := &Window{}
	return w, w.init(f.Application)
}

//f.Menu.Application.New = glib.SimpleActionNew("new_conn", nil)
//f.Menu.Application.Open = glib.SimpleActionNew("open", nil)
//f.Menu.Application.Save = glib.SimpleActionNew("save", nil)
//f.Menu.Application.Preferences = glib.SimpleActionNew("prefernces", nil)

//menu := glib.MenuNew()
//menu.Append("New Connection", "app.new_conn")
//menu.Append("Open", "app.open")
//menu.Append("Save", "app.save")
//menu.Append("Preferences", "app.preferences")

//f.AddAction(f.Menu.Application.New)
//f.AddAction(f.Menu.Application.Open)
//f.AddAction(f.Menu.Application.Save)
//f.AddAction(f.Menu.Application.Preferences)

func (f *App) headerMenu() (*gtk.HeaderBar, error) {
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

	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append("New tab", "app.new_conn")
	//menu.Append("Open", "app.open")
	//menu.Append("Save", "app.save")
	menu.Append("Preferences", "app.preferences")

	f.AddAction(f.Menu.Application.New)
	//f.AddAction(f.Menu.Application.Open)
	//f.AddAction(f.Menu.Application.Save)
	f.AddAction(f.Menu.Application.Preferences)

	mbtn.SetMenuModel(&menu.MenuModel)

	// add the menu button to the header
	header.PackStart(mbtn)

	// Assemble the window
	return header, nil
}
