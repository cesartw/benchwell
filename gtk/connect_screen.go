package gtk

import (
	"fmt"
	"strconv"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk/controls"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func newConnectScreen() (*ConnectScreen, error) {
	cs := &ConnectScreen{}
	return cs, cs.init()
}

type ConnectScreen struct {
	*gtk.Paned
	connectionList *controls.List
	activeForm     *stdform

	activeConnectionIndex controls.MVar
	btnSave               *gtk.Button
	btnConnect            *gtk.Button
	btnTest               *gtk.Button
}

func (c *ConnectScreen) init() error {
	var err error

	c.Paned, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return err
	}

	c.Paned.SetHExpand(true)
	c.Paned.SetVExpand(true)

	frame1, err := gtk.FrameNew("")
	if err != nil {
		return err
	}

	frame2, err := gtk.FrameNew("")
	if err != nil {
		return err
	}

	frame1.SetShadowType(gtk.SHADOW_IN)
	frame1.SetSizeRequest(50, -1)
	frame2.SetShadowType(gtk.SHADOW_IN)
	frame2.SetSizeRequest(50, -1)

	c.connectionList, err = controls.NewList(controls.ListOptions{SelectOnRightClick: true})
	if err != nil {
		return err
	}

	c.connectionList.OnButtonPress(c.onConnectListButtonPress)

	c.connectionList.SetHExpand(true)
	c.connectionList.SetVExpand(true)
	frame1.Add(c.connectionList)

	forms, err := c.forms()
	if err != nil {
		return err
	}

	btnBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return err
	}
	btnBox.SetLayout(gtk.BUTTONBOX_EDGE)

	c.btnConnect, err = gtk.ButtonNew()
	if err != nil {
		return err
	}
	c.btnSave, err = gtk.ButtonNew()
	if err != nil {
		return err
	}
	c.btnTest, err = gtk.ButtonNew()
	if err != nil {
		return err
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

	c.Paned.ShowAll()

	return nil
}

func (c *ConnectScreen) onConnectListButtonPress(_ *gtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	m, err := gtk.MenuNew()
	if err != nil {
		return
	}

	mi, err := gtk.MenuItemNewWithLabel("New")
	if err != nil {
		return
	}
	mi.Connect("activate", func() {
		c.activeForm.Clear()
		c.activeForm.GrabFocus()
	})
	m.Add(mi)

	mi, err = gtk.MenuItemNewWithLabel("Connect")
	if err != nil {
		return
	}

	mi.Connect("activate", func() {
		index, ok := c.connectionList.SelectedItemIndex()
		if !ok {
			return
		}

		c.onConnect(index)
	})

	m.Add(mi)

	mi, err = gtk.MenuItemNewWithLabel("Test")
	if err != nil {
		return
	}
	mi.Connect("activate", func() {
		c.btnTest.Emit("clicked")
	})
	m.Add(mi)

	mi, err = gtk.MenuItemNewWithLabel("Delete")
	if err != nil {
		return
	}
	m.Add(mi)

	m.ShowAll()
	m.PopupAtPointer(e)
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

	c.connectionList.UpdateItems(names)
}

func (c *ConnectScreen) SetFormConnection(conn *config.Connection) {
	c.activeForm.SetConnection(conn)
	if conn.Valid() {
		c.btnConnect.SetSensitive(true)
		c.btnTest.SetSensitive(true)
		c.btnSave.SetSensitive(true)
	}
}

func (c *ConnectScreen) OnConnectionSelected(fn interface{}) {
	c.connectionList.Connect("row-selected", fn)
}

func (c *ConnectScreen) OnConnectionActivated(fn interface{}) {
	c.connectionList.Connect("row-activated", fn)
}

func (c *ConnectScreen) onConnect(index int) {
	c.btnConnect.Emit("activate")
}

func (c *ConnectScreen) OnConnect(f interface{}) {
	c.btnConnect.Connect("clicked", f)
}

func (c *ConnectScreen) OnSave(f interface{}) {
	c.btnSave.Connect("clicked", f)
}

func (c *ConnectScreen) OnTest(f interface{}) {
	c.btnTest.Connect("clicked", f)
}

func (c *ConnectScreen) ActiveConnectionIndex() int {
	if c.activeConnectionIndex.Get() == nil {
		return -1
	}
	return c.activeConnectionIndex.Get().(int)
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
			if keyEvent.KeyVal() == 65293 && keyEvent.State()&gdk.GDK_CONTROL_MASK > 0 {
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
	return c.activeForm.GetConnection()
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
