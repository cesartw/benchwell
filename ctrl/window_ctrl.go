package ctrl

import (
	"io/ioutil"
	"os"

	ggtk "github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
)

type tab_type int

const (
	TAB_TYPE_DB tab_type = iota
	TAB_TYPE_HTTP
)

type tabCtrl interface {
	Close() bool
	Removed()
	Title() string
	Content() ggtk.IWidget
	SetFileText(string)
	Config() *config.Config
	AddTab() error
	OnCloseTab()
}

type WindowCtrl struct {
	*AppCtrl
	window *gtk.Window
}

func (c WindowCtrl) Init(parent *AppCtrl) (*WindowCtrl, error) {
	var err error
	ctrl := &c
	ctrl.AppCtrl = parent
	ctrl.window, err = gtk.Window{}.Init(parent.App.Application, &c)
	if err != nil {
		return nil, err
	}

	return ctrl, ctrl.AddTab(TAB_TYPE_DB)
}

func (c *WindowCtrl) OnSaveQuery(query, path string) {
	err := ioutil.WriteFile(path, []byte(query), os.FileMode(666))
	if err != nil {
		c.window.PushStatus("failed to save file: %#v", err)
	}
}

func (c *WindowCtrl) OnFileSelected(filepath string) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		c.Config().Error("reading file", err)
		return
	}

	c.currentWindowTab().SetFileText(string(bytes))
}

func (c *WindowCtrl) OnNewSubTab() {
	err := c.currentWindowTab().AddTab()
	if err != nil {
		c.Config().Error(err)
		return
	}
	c.window.PushStatus("Ready")
}

func (c *WindowCtrl) OnNewTab() {
	err := c.AddTab(TAB_TYPE_DB)
	if err != nil {
		c.Config().Error(err)
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

func (c *WindowCtrl) AddTab(t tab_type) error {
	var (
		err  error
		ctrl tabCtrl
	)

	switch t {
	case TAB_TYPE_DB:
		ctrl, err = DbTabCtrl{}.Init(c)
		if err != nil {
			return err
		}
	}

	tab, err := gtk.ToolTab{}.Init(c.window)
	if err != nil {
		return err
	}
	tab.SetContent(gtk.ToolTabOptions{
		Content: ctrl.Content(),
		Title:   ctrl.Title(),
		Ctrl:    ctrl,
	})

	c.window.AddToolTab(tab)

	return nil
}

func (c *WindowCtrl) OnCloseTab() {
	c.currentWindowTab().Close()
	c.window.RemoveCurrentPage()
}

func (c *WindowCtrl) ChangeTitle(title string) {
	c.currentWindowTab().SetTitle(title)
}

func (c *WindowCtrl) currentWindowTab() *gtk.ToolTab {
	return c.window.CurrentTab()
}

// TODO: not used
//func (c *WindowCtrl) onNotebookDoubleClick(_ *ggtk.ListBox, e *gdk.Event) {
//keyEvent := gdk.EventButtonNewFromEvent(e)

//if keyEvent.Button() != gdk.BUTTON_PRIMARY {
//return
//}
//if keyEvent.Type() != gdk.EVENT_2BUTTON_PRESS {
//return
//}

//if err := c.AddTab(); err != nil {
//c.Config().Error(err)
//}
//}
