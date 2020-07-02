package gtk

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gtk"
)

type socketForm struct {
	*gtk.Box
	fields []string
	conn   *config.Connection

	entryName     *gtk.Entry
	entrySocket   *gtk.Entry
	entryUser     *gtk.Entry
	entryPassword *gtk.Entry
	entryDatabase *gtk.Entry

	labelName     *gtk.Label
	labelSocket   *gtk.Label
	labelUser     *gtk.Label
	labelPassword *gtk.Label
	labelDatabase *gtk.Label
	queries       []config.Query
}

func (f socketForm) Init() (*socketForm, error) {
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

	f.labelSocket, err = gtk.LabelNew("Socket")
	if err != nil {
		return nil, err
	}
	f.labelSocket.SetHAlign(gtk.ALIGN_START)

	f.entrySocket, err = gtk.EntryNew()
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
	labelsBox.PackStart(f.labelSocket, true, true, 5)
	inputsBox.PackStart(f.entrySocket, true, true, 5)
	labelsBox.PackStart(f.labelUser, true, true, 5)
	inputsBox.PackStart(f.entryUser, true, true, 5)
	labelsBox.PackStart(f.labelPassword, true, true, 5)
	inputsBox.PackStart(f.entryPassword, true, true, 5)
	labelsBox.PackStart(f.labelDatabase, true, true, 5)
	inputsBox.PackStart(f.entryDatabase, true, true, 5)

	return &f, nil
}

func (f *socketForm) Clear() {
	f.entryName.SetText("")
	f.entrySocket.SetText("")
	f.entryUser.SetText("")
	f.entryPassword.SetText("")
	f.entryDatabase.SetText("")
}

func (f *socketForm) GrabFocus() {
	f.entryName.GrabFocus()
}

func (f *socketForm) SetConnection(conn *config.Connection) {
	f.conn = conn
	f.entryName.SetText(conn.Name)
	f.entrySocket.SetText(conn.Socket)
	f.entryUser.SetText(conn.User)
	f.entryPassword.SetText(conn.Password)
	f.entryDatabase.SetText(conn.Database)
}

func (f *socketForm) GetConnection() *config.Connection {
	f.conn.Type = "tcp"
	f.conn.Name, _ = f.entryName.GetText()
	f.conn.Socket, _ = f.entrySocket.GetText()
	f.conn.User, _ = f.entryUser.GetText()
	f.conn.Password, _ = f.entryPassword.GetText()
	f.conn.Database, _ = f.entryDatabase.GetText()

	return f.conn
}

func (f *socketForm) onChange(fn interface{}) {
	f.entryName.Connect("key-release-event", fn)
	f.entrySocket.Connect("key-release-event", fn)
	f.entryUser.Connect("key-release-event", fn)
	f.entryPassword.Connect("key-release-event", fn)
	f.entryDatabase.Connect("key-release-event", fn)
}
