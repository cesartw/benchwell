package gtk

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type Window struct {
	*gtk.ApplicationWindow
	nb          *gtk.Notebook
	box         *gtk.Box // holds nb and statusbar
	statusBar   *gtk.Statusbar
	statusBarID uint
}

func (w *Window) init(app *gtk.Application) (err error) {
	w.ApplicationWindow, err = gtk.ApplicationWindowNew(app)
	w.SetTitle("SQLHero")

	w.nb, err = gtk.NotebookNew()
	if err != nil {
		return err
	}

	w.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	w.box.SetVExpand(true)
	w.box.SetHExpand(true)

	w.statusBar, err = gtk.StatusbarNew()
	if err != nil {
		return err
	}

	w.box.PackStart(w.nb, true, true, 0)
	w.box.PackEnd(w.statusBar, false, false, 0)

	w.statusBarID = w.statusBar.GetContextId("main")

	w.ApplicationWindow.Add(w.box)

	return nil
}

func (w *Window) OnTabClick(f interface{}) {
	w.nb.Connect("button-press-event", f)
}

func (w *Window) AddTab(label *gtk.Label, wd gtk.IWidget) error {
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
		if index == -1 {
			return
		}
		w.nb.RemovePage(index)
	})

	return nil
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
