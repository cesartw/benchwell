package ctrl

import (
	"errors"

	"github.com/gotk3/gotk3/gdk"
	ggtk "github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
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
	err := c.AddTab()
	if err != nil {
		c.factory.PushStatus(err.Error())
	} else {
		c.factory.PushStatus("Ready")
	}

	c.factory.Menu.Application.New.Connect("activate", func() {
		err := c.AddTab()
		if err != nil {
			c.factory.PushStatus(err.Error())
		}
	})

	c.factory.OnTabClick(c.onNotebookDoubleClick)

	c.factory.Show()
}

func (c *MainCtrl) onNotebookDoubleClick(_ *ggtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_PRIMARY {
		return
	}
	if keyEvent.Type() != gdk.EVENT_2BUTTON_PRESS {
		return
	}

	if err := c.AddTab(); err != nil {
		c.log.Error(err)
	}
}

func (c *MainCtrl) AddTab() error {
	tab, err := TabCtrl{}.init(c)
	if err != nil {
		return err
	}

	tab.tab.Show()
	c.factory.AddTab(tab.tabLabel, tab.tab)

	return nil
}
