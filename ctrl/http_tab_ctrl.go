package ctrl

import (
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type HTTPTabCtrl struct {
	*WindowCtrl
	scr                *gtk.HTTPScreen
	client             *http.Client
	selectedCollection *config.HTTPCollection
}

func (c HTTPTabCtrl) Init(p *WindowCtrl) (*HTTPTabCtrl, error) {
	c.WindowCtrl = p
	c.client = &http.Client{}

	var err error
	c.scr, err = gtk.HTTPScreen{}.Init(c.window, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *HTTPTabCtrl) OnCollectionSelected() {
	id, err := c.scr.GetSelectedCollectionID()
	if err != nil {
		c.window.PushStatus("getting collection: " + err.Error())
		return
	}

	if id == 0 {
		return
	}

	for _, collection := range c.config.Collections {
		if collection.ID != id {
			continue
		}
		c.selectedCollection = collection
		break
	}

	if c.selectedCollection == nil {
		c.window.PushStatus("collection not found")
		return
	}

	err = c.selectedCollection.LoadRootItems()
	if err != nil {
		c.window.PushStatus("Error loading collection: " + err.Error())
		return
	}

	err = c.scr.LoadCollection(c.selectedCollection.Items)
	if err != nil {
		c.window.PushStatus("Error loading item: " + err.Error())
	}
}

func (c *HTTPTabCtrl) OnLoadItem() {
	itemID, path, err := c.scr.GetSelectedItemID()
	if err != nil {
		c.window.PushStatus("getting iter: " + err.Error())
		return
	}

	var item *config.HTTPItem
	items := c.selectedCollection.Items

	for _, i := range items {
		if found := i.SearchID(itemID); found != nil {
			item = found
			break
		}
	}
	if item == nil {
		c.window.PushStatus("no item found")
		return
	}

	err = item.LoadFull()
	if err != nil {
		c.window.PushStatus("loading item: " + err.Error())
		return
	}

	if item.IsFolder {
		c.scr.LoadFolder(path, item)
	} else {
		c.scr.SetRequest(item)
	}
}

func (c *HTTPTabCtrl) Save()   {}
func (c *HTTPTabCtrl) SaveAs() {}
func (c *HTTPTabCtrl) Send() {
	req, err := c.scr.GetRequest()
	if err != nil {
		c.window.PushStatus("getting request: ", err.Error())
		return
	}

	httpreq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		c.window.PushStatus("building request: ", err.Error())
		return
	}

	for header, values := range req.Headers {
		for _, value := range values {
			httpreq.Header.Add(header, value)
		}
	}

	now := time.Now()
	resp, err := c.client.Do(httpreq)
	if err != nil {
		c.window.PushStatus("failed: ", err.Error())
		return
	}
	duration := time.Since(now)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.window.PushStatus("reading body: ", err.Error())
		return
	}

	c.scr.SetResponse(string(b), resp.Header, duration)
}

func (c *HTTPTabCtrl) Close()                    {}
func (c *HTTPTabCtrl) Removed()                  {}
func (c *HTTPTabCtrl) Title() string             { return "HTTP" }
func (c *HTTPTabCtrl) Content() ggtk.IWidget     { return c.scr }
func (c *HTTPTabCtrl) SetFileText(string)        {}
func (c *HTTPTabCtrl) OnCloseTab()               {}
func (c *HTTPTabCtrl) SetWindowCtrl(interface{}) {}
