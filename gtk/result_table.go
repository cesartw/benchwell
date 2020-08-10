package gtk

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
)

var once = sync.Once{}

type MODE int

// Grid modes
const (
	MODE_RAW = iota
	MODE_DEF
)

// Row status
const (
	STATUS_NEW = iota
	STATUS_CHANGED
	STATUS_PRISTINE
)

type parser func(driver.ColDef, string) (interface{}, error)
type Result struct {
	*gtk.TreeView
	cols           []fmt.Stringer
	data           [][]interface{}
	store          *gtk.ListStore
	updateCallback func([]driver.ColDef, []interface{}) error
	createCallback func([]driver.ColDef, []interface{}) ([]interface{}, error)

	mode   MODE
	parser parser

	ddMenu struct {
		*gtk.Menu
		clone    *gtk.MenuItem
		cpInsert *gtk.MenuItem
		cp       *gtk.MenuItem
	}

	rowAtCursor *gtk.TreePath
	colAtCursor *gtk.TreeViewColumn

	onCopyInsertFn func([]driver.ColDef, []interface{})
	ctrl           resultCtrl
}

type resultCtrl interface {
	OnUpdateRecord([]driver.ColDef, []interface{}) error
	OnCreateRecord([]driver.ColDef, []interface{}) ([]interface{}, error)
	OnCopyInsert([]driver.ColDef, []interface{})
	Config() *config.Config
}

func (u Result) Init(_ *Window, ctrl resultCtrl, parser parser) (*Result, error) {
	var err error
	u.parser = parser
	u.ctrl = ctrl

	u.TreeView, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	sel, err := u.TreeView.GetSelection()
	if err != nil {
		return nil, err
	}
	sel.SetMode(gtk.SELECTION_MULTIPLE)

	u.TreeView.SetProperty("rubber-banding", true)
	u.TreeView.SetProperty("enable-grid-lines", gtk.TREE_VIEW_GRID_LINES_HORIZONTAL)
	u.TreeView.SetProperty("activate-on-single-click", true)

	u.TreeView.SetEnableSearch(true)

	u.TreeView.Connect("key-press-event", u.onTreeViewKeyPress)
	u.TreeView.Connect("button-press-event", u.onTreeViewButtonPress)

	//if len(cols) > 0 {
	//u.UpdateColumns(cols)
	//u.UpdateData(data)
	//}

	u.ddMenu.Menu, err = gtk.MenuNew()
	if err != nil {
		return nil, err
	}

	u.ddMenu.clone, err = BWMenuItemWithImage("Clone Row", "gtk-convert")
	if err != nil {
		return nil, err
	}
	u.ddMenu.clone.Connect("activate", u.onCloneRow)
	u.ddMenu.Add(u.ddMenu.clone)

	u.ddMenu.cpInsert, err = BWMenuItemWithImage("Copy Insert", "gtk-page-setup")
	if err != nil {
		return nil, err
	}
	u.ddMenu.cpInsert.Connect("activate", u.onCopyInsert)
	u.ddMenu.Add(u.ddMenu.cpInsert)

	u.ddMenu.cp, err = BWMenuItemWithImage("Copy Field", "gtk-copy")
	if err != nil {
		return nil, err
	}
	u.ddMenu.cp.Connect("activate", u.onCopy)
	u.ddMenu.Add(u.ddMenu.cp)

	u.updateCallback = ctrl.OnUpdateRecord
	u.createCallback = ctrl.OnCreateRecord
	u.onCopyInsertFn = ctrl.OnCopyInsert

	return &u, nil
}

func (u *Result) UpdateColumns(cols []driver.ColDef) error {
	u.mode = MODE_DEF
	// columns shift to the left
	for _ = range u.cols {
		u.TreeView.RemoveColumn(u.TreeView.GetColumn(0))
	}

	u.cols = colDefSliceToStringerSlice(cols)
	u.data = nil

	columns := make([]glib.Type, len(u.cols)+1) // +1 internal status col
	for i, col := range u.cols {
		c, err := u.createColumn(col.String(), i, u.cols[i].(driver.ColDef).Type == driver.ColTypeLongString)
		if err != nil {
			return err
		}

		u.TreeView.InsertColumn(c, i)
		// default type
		columns[i] = glib.TYPE_STRING
	}

	columns[len(u.cols)] = glib.TYPE_INT

	var err error
	u.store, err = gtk.ListStoreNew(columns...)
	if err != nil {
		return err
	}
	u.TreeView.SetModel(u.store)

	return nil
}

func (u *Result) UpdateData(data [][]interface{}) error {
	u.mode = MODE_DEF
	u.data = data
	u.store.Clear()

	for _, row := range data {
		u.AddRow(row)
	}

	return nil
}

func (u *Result) UpdateRawData(cols []string, data [][]interface{}) error {
	// columns shift to the left
	u.mode = MODE_RAW
	for _ = range u.cols {
		u.TreeView.RemoveColumn(u.TreeView.GetColumn(0))
	}

	u.cols = stringSliceToStringerSlice(cols)
	u.data = data

	columns := make([]glib.Type, len(u.cols))
	for i, col := range u.cols {
		gtkc, err := u.createColumn(col.String(), i, false)
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

	return nil
}

func (u *Result) AddEmptyRow() (err error) {
	if u.mode == MODE_RAW {
		return nil
	}

	var iter *gtk.TreeIter

	sel, err := u.TreeView.GetSelection()
	if err != nil {
		return err
	}

	var storeSelected *gtk.TreeIter
	sel.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, _ *gtk.TreeIter, userData ...interface{}) {
		if storeSelected != nil {
			return
		}

		storeSelected, err = u.store.GetIter(path)
		if err != nil {
			return
		}
	})

	if storeSelected != nil {
		iter = u.store.InsertAfter(storeSelected)
	}

	if iter == nil {
		iter = u.store.Prepend()
	}

	p, err := u.store.GetPath(iter)
	if err != nil {
		return err
	}
	sel.UnselectAll()
	sel.SelectPath(p)

	i, _ := strconv.Atoi(p.String())

	row := make([]interface{}, len(u.cols))
	data2 := u.data[:i]
	data2 = append(data2, row)
	data2 = append(data2, u.data[i:]...)
	u.data = data2

	u.TreeView.RowActivated(p, nil)

	// vertically center scroll at new row
	u.TreeView.ScrollToCell(p, nil, true, 0.5, 0)

	columns := make([]int, len(row))
	for i, col := range u.cols {
		def := col.(driver.ColDef)
		columns[i] = i

		row[i], err = u.parser(def, driver.NULL_PATTERN)
		if err != nil || row[i] == nil {
			row[i] = driver.NULL_PATTERN
		}
	}

	// Set the contents of the list store row that the iterator represents
	err = u.store.Set(iter,
		columns,
		row,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *Result) AddRow(originalRow []interface{}) {
	// Get an iterator for a new row at the end of the list store
	iter := u.store.Append()

	if len(originalRow) != len(u.cols) {
		log.Fatal("wrong row length")
	}

	row := make([]interface{}, len(u.cols))

	columns := make([]int, len(originalRow))
	for i, d := range originalRow {
		row[i] = originalRow[i]
		columns[i] = i

		if d == nil {
			row[i] = "<NULL>"
			continue
		}

		if data, ok := d.(string); ok && len(data) > 1000 {
			s := make([]byte, 255)
			copy(s, []byte(d.(string)))
			row[i] = strings.TrimSpace(strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(string(s), "\t"), "\n"), "\r"))
		}

		if s, ok := d.(int64); ok {
			row[i] = s
		}
	}

	if u.mode == MODE_DEF {
		columns = append(columns, len(row))
		row = append(row, STATUS_PRISTINE)
	}

	// Set the contents of the list store row that the iterator represents
	err := u.store.Set(iter, columns, row)
	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}

func (u *Result) GetCurrentIter() (*gtk.TreeIter, error) {
	var storeSelected *gtk.TreeIter

	sel, err := u.TreeView.GetSelection()
	if err != nil {
		return nil, err
	}

	sel.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, _ *gtk.TreeIter, userData ...interface{}) {
		if storeSelected != nil {
			return
		}

		storeSelected, err = u.store.GetIter(path)
		if err != nil {
			return
		}
	})

	return storeSelected, nil
}

func (u *Result) SelectedIsNewRecord() (bool, error) {
	iter, err := u.GetCurrentIter()
	if err != nil {
		return false, err
	}

	if iter == nil {
		return false, nil
	}

	lastColValue, err := u.store.GetValue(iter, len(u.cols))
	if err != nil {
		return false, err
	}
	status, err := lastColValue.GoValue()
	if err != nil {
		return false, err
	}
	return status.(int) == STATUS_NEW, nil
}

func (u *Result) RemoveSelected() error {
	iter, err := u.GetCurrentIter()
	if err != nil {
		return err
	}

	if iter == nil {
		return nil
	}

	u.store.Remove(iter)
	path, err := u.store.GetPath(iter)
	if err != nil {
		return err
	}
	index, _ := strconv.Atoi(path.String())

	u.data = append(u.data[:index], u.data[index+1:]...)
	return nil
}

func (u *Result) GetRowID() ([]driver.ColDef, []interface{}, error) {
	iter, err := u.GetCurrentIter()
	if err != nil {
		return nil, nil, err

	}
	if iter == nil {
		return nil, nil, nil
	}

	path, err := u.store.GetPath(iter)
	if err != nil {
		return nil, nil, err
	}

	row, _ := strconv.Atoi(path.String())

	pkCols := []driver.ColDef{}
	values := []interface{}{}
	for i, col := range u.cols {
		def := col.(driver.ColDef)
		if !def.PK {
			continue
		}

		pkCols = append(pkCols, def)
		values = append(values, u.data[row][i])
	}

	return pkCols, values, nil
}

func (u *Result) ForEachSelected(f func([]driver.ColDef, []interface{})) error {
	selection, err := u.GetSelection()
	if err != nil {
		return err
	}
	// TODO: maybe there's a memory leak here.
	//       using FreeFull causes a 'double free or corruption' error
	selection.GetSelectedRows(u.store).Foreach(func(i interface{}) {
		path := i.(*gtk.TreePath)
		row, _ := strconv.Atoi(path.String())

		pkCols := []driver.ColDef{}
		values := []interface{}{}
		for i, col := range u.cols {
			def := col.(driver.ColDef)
			if !def.PK {
				continue
			}

			pkCols = append(pkCols, def)
			values = append(values, u.data[row][i])
		}

		f(pkCols, values)
	})

	return nil
}

func (u *Result) GetRow() ([]driver.ColDef, []interface{}, error) {
	iter, err := u.GetCurrentIter()
	if err != nil {
		return nil, nil, err

	}
	if iter == nil {
		return nil, nil, nil
	}

	values := make([]interface{}, len(u.cols))
	for i := range u.cols {
		v, err := u.store.GetValue(iter, i)
		if err != nil {
			return nil, nil, err
		}

		n, err := v.GoValue()
		if err != nil {
			return nil, nil, err
		}

		values[i] = n
	}

	return stringerSliceToColDefSlice(u.cols), values, nil
}

func (u *Result) UpdateRow(values []interface{}) error {
	iter, err := u.GetCurrentIter()
	if err != nil {
		return err
	}

	columns := make([]int, len(u.cols))
	for i := range u.cols {
		columns[i] = i
	}

	path, err := u.store.GetPath(iter)
	if err != nil {
		return err
	}
	i, _ := strconv.Atoi(path.String())

	u.data[i] = values

	u.store.SetValue(iter, len(u.cols), STATUS_PRISTINE)

	return u.store.Set(iter,
		columns,
		values)
}

func (u *Result) SortOptions() []driver.SortOption {
	if u.mode == MODE_RAW {
		return nil
	}

	opts := []driver.SortOption{}

	for i, col := range u.cols {
		treeCol := u.GetColumn(i)
		if !treeCol.GetSortIndicator() {
			continue
		}

		if treeCol.GetSortOrder() == gtk.SORT_ASCENDING {
			opts = append(opts, driver.SortOption{Column: col.(driver.ColDef), Direction: driver.SortDirectionAsc})
		} else {
			opts = append(opts, driver.SortOption{Column: col.(driver.ColDef), Direction: driver.SortDirectionDesc})
		}
	}

	return opts
}

func (u *Result) onCopyInsert() {
	if u.mode == MODE_RAW {
		return
	}

	cols, values, err := u.GetRow()
	if err != nil {
		return
	}

	u.onCopyInsertFn(cols, values)
}

func (u *Result) onEdited(_ *gtk.CellRendererText, rowPath string, newValue string, userData interface{}) {
	u.ctrl.Config().Debug("cell edited")
	if u.mode == MODE_RAW {
		return
	}

	row, _ := strconv.Atoi(rowPath)
	u.onSaveCell(row, userData.(int), newValue)
}

func (u *Result) onSaveCell(row, column int, newValue string) {
	tpath, err := gtk.TreePathNewFromString(fmt.Sprintf("%d", row))
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}

	iter, err := u.store.GetIter(tpath)
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}

	err = u.store.SetValue(iter, column, newValue)
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}

	// is a new record
	lastColValue, err := u.store.GetValue(iter, len(u.cols))
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}
	status, err := lastColValue.GoValue()
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}
	if status.(int) == STATUS_NEW {
		return
	}
	/////////////

	pkCols := []driver.ColDef{}
	values := []interface{}{}
	for i, col := range u.cols {
		def := col.(driver.ColDef)
		if !def.PK {
			continue
		}

		pkCols = append(pkCols, def)
		values = append(values, u.data[row][i])
	}

	// PK-LESS rows
	if len(pkCols) == 0 {
		for _, col := range u.cols {
			pkCols = append(pkCols, col.(driver.ColDef))
		}
	}

	affectedCol := u.cols[column].(driver.ColDef)
	pkCols = append(pkCols, affectedCol)
	parsedValue, err := u.parser(affectedCol, newValue)
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}

	values = append(values, parsedValue)

	err = u.updateCallback(pkCols, values)
	if err != nil {
		u.ctrl.Config().Error(err)
		return
	}

	u.data[row][column], _ = u.parser(affectedCol, newValue)
}

func (u *Result) createColumn(title string, id int, useEditModal bool) (*gtk.TreeViewColumn, error) {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}
	cellRenderer.SetProperty("editable", true)
	cellRenderer.SetProperty("xpad", 10)
	cellRenderer.SetProperty("height", 23)
	cellRenderer.SetProperty("width-chars", u.ctrl.Config().GUI.CellWidth.String())
	cellRenderer.SetProperty("max-width-chars", u.ctrl.Config().GUI.CellWidth.String())
	cellRenderer.SetProperty("wrap-width", u.ctrl.Config().GUI.CellWidth.String())

	// i think "text" refers to a property of the column.
	// `"text", id` means that the text source for the column should come from
	// the listore column with id = `id`
	// NOTE: single _ is not display, maybe it's an issue with my system
	column, err := gtk.TreeViewColumnNewWithAttribute(underscoreFix(title), cellRenderer, "text", id)
	if err != nil {
		return nil, err
	}
	column.SetResizable(true)
	// TODO: this limits resizing
	//column.SetMaxWidth(300)

	if u.mode != MODE_DEF {
		return column, nil
	}

	cellRenderer.Connect("edited", func(cell *gtk.CellRendererText, path string, newValue string, userData interface{}) {
		if !useEditModal {
			u.onEdited(cell, path, newValue, userData)
		}
	}, id)

	cellRenderer.Connect("editing-started", func() {
		if !useEditModal {
			return
		}

		u.editDialog(u.dataAtCursor(), func(newValue string) {
			row, _ := strconv.Atoi(u.rowAtCursor.String())
			var at int
			for i, col := range u.cols {
				if col.String() == underscoreUnfix(u.colAtCursor.GetTitle()) {
					at = i
					break
				}
			}

			u.onSaveCell(row, at, newValue)
		})
	})

	column.SetClickable(true)
	column.Connect("clicked", u.onSort(column))

	return column, nil
}

func (u *Result) onSort(column *gtk.TreeViewColumn) func() {
	return func() {
		if !column.GetSortIndicator() {
			column.SetSortIndicator(true)
			column.SetSortOrder(gtk.SORT_ASCENDING)
			return
		}

		if column.GetSortOrder() == gtk.SORT_ASCENDING {
			column.SetSortOrder(gtk.SORT_DESCENDING)
		} else {
			column.SetSortIndicator(false)
		}
	}
}

func (u *Result) editDialog(data string, done func(string)) {
	modal, err := gtk.DialogNewWithButtons("Edit", nil,
		gtk.DIALOG_DESTROY_WITH_PARENT|gtk.DIALOG_MODAL,
		[]interface{}{"Save", gtk.RESPONSE_ACCEPT},
		[]interface{}{"Cancel", gtk.RESPONSE_CANCEL},
	)
	if err != nil {
		return
	}

	modal.SetDefaultSize(400, 400)
	content, err := modal.GetContentArea()
	if err != nil {
		return
	}

	sw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return
	}
	sw.SetVExpand(true)
	sw.SetHExpand(true)

	textView, err := gtk.TextViewNew()
	if err != nil {
		return
	}
	textView.SetVExpand(true)
	textView.SetHExpand(true)

	buff, err := textView.GetBuffer()
	if err != nil {
		return
	}

	buff.Insert(buff.GetStartIter(), data)

	sw.Add(textView)
	content.Add(sw)
	content.ShowAll()

	resp := modal.Run()
	defer modal.Destroy()

	if resp != gtk.RESPONSE_ACCEPT {
		return
	}

	data, err = buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
	if err != nil {
		u.ctrl.Config().Error(err)
	}

	done(data)
}

func (u *Result) onTreeViewButtonPress(_ *gtk.TreeView, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	path, col, _, _, ok := u.TreeView.GetPathAtPos(int(keyEvent.X()), int(keyEvent.Y()))
	if !ok {
		return
	}

	u.rowAtCursor = path
	u.colAtCursor = col

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}
	u.ddMenu.ShowAll()
	u.ddMenu.PopupAtPointer(e)
}

func (u *Result) onTreeViewKeyPress(_ *gtk.TreeView, e *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(e)
	if keyEvent.KeyVal() == gdk.KEY_F2 {
		path, col := u.TreeView.GetCursor()
		u.TreeView.SetCursor(path, col, true)
	}
}

func (u *Result) onCloneRow() {
	cols, data, err := u.GetRow()
	if err != nil {
		u.ctrl.Config().Error("getting current row", err)
		return
	}

	err = u.AddEmptyRow()
	if err != nil {
		u.ctrl.Config().Error("adding empty row", err)
		return
	}

	iter, err := u.GetCurrentIter()
	if err != nil {
		u.ctrl.Config().Errorf("getting new row iter: %s", err)
		return

	}
	if iter == nil {
		u.ctrl.Config().Debug("no row selected")
		return
	}

	for i := range data {
		if cols[i].PK {
			continue
		}

		err = u.store.SetValue(iter, i, data[i])
		if err != nil {
			u.ctrl.Config().Errorf("setting value: %s", err)
			return
		}
	}
}

func (u *Result) onCopy() {
	ClipboardCopy(u.dataAtCursor())
}

func (u *Result) dataAtCursor() string {
	var at int
	for i, col := range u.cols {
		if col.String() == underscoreUnfix(u.colAtCursor.GetTitle()) {
			at = i
			break
		}
	}

	// NOTE: I think FreeFull is buggy
	//u.TreeView.GetColumns().FreeFull(func(c interface{}) {
	//if c.(*gtk.TreeViewColumn).GetTitle() == u.colAtCursor.GetTitle() {
	//at = colIndex
	//}
	//colIndex++
	//})

	row, _ := strconv.Atoi(u.rowAtCursor.String())

	return string(u.data[row][at].(string))
}

func underscoreFix(s string) string {
	return strings.Replace(s, "_", "__", -1)
}

func underscoreUnfix(s string) string {
	return strings.Replace(s, "__", "_", -1)
}
