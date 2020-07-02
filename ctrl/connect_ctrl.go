package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
)

type ConnectCtrl struct {
	*WindowTabCtrl
	scr *gtk.ConnectScreen
}

func (c ConnectCtrl) Init(p *WindowTabCtrl) (*ConnectCtrl, error) {
	c.WindowTabCtrl = p

	var err error
	c.scr, err = gtk.ConnectScreen{}.Init(&c)
	if err != nil {
		return nil, err
	}

	c.scr.SetConnections(config.Env.Connections)

	return &c, nil
}

func (c *ConnectCtrl) AddEmptyTab() error {
	return nil
}

func (c *ConnectCtrl) SetFileText(s string) {
}

func (c *ConnectCtrl) OnTest() {
	var conn *config.Connection
	index := c.scr.ActiveConnectionIndex()
	if index > 0 {
		conn = config.Env.Connections[index]
	} else {
		conn = c.scr.GetFormConnection()
	}

	ctx, err := c.Engine.Connect(nil, *conn)
	if err != nil {
		config.Env.Log.Error(err)
		c.window.PushStatus("Fail connection `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}

	c.window.PushStatus("Connection to `%s`(%s) was successful", conn.Name, conn.Host)
	c.Engine.Disconnect(ctx)
}

func (c *ConnectCtrl) OnSave() {
	index := c.scr.ActiveConnectionIndex()
	if index == -1 || index >= len(config.Env.Connections) {
		config.Env.Connections = append(config.Env.Connections, c.scr.GetFormConnection())
	} else {
		config.Env.Connections[index] = c.scr.GetFormConnection()
	}

	config.Env.Save(c.window.ApplicationWindow)
	c.scr.SetConnections(config.Env.Connections)
	c.scr.ConnectionList.SelectRow(c.scr.ConnectionList.GetRowAtIndex(index))

	c.window.PushStatus("Saved")
}

func (c *ConnectCtrl) OnDeleteConnection() {
	index := c.scr.ActiveConnectionIndex()
	if index == -1 {
		return
	}

	// deleting an saved connection
	if index < len(config.Env.Connections) {
		config.Env.Connections = append(config.Env.Connections[:index], config.Env.Connections[index+1:]...)
	}

	config.Env.Save(c.window.ApplicationWindow)
	c.scr.SetConnections(config.Env.Connections)
	c.scr.ClearForm()

	c.window.PushStatus("Deleted")
}

func (c *ConnectCtrl) OnNewConnection() {
	row, err := c.scr.ConnectionList.AppendItem(gtk.Stringer("New Connection"))
	if err != nil {
		config.Env.Log.Error(err)
	}
	c.scr.ConnectionList.SelectRow(row)
	c.scr.ClearForm()
	c.scr.FocusForm()
	c.scr.SetConnection(&config.Connection{Name: "New Connection", Port: 3306})
}

func (c *ConnectCtrl) OnConnectionSelected() {
	row := c.scr.ConnectionList.GetSelectedRow()
	if row.GetIndex() == -1 {
		return
	}

	if row.GetIndex() >= len(config.Env.Connections) {
		c.scr.ClearForm()
		c.scr.FocusForm()
		return
	}

	err := config.Env.Connections[row.GetIndex()].Decrypt(c.window.ApplicationWindow)
	if err != nil {
		c.window.PushStatus("Fail to decrypt password")
		c.scr.ConnectionList.ClearSelection()
		return
	}

	c.scr.SetConnection(config.Env.Connections[row.GetIndex()])
}

func (c *ConnectCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectCtrl) Close() bool {
	return false
}

func (c *ConnectCtrl) Connecting(cancel func()) {
	c.scr.Connecting(cancel)
}

func (c *ConnectCtrl) CancelConnecting() {
	c.scr.CancelConnecting()
}
