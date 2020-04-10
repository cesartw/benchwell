package ctrl

import (
	"context"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type ConnectCtrl struct {
	*TabCtrl
	scr *gtk.ConnectScreen
}

func (c ConnectCtrl) init(p *TabCtrl) (*ConnectCtrl, error) {
	c.TabCtrl = p

	var err error
	c.scr, err = c.factory.NewConnectScreen()
	if err != nil {
		return nil, err
	}

	c.scr.SetConnections(config.Env.Connections)

	c.scr.OnConnectionSelected(func(list *ggtk.ListBox) {
		row := list.GetSelectedRow()

		c.scr.SetFormConnection(config.Env.Connections[row.GetIndex()])
	})

	c.scr.OnTest(c.onTest)
	c.scr.OnSave(c.onSave)

	return &c, nil
}

func (c *ConnectCtrl) onTest() {
	var conn *config.Connection
	index := c.scr.ActiveConnectionIndex()
	if index > 0 {
		conn = config.Env.Connections[index]
	} else {
		conn = c.scr.GetFormConnection()
	}

	ctx, err := c.engine.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
	if err != nil {
		config.Env.Log.Error(err)
		c.factory.PushStatus("Fail connection `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}

	c.factory.PushStatus("Connection to `%s`(%s) was successful", conn.Name, conn.Host)
	c.engine.Disconnect(ctx)
}

func (c *ConnectCtrl) onSave() {
	index := c.scr.ActiveConnectionIndex()
	if index == -1 {
		config.Env.Connections = append(config.Env.Connections, c.scr.GetFormConnection())
	}

	config.Env.Save()

	c.factory.PushStatus("Saved")
}

func (c *ConnectCtrl) Screen() interface{} {
	return c.scr
}
