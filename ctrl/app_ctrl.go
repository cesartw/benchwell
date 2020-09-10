package ctrl

import (
	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
	"bitbucket.org/goreorto/benchwell/sqlengine"
)

type AppCtrl struct {
	App    *gtk.Application
	Engine *sqlengine.Engine

	windows []*WindowCtrl
}

func (c AppCtrl) Init(eng *sqlengine.Engine) *AppCtrl {
	defer config.LogStart("AppCtrl.Init", nil)()

	c.Engine = eng
	return &c
}

func (c *AppCtrl) AppID() string {
	defer config.LogStart("AppCtrl.AppID", nil)()

	return config.AppID
}

func (c *AppCtrl) OnActivate() {
	defer config.LogStart("AppCtrl.OnActivate", nil)()

	c.OnNewWindow()

	// TODO: every double click is triggering this handler
	//c.factory.OnTabClick(c.onNotebookDoubleClick)

	//c.factory.Show()
}

func (c *AppCtrl) OnNewWindow() {
	defer config.LogStart("AppCtrl.OnNewWindow", nil)()

	err := c.createWindow()
	if err != nil {
		panic(err)
	}
}
func (c *AppCtrl) OnPreferences() {
	defer config.LogStart("AppCtrl.OnPreferences", nil)()
}

func (c *AppCtrl) OnShutdown() {
	defer config.LogStart("AppCtrl.OnShutdown", nil)()
}

func (c *AppCtrl) OnStartup() {
	defer config.LogStart("AppCtrl.OnStartup", nil)()
}

func (c *AppCtrl) createWindow() error {
	defer config.LogStart("AppCtrl.createWindow", nil)()

	window, err := WindowCtrl{}.Init(c)
	if err != nil {
		return err
	}
	c.windows = append(c.windows, window)

	window.Show()
	c.App.AddWindow(window.window)

	return nil
}

func (c *AppCtrl) ShowAll() {
	defer config.LogStart("AppCtrl.ShowAll", nil)()

	for _, w := range c.windows {
		w.Show()
	}
}

func (c *AppCtrl) HideAll() {
	defer config.LogStart("AppCtrl.HideAll", nil)()

	for _, w := range c.windows {
		w.Hide()
	}
}
