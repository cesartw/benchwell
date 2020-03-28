package gtk

import "github.com/gotk3/gotk3/gtk"

type Tab struct {
	*gtk.Box
}

func NewTab() (*Tab, error) {
	t := &Tab{}
	var err error

	t.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	t.Box.SetVExpand(true)
	t.Box.SetHExpand(true)

	return t, nil
}
