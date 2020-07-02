package gtk

import (
	"fmt"
	"strconv"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gtk"
)

type tcpForm struct {
	*gtk.Box
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

func (f *tcpForm) GetConnection() *config.Connection {
	f.conn.Type = "socket"
	f.conn.Name, _ = f.entryName.GetText()
	f.conn.Host, _ = f.entryHost.GetText()
	portS, _ := f.entryPort.GetText()
	if portS == "" {
		f.conn.Port = 3306
	} else {
		f.conn.Port, _ = strconv.Atoi(portS)
	}
	f.conn.User, _ = f.entryUser.GetText()
	f.conn.Password, _ = f.entryPassword.GetText()
	f.conn.Database, _ = f.entryDatabase.GetText()

	return f.conn
}

func (f *tcpForm) onChange(fn interface{}) {
	f.entryName.Connect("key-release-event", fn)
	f.entryHost.Connect("key-release-event", fn)
	f.entryPort.Connect("key-release-event", fn)
	f.entryUser.Connect("key-release-event", fn)
	f.entryPassword.Connect("key-release-event", fn)
	f.entryDatabase.Connect("key-release-event", fn)
}
