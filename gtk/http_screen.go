package gtk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

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
	save    *OptionButton

	// request body
	body         *SourceView
	mime         *gtk.ComboBoxText
	bodySize     *gtk.Label
	headers      *KeyValues
	params       *KeyValues
	buildingAddr bool

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
	Save()
	SaveAs()
	Send()
	OnLoadItem()
	OnCollectionSelected()
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
	h.collection, err = HTTPCollection{}.Init(w, &h, ctrl)
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

func (h *HTTPScreen) GetSelectedItemID() (int64, string, error) {
	return h.collection.GetSelectedItemID()
}

func (h *HTTPScreen) GetSelectedCollectionID() (int64, error) {
	return h.collection.GetSelectedCollectionID()
}

func (h *HTTPScreen) LoadCollection(items []*config.HTTPItem) error {
	return h.collection.LoadCollection(items)
}

func (h *HTTPScreen) LoadFolder(at string, item *config.HTTPItem) error {
	return h.collection.LoadFolder(at, item)
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
		h.method.Append(strings.ToLower(m), m)
	}
	h.method.SetActive(0)

	h.address, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	h.address.SetPlaceholderText("http://localhost/path.json")
	h.address.Connect("key-release-event", h.onAddressChange)

	h.send, err = gtk.ButtonNewWithLabel("SEND")
	if err != nil {
		return nil, err
	}
	h.send.Connect("clicked", h.ctrl.Send)

	h.save, err = BWOptionButtonNew("SAVE", h.w, []string{"Save as", "win.saveas"})
	if err != nil {
		return nil, err
	}
	h.save.ConnectAction("win.saveas", h.ctrl.SaveAs)
	h.save.btn.Connect("activate", h.ctrl.Save)

	box.PackStart(h.method, false, false, 0)
	box.PackStart(h.address, true, true, 0)
	box.PackEnd(h.save, false, false, 0)
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

	h.headers, err = KeyValues{}.Init(func() {})
	if err != nil {
		return nil, err
	}

	h.params, err = KeyValues{}.Init(func() {
		if h.buildingAddr {
			return
		}
		h.buildingAddr = true
		defer func() { h.buildingAddr = false }()

		params, _ := h.params.Collect()
		address, _ := h.address.GetText()
		u, err := url.Parse(address)
		if err != nil {
			u = &url.URL{}
			u.Host = address
		}
		u.RawQuery = url.Values(params).Encode()

		h.address.SetText(u.String())
	})
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

	h.mime, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	h.mime.Append("json", "JSON")
	h.mime.Append("plain", "PLAIN")
	h.mime.Append("html", "HTML")
	h.mime.Append("xml", "XML")
	h.mime.Append("yaml", "YAML")
	h.mime.Connect("changed", func() {
		mime := h.mime.GetActiveID()

		if mime == "plain" {
			h.body.SetLanguage("")
			return
		}

		err := h.body.SetLanguage(mime)
		if err != nil {
			config.Error(err, "setting language")
		}
	})

	h.mime.SetActive(0)

	mimeOptions.PackStart(h.mime, false, false, 0)
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
	responseSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	responseSW.Add(h.response)

	details.PackStart(h.status, false, false, 0)
	details.PackStart(h.duration, false, false, 0)
	details.PackEnd(h.respSize, false, false, 0)

	box.PackStart(details, false, false, 0)
	box.PackStart(responseSW, true, true, 0)

	return box, nil
}

func (h *HTTPScreen) onAddressChange() {
	if h.buildingAddr {
		return
	}
	h.buildingAddr = true
	defer func() { h.buildingAddr = false }()

	address, err := h.address.GetText()
	if err != nil {
		h.w.PushStatus("getting address: " + err.Error())
		return
	}

	u, err := url.Parse(address)
	if err != nil {
		h.w.PushStatus("parsing address: " + err.Error())
		return
	}

	h.params.Clear()

	keys := []string{}
	params := u.Query()
	for key, _ := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		for _, value := range params[key] {
			h.params.AddWithValues(key, value, true)
		}
	}
}

type Request struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string][]string
}

func (h *HTTPScreen) GetRequest() (*Request, error) {
	address, err := h.address.GetText()
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method:  h.method.GetActiveText(),
		URL:     address,
		Headers: map[string][]string{},
	}

	headers, err := h.headers.Collect()
	if err != nil {
		return nil, err
	}

	params, err := h.params.Collect()
	if err != nil {
		return nil, err
	}

	req.URL = req.URL + "?" + url.Values(params).Encode()

	buff, err := h.body.GetBuffer()
	if err != nil {
		return nil, err
	}
	start := buff.GetStartIter()
	end := buff.GetStartIter()

	txt, err := buff.GetText(start, end, false)
	if err != nil {
		return nil, err
	}

	req.Body = bytes.NewReader([]byte(txt))
	req.Headers = headers

	return req, nil
}

func (h *HTTPScreen) SetResponse(body string, headers http.Header, duration time.Duration) {
	switch strings.Split(headers.Get("Content-Type"), ";")[0] {
	case "application/json", "text/json":
		err := h.response.SetLanguage("json")
		if err != nil {
			config.Error(err, "setting language")
		}

		var out bytes.Buffer
		json.Indent(&out, []byte(body), "", "\t")
		body = out.String()
	case "text/html":
		err := h.response.SetLanguage("html")
		if err != nil {
			config.Error(err, "setting language")
		}
	case "text/yaml", "application/x-yaml":
		err := h.response.SetLanguage("yaml")
		if err != nil {
			config.Error(err, "setting language")
		}
	case "text/xml", "application/xml":
		err := h.response.SetLanguage("xml")
		if err != nil {
			config.Error(err, "setting language")
		}
	}

	h.respSize.SetText(fmt.Sprintf("%dKB", len([]byte(body))/1024))
	h.duration.SetText(duration.String())

	buff, _ := h.response.GetBuffer()
	buff.SetText(body)
}

func (h *HTTPScreen) SetRequest(req *config.HTTPItem) {
	h.address.SetText(req.URL)
	h.method.SetActiveID(strings.ToLower(req.Method))
	h.headers.Clear()
	for _, kv := range req.Headers {
		h.headers.AddWithValues(kv.Key, kv.Value, true)
	}
	for _, kv := range req.Params {
		h.params.AddWithValues(kv.Key, kv.Value, kv.Enabled)
	}

	h.mime.SetActiveID(req.Mime)
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
	c.enabled.SetActive(true)

	c.Box.PackStart(c.enabled, false, false, 5)
	c.Box.PackStart(c.key, true, true, 0)
	c.Box.PackStart(c.value, true, true, 0)
	c.Box.PackEnd(c.remove, false, false, 5)

	return &c, nil
}

func (c *KeyValue) Get() (string, string, error) {
	key, err := c.key.GetText()
	if err != nil {
		return "", "", err
	}
	value, err := c.value.GetText()
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

type KeyValues struct {
	*gtk.Box
	keyvalues []*KeyValue
	onChange  func()
}

func (c KeyValues) Init(onChange func()) (*KeyValues, error) {
	var err error
	c.onChange = onChange
	c.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	return &c, c.AddEmpty()
}

func (c *KeyValues) AddEmpty() error {
	_, err := c.add()
	if err != nil {
		return err
	}

	return nil
}

func (c *KeyValues) AddWithValues(key, value string, enabled bool) error {
	for _, kv := range c.keyvalues {
		s, _ := kv.key.GetText()
		if s != "" {
			continue
		}
		kv.key.SetText(key)
		kv.value.SetText(value)
		return nil
	}

	kv, err := c.add()
	if err != nil {
		return err
	}
	kv.enabled.SetActive(enabled)

	kv.key.SetText(key)
	kv.value.SetText(value)

	return nil
}

func (c *KeyValues) add() (*KeyValue, error) {
	kv, err := KeyValue{}.Init()
	if err != nil {
		return kv, err
	}

	kv.remove.Connect("clicked", c.onBlur(kv))

	focused := func() {
		if c.keyvalues[len(c.keyvalues)-1] != kv {
			return
		}

		c.AddEmpty()
	}
	kv.key.Connect("grab-focus", focused)
	kv.value.Connect("grab-focus", focused)
	kv.key.Connect("key-release-event", c.onChange)
	kv.value.Connect("key-release-event", c.onChange)
	kv.enabled.Connect("toggled", c.onChange)
	kv.remove.Connect("clicked", c.onChange)

	c.Box.PackStart(kv, false, false, 0)
	c.keyvalues = append(c.keyvalues, kv)
	kv.ShowAll()

	return kv, nil
}

func (c *KeyValues) onBlur(kv *KeyValue) func() {
	return func() {
		c.Box.Remove(kv)

		for i, v := range c.keyvalues {
			if v != kv {
				continue
			}
			c.keyvalues = append(c.keyvalues[:i], c.keyvalues[i+1:]...)
		}

		if len(c.keyvalues) == 0 {
			c.AddEmpty()
		}
	}
}

func (c *KeyValues) Clear() {
	for _, kv := range c.keyvalues {
		c.Remove(kv)
	}
	c.keyvalues = nil
	c.AddEmpty()
}

func (c KeyValues) Collect() (map[string][]string, error) {
	keyvalues := map[string][]string{}

	for _, kv := range c.keyvalues {
		if !kv.enabled.GetActive() {
			continue
		}

		key, value, err := kv.Get()
		if err != nil {
			return nil, err
		}

		if key == "" {
			continue
		}

		keyvalues[key] = append(keyvalues[key], value)
	}

	return keyvalues, nil
}
