package gtk

import (
	"log"

	"bitbucket.org/goreorto/sqlhero/sqlengine"
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

	f.Connect("activate", func() {
		f.mainWindow, err = f.newMainScreen()
		if err != nil {
			log.Fatal(err)
		}

		f.AddWindow(f.mainWindow)
	})

	return f, nil
}

func (f *Factory) NewConnectScreen() (*ConnectScreen, error) {
	return newConnectScreen()
}

func (f *Factory) NewConnectionScreen(ctx sqlengine.Context) (*ConnectionScreen, error) {
	return newConnectionScreen(ctx)
}

func (f *Factory) newMainScreen() (*Window, error) {
	w := &Window{}
	return w, w.init()
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
