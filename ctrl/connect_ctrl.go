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
	c.DbTabCtrl = p

	var err error
	c.scr, err = gtk.ConnectScreen{}.Init(c.window, &c)
	if err != nil {
		return nil, err
	}

	c.scr.SetConnections(c.Config().Connections)

	return &c, nil
}

func (c *ConnectCtrl) Title() string {
	return "Connect"
}

func (c *ConnectCtrl) Content() ggtk.IWidget {
	return c.scr
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
		conn = c.Config().Connections[index]
	} else {
		conn = c.scr.GetFormConnection()
	}

	ctx, err := c.Engine.Connect(context.Background(), *conn)
	if err != nil {
		c.Config().Error(err)
		c.window.PushStatus("Fail connection `%s`(%s): %s", conn.Name, conn.Host, err.Error())
		return
	}

	c.window.PushStatus("Connection to `%s`(%s) was successful", conn.Name, conn.Host)
	c.Engine.Disconnect(ctx)
}

func (c *ConnectCtrl) OnSave() {
	conn := c.scr.GetFormConnection()
	err := c.Config().SaveConnection(conn)
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	c.scr.SetConnections(c.Config().Connections)
	for i, co := range c.Config().Connections {
		if co.ID == conn.ID {
			c.scr.ConnectionList.SelectRow(c.scr.ConnectionList.GetRowAtIndex(i))
			break
		}
	}

	c.window.PushStatus("Saved")
}

func (c *ConnectCtrl) OnDeleteConnection() {
	index := c.scr.ActiveConnectionIndex()
	if index == -1 {
		return
	}

	err := c.Config().DeleteConnection(c.Config().Connections[index])
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	//c.Config().Save(c.window.ApplicationWindow)
	c.scr.SetConnections(c.Config().Connections)
	c.scr.ClearForm()

	c.window.PushStatus("Deleted")
}

func (c *ConnectCtrl) OnNewConnection() {
	row, err := c.scr.ConnectionList.AppendItem(gtk.Stringer("New Connection"))
	if err != nil {
		c.Config().Error(err)
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

	if row.GetIndex() >= len(c.Config().Connections) {
		c.scr.ClearForm()
		c.scr.FocusForm()
		return
	}

	err := c.Config().Connections[row.GetIndex()].Decrypt(c.window.ApplicationWindow)
	if err != nil {
		c.window.PushStatus("Fail to decrypt password")
		c.scr.ConnectionList.ClearSelection()
		return
	}

	c.scr.SetConnection(c.Config().Connections[row.GetIndex()])
}

func (c *ConnectCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectCtrl) Close() bool {
	return false
}

func (c *ConnectCtrl) FullClose() {
}

func (c *ConnectCtrl) Connecting(cancel func()) {
	c.scr.Connecting(cancel)
}

func (c *ConnectCtrl) CancelConnecting() {
	c.scr.CancelConnecting()
}
