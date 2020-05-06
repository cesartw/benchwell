package ctrl

import (
	"context"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
)

type ConnectCtrl struct {
	*TabCtrl
	scr *gtk.ConnectScreen
}

func (c ConnectCtrl) Init(p *TabCtrl) (*ConnectCtrl, error) {
	c.TabCtrl = p

	var err error
	c.scr, err = gtk.NewConnectScreen()
	if err != nil {
		return nil, err
	}

	c.scr.SetConnections(config.Env.Connections)

	c.scr.OnConnectionSelected(c.onConnectionSelected)

	c.scr.OnTest(c.onTest)
	c.scr.OnSave(c.onSave)

	c.scr.OnNewConnection(c.onNewConnection)
	c.scr.OnDeleteConnection(c.onDeleteConnection)

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
		c.window.PushStatus("Fail connection `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}

	c.window.PushStatus("Connection to `%s`(%s) was successful", conn.Name, conn.Host)
	c.engine.Disconnect(ctx)
}

func (c *ConnectCtrl) onSave() {
	index := c.scr.ActiveConnectionIndex()
	if index == -1 || index >= len(config.Env.Connections) {
		config.Env.Connections = append(config.Env.Connections, c.scr.GetFormConnection())
	} else {
		config.Env.Connections[index] = c.scr.GetFormConnection()
	}

	config.Env.Save()
	c.scr.SetConnections(config.Env.Connections)
	c.scr.ConnectionList.SelectRow(c.scr.ConnectionList.GetRowAtIndex(index))

	c.window.PushStatus("Saved")
}

func (c *ConnectCtrl) onDeleteConnection() {
	index := c.scr.ActiveConnectionIndex()
	if index == -1 {
		return
	}

	config.Env.Connections = append(config.Env.Connections[:index], config.Env.Connections[index+1:]...)

	config.Env.Save()
	c.scr.SetConnections(config.Env.Connections)
	c.scr.ClearForm()

	c.window.PushStatus("Deleted")
}

func (c *ConnectCtrl) onNewConnection() {
	row, err := c.scr.ConnectionList.AddItem("New Connection")
	if err != nil {
		config.Env.Log.Error(err)
	}
	c.scr.ConnectionList.SelectRow(row)
	c.scr.ClearForm()
	c.scr.FocusForm()
	c.scr.SetFormConnection(&config.Connection{Name: "New Connection", Port: 3306})
}

func (c *ConnectCtrl) onConnectionSelected() {
	row := c.scr.ConnectionList.GetSelectedRow()
	if row.GetIndex() == -1 {
		return
	}

	if row.GetIndex() >= len(config.Env.Connections) {
		return
	}

	err := config.Env.Connections[row.GetIndex()].Decrypt()
	if err != nil {
		c.window.PushStatus("Fail to decrypt password")
		return
	}

	c.scr.SetFormConnection(config.Env.Connections[row.GetIndex()])
}

func (c *ConnectCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectCtrl) Close() bool {
	return false
}

func (c *ConnectCtrl) AddTab() error {
	return nil
}
