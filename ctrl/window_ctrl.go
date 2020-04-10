package ctrl

import (
	"github.com/gotk3/gotk3/gdk"
	ggtk "github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
)

type WindowCtrl struct {
	*AppCtrl
	window *gtk.Window
}

func (c WindowCtrl) Init(parent *AppCtrl) (*WindowCtrl, error) {
	var err error
	c.AppCtrl = parent
	c.window, err = gtk.Window{}.Init(parent.app.Application)
	if err != nil {
		return nil, err
	}

	c.window.Menu.NewTab.Connect("activate", func() {
		err := c.AddTab()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}
		c.window.PushStatus("Ready")
	})
	c.window.Menu.CloseTab.Connect("activate", func() {
		c.window.RemoveCurrentPage()
	})

	return &c, c.AddTab()
}

func (c *WindowCtrl) Show() {
	c.window.Show()
}

func (c *WindowCtrl) OnActivate() {
	err := c.AddTab()
	if err != nil {
		c.window.PushStatus(err.Error())
	} else {
		c.window.PushStatus("Ready")
	}

	c.window.Show()
}

func (c *WindowCtrl) onNotebookDoubleClick(_ *ggtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_PRIMARY {
		return
	}
	if keyEvent.Type() != gdk.EVENT_2BUTTON_PRESS {
		return
	}

	if err := c.AddTab(); err != nil {
		config.Env.Log.Error(err)
	}
}

func (c *WindowCtrl) AddTab() error {
	tab, err := TabCtrl{}.Init(c)
	if err != nil {
		return err
	}

	tab.tab.Show()
	c.window.AddTab(tab.tabLabel, tab.tab)

	return nil
}
