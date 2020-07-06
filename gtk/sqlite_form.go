package gtk

import (
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
)

type sqliteForm struct {
	*gtk.Grid
	fields []string
	conn   *config.Connection

	entryName   *gtk.Entry
	btnFile     *gtk.FileChooserButton
	fileChooser *gtk.FileChooserDialog
	filename    string
	//entryUser     *gtk.Entry
	//entryPassword *gtk.Entry
	//entryDatabase *gtk.Entry

	labelName *gtk.Label
	labelFile *gtk.Label
	//labelUser     *gtk.Label
	//labelPassword *gtk.Label
	//labelDatabase *gtk.Label
	queries []config.Query
}

func (f sqliteForm) Init(w *Window) (*sqliteForm, error) {
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

	f.labelFile, err = gtk.LabelNew("File")
	if err != nil {
		return nil, err
	}
	f.labelFile.SetHAlign(gtk.ALIGN_START)

	f.fileChooser, err = gtk.FileChooserDialogNewWith2Buttons("Select", &w.Window, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Ok", gtk.RESPONSE_ACCEPT, "Cancel", gtk.RESPONSE_CANCEL)
	if err != nil {
		return nil, err
	}

	f.btnFile, err = gtk.FileChooserButtonNewWithDialog(f.fileChooser)
	if err != nil {
		return nil, err
	}
	f.fileChooser.Connect("file-activated", f.onFileSet)

	/*
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
	*/

	f.Attach(f.labelName, 0, 1, 1, 1)
	f.Attach(f.entryName, 1, 1, 2, 1)

	f.Attach(f.labelFile, 0, 2, 1, 1)
	f.Attach(f.btnFile, 1, 2, 2, 1)

	/*
			f.Attach(f.labelUser, 0, 3, 1, 1)
			f.Attach(f.entryUser, 1, 3, 2, 1)

			f.Attach(f.labelPassword, 0, 4, 1, 1)
			f.Attach(f.entryPassword, 1, 4, 2, 1)

		f.Attach(f.labelDatabase, 0, 5, 1, 1)
		f.Attach(f.entryDatabase, 1, 5, 2, 1)
	*/

	return &f, nil
}

func (f *sqliteForm) Clear() {
	f.conn = nil
	f.entryName.SetText("")
	f.filename = ""
	/*
		f.entryUser.SetText("")
		f.entryPassword.SetText("")
		f.entryDatabase.SetText("")
	*/
}

func (f *sqliteForm) GrabFocus() {
	f.entryName.GrabFocus()
}

func (f *sqliteForm) SetConnection(conn *config.Connection) {
	f.conn = conn
	f.entryName.SetText(conn.Name)
	f.btnFile.SelectFilename(conn.File)
	f.filename = conn.File
	/*
		f.entryUser.SetText(conn.User)
		f.entryPassword.SetText(conn.Password)
		f.entryDatabase.SetText(conn.Database)
	*/
}

func (f *sqliteForm) GetConnection() (*config.Connection, bool) {
	var newConn bool
	conn := f.conn
	if conn == nil {
		newConn = true
		conn = &config.Connection{}
	}
	conn.Adapter = "sqlite"
	conn.Database = filepath.Base(f.filename)
	conn.Name, _ = f.entryName.GetText()
	conn.File = f.filename
	/*
		conn.User, _ = f.entryUser.GetText()
		conn.Password, _ = f.entryPassword.GetText()
		conn.Database, _ = f.entryDatabase.GetText()
	*/

	return conn, newConn
}

func (f *sqliteForm) onChange(fn func(form)) {
	ff := func() { fn(f) }
	f.entryName.Connect("key-release-event", ff)
	f.fileChooser.Connect("file-activated", ff)
	/*
		f.entryUser.Connect("key-release-event", fn)
		f.entryPassword.Connect("key-release-event", fn)
		f.entryDatabase.Connect("key-release-event", fn)
	*/
}

func (f *sqliteForm) onFileSet() {
	f.filename = f.fileChooser.GetFilename()
}
