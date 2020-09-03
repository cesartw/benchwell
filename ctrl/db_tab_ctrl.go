package ctrl

import (
	"context"

	ggtk "github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
	"bitbucket.org/goreorto/benchwell/sqlengine"
)

type DbTabCtrl struct {
	*WindowCtrl
	screenHolder   *gtk.DB
	connectCtrl    *ConnectCtrl
	connectionCtrl *ConnectionCtrl

	currentCtrl interface {
		Close()
		SetFileText(string)
		Content() ggtk.IWidget
		Title() string
	}
}

func (c DbTabCtrl) Init(p *WindowCtrl) (*DbTabCtrl, error) {
	defer config.LogStart("DbTabCtrl.Init", nil)()

	c.WindowCtrl = p

	var err error
	c.screenHolder, err = gtk.DB{}.Init(c.window)
	if err != nil {
		return nil, err
	}

	c.launchConnect()

	return &c, nil
}

func (c *DbTabCtrl) Close() {
	c.currentCtrl.Close()
}

func (c *DbTabCtrl) Title() string {
	defer config.LogStart("DbTabCtrl.Title", nil)()

	return c.currentCtrl.Title()
}

func (c *DbTabCtrl) Content() ggtk.IWidget {
	defer config.LogStart("DbTabCtrl.Content", nil)()

	return c.screenHolder
}

func (c *DbTabCtrl) SetWindowCtrl(i interface{}) {
	defer config.LogStart("DbTabCtrl.SetWindowCtrl", nil)()

	c.WindowCtrl = i.(*WindowCtrl)
}

func (c *DbTabCtrl) SetFileText(s string) {
	defer config.LogStart("DbTabCtrl.SetFileText", nil)()

	c.currentCtrl.SetFileText(s)
}

func (c *DbTabCtrl) Show() {
	defer config.LogStart("DbTabCtrl.Show", nil)()

	//c.tab.Show()
}

func (c *DbTabCtrl) Removed() {
	defer config.LogStart("DbTabCtrl.Removed", nil)()

	if c.connectionCtrl != nil {
		c.Engine.Disconnect(c.connectionCtrl.ctx)
		c.window.PushStatus("Disconnected")
	}
}

func (c *DbTabCtrl) launchConnect() {
	defer config.LogStart("DbTabCtrl.launchConnect", nil)()

	var err error
	c.connectCtrl, err = ConnectCtrl{}.Init(c)
	if err != nil {
		config.Error(err)
		return
	}
	c.currentCtrl = c.connectCtrl
	c.screenHolder.SetContent(c.connectCtrl.scr)
}

func (c *DbTabCtrl) launchConnection(ctx *sqlengine.Context, conn *config.Connection) {
	defer config.LogStart("DbTabCtrl.launchConnection", nil)()

	var err error
	c.connectionCtrl, err = ConnectionCtrl{}.Init(ctx, c, conn)
	if err != nil {
		config.Error(err)
		return
	}
	c.currentCtrl = c.connectionCtrl
	c.screenHolder.SetContent(c.connectionCtrl.scr)
}

func (c *DbTabCtrl) OnConnect() {
	defer config.LogStart("DbTabCtrl.OnConnect", nil)()

	c.onConnect(nil)
}

func (c *DbTabCtrl) onConnect(f func()) {
	defer config.LogStart("DbTabCtrl.OnConnect", nil)()

	conn := c.connectCtrl.scr.GetFormConnection()

	cancel := c.window.Go(func(ctx context.Context) func() {
		engineCtx, err := c.Engine.Connect(ctx, *conn)
		if err != nil {
			return func() {
				config.Error(err)
				c.window.PushStatus("Failed connect to `%s`(%s): %s", conn.Name, conn.Host, err.Error())
				c.connectCtrl.CancelConnecting()
			}
		}

		return func() {
			c.window.PushStatus("Connected to `%s`(%s)", conn.Name, conn.Host)
			c.connectCtrl.CancelConnecting()
			c.launchConnection(engineCtx, conn)
			if f != nil {
				f()
			}
		}
	})

	c.connectCtrl.Connecting(cancel)
}
