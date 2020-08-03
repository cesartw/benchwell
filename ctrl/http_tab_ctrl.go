package ctrl

import (
	"bitbucket.org/goreorto/benchwell/gtk"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type HTTPTabCtrl struct {
	*WindowCtrl
	scr *gtk.HTTPScreen
}

func (c HTTPTabCtrl) Init(p *WindowCtrl) (*HTTPTabCtrl, error) {
	c.WindowCtrl = p

	var err error
	c.scr, err = gtk.HTTPScreen{}.Init(c.window, c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
func (c *HTTPTabCtrl) Close()                    {}
func (c *HTTPTabCtrl) Removed()                  {}
func (c *HTTPTabCtrl) Title() string             { return "HTTP" }
func (c *HTTPTabCtrl) Content() ggtk.IWidget     { return c.scr }
func (c *HTTPTabCtrl) SetFileText(string)        {}
func (c *HTTPTabCtrl) OnCloseTab()               {}
func (c *HTTPTabCtrl) SetWindowCtrl(interface{}) {}
