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
	//selectedItem       *config.HTTPItem
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

	var err error
	c.selectedCollection = c.scr.GetSelectedCollection()
	if err != nil {
		return
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

	item := c.scr.GetSelectedItem()
	if item == nil {
		c.window.PushStatus("no item found")
		return
	}

	err := item.LoadFull()
	if err != nil {
		c.window.PushStatus("loading item: " + err.Error())
		return
	}

	if item.IsFolder {
		if item.Loaded {
			return
		}
		item.Loaded = true

		path, err := c.scr.GetSelectedPath()
		if err != nil {
			return
		}

		c.scr.LoadFolder(path, item)
	} else {
		c.scr.SetRequest(item)
	}
}

func (c *HTTPTabCtrl) OnDeleteItem() {
	defer config.LogStart("HTTPTabCtrl.OnDeleteItem", nil)()
	item := c.scr.GetSelectedItem()
	if item == nil {
		return
	}

	path, err := c.scr.GetSelectedPath()
	if err != nil {
		return
	}

	err = item.Delete()
	if err != nil {
		c.window.PushStatus("couldn't delete item: " + err.Error())
		return
	}

	c.scr.RemoveItem(path)
}

func (c *HTTPTabCtrl) OnNewFolder() {
}

func (c *HTTPTabCtrl) OnNewRequest() {
	c.OnLoadItem()

	path, err := c.scr.GetSelectedPath()
	if err != nil {
		return
	}

	selectedItem := c.scr.GetSelectedItem()

	item := &config.HTTPItem{HTTPRequest: config.HTTPRequest{Method: "GET"}, HTTPCollectionID: c.selectedCollection.ID}
	if selectedItem != nil && selectedItem.IsFolder {
		item.ParentID = selectedItem.ID
		selectedItem.Items = append(selectedItem.Items, item)
	}

	if err := item.Save(); err != nil {
		c.window.PushStatus("failed to save item: %s", err.Error())
		return
	}

	c.scr.SetRequest(item)
	c.scr.AddItem(item, path)
}

func (c *HTTPTabCtrl) OnSave(item *config.HTTPItem) error {
	defer config.LogStart("HTTPTabCtrl.OnSave", nil)()

	if err := item.Save(); err != nil {
		c.window.PushStatus("failed to save item: %s", err.Error())
		return err
	}

	c.window.PushStatus("Saved")

	return nil
}

func (c *HTTPTabCtrl) OnSaveAs() error {
	defer config.LogStart("HTTPTabCtrl.OnSaveAs", nil)()
	return nil
}

func (c *HTTPTabCtrl) OnSend() {
	defer config.LogStart("HTTPTabCtrl.OnSend", nil)()

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
