package ctrl

import (
	"context"

	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
)

type ConnectCtrl struct {
	*MainCtrl
	scr *gtk.ConnectScreen
}

func (c ConnectCtrl) init(p *MainCtrl) (*ConnectCtrl, error) {
	c.MainCtrl = p

	var err error
	c.scr, err = c.factory.NewConnectScreen()
	if err != nil {
		return nil, err
	}

	c.scr.SetConnections(c.config.Connection)
	c.scr.OnConnect(c.onConnect)
	c.scr.OnTest(c.onTest)
	c.scr.OnSave(c.onSave)

	return &c, nil
}

func (c *ConnectCtrl) onConnect() {
	conn := c.scr.ActiveConnection()
	ctx, err := c.engine.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
	if err != nil {
		c.log.Error(err)
		return
	}
	c.factory.PushStatus("Connected to `%s`(%s)", conn.Name, conn.Host)

	c.launchConnection(ctx, conn)
}

func (c *ConnectCtrl) onTest() {
	conn := c.scr.ActiveConnection()
	ctx, err := c.engine.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
	if err != nil {
		c.log.Error(err)
		return
	}

	c.factory.PushStatus("Connection to `%s`(%s) was successful", conn.Name, conn.Host)
	c.engine.Disconnect(ctx)
}

func (c *ConnectCtrl) onSave() {
	c.factory.PushStatus("Saved")
}

func (c *ConnectCtrl) Screen() interface{} {
	return c.scr
}
