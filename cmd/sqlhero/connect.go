package main

import (
	"fmt"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/ui/controls"
	"github.com/gotk3/gotk3/gtk"
)

type ConnectScreen struct {
	*gtk.Paned
	connectionList *controls.List
	connections    []*config.Connection
	activeForm     *Form
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

	c.connectionList, err = controls.NewList(nil)
	if err != nil {
		return err
	}

	c.connectionList.Connect("row-activated", c.onConnectionActivated)
	c.connectionList.Connect("row-selected", c.onConnectionSelected)

	c.connectionList.SetHExpand(true)
	c.connectionList.SetVExpand(true)
	frame1.Add(c.connectionList)

	c.Paned.Pack1(frame1, false, true)
	c.Paned.Pack2(frame2, false, false)

	forms, err := c.forms()
	if err != nil {
		return err
	}

	frame2.Add(forms)
	c.Paned.ShowAll()

	return nil
}

type Form struct {
	inputs map[string]*gtk.Entry
	labels map[string]*gtk.Label
}

func (f *Form) New(fields []string) (*Form, error) {
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

func (c *ConnectScreen) forms() (*gtk.Notebook, error) {
	nb, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}

	nb.SetSizeRequest(300, 200)
	nb.SetVAlign(gtk.ALIGN_CENTER)
	nb.SetHAlign(gtk.ALIGN_CENTER)
	nb.SetShowBorder(true)
	nb.SetCanFocus(true)

	forms := []*Form{}
	content, label, form, err := c.nbPage("Standard", []string{"Name", "Host", "Port", "User", "Password", "Database"})
	if err != nil {
		return nil, err
	}
	c.activeForm = form

	forms = append(forms, form)
	nb.AppendPage(content, label)

	return nb, nil
}

func (c *ConnectScreen) nbPage(title string, fields []string) (gtk.IWidget, gtk.IWidget, *Form, error) {
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

	frm, err := (&Form{}).New(fields)
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

	conn := c.connections[index]
	c.activeForm.inputs["Name"].SetText(conn.Name)
	c.activeForm.inputs["Host"].SetText(conn.Host)
	c.activeForm.inputs["Port"].SetText(fmt.Sprintf("%d", conn.Port))
	c.activeForm.inputs["User"].SetText(conn.Username)
	c.activeForm.inputs["Password"].SetText(conn.Password)
	c.activeForm.inputs["Database"].SetText(conn.Database)
}

func (c *ConnectScreen) UpdateForm(conn config.Connection) {
}
