package ctrl

import (
	"context"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type TabCtrl struct {
	*MainCtrl
	tab           *gtk.Tab
	tabLabel      *ggtk.Label
	connectScr    *gtk.ConnectScreen
	connectionScr *gtk.ConnectionScreen
}

func (c TabCtrl) init(p *MainCtrl) (*TabCtrl, error) {
	var err error
	c.MainCtrl = p

	c.tab, err = gtk.NewTab()
	if err != nil {
		return nil, err
	}

	c.tabLabel, err = ggtk.LabelNew("New Connection")
	if err != nil {
		return nil, err
	}

	c.launchConnect()

	return &c, nil
}

func (c *TabCtrl) launchConnect() {
	ctl, err := ConnectCtrl{}.init(c)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	if c.connectionScr != nil {
		c.tab.Remove(c.connectionScr)
	}

	c.connectScr = ctl.scr
	c.connectScr.OnConnect(c.onConnect)
	c.tab.PackStart(c.connectScr, true, true, 0)
}

func (c *TabCtrl) launchConnection(ctx sqlengine.Context, conn *config.Connection) {
	ctl, err := ConnectionCtrl{}.init(ctx, c, conn)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.tabLabel.SetText(conn.Name)

	if c.connectScr != nil {
		c.tab.Remove(c.connectScr)
	}

	c.connectionScr = ctl.scr
	c.tab.PackStart(c.connectionScr, true, true, 0)
}

func (c *TabCtrl) onConnect() {
	conn, ok := c.connectScr.ActiveConnection()
	if !ok {
		return
	}
	ctx, err := c.engine.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
	if err != nil {
		config.Env.Log.Error(err)
		c.factory.PushStatus("Failed connect to `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}
	c.factory.PushStatus("Connected to `%s`(%s)", conn.Name, conn.Host)

	c.launchConnection(ctx, conn)
}
