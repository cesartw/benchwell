package gtk

import (
	"fmt"
	"strings"

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

type HTTPScreen struct {
	*gtk.Paned
	w *Window

	// Address bar
	method  *gtk.ComboBoxText
	address *gtk.Entry
	send    *gtk.Button

	// request body
	body     *SourceView
	bodySize *gtk.Label
	headers  *KeyValues
	params   *KeyValues
	bodyMime string

	// response
	response *SourceView
	status   *gtk.Label
	duration *gtk.Label
	respSize *gtk.Label

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
	main.SetHExpand(true)
	main.SetVExpand(true)

	addressBar, err := h.buildAddressBar()
	if err != nil {
		return nil, err
	}

	request, err := h.buildRequest()
	if err != nil {
		return nil, err
	}
	response, err := h.buildResponse()
	if err != nil {
		return nil, err
	}

	vPaned, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	vPaned.SetVExpand(true)
	vPaned.SetHExpand(true)
	vPaned.SetWideHandle(true)

	vPaned.Add1(request)
	vPaned.Add2(response)

	main.PackStart(addressBar, false, false, 0)
	main.PackStart(vPaned, true, true, 0)

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
	box.SetName("AddressBar")

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
	box.PackStart(h.method, false, false, 0)
	box.PackStart(h.address, true, true, 0)
	box.PackEnd(h.send, false, false, 0)

	return box, nil
}

func (h *HTTPScreen) buildRequest() (*gtk.Notebook, error) {
	nb, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}

	h.body, err = SourceView{}.Init(h.w, SourceViewOptions{true, true, "json"}, h.ctrl)
	if err != nil {
		return nil, err
	}
	h.body.SetHExpand(true)
	h.body.SetVExpand(true)
	h.body.SetProperty("highlight-current-line", false)
	h.body.SetProperty("show-line-numbers", true)
	h.body.SetProperty("show-right-margin", true)
	h.body.SetProperty("auto-indent", true) // it just keeps the current indent
	h.body.SetProperty("show-line-marks", true)
	h.body.SetProperty("highlight-current-line", true)
	buff, err := h.body.GetBuffer()
	if err != nil {
		return nil, err
	}
	buff.Connect("changed", func() {
		start := buff.GetStartIter()
		end := buff.GetEndIter()
		txt, err := buff.GetText(start, end, false)
		if err != nil {
			return
		}
		h.bodySize.SetText(fmt.Sprintf("%dKB", len([]byte(txt))/1028))
	})

	bodySW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	bodySW.Add(h.body)

	h.bodySize, err = BWLabelNewWithClass("0KB", "tag")
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

	bodyLabel, err := gtk.LabelNew("Body")
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

	mimeOptions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	mimeOptions.SetName("BodyMimeOpts")

	bodyBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	bodyBox.SetHExpand(true)
	bodyBox.SetVExpand(true)

	bodyBox.PackStart(mimeOptions, false, false, 0)
	bodyBox.PackStart(bodySW, true, true, 0)

	cbMime, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	cbMime.AppendText("JSON")
	cbMime.AppendText("PLAIN")
	cbMime.AppendText("HTML")
	cbMime.AppendText("XML")
	cbMime.AppendText("YAML")
	cbMime.Connect("changed", func() {
		h.bodyMime = strings.ToLower(cbMime.GetActiveText())

		if h.bodyMime == "plain" {
			err = h.body.SetLanguage("")
			return
		}

		err := h.body.SetLanguage(h.bodyMime)
		if err != nil {
			h.ctrl.Config().Error(err, "setting language")
		}
	})

	cbMime.SetActive(0)

	mimeOptions.PackStart(cbMime, false, false, 0)
	mimeOptions.PackEnd(h.bodySize, false, false, 0)

	nb.AppendPage(bodyBox, bodyLabel)
	nb.AppendPage(h.params, paramsLabel)
	nb.AppendPage(h.headers, headersLabel)

	return nb, nil
}

func (h *HTTPScreen) buildResponse() (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	h.response, err = SourceView{}.Init(h.w, SourceViewOptions{true, true, ""}, h.ctrl)
	if err != nil {
		return nil, err
	}
	h.response.SetHExpand(true)
	h.response.SetVExpand(true)
	h.response.SetProperty("highlight-current-line", false)
	h.response.SetProperty("show-line-numbers", true)
	h.response.SetProperty("show-right-margin", true)
	h.response.SetProperty("auto-indent", true) // it just keeps the current indent
	h.response.SetProperty("show-line-marks", true)
	h.response.SetProperty("highlight-current-line", true)
	buff, err := h.response.GetBuffer()
	if err != nil {
		return nil, err
	}
	buff.Connect("insert-text", func() bool {
		return false
	})

	details, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	details.SetName("ResponseDetails")

	h.status, err = BWLabelNewWithClass("200 OK", "tag")
	if err != nil {
		return nil, err
	}
	h.status.SetName("Status")
	BWAddClass(h.status, "success")

	h.duration, err = BWLabelNewWithClass("0ms", "tag")
	if err != nil {
		return nil, err
	}
	h.duration.SetName("Duration")

	h.respSize, err = BWLabelNewWithClass("0KB", "tag")
	if err != nil {
		return nil, err
	}

	details.PackStart(h.status, false, false, 0)
	details.PackStart(h.duration, false, false, 0)
	details.PackEnd(h.respSize, false, false, 0)

	box.PackStart(details, false, false, 0)
	box.PackStart(h.response, true, true, 0)

	return box, nil
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
