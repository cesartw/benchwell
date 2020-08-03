package gtk

import (
	"bitbucket.org/goreorto/benchwell/config"
	"github.com/gotk3/gotk3/gtk"
)

var methods = [7]string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"OPTIONS",
	"HEAD",
}

type KeyValue struct {
	*gtk.Box
	enabled *gtk.CheckButton
	key     *gtk.Entry
	value   *gtk.Entry
	remove  *gtk.Button
}

func (c KeyValue) Init() (*KeyValue, error) {
	var err error
	c.key, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	c.key.SetPlaceholderText("Name")

	c.value, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	c.value.SetPlaceholderText("Value")

	c.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}

	c.remove, err = BWButtonNewFromIconName("close", ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	c.enabled, err = gtk.CheckButtonNew()
	if err != nil {
		return nil, err
	}

	c.Box.PackStart(c.enabled, false, false, 5)
	c.Box.PackStart(c.key, true, true, 0)
	c.Box.PackStart(c.value, true, true, 0)
	c.Box.PackEnd(c.remove, false, false, 5)

	return &c, nil
}

type KeyValues struct {
	*gtk.Box
	keyvalues []*KeyValue
}

func (c KeyValues) Init() (*KeyValues, error) {
	var err error
	c.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	return &c, c.Add()
}

func (c *KeyValues) Add() error {
	kv, err := KeyValue{}.Init()
	if err != nil {
		return err
	}

	kv.remove.Connect("clicked", func() {
		c.Box.Remove(kv)

		for i, v := range c.keyvalues {
			if v != kv {
				continue
			}
			c.keyvalues = append(c.keyvalues[:i], c.keyvalues[i+1:]...)
		}

		if len(c.keyvalues) == 0 {
			c.Add()
		}
	})

	focused := func() {
		if c.keyvalues[len(c.keyvalues)-1] != kv {
			return
		}

		c.Add()
	}
	kv.key.Connect("grab-focus", focused)
	kv.value.Connect("grab-focus", focused)

	c.Box.PackStart(kv, false, false, 0)
	c.keyvalues = append(c.keyvalues, kv)
	kv.ShowAll()

	return nil
}

type HTTPScreen struct {
	*gtk.Paned
	w *Window

	// Address bar
	method  *gtk.ComboBoxText
	address *gtk.Entry
	send    *gtk.Button

	// request body
	body    *SourceView
	headers *KeyValues
	params  *KeyValues

	// side panel
	collection *HTTPCollection

	ctrl httpScreenCtrl
}

type httpScreenCtrl interface {
	Config() *config.Config
}

func (h HTTPScreen) Init(w *Window, ctrl httpScreenCtrl) (*HTTPScreen, error) {
	var err error

	h.w = w
	h.ctrl = ctrl
	h.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	h.Paned.SetWideHandle(true)
	h.collection, err = HTTPCollection{}.Init(w, ctrl)
	if err != nil {
		return nil, err
	}

	main, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	addressBar, err := h.buildAddressBar()
	if err != nil {
		return nil, err
	}

	h.body, err = SourceView{}.Init(h.w, SourceViewOptions{true, true, "json"}, ctrl)
	if err != nil {
		return nil, err
	}

	options, err := h.buildOptions()
	if err != nil {
		return nil, err
	}

	main.PackStart(addressBar, false, false, 0)
	main.PackStart(options, false, false, 5)
	main.PackStart(h.body, true, true, 0)

	h.Paned.Add1(h.collection)
	h.Paned.Add2(main)
	h.Paned.ShowAll()

	return &h, nil
}

func (h *HTTPScreen) buildAddressBar() (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}

	h.method, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	for _, m := range methods {
		h.method.AppendText(m)
	}
	h.method.SetActive(0)

	h.address, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	h.address.SetPlaceholderText("http://localhost/path.json")

	h.send, err = gtk.ButtonNewWithLabel("SEND")
	if err != nil {
		return nil, err
	}

	// PackStart(child IWidget, expand, fill bool, padding uint) {
	box.PackStart(h.method, false, false, 5)
	box.PackStart(h.address, true, true, 0)
	box.PackEnd(h.send, false, false, 5)

	return box, nil
}

func (h *HTTPScreen) buildOptions() (*gtk.Notebook, error) {
	nb, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}

	h.headers, err = KeyValues{}.Init()
	if err != nil {
		return nil, err
	}

	h.params, err = KeyValues{}.Init()
	if err != nil {
		return nil, err
	}

	paramsLabel, err := gtk.LabelNew("Params")
	if err != nil {
		return nil, err
	}
	headersLabel, err := gtk.LabelNew("Headers")
	if err != nil {
		return nil, err
	}

	nb.AppendPage(h.params, paramsLabel)
	nb.AppendPage(h.headers, headersLabel)

	return nb, nil
}
