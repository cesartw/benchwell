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
		Close() bool
		FullClose()
		AddEmptyTab() error
		SetFileText(string)
		Content() ggtk.IWidget
		Title() string
	}
}

func (c *DbTabCtrl) Title() string {
	return c.currentCtrl.Title()
}

func (c *DbTabCtrl) Content() ggtk.IWidget {
	return c.screenHolder
}

func (c DbTabCtrl) Init(p *WindowCtrl) (*DbTabCtrl, error) {
	c.WindowCtrl = p

	var err error
	c.screenHolder, err = gtk.DB{}.Init(c.window)
	if err != nil {
		return nil, err
	}

	c.launchConnect()

	return &c, nil
}

func (c *DbTabCtrl) SetWindowCtrl(i interface{}) {
	c.WindowCtrl = i.(*WindowCtrl)
}

func (c *DbTabCtrl) AddTab() error {
	return c.currentCtrl.AddEmptyTab()
}

func (c *DbTabCtrl) SetFileText(s string) {
	c.currentCtrl.SetFileText(s)
}

func (c *DbTabCtrl) Show() {
	//c.tab.Show()
}

func (c *DbTabCtrl) Removed() {
	if c.connectionCtrl != nil {
		c.Engine.Disconnect(c.connectionCtrl.mainCtx)
		c.window.PushStatus("Disconnected")
	}
}

// Close delegates the close tab action ot connect or connection screen
func (c *DbTabCtrl) Close() {
	// TODO: figure out which screen is open
	c.currentCtrl.FullClose()
}

// Close all tabs
func (c *DbTabCtrl) FullClose() {
	c.currentCtrl.FullClose()
}

func (c *DbTabCtrl) launchConnect() {
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
	var err error
	c.connectionCtrl, err = ConnectionCtrl{}.Init(ctx, c, conn)
	if err != nil {
		config.Error(err)
		return
	}
	c.currentCtrl = c.connectionCtrl
	c.screenHolder.SetContent(c.connectionCtrl.scr)
	c.ChangeTitle(c.currentCtrl.Title())
}

func (c *DbTabCtrl) OnConnect() {
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
		}
	})

	c.connectCtrl.Connecting(cancel)
}
