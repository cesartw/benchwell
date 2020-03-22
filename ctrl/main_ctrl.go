package ctrl

import (
	"errors"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	ggtk "github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Config  *config.Config
	Log     *logrus.Logger
	Factory *gtk.Factory
	Engine  *sqlengine.Engine
}

func (o Options) Valid() error {
	if o.Config == nil {
		return errors.New("Config is required")
	}

	if o.Factory == nil {
		return errors.New("Factory is required")
	}

	if o.Engine == nil {
		return errors.New("Engine is required")
	}

	return nil
}

type MainCtrl struct {
	config  *config.Config
	log     *logrus.Logger
	factory *gtk.Factory
	engine  *sqlengine.Engine

	currentCtrl interface {
		Screen() interface{}
	}
}

func (c MainCtrl) Init(opts Options) (*MainCtrl, error) {
	if err := opts.Valid(); err != nil {
		return nil, err
	}
	c.engine = opts.Engine
	c.factory = opts.Factory
	c.config = opts.Config
	c.log = opts.Log

	if c.log == nil {
		c.log = logrus.New()
	}

	return &c, nil
}

func (c *MainCtrl) OnActivate() {
	c.launchConnect()
	c.factory.PushStatus("Ready")
}

func (c *MainCtrl) launchConnect() {
	ctl, err := ConnectCtrl{}.init(c)
	if err != nil {
		c.log.Error(err)
		return
	}

	if c.currentCtrl != nil {
		c.factory.Remove(c.currentCtrl.Screen().(ggtk.IWidget))
	}
	c.factory.Add(ctl.Screen().(ggtk.IWidget))

	c.currentCtrl = ctl

	c.factory.Show()
}

func (c *MainCtrl) launchConnection(ctx sqlengine.Context, conn *config.Connection) {
	ctl, err := ConnectionCtrl{}.init(ctx, c, conn)
	if err != nil {
		c.log.Error(err)
		return
	}

	c.factory.Remove(c.currentCtrl.Screen().(ggtk.IWidget))
	c.factory.Add(ctl.Screen().(ggtk.IWidget))

	c.currentCtrl = ctl
}
