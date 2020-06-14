package gtk

import (
	"fmt"
	"strconv"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ConnectScreen struct {
	*gtk.Paned
	ConnectionList *List
	activeForm     *stdform

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
	frame1.SetSizeRequest(50, -1)
	frame2.SetShadowType(gtk.SHADOW_IN)
	frame2.SetSizeRequest(50, -1)

	c.ConnectionList, err = List{}.Init(ListOptions{SelectOnRightClick: true, StockIcon: "gtk-connect"})
	if err != nil {
		return nil, err
	}

	c.ConnectionList.OnButtonPress(c.onConnectListButtonPress)
	c.ConnectionList.Connect("row-activated", func() {
		c.btnConnect.Emit("activate")
	})

	c.ConnectionList.SetHExpand(true)
	c.ConnectionList.SetVExpand(true)
	frame1.Add(c.ConnectionList)

	forms, err := c.forms()
	if err != nil {
		return nil, err
	}

	btnBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	btnBox.SetLayout(gtk.BUTTONBOX_EDGE)

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

	btnBox.Add(c.btnConnect)
	btnBox.Add(c.btnTest)
	btnBox.Add(c.btnSave)

	forms.Add(btnBox)
	frame2.Add(forms)

	c.Paned.Pack1(frame1, false, true)
	c.Paned.Pack2(frame2, false, false)

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

func (c *ConnectScreen) nbPage(title string) (gtk.IWidget, *stdform, error) {
	label, err := gtk.LabelNew(title)
	if err != nil {
		return nil, nil, err
	}

	frm, err := stdform{}.init()
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
	c.activeForm.Clear()
	c.btnConnect.SetSensitive(false)
	c.btnTest.SetSensitive(false)
	c.btnSave.SetSensitive(false)
}

func (c *ConnectScreen) FocusForm() {
	c.activeForm.GrabFocus()
}

func (c *ConnectScreen) SetFormConnection(conn *config.Connection) {
	c.activeForm.SetConnection(conn)
	c.activeForm.queries = conn.Queries
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

func (c *ConnectScreen) forms() (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	nb, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	box.Add(nb)

	box.SetSizeRequest(300, 200)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetHAlign(gtk.ALIGN_CENTER)
	nb.SetShowBorder(true)
	nb.SetCanFocus(true)

	forms := []*stdform{}
	label, frm, err := c.nbPage("Standard")
	if err != nil {
		return nil, err
	}
	c.activeForm = frm
	c.activeForm.onChange(func(_ *gtk.Entry, e *gdk.Event) bool {
		conn := c.activeForm.GetConnection()

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

	forms = append(forms, frm)
	nb.AppendPage(frm, label)

	return box, nil
}

func (c *ConnectScreen) GetFormConnection() *config.Connection {
	conn := c.activeForm.GetConnection()
	conn.Queries = c.activeForm.queries
	return conn
}

type stdform struct {
	*gtk.Box
	fields []string

	entryName     *gtk.Entry
	entryHost     *gtk.Entry
	entryPort     *gtk.Entry
	entryUser     *gtk.Entry
	entryPassword *gtk.Entry
	entryDatabase *gtk.Entry

	labelName     *gtk.Label
	labelHost     *gtk.Label
	labelPort     *gtk.Label
	labelUser     *gtk.Label
	labelPassword *gtk.Label
	labelDatabase *gtk.Label
	queries       []config.Query
}

func (f stdform) init() (*stdform, error) {
	var err error

	f.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	labelsBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	inputsBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	f.Box.PackStart(labelsBox, true, false, 5)
	f.Box.PackStart(inputsBox, false, false, 5)

	f.labelName, err = gtk.LabelNew("Name")
	if err != nil {
		return nil, err
	}
	f.labelName.SetHAlign(gtk.ALIGN_START)

	f.entryName, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	f.labelHost, err = gtk.LabelNew("Host")
	if err != nil {
		return nil, err
	}
	f.labelHost.SetHAlign(gtk.ALIGN_START)

	f.entryHost, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	f.labelPort, err = gtk.LabelNew("Port")
	if err != nil {
		return nil, err
	}
	f.labelPort.SetHAlign(gtk.ALIGN_START)

	f.entryPort, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryPort.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	f.labelUser, err = gtk.LabelNew("User")
	if err != nil {
		return nil, err
	}
	f.labelUser.SetHAlign(gtk.ALIGN_START)

	f.entryUser, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	f.labelPassword, err = gtk.LabelNew("Password")
	if err != nil {
		return nil, err
	}
	f.labelPassword.SetHAlign(gtk.ALIGN_START)

	f.entryPassword, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryPassword.SetProperty("input-purpose", gtk.INPUT_PURPOSE_PASSWORD)
	f.entryPassword.SetProperty("visibility", false)

	f.labelDatabase, err = gtk.LabelNew("Database")
	if err != nil {
		return nil, err
	}
	f.labelDatabase.SetHAlign(gtk.ALIGN_START)

	f.entryDatabase, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	labelsBox.PackStart(f.labelName, true, true, 5)
	inputsBox.PackStart(f.entryName, true, true, 5)
	labelsBox.PackStart(f.labelHost, true, true, 5)
	inputsBox.PackStart(f.entryHost, true, true, 5)
	labelsBox.PackStart(f.labelPort, true, true, 5)
	inputsBox.PackStart(f.entryPort, true, true, 5)
	labelsBox.PackStart(f.labelUser, true, true, 5)
	inputsBox.PackStart(f.entryUser, true, true, 5)
	labelsBox.PackStart(f.labelPassword, true, true, 5)
	inputsBox.PackStart(f.entryPassword, true, true, 5)
	labelsBox.PackStart(f.labelDatabase, true, true, 5)
	inputsBox.PackStart(f.entryDatabase, true, true, 5)

	return &f, nil
}

func (f *stdform) Clear() {
	f.entryName.SetText("")
	f.entryHost.SetText("")
	f.entryPort.SetText("")
	f.entryUser.SetText("")
	f.entryPassword.SetText("")
	f.entryDatabase.SetText("")
}

func (f *stdform) GrabFocus() {
	f.entryName.GrabFocus()
}

func (f *stdform) SetConnection(conn *config.Connection) {
	f.entryName.SetText(conn.Name)
	f.entryHost.SetText(conn.Host)
	f.entryPort.SetText(fmt.Sprintf("%d", conn.Port))
	f.entryUser.SetText(conn.User)
	f.entryPassword.SetText(conn.Password)
	f.entryDatabase.SetText(conn.Database)
}

func (f *stdform) GetConnection() *config.Connection {
	conn := &config.Connection{}
	conn.Name, _ = f.entryName.GetText()
	conn.Host, _ = f.entryHost.GetText()
	portS, _ := f.entryPort.GetText()
	if portS == "" {
		conn.Port = 3306
	} else {
		conn.Port, _ = strconv.Atoi(portS)
	}
	conn.User, _ = f.entryUser.GetText()
	conn.Password, _ = f.entryPassword.GetText()
	conn.Database, _ = f.entryDatabase.GetText()

	return conn
}

func (f *stdform) onChange(fn interface{}) {
	f.entryName.Connect("key-release-event", fn)
	f.entryHost.Connect("key-release-event", fn)
	f.entryPort.Connect("key-release-event", fn)
	f.entryUser.Connect("key-release-event", fn)
	f.entryPassword.Connect("key-release-event", fn)
	f.entryDatabase.Connect("key-release-event", fn)
}
