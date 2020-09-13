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
	collection  *HTTPCollection
	currentItem *config.HTTPItem

	ctrl httpScreenCtrl
}

type httpScreenCtrl interface {
	OnSave(*config.HTTPItem) error
	OnSaveAs() error
	OnSend()
	OnLoadItem()
	OnCollectionSelected()
	OnDeleteItem()
	OnNewRequest()
	OnNewFolder()
}

func (h HTTPScreen) Init(w *Window, ctrl httpScreenCtrl) (*HTTPScreen, error) {
	defer config.LogStart("HTTPScreen.Init", nil)()

	var err error

	h.w = w
	h.ctrl = ctrl
	h.currentItem = &config.HTTPItem{
		HTTPRequest: config.HTTPRequest{
			Headers: []*config.HTTPKV{},
			Params:  []*config.HTTPKV{},
		},
	}

	h.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	h.Paned.Show()
	h.Paned.SetWideHandle(true)

	h.collection, err = HTTPCollection{}.Init(w, &h, ctrl)
	if err != nil {
		return nil, err
	}
	h.collection.Show()

	main, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	main.Show()
	main.SetHExpand(true)
	main.SetVExpand(true)

	addressBar, err := h.buildAddressBar()
	if err != nil {
		return nil, err
	}
	addressBar.Show()

	request, err := h.buildRequest()
	if err != nil {
		return nil, err
	}
	request.Show()

	response, err := h.buildResponse()
	if err != nil {
		return nil, err
	}
	response.Show()

	vPaned, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	vPaned.Show()
	vPaned.SetVExpand(true)
	vPaned.SetHExpand(true)
	vPaned.SetWideHandle(true)

	vPaned.Add1(request)
	vPaned.Add2(response)

	main.PackStart(addressBar, false, false, 0)
	main.PackStart(vPaned, true, true, 0)

	h.Paned.Pack1(h.collection, false, false)
	h.Paned.Pack2(main, false, true)

	return &h, nil
}

func (h *HTTPScreen) AddItem(item *config.HTTPItem, path *gtk.TreePath) {
	h.collection.AddItem(item, path)
}

func (h *HTTPScreen) GetSelectedItem() *config.HTTPItem {
	defer config.LogStart("HTTPScreen.GetSelectedItem", nil)()

	return h.collection.GetSelectedItem()
}

func (h *HTTPScreen) GetSelectedCollection() *config.HTTPCollection {
	defer config.LogStart("HTTPScreen.GetSelectedCollection", nil)()
	return h.collection.GetSelectedCollection()
}

func (h *HTTPScreen) GetSelectedPath() (*gtk.TreePath, error) {
	defer config.LogStart("HTTPScreen.GetSelectedPath", nil)()
	return h.collection.GetSelectedPath()
}

func (h *HTTPScreen) LoadCollection(items []*config.HTTPItem) error {
	defer config.LogStart("HTTPScreen.LoadCollection", nil)()

	return h.collection.LoadCollection(items)
}

func (h *HTTPScreen) LoadFolder(path *gtk.TreePath, item *config.HTTPItem) error {
	defer config.LogStart("HTTPScreen.LoadFolder", nil)()

	return h.collection.LoadFolder(path, item)
}

func (h *HTTPScreen) RemoveItem(path *gtk.TreePath) {
	h.collection.RemoveItem(path)
}

func (h *HTTPScreen) buildAddressBar() (*gtk.Box, error) {
	defer config.LogStart("HTTPScreen.buildAddressBar", nil)()

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}
	box.SetName("AddressBar")

	h.method, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	h.method.Show()
	h.method.SetActive(0)
	h.method.Connect("changed", func() {
		if h.currentItem == nil {
			return
		}
		h.currentItem.Method = h.method.GetActiveText()
	})

	for _, m := range methods {
		h.method.Append(strings.ToLower(m), m)
	}

	h.address, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	h.address.Show()
	h.address.SetPlaceholderText("http://localhost/path.json")
	h.address.Connect("key-release-event", h.onAddressChange)

	h.send, err = gtk.ButtonNewWithLabel("SEND")
	if err != nil {
		return nil, err
	}
	h.send.Show()
	h.send.Connect("clicked", h.ctrl.OnSend)
	BWAddClass(h.send, "suggested-action")

	h.save, err = BWOptionButtonNew("SAVE", h.w, []string{"Save as", "win.saveas"})
	if err != nil {
		return nil, err
	}
	h.save.Show()
	h.save.ConnectAction("win.saveas", h.ctrl.OnSaveAs)
	h.save.btn.Connect("clicked", func() {
		err := h.ctrl.OnSave(h.currentItem)
		if err != nil {
			return
		}

		path, _ := h.GetSelectedPath()
		item := h.GetSelectedItem()
		iter, _ := h.collection.store.GetIter(path)
		h.collection.store.SetValue(iter, COLUMN_METHOD, item.Method)
	})

	box.PackStart(h.method, false, false, 0)
	box.PackStart(h.address, true, true, 0)
	box.PackEnd(h.save, false, false, 0)
	box.PackEnd(h.send, false, false, 0)

	return box, nil
}

func (h *HTTPScreen) buildRequest() (*gtk.Notebook, error) {
	defer config.LogStart("HTTPScreen.buildRequest", nil)()

	nb, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}

	h.body, err = SourceView{}.Init(h.w, SourceViewOptions{true, true, "json"}, h.ctrl)
	if err != nil {
		return nil, err
	}
	h.body.Show()
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
	bodySW.Show()
	bodySW.Add(h.body)

	h.bodySize, err = BWLabelNewWithClass("0KB", "tag")
	if err != nil {
		return nil, err
	}
	h.bodySize.Show()

	h.headers, err = KeyValues{}.Init(&h.currentItem.Headers, func() {})
	if err != nil {
		return nil, err
	}
	h.headers.Show()

	h.params, err = KeyValues{}.Init(&h.currentItem.Params, func() {
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
	h.params.Show()

	bodyLabel, err := gtk.LabelNew("Body")
	if err != nil {
		return nil, err
	}
	bodyLabel.Show()

	paramsLabel, err := gtk.LabelNew("Params")
	if err != nil {
		return nil, err
	}
	paramsLabel.Show()

	headersLabel, err := gtk.LabelNew("Headers")
	if err != nil {
		return nil, err
	}
	headersLabel.Show()

	mimeOptions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	mimeOptions.Show()
	mimeOptions.SetName("BodyMimeOpts")

	bodyBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	bodyBox.Show()
	bodyBox.SetHExpand(true)
	bodyBox.SetVExpand(true)

	bodyBox.PackStart(mimeOptions, false, false, 0)
	bodyBox.PackStart(bodySW, true, true, 0)

	h.mime, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	h.mime.Show()
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
	defer config.LogStart("HTTPScreen.buildResponse", nil)()

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	h.response, err = SourceView{}.Init(h.w, SourceViewOptions{true, true, ""}, h.ctrl)
	if err != nil {
		return nil, err
	}
	h.response.Show()
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
	details.Show()
	details.SetName("ResponseDetails")

	h.status, err = BWLabelNewWithClass("200 OK", "tag")
	if err != nil {
		return nil, err
	}
	h.status.Show()
	h.status.SetName("Status")
	BWAddClass(h.status, "success")

	h.duration, err = BWLabelNewWithClass("0ms", "tag")
	if err != nil {
		return nil, err
	}
	h.duration.Show()
	h.duration.SetName("Duration")

	h.respSize, err = BWLabelNewWithClass("0KB", "tag")
	if err != nil {
		return nil, err
	}
	h.respSize.Show()
	responseSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	responseSW.Show()
	responseSW.Add(h.response)

	details.PackStart(h.status, false, false, 0)
	details.PackStart(h.duration, false, false, 0)
	details.PackEnd(h.respSize, false, false, 0)

	box.PackStart(details, false, false, 0)
	box.PackStart(responseSW, true, true, 0)

	return box, nil
}

func (h *HTTPScreen) onAddressChange() {
	defer config.LogStart("HTTPScreen.onAddressChange", nil)()

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
	h.currentItem.URL = address

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
			h.params.Add(&config.HTTPKV{Var: config.Var{Key: key, Value: value, Enabled: true}})
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
	defer config.LogStart("HTTPScreen.GetRequest", nil)()

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
	defer config.LogStart("HTTPScreen.SetResponse", nil)()

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
	defer config.LogStart("HTTPScreen.SetResponse", nil)()
	h.currentItem = req

	h.address.SetText(req.URL)
	h.method.SetActiveID(strings.ToLower(req.Method))
	h.headers.Clear()
	for _, kv := range req.Headers {
		h.headers.Add(kv)
	}
	for _, kv := range req.Params {
		h.params.Add(kv)
	}

	h.mime.SetActiveID(req.Mime)
}

func (h *HTTPScreen) CurrentItem() *config.HTTPItem {
	return h.currentItem
}
