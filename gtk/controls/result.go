package controls

import (
	"log"
	"reflect"

	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Result struct {
	*gtk.TreeView
	cols  []driver.ColDef
	data  [][]interface{}
	store *gtk.ListStore
}

func NewResult(cols []driver.ColDef, data [][]interface{}) (u *Result, err error) {
	u = &Result{}

	u.TreeView, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}

	if len(cols) > 0 {
		u.UpdateData(cols, data)
	}

	return u, nil
}

func (u *Result) UpdateData(cols []driver.ColDef, data [][]interface{}) error {
	for i := range u.cols {
		u.TreeView.RemoveColumn(u.TreeView.GetColumn(i))
	}
	u.cols = cols
	u.data = data

	columns := make([]glib.Type, len(cols))
	for i, col := range cols {
		u.TreeView.InsertColumn(u.createColumn(col.Name, i), i)
		columns[i] = glib.TYPE_STRING
	}

	for _, row := range data {
		for i, col := range row {
			if col == nil {
				continue
			}

			switch reflect.TypeOf(col).Kind() {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
				columns[i] = glib.TYPE_INT64
			case reflect.Bool:
				columns[i] = glib.TYPE_BOOLEAN
			default:
				columns[i] = glib.TYPE_STRING
			}
		}
		break
	}

	var err error
	u.store, err = gtk.ListStoreNew(columns...)
	if err != nil {
		return err
	}
	u.TreeView.SetModel(u.store)

	for _, row := range data {
		u.AddRow(row)
	}

	return nil
}

func (u *Result) AddRow(row []interface{}) {
	// Get an iterator for a new row at the end of the list store
	iter := u.store.Append()

	if len(row) != len(u.cols) {
		log.Fatal("wrong row length")
	}

	columns := make([]int, len(row))
	for i, d := range row {
		columns[i] = i
		if s, ok := d.([]uint8); ok {
			row[i] = string(s)
		}
	}

	// Set the contents of the list store row that the iterator represents
	err := u.store.Set(iter,
		columns,
		row,
	)

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}

func (u *Result) createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	cellRenderer.SetProperty("editable", true)
	cellRenderer.Connect("edited", u.edited, id)

	// i think "text" refers to a property of the column.
	// `"text", id` means that the text source for the column should come from
	// the listore column with id = `id`
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

func (u *Result) edited(cell *gtk.CellRendererText, path string, newText string, userData interface{}) {
	columnIndex := userData.(int)

	tpath, err := gtk.TreePathNewFromString(path)
	if err != nil {
		log.Fatal(err)
	}

	iter, err := u.store.GetIter(tpath)
	if err != nil {
		log.Fatal(err)
	}

	err = u.store.SetValue(iter, columnIndex, newText)
	if err != nil {
		log.Fatal(err)
	}
}
