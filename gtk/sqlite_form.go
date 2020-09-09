package gtk

import (
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
)

type sqliteForm struct {
	*gtk.Grid
	fields []string
	conn   *config.Connection

	entryName   *gtk.Entry
	btnFile     *gtk.FileChooserButton
	fileChooser *gtk.FileChooserDialog
	filename    string

	labelName *gtk.Label
	labelFile *gtk.Label
	queries   []config.Query
}

func (f sqliteForm) Init(w *Window) (*sqliteForm, error) {
	defer config.LogStart("sqliteForm.Init", nil)()

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
	f.labelName.Show()
	f.labelName.SetHAlign(gtk.ALIGN_START)

	f.entryName, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	f.entryName.Show()

	f.labelFile, err = gtk.LabelNew("File")
	if err != nil {
		return nil, err
	}
	f.labelFile.Show()
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
	f.btnFile.Show()

	f.fileChooser.Connect("file-activated", f.onFileSet)

	f.Attach(f.labelName, 0, 1, 1, 1)
	f.Attach(f.entryName, 1, 1, 2, 1)

	f.Attach(f.labelFile, 0, 2, 1, 1)
	f.Attach(f.btnFile, 1, 2, 2, 1)

	return &f, nil
}

func (f *sqliteForm) Clear() {
	defer config.LogStart("sqliteForm.Clear", nil)()

	f.conn = nil
	f.entryName.SetText("")
	f.filename = ""
}

func (f *sqliteForm) GrabFocus() {
	defer config.LogStart("sqliteForm.GrabFocus", nil)()

	f.entryName.GrabFocus()
}

func (f *sqliteForm) SetConnection(conn *config.Connection) {
	defer config.LogStart("sqliteForm.SetConnection", nil)()

	f.conn = conn
	f.entryName.SetText(conn.Name)
	f.btnFile.SelectFilename(conn.File)
	f.filename = conn.File
}

func (f *sqliteForm) GetConnection() (*config.Connection, bool) {
	defer config.LogStart("sqliteForm.GetConnection", nil)()

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

	return conn, newConn
}

func (f *sqliteForm) onChange(fn func(form)) {
	defer config.LogStart("sqliteForm.onChange", nil)()

	ff := func() { fn(f) }
	f.entryName.Connect("key-release-event", ff)
	f.fileChooser.Connect("file-activated", ff)
}

func (f *sqliteForm) onFileSet() {
	defer config.LogStart("sqliteForm.onFileSet", nil)()

	f.filename = f.fileChooser.GetFilename()
}
