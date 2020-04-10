package ctrl

import (
	"errors"

	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
)

type Options struct {
	Config *config.Config
	Log    *logrus.Logger
	App    *gtk.App
	Engine *sqlengine.Engine
}

func (o Options) Valid() error {
	if o.Config == nil {
		return errors.New("Config is required")
	}

	if o.App == nil {
		return errors.New("App is required")
	}

	if o.Engine == nil {
		return errors.New("Engine is required")
	}

	return nil
}

type AppCtrl struct {
	config *config.Config
	app    *gtk.App
	engine *sqlengine.Engine

	windows []*WindowCtrl
}

func (c AppCtrl) Init(opts Options) (*AppCtrl, error) {
	if err := opts.Valid(); err != nil {
		return nil, err
	}
	c.engine = opts.Engine
	c.app = opts.App
	c.config = opts.Config

	return &c, nil
}

func (c *AppCtrl) OnActivate() {
	err := c.CreateWindow()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.app.Menu.Application.NewWindow.Connect("activate", func() {
		err := c.CreateWindow()
		if err != nil {
			config.Env.Log.Error(err)
		}
	})
	c.app.Menu.Application.Preferences.Connect("activate", func() {
		config.Env.Log.Print("launch preferences modal")
	})

	// TODO: every double click is triggering this handler
	//c.factory.OnTabClick(c.onNotebookDoubleClick)

	//c.factory.Show()
}

func (c *AppCtrl) CreateWindow() error {
	window, err := WindowCtrl{}.Init(c)
	if err != nil {
		return err
	}
	c.windows = append(c.windows, window)

	window.Show()

	return nil
}
