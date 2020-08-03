package gtk

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	COLUMN_ICON = iota
	COLUMN_TEXT
)

type HTTPCollection struct {
	*gtk.TreeView
	store *gtk.TreeStore
}

func (h HTTPCollection) Init(w *Window) (*HTTPCollection, error) {
	var err error
	h.TreeView, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	h.TreeView.SetHeadersVisible(false)

	col, err := h.createImageColumn("Icon", COLUMN_ICON)
	if err != nil {
		return nil, err
	}
	h.TreeView.AppendColumn(col)

	col, err = h.createTextColumn("Version", COLUMN_TEXT)
	if err != nil {
		return nil, err
	}
	h.TreeView.AppendColumn(col)

	h.store, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create tree store:", err)
	}
	h.TreeView.SetModel(h.store)

	imageOK, err := BWPixbufFromFile("open", 16)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}

	iter, err := h.addSubRow(nil, imageOK, "test1-1")
	if err != nil {
		return nil, err
	}
	_, err = h.addSubRow(iter, imageOK, "test1-2")
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *HTTPCollection) createTextColumn(title string, id int) (*gtk.TreeViewColumn, error) {
	// In this column we want to show text, hence create a text renderer
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}

	// Tell the renderer where to pick input from. Text renderer understands
	// the "text" property.
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (h *HTTPCollection) createImageColumn(title string, id int) (*gtk.TreeViewColumn, error) {
	// In this column we want to show image data from Pixbuf, hence
	// create a pixbuf renderer
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		return nil, err
	}

	// Tell the renderer where to pick input from. Pixbuf renderer understands
	// the "pixbuf" property.
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", id)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (h *HTTPCollection) addSubRow(
	iter *gtk.TreeIter,
	icon *gdk.Pixbuf,
	text string,
) (*gtk.TreeIter, error) {
	// Get an iterator for a new row at the end of the list store
	i := h.store.Append(iter)

	// Set the contents of the tree store row that the iterator represents
	err := h.store.SetValue(i, COLUMN_ICON, icon)
	if err != nil {
		return nil, err
	}
	err = h.store.SetValue(i, COLUMN_TEXT, text)
	if err != nil {
		return nil, err
	}
	return i, nil
}
