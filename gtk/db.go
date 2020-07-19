package gtk

import "github.com/gotk3/gotk3/gtk"

type DB struct {
	*gtk.Box
	content gtk.IWidget
	w       *Window
}

func (d DB) Init(w *Window) (*DB, error) {
	var err error
	d.w = w
	d.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	d.Box.Show()

	return &d, nil
}

func (d *DB) SetContent(w gtk.IWidget) {
	if d.content != nil {
		d.Remove(d.content)
	}
	d.content = w
	d.PackStart(w, true, true, 0)
}
