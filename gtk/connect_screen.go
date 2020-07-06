package gtk

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type form interface {
	Clear()
	GetConnection() (*config.Connection, bool)
	GrabFocus()
	SetConnection(conn *config.Connection)
}

type ConnectScreen struct {
	w *Window
	*gtk.Paned
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
	Config() *config.Config
}

func (c ConnectScreen) Init(w *Window, ctrl connectScreenCtrl) (*ConnectScreen, error) {
	var err error
	c.w = w
	c.ctrl = ctrl

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

	c.ConnectionList, err = List{}.Init(c.w, &ListOptions{SelectOnRightClick: true, StockIcon: "gtk-connect"}, ctrl)
	if err != nil {
		return nil, err
	}

	c.ConnectionList.OnButtonPress(c.onConnectListButtonPress)
	c.ConnectionList.SetHExpand(true)
	c.ConnectionList.SetVExpand(true)
	frame1.Add(c.ConnectionList)

	mysqlForms, err := c.buildMysqlForms()
	if err != nil {
		return nil, err
	}
	sqliteForm, err := c.buildSqliteForms()
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

	mysqlForms.ShowAll()

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	vbox.SetSizeRequest(300, 200)
	vbox.SetVAlign(gtk.ALIGN_CENTER)
	vbox.SetHAlign(gtk.ALIGN_CENTER)

	vbox.PackEnd(c.btnBox, false, true, 0)

	adapterStack, err := gtk.StackSwitcherNew()
	if err != nil {
		return nil, err
	}
	adapterStack.SetVExpand(true)
	adapterStack.SetHExpand(true)
	c.stack, err = gtk.StackNew()
	if err != nil {
		return nil, err
	}
	c.stack.SetHomogeneous(true)
	c.stack.AddTitled(mysqlForms, "mysql", "Mysql")
	c.stack.AddTitled(sqliteForm, "sqlite", "SQLite")
	c.stack.SetVisibleChildName("mysql")
	c.stack.SetVExpand(true)
	c.stack.SetHExpand(true)

	adapterStack.SetStack(c.stack)

	vbox.PackStart(adapterStack, false, true, 0)
	vbox.PackStart(c.stack, true, false, 0)

	frame2.Add(vbox)
	c.formOverlay, err = CancelOverlay{}.Init(frame2)
	if err != nil {
		return nil, err
	}

	c.Paned.Pack1(frame1, false, true)
	c.Paned.Pack2(c.formOverlay, true, false)
	c.Paned.ShowAll()

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

func (c *ConnectScreen) buildTcpForm() (*gtk.Label, *tcpForm, error) {
	label, err := gtk.LabelNew("TCP/IP")
	if err != nil {
		return nil, nil, err
	}

	frm, err := tcpForm{}.Init(c.w)
	if err != nil {
		return nil, nil, err
	}

	frm.ShowAll()
	label.Show()

	return label, frm, nil
}

func (c *ConnectScreen) buildSocketForm() (*gtk.Label, *socketForm, error) {
	label, err := gtk.LabelNew("Socket")
	if err != nil {
		return nil, nil, err
	}

	frm, err := socketForm{}.Init(c.w)
	if err != nil {
		return nil, nil, err
	}

	frm.ShowAll()
	label.Show()

	return label, frm, nil
}

func (c *ConnectScreen) buildSshForm() (*gtk.Label, *sshForm, error) {
	label, err := gtk.LabelNew("SSH")
	if err != nil {
		return nil, nil, err
	}

	frm, err := sshForm{}.Init(c.w)
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
			c.ctrl.Config().Errorf("invalid connection type '%s'", conn.Type)
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
	return c.ConnectionList.GetSelectedRow().GetIndex()
}

func (c *ConnectScreen) onConnect() {
	c.btnConnect.Emit("activate")
}

func (c *ConnectScreen) Dispose() {
}

func (c *ConnectScreen) buildMysqlForms() (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	c.forms.mysql.notebook, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
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
	frm, err := sqliteForm{}.Init(c.w)
	if err != nil {
		return nil, err
	}

	c.forms.sqlite.form = frm
	frm.onChange(c.onFormChanged)
	return frm, nil
}

func (c *ConnectScreen) GetFormConnection() *config.Connection {
	conn, _ := c.forms.active.GetConnection()
	//conn.Queries = c.activeForm.queries
	return conn
}

func (c *ConnectScreen) onChangeCurrentPage(_ *gtk.Notebook, _ gtk.IWidget, currentPage int) {
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
	conn, _ := c.forms.active.GetConnection()
	if driver.ValidateConnection(*conn) {
		keyEvent := gdk.EventKeyNewFromEvent(e)
		if keyEvent.KeyVal() == 65293 && keyEvent.State()&gdk.CONTROL_MASK > 0 {
			c.btnConnect.Emit("activate")
		}

		return
	}
}
