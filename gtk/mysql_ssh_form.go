package gtk

import (
	"fmt"
	"strconv"

	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
)

type sshForm struct {
	*gtk.Grid
	fields []string
	conn   *config.Connection

	entryName     *gtk.Entry
	entryDbHost   *gtk.Entry
	entryPort     *gtk.Entry
	entryUser     *gtk.Entry
	entryPassword *gtk.Entry
	entryDatabase *gtk.Entry

	labelName     *gtk.Label
	labelDbHost   *gtk.Label
	labelPort     *gtk.Label
	labelUser     *gtk.Label
	labelPassword *gtk.Label
	labelDatabase *gtk.Label

	entrySshHost  *gtk.Entry
	entrySshAgent *gtk.Entry

	labelSshHost  *gtk.Label
	labelSshAgent *gtk.Label
}

func (f sshForm) Init(_ *Window) (*sshForm, error) {
	defer config.LogStart("sshForm.Init", nil)()

	var err error
	f.Grid, err = gtk.GridNew()
	if err != nil {
		return nil, err
	}
	f.SetName("form")
	f.Grid.SetColumnHomogeneous(true)
	f.Grid.SetRowSpacing(5)

	f.labelName, err = gtk.LabelNew("Name")
	if err != nil {
		return nil, err
	}
	f.labelName.Show()
	f.labelName.SetHAlign(gtk.ALIGN_START)

	f.entryName, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryName.Show()

	f.labelDbHost, err = gtk.LabelNew("Host")
	if err != nil {
		return nil, err
	}
	f.labelDbHost.Show()
	f.labelDbHost.SetHAlign(gtk.ALIGN_START)

	f.entryDbHost, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryDbHost.Show()

	f.labelPort, err = gtk.LabelNew("Port")
	if err != nil {
		return nil, err
	}
	f.labelPort.Show()
	f.labelPort.SetHAlign(gtk.ALIGN_START)

	f.entryPort, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryPort.Show()
	f.entryPort.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	f.labelUser, err = gtk.LabelNew("User")
	if err != nil {
		return nil, err
	}
	f.labelUser.Show()
	f.labelUser.SetHAlign(gtk.ALIGN_START)

	f.entryUser, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryUser.Show()

	f.labelPassword, err = gtk.LabelNew("Password")
	if err != nil {
		return nil, err
	}
	f.labelPassword.Show()
	f.labelPassword.SetHAlign(gtk.ALIGN_START)

	f.entryPassword, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryPassword.Show()
	f.entryPassword.SetProperty("input-purpose", gtk.INPUT_PURPOSE_PASSWORD)
	f.entryPassword.SetProperty("visibility", false)

	f.labelDatabase, err = gtk.LabelNew("Database")
	if err != nil {
		return nil, err
	}
	f.labelDatabase.Show()
	f.labelDatabase.SetHAlign(gtk.ALIGN_START)

	f.entryDatabase, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryDatabase.Show()

	// SSH
	sshFrame, err := gtk.FrameNew("SSH")
	if err != nil {
		return nil, err
	}
	sshFrame.Show()
	sshFrame.SetProperty("shadow-type", gtk.SHADOW_NONE)

	sshBox, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	sshBox.Show()

	f.labelSshHost, err = gtk.LabelNew("Host")
	if err != nil {
		return nil, err
	}
	f.labelSshHost.Show()
	f.labelSshHost.SetHAlign(gtk.ALIGN_START)

	f.entrySshHost, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entrySshHost.Show()

	f.labelSshAgent, err = gtk.LabelNew("Agent")
	if err != nil {
		return nil, err
	}
	f.labelSshAgent.Show()
	f.labelSshAgent.SetHAlign(gtk.ALIGN_START)

	f.entrySshAgent, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entrySshAgent.Show()

	sshBox.SetColumnHomogeneous(true)
	sshBox.SetRowSpacing(5)
	sshBox.Attach(f.labelSshHost, 0, 0, 1, 1)
	sshBox.Attach(f.entrySshHost, 1, 0, 2, 1)
	sshBox.Attach(f.labelSshAgent, 0, 1, 1, 1)
	sshBox.Attach(f.entrySshAgent, 1, 1, 2, 1)
	sshFrame.Add(sshBox)

	//f.Grid.SetRowHomogeneous(true)
	f.Attach(f.labelName, 0, 0, 1, 1)
	f.Attach(f.entryName, 1, 0, 2, 1)

	f.Attach(f.labelDbHost, 0, 1, 1, 1)
	f.Attach(f.entryDbHost, 1, 1, 2, 1)

	f.Attach(f.labelPort, 0, 2, 1, 1)
	f.Attach(f.entryPort, 1, 2, 2, 1)

	f.Attach(f.labelUser, 0, 3, 1, 1)
	f.Attach(f.entryUser, 1, 3, 2, 1)

	f.Attach(f.labelPassword, 0, 4, 1, 1)
	f.Attach(f.entryPassword, 1, 4, 2, 1)

	f.Attach(f.labelDatabase, 0, 5, 1, 1)
	f.Attach(f.entryDatabase, 1, 5, 2, 1)

	f.Attach(sshFrame, 0, 6, 3, 2)

	return &f, nil
}

func (f *sshForm) newInput(l string) (*gtk.Label, *gtk.Entry, error) {
	defer config.LogStart("sshForm.newInput", nil)()

	label, err := gtk.LabelNew(l)
	if err != nil {
		return nil, nil, err
	}
	label.Show()
	label.SetHAlign(gtk.ALIGN_START)

	entry, err := gtk.EntryNew()
	if err != nil {
		return nil, nil, err
	}

	return label, entry, nil
}

func (f *sshForm) Clear() {
	defer config.LogStart("sshForm.Clear", nil)()

	f.conn = nil
	f.entryName.SetText("")
	f.entryDbHost.SetText("")
	f.entryPort.SetText("")
	f.entryUser.SetText("")
	f.entryPassword.SetText("")
	f.entryDatabase.SetText("")
	f.entrySshHost.SetText("")
	f.entrySshAgent.SetText("")
}

func (f *sshForm) GrabFocus() {
	defer config.LogStart("sshForm.GrabFocus", nil)()

	f.entryName.GrabFocus()
}

func (f *sshForm) SetConnection(conn *config.Connection) {
	defer config.LogStart("sshForm.SetConnection", nil)()

	f.conn = conn
	f.entryName.SetText(conn.Name)
	f.entryDbHost.SetText(conn.Host)
	f.entryPort.SetText(fmt.Sprintf("%d", conn.Port))
	f.entryUser.SetText(conn.User)
	f.entryPassword.SetText(conn.Password)
	f.entryDatabase.SetText(conn.Database)
	f.entrySshHost.SetText(conn.SshHost)
	f.entrySshAgent.SetText(conn.SshAgent)
}

func (f *sshForm) GetConnection() (*config.Connection, bool) {
	defer config.LogStart("sshForm.GetConnection", nil)()

	var newConn bool
	conn := f.conn
	if conn == nil {
		conn = &config.Connection{}
		newConn = true
	}

	conn.Adapter = "mysql"
	conn.Type = "ssh"
	conn.Name, _ = f.entryName.GetText()
	conn.Host, _ = f.entryDbHost.GetText()
	portS, _ := f.entryPort.GetText()
	if portS == "" {
		conn.Port = 3306
	} else {
		conn.Port, _ = strconv.Atoi(portS)
	}
	conn.User, _ = f.entryUser.GetText()
	conn.Password, _ = f.entryPassword.GetText()
	conn.Database, _ = f.entryDatabase.GetText()
	conn.SshHost, _ = f.entrySshHost.GetText()
	conn.SshAgent, _ = f.entrySshAgent.GetText()

	return conn, newConn
}

func (f *sshForm) onChange(fn func(form)) {
	defer config.LogStart("sshForm.onChange", nil)()

	ff := func() { fn(f) }
	f.entryName.Connect("key-release-event", ff)
	f.entryDbHost.Connect("key-release-event", ff)
	f.entryPort.Connect("key-release-event", ff)
	f.entryUser.Connect("key-release-event", ff)
	f.entryPassword.Connect("key-release-event", ff)
	f.entryDatabase.Connect("key-release-event", ff)
	f.entrySshHost.Connect("key-release-event", ff)
	f.entrySshAgent.Connect("key-release-event", ff)
}
