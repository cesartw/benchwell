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
	tabs   []*WindowTabCtrl
}

func (c WindowCtrl) Init(parent *AppCtrl) (*WindowCtrl, error) {
	var err error
	ctrl := &c
	ctrl.AppCtrl = parent
	ctrl.window, err = gtk.Window{}.Init(parent.App.Application, &c)
	if err != nil {
		return nil, err
	}

	return ctrl, ctrl.AddTab()
}

func (c *WindowCtrl) OnNewSubTab() {
	err := c.currentWindowTab().AddTab()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	c.window.PushStatus("Ready")
}

func (c *WindowCtrl) OnCloseTab() {
	if c.currentWindowTab().Close() {
		return
	}

	i := c.window.CurrentPage()

	c.tabs = append(c.tabs[i:], c.tabs[:i+1]...)
	c.window.RemoveCurrentPage()
}

func (c *WindowCtrl) OnNewTab() {
	err := c.AddTab()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	c.window.PushStatus("Ready")
}

func (c *WindowCtrl) Show() {
	c.window.Show()
}

func (c *WindowCtrl) Hide() {
	c.window.Hide()
}

//func (c *WindowCtrl) OnActivate() {
//err := c.AddTab()
//if err != nil {
//c.window.PushStatus(err.Error())
//} else {
//c.window.PushStatus("Ready")
//}

//c.window.Show()
//}

func (c *WindowCtrl) AddTab() error {
	tab, err := WindowTabCtrl{}.Init(c)
	if err != nil {
		return err
	}

	tab.tab.Show()
	c.window.AddTab(tab.tabLabel, tab.tab, tab.Removed)
	c.tabs = append(c.tabs, tab)

	return nil
}

func (c *WindowCtrl) currentWindowTab() *WindowTabCtrl {
	return c.tabs[c.window.CurrentPage()]
}

// TODO: not used
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
