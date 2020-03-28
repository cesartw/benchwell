package ctrl

import (
	"context"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
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

	c.scr.SetConnections(c.config.Connections)
	c.scr.OnTest(c.onTest)
	c.scr.OnSave(c.onSave)

	return &c, nil
}

func (c *ConnectCtrl) onTest() {
	conn, ok := c.scr.ActiveConnection()
	if !ok {
		return
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
	c.factory.PushStatus("Saved")
}

func (c *ConnectCtrl) Screen() interface{} {
	return c.scr
}
