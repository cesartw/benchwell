package gtk

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Factory struct {
	*gtk.Application
	mainWindow *Window
}

func New(appid string) (*Factory, error) {
	var err error
	f := &Factory{}

	f.Application, err = gtk.ApplicationNew(appid, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}

	f.Connect("startup", func() {
		f.SetAppMenu(&f.menu().MenuModel)
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

// Application
//  -> Preferences
//     -> Tabs
func (f *Factory) menu() *glib.Menu {
	open := glib.MenuItemNewWithLabel("Open")
	save := glib.MenuItemNewWithLabel("Save")
	prefs := glib.MenuItemNewWithLabel("Preferences")

	main := glib.MenuNew()
	main.AppendItem(open)
	main.AppendItem(save)
	main.AppendItem(prefs)

	return main
}

func (f *Factory) NewConnectScreen() (*ConnectScreen, error) {
	return newConnectScreen()
}

func (f *Factory) newMainScreen() (*Window, error) {
	w := &Window{}
	return w, w.init(f.Application)
}

func (f *Factory) Add(w gtk.IWidget) {
	f.mainWindow.Add(w)
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
