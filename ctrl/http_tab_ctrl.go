package ctrl

import (
	"net/http"

	"bitbucket.org/goreorto/benchwell/gtk"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type HTTPTabCtrl struct {
	*WindowCtrl
	scr    *gtk.HTTPScreen
	client *http.Client
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

func (c *HTTPTabCtrl) Save()   {}
func (c *HTTPTabCtrl) SaveAs() {}
func (c *HTTPTabCtrl) Send() {
	req, err := c.scr.GetRequest()
	if err != nil {
		c.window.PushStatus("getting request: ", err.Error())
		return
	}

	//func NewRequest(method, url string, body io.Reader) (*Request, error) {
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

}

func (c *HTTPTabCtrl) Close()                    {}
func (c *HTTPTabCtrl) Removed()                  {}
func (c *HTTPTabCtrl) Title() string             { return "HTTP" }
func (c *HTTPTabCtrl) Content() ggtk.IWidget     { return c.scr }
func (c *HTTPTabCtrl) SetFileText(string)        {}
func (c *HTTPTabCtrl) OnCloseTab()               {}
func (c *HTTPTabCtrl) SetWindowCtrl(interface{}) {}
