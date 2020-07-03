package gtk

import (
	"fmt"
	"strconv"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gtk"
)

type tcpForm struct {
	*gtk.Grid
	fields []string
	conn   *config.Connection

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

func (f tcpForm) Init() (*tcpForm, error) {
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

	f.Attach(f.labelName, 0, 0, 1, 1)
	f.Attach(f.entryName, 1, 0, 2, 1)

	f.Attach(f.labelHost, 0, 1, 1, 1)
	f.Attach(f.entryHost, 1, 1, 2, 1)

	f.Attach(f.labelPort, 0, 2, 1, 1)
	f.Attach(f.entryPort, 1, 2, 2, 1)

	f.Attach(f.labelUser, 0, 3, 1, 1)
	f.Attach(f.entryUser, 1, 3, 2, 1)

	f.Attach(f.labelPassword, 0, 4, 1, 1)
	f.Attach(f.entryPassword, 1, 4, 2, 1)

	f.Attach(f.labelDatabase, 0, 5, 1, 1)
	f.Attach(f.entryDatabase, 1, 5, 2, 1)

	return &f, nil
}

func (f *tcpForm) Clear() {
	f.entryName.SetText("")
	f.entryHost.SetText("")
	f.entryPort.SetText("")
	f.entryUser.SetText("")
	f.entryPassword.SetText("")
	f.entryDatabase.SetText("")
}

func (f *tcpForm) GrabFocus() {
	f.entryName.GrabFocus()
}

func (f *tcpForm) SetConnection(conn *config.Connection) {
	f.conn = conn
	f.entryName.SetText(conn.Name)
	f.entryHost.SetText(conn.Host)
	f.entryPort.SetText(fmt.Sprintf("%d", conn.Port))
	f.entryUser.SetText(conn.User)
	f.entryPassword.SetText(conn.Password)
	f.entryDatabase.SetText(conn.Database)
}

func (f *tcpForm) GetConnection() (*config.Connection, bool) {
	var newConn bool
	conn := f.conn
	if conn == nil {
		newConn = true
		conn = &config.Connection{Adapter: "mysql"}
	}

	conn.Type = "tcp"
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

	return conn, newConn
}

func (f *tcpForm) onChange(fn interface{}) {
	f.entryName.Connect("key-release-event", fn)
	f.entryHost.Connect("key-release-event", fn)
	f.entryPort.Connect("key-release-event", fn)
	f.entryUser.Connect("key-release-event", fn)
	f.entryPassword.Connect("key-release-event", fn)
	f.entryDatabase.Connect("key-release-event", fn)
}
