package gtk

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type Window struct {
	*gtk.ApplicationWindow
	//builder     *gtk.Builder
	box         *gtk.Box
	statusBar   *gtk.Statusbar
	statusBarID uint
}

func (w *Window) init(app *gtk.Application) (err error) {
	//w.builder, err = gtk.BuilderNewFromFile("main.glade")
	//if err != nil {
	//return err
	//}

	w.ApplicationWindow, err = gtk.ApplicationWindowNew(app)
	w.SetTitle("SQLHero")
	w.SetShowMenubar(true)

	//w.builder.ConnectSignals(signals)
	//obj, err := w.builder.GetObject("MainWindow")
	//if err != nil {
	//return err
	//}

	//w.Window = obj.(*gtk.Window)

	//obj, err = w.builder.GetObject("MainWindowBox")
	//if err != nil {
	//return err
	//}

	w.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}

	w.statusBar, err = gtk.StatusbarNew()
	if err != nil {
		return err
	}
	w.box.PackEnd(w.statusBar, false, false, 0)

	w.statusBarID = w.statusBar.GetContextId("main")

	w.ApplicationWindow.Add(w.box)

	return nil
}

func (w *Window) Add(wd gtk.IWidget) {
	w.box.Add(wd)
}

func (w *Window) Remove(wd gtk.IWidget) {
	w.box.Remove(wd)
}

func (w Window) PushStatus(format string, args ...interface{}) {
	w.statusBar.Push(w.statusBarID, fmt.Sprintf(format, args...))
}
