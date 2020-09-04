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
	Close()
	Title() string
	Removed()
	Content() ggtk.IWidget
	SetFileText(string)
	OnCloseTab(string)
	SetWindowCtrl(interface{})
}

type WindowCtrl struct {
	*AppCtrl
	window *gtk.Window
}

func (c WindowCtrl) Init(parent *AppCtrl) (*WindowCtrl, error) {
	defer config.LogStart("WindowCtrl.Init", nil)()

	var err error
	c.AppCtrl = parent
	c.window, err = gtk.Window{}.Init(parent.App.Application, &c)
	if err != nil {
		return nil, err
	}
	_, err = c.AddTab(TAB_TYPE_HTTP)
	if err != nil {
		return nil, err
	}
	_, err = c.AddTab(TAB_TYPE_DB)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *WindowCtrl) OnSaveQuery(query, path string) {
	defer config.LogStart("WindowCtrl.OnSaveQuery", nil)()

	err := ioutil.WriteFile(path, []byte(query), os.FileMode(666))
	if err != nil {
		c.window.PushStatus("failed to save file: %#v", err)
	}
}

func (c *WindowCtrl) OnFileSelected(filepath string) {
	defer config.LogStart("WindowCtrl.OnFileSelected", nil)()

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		config.Error("reading file", err)
		return
	}

	c.currentWindowTab().SetFileText(string(bytes))
}

func (c *WindowCtrl) OnNewDatabaseTab() {
	defer config.LogStart("WindowCtrl.OnNewDatabaseTab", nil)()

	_, err := c.AddTab(TAB_TYPE_DB)
	if err != nil {
		config.Error(err)
		return
	}
	c.window.PushStatus("Ready")
}

func (c *WindowCtrl) OnNewHTTPTab() {
	defer config.LogStart("WindowCtrl.OnNewHTTPTab", nil)()

	_, err := c.AddTab(TAB_TYPE_HTTP)
	if err != nil {
		config.Error(err)
		return
	}
	c.window.PushStatus("Ready")
}

func (c *WindowCtrl) Show() {
	defer config.LogStart("WindowCtrl.Show", nil)()

	c.window.Show()
}

func (c *WindowCtrl) Hide() {
	defer config.LogStart("WindowCtrl.Hide", nil)()

	c.window.Hide()
}

func (c *WindowCtrl) AddTab(t tab_type) (interface{}, error) {
	defer config.LogStart("WindowCtrl.AddTab", nil)()

	var (
		err  error
		ctrl tabCtrl
	)

	switch t {
	case TAB_TYPE_DB:
		ctrl, err = DbTabCtrl{}.Init(c)
		if err != nil {
			return nil, err
		}
	case TAB_TYPE_HTTP:
		ctrl, err = HTTPTabCtrl{}.Init(c)
		if err != nil {
			return nil, err
		}
	}

	tab, err := gtk.ToolTab{}.Init(c.window)
	if err != nil {
		return nil, err
	}
	tab.SetContent(gtk.ToolTabOptions{
		Content: ctrl.Content(),
		Ctrl:    ctrl,
	})

	c.window.AddToolTab(tab)

	return ctrl, nil
}

// click on main tab close
func (c *WindowCtrl) OnCloseTab(id string) {
	defer config.LogStart("WindowCtrl.OnCloseTab", nil)()

	// tell the tool tab that we closing it
	c.currentWindowTab().Close()
	c.window.RemovePage(id)
}

func (c *WindowCtrl) ChangeTitle(title string) {
	defer config.LogStart("WindowCtrl.ChangeTitle", nil)()

	c.currentWindowTab().SetTitle(title)
}

func (c *WindowCtrl) currentWindowTab() *gtk.ToolTab {
	defer config.LogStart("WindowCtrl.currentWindowTab", nil)()

	return c.window.CurrentTab()
}
