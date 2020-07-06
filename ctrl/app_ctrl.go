package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
)

type AppCtrl struct {
	config *config.Config
	App    *gtk.Application
	Engine *sqlengine.Engine

	windows []*WindowCtrl
}

func (c AppCtrl) Init(cfg *config.Config, eng *sqlengine.Engine) *AppCtrl {
	c.config = cfg
	c.Engine = eng
	return &c
}

func (c *AppCtrl) AppID() string {
	return config.AppID
}

func (c *AppCtrl) Config() *config.Config {
	return c.config
}

func (c *AppCtrl) OnActivate() {
	c.OnNewWindow()

	// TODO: every double click is triggering this handler
	//c.factory.OnTabClick(c.onNotebookDoubleClick)

	//c.factory.Show()
}

func (c *AppCtrl) OnNewWindow() {
	err := c.createWindow()
	if err != nil {
		panic(err)
	}
}
func (c *AppCtrl) OnPreferences() {
	c.config.Print("launch preferences modal")
}

func (c *AppCtrl) OnShutdown() {
	c.config.Debug("application shutdown")
}

func (c *AppCtrl) OnStartup() {
	c.config.Debug("application started")
}

func (c *AppCtrl) createWindow() error {
	window, err := WindowCtrl{}.Init(c)
	if err != nil {
		return err
	}
	c.windows = append(c.windows, window)

	window.Show()

	return nil
}

func (c *AppCtrl) ShowAll() {
	for _, w := range c.windows {
		w.Show()
	}
}
func (c *AppCtrl) HideAll() {
	for _, w := range c.windows {
		w.Hide()
	}
}
