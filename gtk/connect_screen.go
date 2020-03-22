package gtk

import (
	"fmt"

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
	connections    []*config.Connection
	activeForm     *form

	activeConnection controls.MVar
	btnSave          *gtk.Button
	btnConnect       *gtk.Button
	btnTest          *gtk.Button
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

	c.connectionList.Connect("row-activated", c.onConnectionActivated)
	c.connectionList.Connect("row-selected", c.onConnectionSelected)

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
		c.activeForm.inputs["Name"].SetText("")
		c.activeForm.inputs["Host"].SetText("")
		c.activeForm.inputs["Port"].SetText("")
		c.activeForm.inputs["User"].SetText("")
		c.activeForm.inputs["Password"].SetText("")
		c.activeForm.inputs["Database"].SetText("")

		c.activeForm.inputs["Name"].GrabFocus()
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

		c.connect(c.connections[index])
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

func (c *ConnectScreen) nbPage(title string, fields []string) (gtk.IWidget, gtk.IWidget, *form, error) {
	label, err := gtk.LabelNew(title)
	if err != nil {
		return nil, nil, nil, err
	}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	labelsBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	inputsBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, nil, nil, err
	}
	hbox.PackStart(labelsBox, true, false, 5)
	hbox.PackStart(inputsBox, false, false, 5)

	frm, err := (&form{}).new(fields)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, field := range fields {
		labelsBox.PackStart(frm.labels[field], true, true, 5)
		inputsBox.PackStart(frm.inputs[field], true, true, 5)
	}

	hbox.ShowAll()
	label.Show()
	return hbox, label, frm, nil
}

func (c *ConnectScreen) SetConnections(connections []*config.Connection) {
	c.connections = connections
	names := make([]string, len(connections))
	for i, con := range c.connections {
		names[i] = con.Name
	}

	c.connectionList.UpdateItems(names)
}

func (c *ConnectScreen) onConnectionSelected() {
	index, ok := c.connectionList.SelectedItemIndex()
	if !ok {
		return
	}

	conn := c.connections[index]
	c.activeConnection.Set(conn)

	c.activeForm.inputs["Name"].SetText(conn.Name)
	c.activeForm.inputs["Host"].SetText(conn.Host)
	c.activeForm.inputs["Port"].SetText(fmt.Sprintf("%d", conn.Port))
	c.activeForm.inputs["User"].SetText(conn.Username)
	c.activeForm.inputs["Password"].SetText(conn.Password)
	c.activeForm.inputs["Database"].SetText(conn.Database)
}

func (c *ConnectScreen) onConnectionActivated() {
	index, ok := c.connectionList.ActiveItemIndex()
	if !ok {
		return
	}

	c.connect(c.connections[index])
}

func (c *ConnectScreen) connect(conn *config.Connection) {
	c.activeConnection.Set(conn)
	c.activeForm.inputs["Name"].SetText(conn.Name)
	c.activeForm.inputs["Host"].SetText(conn.Host)
	c.activeForm.inputs["Port"].SetText(fmt.Sprintf("%d", conn.Port))
	c.activeForm.inputs["User"].SetText(conn.Username)
	c.activeForm.inputs["Password"].SetText(conn.Password)
	c.activeForm.inputs["Database"].SetText(conn.Database)
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

func (c *ConnectScreen) ActiveConnection() *config.Connection {
	return c.activeConnection.Get().(*config.Connection)
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

	forms := []*form{}
	content, label, frm, err := c.nbPage("Standard", []string{"Name", "Host", "Port", "User", "Password", "Database"})
	if err != nil {
		return nil, err
	}
	c.activeForm = frm

	forms = append(forms, frm)
	nb.AppendPage(content, label)

	return box, nil
}

type form struct {
	inputs map[string]*gtk.Entry
	labels map[string]*gtk.Label
}

func (f *form) new(fields []string) (*form, error) {
	f.inputs = map[string]*gtk.Entry{}
	f.labels = map[string]*gtk.Label{}

	for _, field := range fields {
		l, err := gtk.LabelNew(field)
		if err != nil {
			return nil, err
		}
		l.SetHAlign(gtk.ALIGN_START)

		e, err := gtk.EntryNew()
		if err != nil {
			return nil, err
		}

		f.inputs[field] = e
		f.labels[field] = l
	}

	return f, nil
}
