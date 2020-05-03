package ctrl

import (
	"context"
	"time"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type TabCtrl struct {
	*WindowCtrl
	tab            *gtk.Tab
	tabLabel       *ggtk.Label
	connectCtrl    *ConnectCtrl
	connectionCtrl *ConnectionCtrl

	currentCtrl interface {
		Close() bool
		AddTab() error
	}
}

func (c TabCtrl) Init(p *WindowCtrl) (*TabCtrl, error) {
	var err error
	c.WindowCtrl = p

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

func (c *TabCtrl) AddTab() error {
	return c.currentCtrl.AddTab()
}

func (c *TabCtrl) Show() {
	c.tab.Show()
}

// Close delegates the close tab action ot connect or connection screen
func (c *TabCtrl) Close() bool {
	// TODO: figure out which screen is open
	return c.currentCtrl.Close()
}

func (c *TabCtrl) launchConnect() {
	var err error
	c.connectCtrl, err = ConnectCtrl{}.Init(c)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	if c.connectionCtrl != nil {
		c.tab.Remove(c.connectionCtrl.scr)
	}

	c.currentCtrl = c.connectCtrl

	c.connectCtrl.scr.OnConnect(c.onConnect)
	c.tab.PackStart(c.connectCtrl.scr, true, true, 0)
}

func (c *TabCtrl) launchConnection(ctx sqlengine.Context, conn *config.Connection) {
	var err error
	c.connectionCtrl, err = ConnectionCtrl{}.Init(ctx, c, conn)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	if conn.Name != "" {
		c.tabLabel.SetText(conn.Name)
	} else {
		c.tabLabel.SetText(conn.Host)
	}

	if c.connectCtrl != nil {
		c.tab.Remove(c.connectCtrl.scr)
	}

	c.currentCtrl = c.connectionCtrl

	c.tab.PackStart(c.connectionCtrl.scr, true, true, 0)
}

func (c *TabCtrl) onConnect() {
	var conn *config.Connection
	index := c.connectCtrl.scr.ActiveConnectionIndex()
	if index == -1 {
		conn = c.connectCtrl.scr.GetFormConnection()
	} else {
		conn = config.Env.Connections[index]
	}

	ctx, done := context.WithTimeout(context.TODO(), time.Second*5)
	defer done()

	ctx, err := c.engine.Connect(sqlengine.Context(ctx), conn.GetDSN())
	if err != nil {
		config.Env.Log.Error(err)
		c.window.PushStatus("Failed connect to `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}
	c.window.PushStatus("Connected to `%s`(%s)", conn.Name, conn.Host)

	c.launchConnection(ctx, conn)
}
