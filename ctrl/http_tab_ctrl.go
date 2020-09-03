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
	defer config.LogStart("HTTPTabCtrl.Init", nil)()

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
	defer config.LogStart("HTTPTabCtrl.OnCollectionSelected", nil)()

	id, err := c.scr.GetSelectedCollectionID()
	if err != nil {
		c.window.PushStatus("getting collection: " + err.Error())
		return
	}

	if id == 0 {
		return
	}

	for _, collection := range config.Collections {
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
	defer config.LogStart("HTTPTabCtrl.OnLoadItem", nil)()

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
		if item.Loaded {
			return
		}
		item.Loaded = true
		c.scr.LoadFolder(path, item)
	} else {
		c.scr.SetRequest(item)
	}
}

func (c *HTTPTabCtrl) Save() {
	defer config.LogStart("HTTPTabCtrl.Save", nil)()

}

func (c *HTTPTabCtrl) SaveAs() {
	defer config.LogStart("HTTPTabCtrl.SaveAs", nil)()

}

func (c *HTTPTabCtrl) Send() {
	defer config.LogStart("HTTPTabCtrl.Send", nil)()

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

func (c *HTTPTabCtrl) Close() {
	defer config.LogStart("HTTPTabCtrl.Close", nil)()

}

func (c *HTTPTabCtrl) Removed() {
	defer config.LogStart("HTTPTabCtrl.Removed", nil)()

}

func (c *HTTPTabCtrl) Title() string {
	defer config.LogStart("HTTPTabCtrl.Title", nil)()

	return "HTTP"
}

func (c *HTTPTabCtrl) Content() ggtk.IWidget {
	defer config.LogStart("HTTPTabCtrl.Content", nil)()

	return c.scr
}

func (c *HTTPTabCtrl) SetFileText(string) {
	defer config.LogStart("HTTPTabCtrl.SetFileText", nil)()

}

func (c *HTTPTabCtrl) SetWindowCtrl(interface{}) {
	defer config.LogStart("HTTPTabCtrl.SetWindowCtrl", nil)()

}
