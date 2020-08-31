package gtk

import (
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
)

type socketForm struct {
	*gtk.Grid
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
}

func (f socketForm) Init(_ *Window) (*socketForm, error) {
	var err error

	f.Grid, err = gtk.GridNew()
	if err != nil {
		return nil, err
	}
	f.SetName("form")
	f.SetColumnHomogeneous(true)
	f.SetRowSpacing(5)

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

	f.Attach(f.labelName, 0, 1, 1, 1)
	f.Attach(f.entryName, 1, 1, 2, 1)

	f.Attach(f.labelSocket, 0, 2, 1, 1)
	f.Attach(f.entrySocket, 1, 2, 2, 1)

	f.Attach(f.labelUser, 0, 3, 1, 1)
	f.Attach(f.entryUser, 1, 3, 2, 1)

	f.Attach(f.labelPassword, 0, 4, 1, 1)
	f.Attach(f.entryPassword, 1, 4, 2, 1)

	f.Attach(f.labelDatabase, 0, 5, 1, 1)
	f.Attach(f.entryDatabase, 1, 5, 2, 1)

	return &f, nil
}

func (f *socketForm) Clear() {
	f.conn = nil
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

func (f *socketForm) GetConnection() (*config.Connection, bool) {
	var newConn bool
	conn := f.conn
	if conn == nil {
		newConn = true
		conn = &config.Connection{}
	}
	conn.Adapter = "mysql"
	conn.Type = "socket"
	conn.Name, _ = f.entryName.GetText()
	conn.Socket, _ = f.entrySocket.GetText()
	conn.User, _ = f.entryUser.GetText()
	conn.Password, _ = f.entryPassword.GetText()
	conn.Database, _ = f.entryDatabase.GetText()

	return conn, newConn
}

func (f *socketForm) onChange(fn func(form)) {
	ff := func() { fn(f) }
	f.entryName.Connect("key-release-event", ff)
	f.entrySocket.Connect("key-release-event", ff)
	f.entryUser.Connect("key-release-event", ff)
	f.entryPassword.Connect("key-release-event", ff)
	f.entryDatabase.Connect("key-release-event", ff)
}
