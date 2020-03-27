package gtk

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Factory struct {
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

func New(appid string) (*Factory, error) {
	var err error
	f := &Factory{}

	f.Application, err = gtk.ApplicationNew(appid, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}

	f.Connect("startup", func() {
		//if f.PrefersAppMenu() {
		//}
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

		f.AddWindow(f.mainWindow)
		f.mainWindow.ShowAll()
	})

	return f, nil
}

func (f *Factory) NewConnectScreen() (*ConnectScreen, error) {
	return newConnectScreen()
}

func (f *Factory) newMainScreen() (*Window, error) {
	w := &Window{}
	return w, w.init(f.Application)
}

func (f *Factory) AddTab(label *gtk.Label, w gtk.IWidget) {
	f.mainWindow.AddTab(label, w)
}

func (f *Factory) OnPageRemoved(fn interface{}) {
	f.mainWindow.OnPageRemoved(fn)
}

func (f *Factory) OnTabClick(fn interface{}) {
	f.mainWindow.OnTabClick(fn)
}

func (f *Factory) PageCount() int {
	return f.mainWindow.PageCount()
}

func (f *Factory) Remove(w gtk.IWidget) {
	f.mainWindow.Remove(w)
}

func (f *Factory) Show() {
	f.mainWindow.Show()
}

func (f *Factory) PushStatus(status string, args ...interface{}) {
	f.mainWindow.PushStatus(status, args...)
}
