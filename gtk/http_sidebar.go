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
	COLUMN_METHOD
)

var colors = map[string]string{
	"POST":    "#2bb6c7", // gdk.NewRGBA(43.0, 182.0, 199.0, 1),
	"PATCH":   "#e68700", // gdk.NewRGBA(230.0, 135.0, 0.0, 1),
	"PUT":     "#de8f00", // gdk.NewRGBA(222.0, 143.0, 0.0, 1),
	"DELETE":  "#ff0000", // gdk.NewRGBA(255.0, 0, 0.0, 1),
	"GET":     "#27d000", // gdk.NewRGBA(39.0, 208.0, 0.0, 1),
	"HEAD":    "#19a800", // gdk.NewRGBA(25.0, 168.0, 0.0, 1),
	"OPTIONS": "#298700", // gdk.NewRGBA(41.0, 135.0, 0.0, 1),
}

type HTTPCollection struct {
	*gtk.Box
	w     *Window
	h     *HTTPScreen
	tree  *gtk.TreeView
	store *gtk.TreeStore

	menu struct {
		*gtk.Menu
		newRequestMenu *gtk.MenuItem
		newFolderMenu  *gtk.MenuItem
		deleteMenu     *gtk.MenuItem
		editMenu       *gtk.MenuItem
	}

	collectioncb *gtk.ComboBoxText

	selectedCollection *config.HTTPCollection
	selectedItem       *config.HTTPItem

	ctrl ctrlHTTPCollection

	nameRenderer *gtk.CellRendererText
}

type ctrlHTTPCollection interface {
	OnLoadItem()
	OnCollectionSelected()
	OnNewRequest()
	OnNewFolder()
	OnDeleteItem()
	OnSave(*config.HTTPItem) error
}

func (h HTTPCollection) Init(
	w *Window,
	hs *HTTPScreen,
	ctrl ctrlHTTPCollection,
) (*HTTPCollection, error) {
	defer config.LogStart("HTTPCollection.Init", nil)()

	var err error

	h.ctrl = ctrl
	h.w = w
	h.h = hs

	h.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	h.tree, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	h.tree.Show()
	h.tree.SetHeadersVisible(false)
	h.tree.SetProperty("show-expanders", false)
	h.tree.Connect("row-activated", ctrl.OnLoadItem)

	h.tree.SetEnableTreeLines(true)
	selection, err := h.tree.GetSelection()
	if err != nil {
		return nil, err
	}
	selection.Connect("changed", func() {
		_, iter, ok := selection.GetSelected()
		if !ok {
			return
		}

		value, _ := h.store.GetValue(iter, COLUMN_ID)
		v, _ := value.GoValue()

		id := v.(int64)
		h.selectedItem = h.searchItem(id)
	})

	h.tree.Connect("button-release-event", func(_ *gtk.TreeView, e *gdk.Event) {
		keyEvent := gdk.EventButtonNewFromEvent(e)
		if keyEvent.Button() != gdk.BUTTON_SECONDARY {
			return
		}

		h.menu.PopupAtPointer(e)
	})

	err = h.buildMenu()
	if err != nil {
		return nil, err
	}
	h.menu.deleteMenu.Connect("activate", ctrl.OnDeleteItem)
	h.menu.newFolderMenu.Connect("activate", ctrl.OnNewFolder)
	h.menu.newRequestMenu.Connect("activate", ctrl.OnNewRequest)
	h.menu.editMenu.Connect("activate", h.onEdit)

	col, err := h.createImageColumn("", COLUMN_ICON)
	if err != nil {
		return nil, err
	}
	h.tree.AppendColumn(col)
	h.tree.SetExpanderColumn(col)

	col, h.nameRenderer, err = h.createTextColumn("", COLUMN_TEXT)
	if err != nil {
		return nil, err
	}

	h.nameRenderer.Connect("edited", func(cell *gtk.CellRendererText, pathS, text string) {
		h.nameRenderer.SetProperty("editable", false)
		path, _ := gtk.TreePathNewFromString(pathS)
		iter, _ := h.store.GetIter(path)

		h.store.SetValue(iter, COLUMN_TEXT, text)
		h.selectedItem.Name = text
		h.ctrl.OnSave(h.selectedItem)
	})

	h.nameRenderer.Connect("editing-canceled", func() {
		h.nameRenderer.SetProperty("editable", false)
	})

	h.tree.AppendColumn(col)

	col, _, err = h.createTextColumn("", COLUMN_METHOD, func(cell *gtk.CellRenderer, iter *gtk.TreeIter) {
		value, _ := h.store.GetValue(iter, COLUMN_METHOD)
		v, _ := value.GoValue()

		if method, ok := v.(string); ok && method != "" {
			cell.SetProperty("markup",
				fmt.Sprintf(`<span foreground="%s">%s</span>`, colors[method], method))
		}
	})
	if err != nil {
		return nil, err
	}
	h.tree.AppendColumn(col)

	h.store, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_INT64, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create tree store:", err)
	}
	h.tree.SetModel(h.store)

	collectionSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	collectionSW.Show()
	collectionSW.Add(h.tree)

	h.collectioncb, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	h.collectioncb.Show()
	h.collectioncb.Append("", "")
	h.collectioncb.Connect("changed", func() {
		if h.collectioncb.GetActiveID() == "" {
			return
		}

		id, _ := strconv.ParseInt(h.collectioncb.GetActiveID(), 10, 64)

		for _, collection := range config.Collections {
			if collection.ID == id {
				h.selectedCollection = collection
				break
			}
		}

		ctrl.OnCollectionSelected()
	})

	for _, collection := range config.Collections {
		h.collectioncb.Append(fmt.Sprintf("%d", collection.ID), collection.Name)
	}

	h.Box.PackStart(h.collectioncb, false, false, 0)
	h.Box.PackStart(collectionSW, true, true, 0)

	return &h, nil
}

func (h *HTTPCollection) RemoveItem(path *gtk.TreePath) {
	iter, err := h.store.GetIter(path)
	if err != nil {
		return
	}

	h.store.Remove(iter)
}

func (h *HTTPCollection) GetSelectedCollection() *config.HTTPCollection {
	defer config.LogStart("HTTPCollection.GetSelectedCollection", nil)()
	return h.selectedCollection
}

func (h *HTTPCollection) GetSelectedItem() *config.HTTPItem {
	defer config.LogStart("HTTPCollection.GetSelectedItem", nil)()
	return h.selectedItem
}

func (h *HTTPCollection) GetSelectedPath() (*gtk.TreePath, error) {
	sel, _ := h.tree.GetSelection()
	_, iter, _ := sel.GetSelected()

	path, err := h.store.GetPath(iter)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (h *HTTPCollection) LoadCollection(items []*config.HTTPItem) error {
	defer config.LogStart("HTTPCollection.LoadCollection", nil)()

	h.store.Clear()
	return h.buildTree(nil, items)
}

func (h *HTTPCollection) LoadFolder(path *gtk.TreePath, item *config.HTTPItem) error {
	defer config.LogStart("HTTPCollection.LoadFolder", nil)()

	iter, err := h.store.GetIter(path)
	if err != nil {
		return err
	}

	err = h.buildTree(iter, item.Items)
	if err != nil {
		return err
	}

	h.tree.ExpandRow(path, false)

	return nil
}

func (h *HTTPCollection) buildTree(iter *gtk.TreeIter, items []*config.HTTPItem) error {
	defer config.LogStart("HTTPCollection.buildTree", nil)()

	for _, item := range items {
		iter, err := h.addRow(iter, item)
		if err != nil {
			return err
		}

		if item.IsFolder {
			err = h.buildTree(iter, item.Items)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *HTTPCollection) createTextColumn(
	title string,
	id int,
	format ...func(*gtk.CellRenderer, *gtk.TreeIter),
) (*gtk.TreeViewColumn, *gtk.CellRendererText, error) {
	defer config.LogStart("HTTPCollection.createTextColumn", nil)()

	// In this column we want to show text, hence create a text renderer
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, nil, err
	}

	// Tell the renderer where to pick input from. Text renderer understands
	// the "text" property.
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		return nil, nil, err
	}
	column.SetCellDataFunc(&cellRenderer.CellRenderer, func(
		tree_column *gtk.TreeViewColumn,
		cell *gtk.CellRenderer,
		tree_model *gtk.TreeModel,
		iter *gtk.TreeIter,
		userData ...interface{}) {
		if len(format) > 0 {
			format[0](cell, iter)
		}
	})

	return column, cellRenderer, nil
}

func (h *HTTPCollection) createImageColumn(title string, id int) (*gtk.TreeViewColumn, error) {
	defer config.LogStart("HTTPCollection.createImageColumn", nil)()

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
	item *config.HTTPItem,
) (*gtk.TreeIter, error) {
	defer config.LogStart("HTTPCollection.addRow", nil)()

	// Get an iterator for a new row at the end of the list store
	iter := h.store.Append(parentIter)

	var err error
	if item.IsFolder {
		image, err := BWPixbufFromFile("directory", "orange", 16)
		if err != nil {
			return nil, err
		}
		err = h.store.SetValue(iter, COLUMN_ICON, image)
		if err != nil {
			return nil, err
		}
	} else {
		// NOTE: gtk_tree_view_column_set_cell_data_func  to format the method name
		err = h.store.SetValue(iter, COLUMN_METHOD, item.Method)
		if err != nil {
			return nil, err
		}
	}

	err = h.store.SetValue(iter, COLUMN_TEXT, item.Name)
	if err != nil {
		return nil, err
	}

	err = h.store.SetValue(iter, COLUMN_ID, item.ID)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (h *HTTPCollection) AddItem(item *config.HTTPItem, path *gtk.TreePath) {
	iter, err := h.store.GetIter(path)
	if err != nil {
		return
	}

	_, _ = h.addRow(iter, item)
}

func (h *HTTPCollection) buildMenu() error {
	defer config.LogStart("HTTPCollection.buildMenu", nil)()

	var err error
	h.menu.Menu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	h.menu.newFolderMenu, err = BWMenuItemWithImage("New folder", "add")
	if err != nil {
		return err
	}
	h.menu.newRequestMenu, err = BWMenuItemWithImage("New request", "add")
	if err != nil {
		return err
	}
	h.menu.deleteMenu, err = BWMenuItemWithImage("Delete", "close")
	if err != nil {
		return err
	}
	h.menu.editMenu, err = BWMenuItemWithImage("Edit", "config")
	if err != nil {
		return err
	}

	h.menu.Add(h.menu.newRequestMenu)
	h.menu.Add(h.menu.newFolderMenu)
	h.menu.Add(h.menu.editMenu)
	h.menu.Add(h.menu.deleteMenu)

	return nil
}

func (h *HTTPCollection) onEdit() {
	h.nameRenderer.SetProperty("editable", true)
	path, _ := h.GetSelectedPath()
	col := h.tree.GetColumn(COLUMN_TEXT)
	h.tree.SetCursorOnCell(path, col, &h.nameRenderer.CellRenderer, true)
}

func (h *HTTPCollection) searchItem(id int64) *config.HTTPItem {
	for _, item := range h.selectedCollection.Items {
		if item.ID == id {
			return item

		}
		for _, subitem := range item.Items {
			if subitem.ID == id {
				return subitem
			}
		}
	}

	return nil
}
