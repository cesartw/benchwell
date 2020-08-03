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
	c.key.SetPlaceholderText("Key")

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

func (c HTTPScreen) Init(w *Window, ctrl httpScreenCtrl) (*HTTPScreen, error) {
	var err error

	c.w = w
	c.ctrl = ctrl
	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	c.Paned.SetWideHandle(true)
	c.collection, err = HTTPCollection{}.Init(w)
	if err != nil {
		return nil, err
	}

	main, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	addressBar, err := c.buildAddressBar()
	if err != nil {
		return nil, err
	}

	frameParams, err := gtk.FrameNew("Params")
	if err != nil {
		return nil, err
	}

	frameHeaders, err := gtk.FrameNew("Headers")
	if err != nil {
		return nil, err
	}

	c.body, err = SourceView{}.Init(c.w, SourceViewOptions{true, true, "json"}, ctrl)
	if err != nil {
		return nil, err
	}

	c.headers, err = KeyValues{}.Init()
	if err != nil {
		return nil, err
	}
	frameParams.Add(c.headers)

	c.params, err = KeyValues{}.Init()
	if err != nil {
		return nil, err
	}
	frameHeaders.Add(c.params)

	main.PackStart(addressBar, false, false, 0)
	main.PackStart(frameParams, false, false, 5)
	main.PackStart(frameHeaders, false, false, 5)
	main.PackStart(c.body, true, true, 0)

	c.Paned.Add1(c.collection)
	c.Paned.Add2(main)
	c.Paned.ShowAll()

	return &c, nil
}

func (c *HTTPScreen) buildAddressBar() (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}

	c.method, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	for _, m := range methods {
		c.method.AppendText(m)
	}
	c.method.SetActive(0)

	c.address, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	c.address.SetPlaceholderText("http://localhost/path.json")

	c.send, err = gtk.ButtonNewWithLabel("SEND")
	if err != nil {
		return nil, err
	}

	// PackStart(child IWidget, expand, fill bool, padding uint) {
	box.PackStart(c.method, false, false, 5)
	box.PackStart(c.address, true, true, 0)
	box.PackEnd(c.send, false, false, 5)

	return box, nil
}
