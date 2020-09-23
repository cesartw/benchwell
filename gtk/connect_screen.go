package gtk

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
)

type form interface {
	Clear()
	GetConnection() (*config.Connection, bool)
	GrabFocus()
	SetConnection(conn *config.Connection)
}

type ConnectScreen struct {
	*gtk.Paned
	w              *Window
	ConnectionList *List
	forms          struct {
		active form
		mysql  struct {
			notebook   *gtk.Notebook
			tcpForm    *tcpForm
			socketForm *socketForm
			sshForm    *sshForm
		}
		sqlite struct {
			form *sqliteForm
		}
	}
	formOverlay *CancelOverlay
	stack       *gtk.Stack

	btnBox     *gtk.ButtonBox
	btnSave    *gtk.Button
	btnConnect *gtk.Button
	btnTest    *gtk.Button

	contextMenu                             *gtk.Menu
	menuNew, menuTest, menuConnect, menuDel *gtk.MenuItem

	ctrl connectScreenCtrl
}

type connectScreenCtrl interface {
	OnConnectionSelected()
	OnTest()
	OnSave()
	OnNewConnection()
	OnDeleteConnection()
	OnConnect()
}

func (c ConnectScreen) Init(w *Window, ctrl connectScreenCtrl) (*ConnectScreen, error) {
	defer config.LogStart("ConnectScreen.Init", nil)()

	var err error
	c.w = w
	c.ctrl = ctrl

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	c.Paned.Show()
	c.Paned.SetHExpand(true)
	c.Paned.SetVExpand(true)
	c.Paned.SetWideHandle(true)

	frame2, err := gtk.FrameNew("")
	if err != nil {
		return nil, err
	}
	frame2.Show()
	frame2.SetShadowType(gtk.SHADOW_IN)

	c.ConnectionList, err = List{}.Init(c.w, &ListOptions{
		SelectOnRightClick: true,
		IconFunc: func(name fmt.Stringer) (string, int) {
			return "connection", ICON_SIZE_BUTTON
		},
	}, ctrl)
	if err != nil {
		return nil, err
	}
	c.ConnectionList.Show()
	c.ConnectionList.OnButtonPress(c.onConnectListButtonPress)
	c.ConnectionList.SetHExpand(true)
	c.ConnectionList.SetVExpand(true)

	connectionListSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	connectionListSW.Add(c.ConnectionList)
	connectionListSW.Show()

	mysqlForms, err := c.buildMysqlForms()
	if err != nil {
		return nil, err
	}
	mysqlForms.Show()

	sqliteForm, err := c.buildSqliteForms()
	if err != nil {
		return nil, err
	}
	sqliteForm.Show()

	c.btnBox, err = gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	c.btnBox.Show()
	c.btnBox.SetLayout(gtk.BUTTONBOX_EDGE)

	c.btnConnect, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	c.btnConnect.Show()
	c.btnSave, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	c.btnSave.Show()
	c.btnTest, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	c.btnTest.Show()

	c.btnConnect.SetLabel("Connect")
	c.btnTest.SetLabel("Test")
	c.btnSave.SetLabel("Save")

	c.btnConnect.SetSensitive(false)
	c.btnTest.SetSensitive(false)
	c.btnSave.SetSensitive(false)

	c.btnBox.Add(c.btnConnect)
	c.btnBox.Add(c.btnTest)
	c.btnBox.Add(c.btnSave)
	c.btnBox.Show()

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	vbox.Show()
	vbox.SetSizeRequest(300, 200)
	vbox.SetVAlign(gtk.ALIGN_CENTER)
	vbox.SetHAlign(gtk.ALIGN_CENTER)

	adapterStack, err := gtk.StackSwitcherNew()
	if err != nil {
		return nil, err
	}
	adapterStack.Show()
	adapterStack.SetVExpand(true)
	adapterStack.SetHExpand(true)

	c.stack, err = gtk.StackNew()
	if err != nil {
		return nil, err
	}
	c.stack.Show()
	c.stack.SetHomogeneous(true)
	c.stack.AddTitled(mysqlForms, "mysql", "Mysql")
	c.stack.AddTitled(sqliteForm, "sqlite", "SQLite")
	c.stack.SetVisibleChildName("mysql")
	c.stack.SetVExpand(true)
	c.stack.SetHExpand(true)

	adapterStack.SetStack(c.stack)

	vbox.PackEnd(c.btnBox, false, true, 0)
	vbox.PackStart(adapterStack, false, true, 0)
	vbox.PackStart(c.stack, true, false, 0)

	frame2.Add(vbox)
	c.formOverlay, err = CancelOverlay{}.Init(frame2)
	if err != nil {
		return nil, err
	}
	c.formOverlay.Show()

	c.Paned.Pack1(connectionListSW, true, true)
	c.Paned.Pack2(c.formOverlay, true, false)
	c.Paned.Show()

	err = c.initMenu()
	if err != nil {
		return nil, err
	}

	c.menuNew.Connect("activate", ctrl.OnNewConnection)
	c.menuDel.Connect("activate", ctrl.OnDeleteConnection)
	c.btnTest.Connect("clicked", ctrl.OnTest)
	c.btnSave.Connect("clicked", ctrl.OnSave)
	c.btnConnect.Connect("clicked", ctrl.OnConnect)
	c.ConnectionList.Connect("row-selected", ctrl.OnConnectionSelected)
	c.ConnectionList.Connect("row-activated", ctrl.OnConnect)

	return &c, nil
}

func (c *ConnectScreen) onConnectListButtonPress(_ *gtk.ListBox, e *gdk.Event) bool {
	defer config.LogStart("ConnectScreen.onConnectListButtonPress", nil)()

	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return false
	}

	c.contextMenu.Show()
	c.contextMenu.PopupAtPointer(e)
	return true
}

func (c *ConnectScreen) initMenu() error {
	defer config.LogStart("ConnectScreen.initMenu", nil)()

	var err error
	c.contextMenu, err = gtk.MenuNew()
	if err != nil {
		return err
	}

	c.menuNew, err = BWMenuItemWithImage("New", "connection")
	if err != nil {
		return err
	}
	c.menuNew.Show()
	c.contextMenu.Add(c.menuNew)

	c.menuConnect, err = BWMenuItemWithImage("Connect", "next")
	if err != nil {
		return err
	}
	c.menuConnect.Show()
	c.menuConnect.Connect("activate", func() {
		c.onConnect()
	})

	c.menuTest, err = BWMenuItemWithImage("Test", "refresh")
	if err != nil {
		return err
	}
	c.menuTest.Show()
	c.menuTest.Connect("activate", func() {
		c.btnTest.Emit("clicked")
	})

	c.menuDel, err = BWMenuItemWithImage("Delete", "close")
	if err != nil {
		return err
	}
	c.menuDel.Show()

	c.contextMenu.Add(c.menuConnect)
	c.contextMenu.Add(c.menuTest)
	c.contextMenu.Add(c.menuDel)

	return nil
}

func (c *ConnectScreen) buildTcpForm() (*gtk.Label, *tcpForm, error) {
	defer config.LogStart("ConnectScreen.buildTcpForm", nil)()

	label, err := gtk.LabelNew("TCP/IP")
	if err != nil {
		return nil, nil, err
	}
	label.Show()

	frm, err := tcpForm{}.Init(c.w)
	if err != nil {
		return nil, nil, err
	}
	frm.Show()

	return label, frm, nil
}

func (c *ConnectScreen) buildSocketForm() (*gtk.Label, *socketForm, error) {
	defer config.LogStart("ConnectScreen.buildSocketForm", nil)()

	label, err := gtk.LabelNew("Socket")
	if err != nil {
		return nil, nil, err
	}
	label.Show()

	frm, err := socketForm{}.Init(c.w)
	if err != nil {
		return nil, nil, err
	}
	frm.Show()

	return label, frm, nil
}

func (c *ConnectScreen) buildSshForm() (*gtk.Label, *sshForm, error) {
	defer config.LogStart("ConnectScreen.buildSshForm", nil)()

	label, err := gtk.LabelNew("SSH")
	if err != nil {
		return nil, nil, err
	}
	label.Show()

	frm, err := sshForm{}.Init(c.w)
	if err != nil {
		return nil, nil, err
	}
	frm.Show()

	return label, frm, nil
}

func (c *ConnectScreen) SetConnections(connections []*config.Connection) {
	defer config.LogStart("ConnectScreen.SetConnections", nil)()

	//c.connections = connections
	names := make([]string, len(connections))
	for i, con := range connections {
		names[i] = con.Name
	}

	c.ConnectionList.UpdateItems(StringSliceToStringers(names))
}

func (c *ConnectScreen) ClearForm() {
	defer config.LogStart("ConnectScreen.ClearForm", nil)()

	c.forms.active.Clear()
	c.btnConnect.SetSensitive(false)
	c.btnTest.SetSensitive(false)
	c.btnSave.SetSensitive(false)
}

func (c *ConnectScreen) Connecting(cancel func()) {
	defer config.LogStart("ConnectScreen.Connecting", nil)()

	c.formOverlay.Run(cancel)
	c.ConnectionList.SetSensitive(false)
}

func (c *ConnectScreen) CancelConnecting() {
	defer config.LogStart("ConnectScreen.CancelConnecting", nil)()

	c.formOverlay.Stop()
	c.ConnectionList.SetSensitive(true)
}

func (c *ConnectScreen) FocusForm() {
	defer config.LogStart("ConnectScreen.FocusForm", nil)()

	c.forms.active.GrabFocus()
}

func (c *ConnectScreen) SetConnection(conn *config.Connection) {
	defer config.LogStart("ConnectScreen.SetConnection", nil)()

	c.forms.mysql.tcpForm.Clear()
	c.forms.mysql.socketForm.Clear()
	c.forms.mysql.sshForm.Clear()

	switch conn.Adapter {
	case "sqlite":
		c.forms.active = c.forms.sqlite.form
		c.stack.SetVisibleChildName("sqlite")
	case "mysql":
		c.stack.SetVisibleChildName("mysql")
		switch conn.Type {
		case "tcp":
			c.forms.active = c.forms.mysql.tcpForm
			c.forms.mysql.notebook.SetCurrentPage(0)
		case "socket":
			c.forms.active = c.forms.mysql.socketForm
			c.forms.mysql.notebook.SetCurrentPage(1)
		case "ssh":
			c.forms.active = c.forms.mysql.sshForm
			c.forms.mysql.notebook.SetCurrentPage(2)
		default:
			config.Errorf("invalid connection type '%s'", conn.Type)
			return
		}
	default:
		return
	}

	c.forms.active.SetConnection(conn)

	valid := driver.ValidateConnection(*conn)
	c.btnConnect.SetSensitive(valid)
	c.btnTest.SetSensitive(valid)
	c.btnSave.SetSensitive(valid)
}

func (c *ConnectScreen) ActiveConnectionIndex() int {
	defer config.LogStart("ConnectScreen.ActiveConnectionIndex", nil)()

	return c.ConnectionList.GetSelectedRow().GetIndex()
}

func (c *ConnectScreen) onConnect() {
	defer config.LogStart("ConnectScreen.onConnect", nil)()

	c.btnConnect.Emit("activate")
}

func (c *ConnectScreen) buildMysqlForms() (*gtk.Box, error) {
	defer config.LogStart("ConnectScreen.buildMysqlForms", nil)()

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	box.Show()

	c.forms.mysql.notebook, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	c.forms.mysql.notebook.Show()
	c.forms.mysql.notebook.SetCanFocus(true)
	box.Add(c.forms.mysql.notebook)

	box.SetSizeRequest(300, 200)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetHAlign(gtk.ALIGN_CENTER)

	// TCP
	{
		label, frm, err := c.buildTcpForm()
		if err != nil {
			return nil, err
		}
		frm.onChange(c.onFormChanged)

		c.forms.mysql.notebook.AppendPage(frm, label)
		c.forms.mysql.tcpForm = frm
	}

	// SOCKET
	{
		label, frm, err := c.buildSocketForm()
		if err != nil {
			return nil, err
		}
		frm.onChange(c.onFormChanged)

		c.forms.mysql.notebook.AppendPage(frm, label)
		c.forms.mysql.socketForm = frm
	}

	// SSH
	{
		label, frm, err := c.buildSshForm()
		if err != nil {
			return nil, err
		}
		frm.onChange(c.onFormChanged)

		c.forms.mysql.notebook.AppendPage(frm, label)
		c.forms.mysql.sshForm = frm
	}

	c.forms.active = c.forms.mysql.tcpForm
	// forms need to be built before can handle this signal
	c.forms.mysql.notebook.Connect("switch-page", c.onChangeCurrentPage)

	return box, nil
}

func (c *ConnectScreen) buildSqliteForms() (*sqliteForm, error) {
	defer config.LogStart("ConnectScreen.buildSqliteForms", nil)()

	frm, err := sqliteForm{}.Init(c.w)
	if err != nil {
		return nil, err
	}
	frm.Show()

	c.forms.sqlite.form = frm
	frm.onChange(c.onFormChanged)
	return frm, nil
}

func (c *ConnectScreen) GetFormConnection() *config.Connection {
	defer config.LogStart("ConnectScreen.GetFormConnection", nil)()

	conn, _ := c.forms.active.GetConnection()
	return conn
}

func (c *ConnectScreen) onChangeCurrentPage(_ *gtk.Notebook, _ gtk.IWidget, currentPage int) {
	defer config.LogStart("ConnectScreen.onChangeCurrentPage", nil)()

	switch currentPage {
	case 0:
		c.forms.active = c.forms.mysql.tcpForm
	case 1:
		c.forms.active = c.forms.mysql.socketForm
	case 2:
		c.forms.active = c.forms.mysql.sshForm
	}

	conn, isNew := c.forms.active.GetConnection()
	if isNew {
		c.btnConnect.SetSensitive(false)
		c.btnTest.SetSensitive(false)
		c.btnSave.SetSensitive(false)
		return
	}

	c.btnConnect.SetSensitive(driver.ValidateConnection(*conn))
	c.btnTest.SetSensitive(driver.ValidateConnection(*conn))
	c.btnSave.SetSensitive(driver.ValidateConnection(*conn))
}

func (c *ConnectScreen) onFormChanged(f form) {
	defer config.LogStart("ConnectScreen.onFormChanged", nil)()

	// NOTE: kinda of a hack
	//       gtk.Stack doesn't emit a signal when the selected item changes
	//       so if a forms is being changed, it means it's active :shrug:
	c.forms.active = f

	conn, _ := c.forms.active.GetConnection()
	valid := driver.ValidateConnection(*conn)
	c.btnConnect.SetSensitive(valid)
	c.btnTest.SetSensitive(valid)
	c.btnSave.SetSensitive(valid)
}

func (c *ConnectScreen) onSubmit(_ *gtk.Entry, e *gdk.Event) {
	defer config.LogStart("ConnectScreen.onSubmit", nil)()

	conn, _ := c.forms.active.GetConnection()
	if driver.ValidateConnection(*conn) {
		keyEvent := gdk.EventKeyNewFromEvent(e)
		if keyEvent.KeyVal() == 65293 && keyEvent.State()&gdk.CONTROL_MASK > 0 {
			c.btnConnect.Emit("activate")
		}

		return
	}
}
