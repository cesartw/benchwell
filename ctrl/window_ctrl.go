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
	tabs   []*ConnectionTabCtrl
}

func (c WindowCtrl) Init(parent *AppCtrl) (*WindowCtrl, error) {
	var err error
	ctrl := &c
	ctrl.AppCtrl = parent
	ctrl.window, err = gtk.Window{}.Init(parent.app.Application)
	if err != nil {
		return nil, err
	}

	// add main tab
	c.window.Menu.NewTab.Connect("activate", func() {
		err := c.AddTab()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}
		c.window.PushStatus("Ready")
	})

	// action menu for sub tabs
	ctrl.window.Menu.NewSubTab.Connect("activate", func() {
		err := ctrl.currentTab().AddTab()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}
		ctrl.window.PushStatus("Ready")
	})

	ctrl.window.Menu.CloseTab.Connect("activate", func() {
		if ctrl.currentTab().Close() {
			return
		}

		i := ctrl.window.CurrentPage()

		ctrl.tabs = append(ctrl.tabs[i:], ctrl.tabs[:i+1]...)
		ctrl.window.RemoveCurrentPage()
	})

	return ctrl, ctrl.AddTab()
}

func (c *WindowCtrl) Show() {
	c.window.Show()
}

func (c *WindowCtrl) Hide() {
	c.window.Hide()
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

func (c *WindowCtrl) AddTab() error {
	tab, err := ConnectionTabCtrl{}.Init(c)
	if err != nil {
		return err
	}

	tab.tab.Show()
	c.window.AddTab(tab.tabLabel, tab.tab, tab.Removed)
	c.tabs = append(c.tabs, tab)

	return nil
}

func (c *WindowCtrl) currentTab() *ConnectionTabCtrl {
	return c.tabs[c.window.CurrentPage()]
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
