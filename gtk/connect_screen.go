package gtk

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ConnectScreen struct {
	*gtk.Paned
	ConnectionList *List
	forms          struct {
		notebook   *gtk.Notebook
		tcpForm    *tcpForm
		socketForm *socketForm
		active     interface {
			Clear()
			GetConnection() *config.Connection
			GrabFocus()
			SetConnection(conn *config.Connection)
		}
	}
	formOverlay *CancelOverlay

	btnBox     *gtk.ButtonBox
	btnSave    *gtk.Button
	btnConnect *gtk.Button
	btnTest    *gtk.Button

	contextMenu                             *gtk.Menu
	menuNew, menuTest, menuConnect, menuDel *gtk.MenuItem
}

func (c ConnectScreen) Init(ctrl interface {
	OnConnectionSelected()
	OnTest()
	OnSave()
	OnNewConnection()
	OnDeleteConnection()
	OnConnect()
}) (*ConnectScreen, error) {
	var err error

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}

	c.Paned.SetHExpand(true)
	c.Paned.SetVExpand(true)

	frame1, err := gtk.FrameNew("")
	if err != nil {
		return nil, err
	}

	frame2, err := gtk.FrameNew("")
	if err != nil {
		return nil, err
	}

	frame1.SetShadowType(gtk.SHADOW_IN)
	frame2.SetShadowType(gtk.SHADOW_IN)

	c.ConnectionList, err = List{}.Init(&ListOptions{SelectOnRightClick: true, StockIcon: "gtk-connect"})
	if err != nil {
		return nil, err
	}

	c.ConnectionList.OnButtonPress(c.onConnectListButtonPress)
	c.ConnectionList.SetHExpand(true)
	c.ConnectionList.SetVExpand(true)
	frame1.Add(c.ConnectionList)

	forms, err := c.buildForms()
	if err != nil {
		return nil, err
	}

	c.btnBox, err = gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	c.btnBox.SetLayout(gtk.BUTTONBOX_EDGE)

	c.btnConnect, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	c.btnSave, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	c.btnTest, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}

	c.btnConnect.SetLabel("Connect")
	c.btnTest.SetLabel("Test")
	c.btnSave.SetLabel("Save")

	c.btnConnect.SetSensitive(false)
	c.btnTest.SetSensitive(false)
	c.btnSave.SetSensitive(false)

	c.btnBox.Add(c.btnConnect)
	c.btnBox.Add(c.btnTest)
	c.btnBox.Add(c.btnSave)

	forms.Add(c.btnBox)
	frame2.Add(forms)
	c.formOverlay, err = CancelOverlay{}.Init(frame2)

	c.Paned.Pack1(frame1, false, true)
	c.Paned.Pack2(c.formOverlay, false, false)

	err = c.initMenu()
	if err != nil {
		return nil, err
	}

	c.Paned.ShowAll()

	c.menuNew.Connect("activate", ctrl.OnNewConnection)
	c.menuDel.Connect("activate", ctrl.OnDeleteConnection)
	c.btnTest.Connect("clicked", ctrl.OnTest)
	c.btnSave.Connect("clicked", ctrl.OnSave)
	c.btnConnect.Connect("clicked", ctrl.OnConnect)
	c.ConnectionList.Connect("row-selected", ctrl.OnConnectionSelected)
	c.ConnectionList.Connect("row-activated", ctrl.OnConnect)

	return &c, nil
}

func (c *ConnectScreen) onConnectListButtonPress(_ *gtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	c.contextMenu.ShowAll()
	c.contextMenu.PopupAtPointer(e)
}

func (c *ConnectScreen) initMenu() error {
	var err error
	c.contextMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.menuNew, err = menuItemWithImage("New", "gtk-new")
	if err != nil {
		return err
	}
	c.contextMenu.Add(c.menuNew)

	c.menuConnect, err = menuItemWithImage("Connect", "gtk-connect")
	if err != nil {
		return err
	}
	c.menuConnect.Connect("activate", func() {
		c.onConnect()
	})
	c.contextMenu.Add(c.menuConnect)

	c.menuTest, err = menuItemWithImage("Test", "gtk-play")
	if err != nil {
		return err
	}
	c.menuTest.Connect("activate", func() {
		c.btnTest.Emit("clicked")
	})
	c.contextMenu.Add(c.menuTest)

	c.menuDel, err = menuItemWithImage("Delete", "gtk-delete")
	if err != nil {
		return err
	}
	c.contextMenu.Add(c.menuDel)

	return nil
}

func (c *ConnectScreen) buildTcpForm() (gtk.IWidget, *tcpForm, error) {
	label, err := gtk.LabelNew("TCP/IP")
	if err != nil {
		return nil, nil, err
	}

	frm, err := tcpForm{}.Init()
	if err != nil {
		return nil, nil, err
	}

	frm.ShowAll()
	label.Show()

	return label, frm, nil
}

func (c *ConnectScreen) buildSocketForm() (gtk.IWidget, *socketForm, error) {
	label, err := gtk.LabelNew("Socket")
	if err != nil {
		return nil, nil, err
	}

	frm, err := socketForm{}.Init()
	if err != nil {
		return nil, nil, err
	}

	frm.ShowAll()
	label.Show()

	return label, frm, nil
}

func (c *ConnectScreen) SetConnections(connections []*config.Connection) {
	//c.connections = connections
	names := make([]string, len(connections))
	for i, con := range connections {
		names[i] = con.Name
	}

	c.ConnectionList.UpdateItems(StringSliceToStringers(names))
}

func (c *ConnectScreen) ClearForm() {
	c.forms.active.Clear()
	c.btnConnect.SetSensitive(false)
	c.btnTest.SetSensitive(false)
	c.btnSave.SetSensitive(false)
}

func (c *ConnectScreen) Connecting(cancel func()) {
	c.formOverlay.Run(cancel)
	c.ConnectionList.SetSensitive(false)
}

func (c *ConnectScreen) CancelConnecting() {
	c.formOverlay.Stop()
	c.ConnectionList.SetSensitive(true)
}

func (c *ConnectScreen) FocusForm() {
	c.forms.active.GrabFocus()
}

func (c *ConnectScreen) SetConnection(conn *config.Connection) {
	switch conn.Type {
	case "tcp":
		c.forms.active = c.forms.tcpForm
		c.forms.notebook.SetCurrentPage(0)
	case "socket":
		c.forms.active = c.forms.socketForm
		c.forms.notebook.SetCurrentPage(1)
	default:
		return
	}

	c.forms.active.SetConnection(conn)

	if conn.Valid() {
		c.btnConnect.SetSensitive(true)
		c.btnTest.SetSensitive(true)
		c.btnSave.SetSensitive(true)
	}
}

func (c *ConnectScreen) ActiveConnectionIndex() int {
	return c.ConnectionList.GetSelectedRow().GetIndex()
}

func (c *ConnectScreen) onConnect() {
	c.btnConnect.Emit("activate")
}

func (c *ConnectScreen) Dispose() {
}

func (c *ConnectScreen) buildForms() (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	c.forms.notebook, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	box.Add(c.forms.notebook)

	box.SetSizeRequest(300, 200)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetHAlign(gtk.ALIGN_CENTER)
	c.forms.notebook.SetShowBorder(true)
	c.forms.notebook.SetCanFocus(true)

	// TCP
	{
		label, frm, err := c.buildTcpForm()
		if err != nil {
			return nil, err
		}
		frm.onChange(func(_ *gtk.Entry, e *gdk.Event) bool {
			conn := c.forms.active.GetConnection()

			if conn.Valid() {
				c.btnConnect.SetSensitive(true)
				c.btnTest.SetSensitive(true)
				c.btnSave.SetSensitive(true)

				keyEvent := gdk.EventKeyNewFromEvent(e)
				if keyEvent.KeyVal() == 65293 && keyEvent.State()&gdk.CONTROL_MASK > 0 {
					c.btnConnect.Emit("activate")
					return false
				}

				return false
			}

			c.btnConnect.SetSensitive(false)
			c.btnTest.SetSensitive(false)
			c.btnSave.SetSensitive(false)
			return false
		})

		c.forms.notebook.AppendPage(frm, label)
		c.forms.tcpForm = frm
	}

	// SOCKET
	{
		label, frm, err := c.buildSocketForm()
		if err != nil {
			return nil, err
		}
		frm.onChange(func(_ *gtk.Entry, e *gdk.Event) bool {
			conn := c.forms.active.GetConnection()

			if conn.Valid() {
				c.btnConnect.SetSensitive(true)
				c.btnTest.SetSensitive(true)
				c.btnSave.SetSensitive(true)

				keyEvent := gdk.EventKeyNewFromEvent(e)
				if keyEvent.KeyVal() == 65293 && keyEvent.State()&gdk.CONTROL_MASK > 0 {
					c.btnConnect.Emit("activate")
					return false
				}

				return false
			}

			c.btnConnect.SetSensitive(false)
			c.btnTest.SetSensitive(false)
			c.btnSave.SetSensitive(false)
			return false
		})

		c.forms.notebook.AppendPage(frm, label)
		c.forms.socketForm = frm
	}

	c.forms.active = c.forms.tcpForm

	return box, nil
}

func (c *ConnectScreen) GetFormConnection() *config.Connection {
	conn := c.forms.active.GetConnection()
	//conn.Queries = c.activeForm.queries
	return conn
}
