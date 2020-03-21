package gtk

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Window struct {
	*gtk.Window
	builder     *gtk.Builder
	box         *gtk.Box
	statusBar   *gtk.Statusbar
	statusBarID uint
}

func (w *Window) init() (err error) {
	w.builder, err = gtk.BuilderNewFromFile("main.glade")
	if err != nil {
		return err
	}

	signals := map[string]interface{}{
		"on_main_window_destroy": w.onMainWindowDestroy,
	}

	w.builder.ConnectSignals(signals)
	obj, err := w.builder.GetObject("MainWindow")
	if err != nil {
		return err
	}

	w.Window = obj.(*gtk.Window)

	obj, err = w.builder.GetObject("MainWindowBox")
	if err != nil {
		return err
	}

	w.box = obj.(*gtk.Box)

	obj, err = w.builder.GetObject("MainStatusBar")
	if err != nil {
		return err
	}

	w.statusBar = obj.(*gtk.Statusbar)
	w.statusBarID = w.statusBar.GetContextId("main")

	return nil
}

func (w *Window) Add(wd gtk.IWidget) {
	w.box.Add(wd)
}

func (w *Window) Remove(wd gtk.IWidget) {
	w.box.Remove(wd)
}

func (w Window) onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
}

func (w Window) PushStatus(format string, args ...interface{}) {
	w.statusBar.Push(w.statusBarID, fmt.Sprintf(format, args...))
}
