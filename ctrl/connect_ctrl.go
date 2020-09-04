package ctrl

import (
	"context"

	ggtk "github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
)

type ConnectCtrl struct {
	*DbTabCtrl
	scr *gtk.ConnectScreen
}

func (c ConnectCtrl) Init(p *DbTabCtrl) (*ConnectCtrl, error) {
	defer config.LogStart("ConnectCtrl.Init", nil)()

	c.DbTabCtrl = p

	var err error
	c.scr, err = gtk.ConnectScreen{}.Init(c.window, &c)
	if err != nil {
		return nil, err
	}

	c.scr.SetConnections(config.Connections)

	return &c, nil
}

func (c *ConnectCtrl) Title() string {
	return "Connect"
}

func (c *ConnectCtrl) Content() ggtk.IWidget {
	defer config.LogStart("ConnectCtrl.Content", nil)()

	return c.scr
}

func (c *ConnectCtrl) SetFileText(s string) {
	defer config.LogStart("ConnectCtrl.SetFileText", nil)()
}

func (c *ConnectCtrl) OnTest() {
	defer config.LogStart("ConnectCtrl.OnTest", nil)()
	var conn *config.Connection
	index := c.scr.ActiveConnectionIndex()
	if index > 0 {
		conn = config.Connections[index]
	} else {
		conn = c.scr.GetFormConnection()
	}

	ctx, err := c.Engine.Connect(context.Background(), *conn)
	if err != nil {
		config.Error(err)
		c.window.PushStatus("Fail connection `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}

	c.window.PushStatus("Connection to `%s`(%s) was successful", conn.Name, conn.Host)
	c.Engine.Disconnect(ctx)
}

func (c *ConnectCtrl) OnSave() {
	defer config.LogStart("ConnectCtrl.OnSave", nil)()

	conn := c.scr.GetFormConnection()
	err := config.SaveConnection(c.window.ApplicationWindow, conn)
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	c.scr.SetConnections(config.Connections)
	for i, co := range config.Connections {
		if co.ID == conn.ID {
			c.scr.ConnectionList.SelectRow(c.scr.ConnectionList.GetRowAtIndex(i))
			break
		}
	}

	c.window.PushStatus("Saved")
}

func (c *ConnectCtrl) OnDeleteConnection() {
	defer config.LogStart("ConnectCtrl.OnDeleteConnection", nil)()
	index := c.scr.ActiveConnectionIndex()
	if index == -1 {
		return
	}

	err := config.DeleteConnection(config.Connections[index])
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	//config.Save(c.window.ApplicationWindow)
	c.scr.SetConnections(config.Connections)
	c.scr.ClearForm()

	c.window.PushStatus("Deleted")
}

func (c *ConnectCtrl) OnNewConnection() {
	defer config.LogStart("ConnectCtrl.OnNewConnection", nil)()
	row, err := c.scr.ConnectionList.AppendItem(gtk.Stringer("New Connection"))
	if err != nil {
		config.Error(err)
	}
	c.scr.ConnectionList.SelectRow(row)
	c.scr.ClearForm()
	c.scr.FocusForm()
	c.scr.SetConnection(&config.Connection{Name: "New Connection", Port: 3306})
}

func (c *ConnectCtrl) OnConnectionSelected() {
	defer config.LogStart("ConnectCtrl.OnConnectionSelected", nil)()
	row := c.scr.ConnectionList.GetSelectedRow()
	if row.GetIndex() == -1 {
		return
	}

	if row.GetIndex() >= len(config.Connections) {
		c.scr.ClearForm()
		c.scr.FocusForm()
		return
	}

	conn := config.Connections[row.GetIndex()]
	err := conn.Decrypt(c.window.ApplicationWindow)
	if err != nil {
		conn.Encrypted = false
		c.window.PushStatus("Fail to decrypt password: %s", err.Error())
		c.scr.ConnectionList.ClearSelection()
		return
	}
	c.scr.SetConnection(conn)
}

func (c *ConnectCtrl) Screen() interface{} {
	defer config.LogStart("ConnectCtrl.Screen", nil)()
	return c.scr
}

func (c *ConnectCtrl) Close() {
	defer config.LogStart("ConnectCtrl.Close", nil)()
}

func (c *ConnectCtrl) Connecting(cancel func()) {
	defer config.LogStart("ConnectCtrl.Connecting", nil)()
	c.scr.Connecting(cancel)
}

func (c *ConnectCtrl) CancelConnecting() {
	defer config.LogStart("ConnectCtrl.CancelConnecting", nil)()
	c.scr.CancelConnecting()
}
