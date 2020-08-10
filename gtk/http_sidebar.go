package gtk

import (
	"fmt"
	"log"
	"strconv"

	"bitbucket.org/goreorto/benchwell/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	COLUMN_ICON = iota
	COLUMN_TEXT
	COLUMN_ID
)

type HTTPCollection struct {
	*gtk.Box
	w     *Window
	tree  *gtk.TreeView
	store *gtk.TreeStore

	colBox *gtk.ComboBoxText

	selectedCollection *config.HTTPCollection

	ctrl ctrlHTTPCollection
}

type ctrlHTTPCollection interface {
	Config() *config.Config
}

func (h HTTPCollection) Init(
	w *Window,
	ctrl ctrlHTTPCollection,
) (*HTTPCollection, error) {
	var err error

	h.ctrl = ctrl
	h.w = w

	h.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	h.tree, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	h.tree.SetHeadersVisible(false)
	h.tree.Connect("row-activated", func(_ *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
		iter, err := h.store.GetIter(path)
		if err != nil {
			h.w.PushStatus("getting iter: " + err.Error())
			return
		}

		storeValue, err := h.store.GetValue(iter, COLUMN_ID)
		if err != nil {
			h.w.PushStatus("getting value: " + err.Error())
			return
		}

		goval, err := storeValue.GoValue()
		if err != nil {
			h.w.PushStatus("getting go value: " + err.Error())
			return
		}
		itemID := goval.(int64)
		for _, item := range h.selectedCollection.Items {
			if item.ID != itemID {
				continue
			}

			err = item.LoadFull()
			if err != nil {
				h.w.PushStatus("loading item: " + err.Error())
				return
			}

			if item.IsFolder {
				err = h.buildTree(iter, item.Items)
				if err != nil {
					h.w.PushStatus("add item: " + err.Error())
					return
				}

				h.tree.ExpandRow(path, false)
			}
		}
	})

	h.colBox, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	h.colBox.Connect("changed", h.onCollectionSelected)

	for _, collection := range h.ctrl.Config().Collections {
		h.colBox.Append(fmt.Sprintf("%d", collection.ID), collection.Name)
	}

	col, err := h.createImageColumn("", COLUMN_ICON)
	if err != nil {
		return nil, err
	}
	h.tree.AppendColumn(col)

	col, err = h.createTextColumn("", COLUMN_TEXT)
	if err != nil {
		return nil, err
	}
	h.tree.AppendColumn(col)

	h.store, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_INT64)
	if err != nil {
		log.Fatal("Unable to create tree store:", err)
	}
	h.tree.SetModel(h.store)

	h.Box.PackStart(h.colBox, false, false, 0)
	h.Box.PackStart(h.tree, true, true, 0)

	if len(h.ctrl.Config().Collections) > 0 {
		h.colBox.SetActive(0)
	}

	return &h, nil
}

func (h *HTTPCollection) buildTree(iter *gtk.TreeIter, items []*config.HTTPItem) error {
	for _, item := range items {
		switch item.IsFolder {
		case true:
			imageOK, err := BWPixbufFromFile("directory", 16)
			if err != nil {
				return err
			}

			iter, err := h.addRow(iter, imageOK, item.ID, item.Name)
			if err != nil {
				return err
			}

			err = h.buildTree(iter, item.Items)
			if err != nil {
				return err
			}
		case false:
			_, err := h.addRow(iter, nil, item.ID, item.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
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

func (h *HTTPCollection) addRow(
	parentIter *gtk.TreeIter,
	icon *gdk.Pixbuf,
	id int64,
	text string,
) (*gtk.TreeIter, error) {
	// Get an iterator for a new row at the end of the list store
	iter := h.store.Append(parentIter)

	var err error
	// Set the contents of the tree store row that the iterator represents
	if icon != nil {
		err = h.store.SetValue(iter, COLUMN_ICON, icon)
		if err != nil {
			return nil, err
		}
	}

	err = h.store.SetValue(iter, COLUMN_TEXT, text)
	if err != nil {
		return nil, err
	}
	err = h.store.SetValue(iter, COLUMN_ID, id)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (h *HTTPCollection) onCollectionSelected() {
	if h.colBox.GetActiveID() == "" {
		return
	}

	id, _ := strconv.ParseInt(h.colBox.GetActiveID(), 10, 64)

	h.store.Clear()

	for _, c := range h.ctrl.Config().Collections {
		if c.ID != id {
			continue
		}
		h.selectedCollection = c
		break
	}
	if h.selectedCollection == nil {
		h.w.PushStatus("collection not found")
		return
	}

	err := h.selectedCollection.LoadRootItems()
	if err != nil {
		h.w.PushStatus("Error loading collection: " + err.Error())
		return
	}

	h.buildTree(nil, h.selectedCollection.Items)
}
