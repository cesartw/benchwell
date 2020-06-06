package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
)

type AppCtrl struct {
	Config *config.Config
	App    *gtk.Application
	Engine *sqlengine.Engine

	windows []*WindowCtrl
}

func (c *AppCtrl) AppID() string {
	return config.AppID
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
	config.Env.Log.Print("launch preferences modal")
}

func (c *AppCtrl) OnShutdown() {
	config.Env.Log.Debug("application shutdown")
}

func (c *AppCtrl) OnStartup() {
	config.Env.Log.Debug("application started")
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
