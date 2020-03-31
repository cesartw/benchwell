package controls

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type MODE int

const (
	MODE_RAW = iota
	MODE_DEF
)

type Result struct {
	*gtk.TreeView
	cols            []fmt.Stringer
	data            [][]interface{}
	store           *gtk.ListStore
	updateCallbacks []func([]driver.ColDef, []interface{}, []interface{}, string, int, int)

	mode MODE
}

func NewResult(cols []driver.ColDef, data [][]interface{}) (u *Result, err error) {
	u = &Result{}

	u.TreeView, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}

	u.TreeView.SetProperty("rubber-banding", true)
	u.TreeView.SetProperty("enable-grid-lines", gtk.TREE_VIEW_GRID_LINES_HORIZONTAL)

	if len(cols) > 0 {
		u.UpdateData(cols, data)
	}

	return u, nil
}

func (u *Result) UpdateData(cols []driver.ColDef, data [][]interface{}) error {
	for i := range u.cols {
		u.TreeView.RemoveColumn(u.TreeView.GetColumn(i))
	}
	u.cols = colDefSliceToStringerSlice(cols)
	u.data = data

	columns := make([]glib.Type, len(u.cols))
	for i, col := range u.cols {
		gtkc, err := u.createColumn(col.String(), i)
		if err != nil {
			return err
		}

		u.TreeView.InsertColumn(gtkc, i)
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

	u.mode = MODE_DEF
	return nil
}

type stringer string

func (s *stringer) String() string {
	return string(*s)
}

func stringSliceToStringerSlice(sc []string) (r []fmt.Stringer) {
	for _, str := range sc {
		st := stringer(str)
		r = append(r, &st)
	}

	return r
}

func colDefSliceToStringerSlice(sc []driver.ColDef) (r []fmt.Stringer) {
	for _, str := range sc {
		r = append(r, &str)
	}

	return r
}

func stringerSliceToColDefSlice(sc []fmt.Stringer) (r []driver.ColDef) {
	for _, str := range sc {
		col := str.(*driver.ColDef)
		r = append(r, *col)
	}

	return r
}

func (u *Result) UpdateRawData(cols []string, data [][]interface{}) error {
	for i := range u.cols {
		u.TreeView.RemoveColumn(u.TreeView.GetColumn(i))
	}

	u.cols = stringSliceToStringerSlice(cols)
	u.data = data

	columns := make([]glib.Type, len(u.cols))
	for i, col := range u.cols {
		gtkc, err := u.createColumn(col.String(), i)
		if err != nil {
			return err
		}

		u.TreeView.InsertColumn(gtkc, i)
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

	u.mode = MODE_RAW
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

func (u *Result) createColumn(title string, id int) (*gtk.TreeViewColumn, error) {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}
	cellRenderer.SetProperty("editable", true)
	cellRenderer.Connect("edited", u.onEdited, id)

	// i think "text" refers to a property of the column.
	// `"text", id` means that the text source for the column should come from
	// the listore column with id = `id`
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (u *Result) onEdited(cell *gtk.CellRendererText, path string, newValue string, userData interface{}) {
	if u.mode == MODE_RAW {
		return
	}

	column := userData.(int)
	row, _ := strconv.Atoi(path)

	tpath, err := gtk.TreePathNewFromString(path)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	iter, err := u.store.GetIter(tpath)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	err = u.store.SetValue(iter, column, newValue)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	oldRow := u.data[row]

	dataLen := u.store.GetNColumns()
	newRow := make([]interface{}, dataLen)
	for i := 0; i < dataLen; i++ {
		val, err := u.store.GetValue(iter, column)
		if err != nil {
			return
		}

		gVal, err := val.GoValue()
		if err != nil {
			return
		}

		switch reflect.TypeOf(gVal).Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			newRow[i] = gVal.(int64)
		case reflect.Bool:
			newRow[i] = gVal.(bool)
		default:
			newRow[i] = gVal.(string)
		}
	}

	for _, fn := range u.updateCallbacks {
		fn(stringerSliceToColDefSlice(u.cols), oldRow, newRow, newValue, row, column)
	}
}

func (u *Result) OnEdited(fn func([]driver.ColDef, []interface{}, []interface{}, string, int, int)) {
	u.updateCallbacks = append(u.updateCallbacks, fn)
}
