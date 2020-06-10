package gtk

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// ResultGrid is a table result tab content
type ResultGrid struct {
	*gtk.Paned

	textView *gtk.TextView
	prevText string
	offset   int64

	conditions *Conditions

	result         *Result
	btnPrev        *gtk.Button
	btnNext        *gtk.Button
	btnRsh         *gtk.Button
	btnShowFilters *gtk.Button
	perPage        *gtk.Entry
	//offset    *gtk.Entry
	pagerMenu *gtk.MenuButton

	btnAddRow    *gtk.Button
	btnDeleteRow *gtk.Button
	btnCreateRow *gtk.Button

	btnSaveQuery *gtk.Button
	btnLoadQuery *gtk.Button

	colFilter *gtk.SearchEntry

	submitCallback func(string)

	//query type
	isDML, isDDL bool
}

func (v ResultGrid) Init(
	w *Window,
	ctrl interface {
		OnUpdateRecord([]driver.ColDef, []interface{}) error
		OnCreateRecord([]driver.ColDef, []interface{}) ([]interface{}, error)
		OnExecQuery(string)
		OnRefresh()
		OnDelete()
		OnCreate()
		OnCopyInsert([]driver.ColDef, []interface{})
		OnFileSelected(string)
		OnSaveQuery(string, string)
	},
	parser parser,
) (*ResultGrid, error) {
	var err error

	v.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}

	v.textView, err = gtk.TextViewNew()
	if err != nil {
		return nil, err
	}
	v.textView.SetName("query")
	v.textView.SetHExpand(true)
	v.textView.SetVExpand(true)

	v.textView.Connect("key-release-event", v.onTextViewKeyRelease) // highlighting
	v.textView.Connect("key-press-event", v.onTextViewKeyPress)     // ctrl+enter exec query

	var resultSW, textViewSW *gtk.ScrolledWindow
	resultSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	actionbar, err := v.actionbar()
	if err != nil {
		return nil, err
	}

	v.conditions, err = Conditions{}.Init()
	if err != nil {
		return nil, err
	}

	v.btnAddRow.Connect("clicked", func() {
		v.result.AddEmptyRow()
	})

	v.btnNext.Connect("clicked", func() {
		v.offset = v.offset + v.PageSize()
	})
	v.btnPrev.Connect("clicked", func() {
		v.offset = v.offset - v.PageSize()
		if v.offset < 0 {
			v.offset = 0
		}
	})
	v.btnShowFilters.Connect("clicked", func() {
		v.conditions.Show()
	})
	v.colFilter.Connect("search-changed", v.onColFilterSearchChanged)

	btnGridBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	btnGridBox.PackStart(actionbar, false, false, 0)
	btnGridBox.PackStart(v.conditions, false, false, 5)
	btnGridBox.PackEnd(resultSW, true, true, 0)

	resultBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	resultBox.PackStart(btnGridBox, true, true, 0)

	tvBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	tvActionBar, err := gtk.ActionBarNew()
	if err != nil {
		return nil, err
	}
	v.btnSaveQuery, err = gtk.ButtonNewFromIconName("gtk-save", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	v.btnLoadQuery, err = gtk.ButtonNewFromIconName("gtk-open", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	v.btnLoadQuery.Connect("clicked", w.OnOpenFile(ctrl.OnFileSelected))
	v.btnSaveQuery.Connect("clicked", v.onSaveQuery(w.OnSaveQuery, ctrl.OnSaveQuery))

	tvActionBar.PackEnd(v.btnSaveQuery)
	tvActionBar.PackEnd(v.btnLoadQuery)
	tvActionBar.SetName("queryactionbar")

	textViewSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	//textViewSW.SetSizeRequest(-1, 200)

	v.result, err = Result{}.Init(ctrl, parser)
	if err != nil {
		return nil, err
	}

	v.result.TreeView.Connect("row-activated", v.onRowActivated)

	v.textView.SetProperty("accepts-tab", true)
	v.Paned.SetProperty("wide-handle", true)
	v.textView.SetLeftMargin(10)
	// this naming mess
	v.textView.SetProperty("top-margin", 10)

	v.textView.SetProperty("wrap-mode", map[string]gtk.WrapMode{
		"none":      gtk.WRAP_NONE,
		"char":      gtk.WRAP_CHAR,
		"word":      gtk.WRAP_WORD,
		"word_char": gtk.WRAP_WORD_CHAR,
	}[config.Env.GUI.Editor.WordWrap])

	resultSW.Add(v.result)
	textViewSW.Add(v.textView)

	tvBox.PackStart(textViewSW, false, true, 0)
	tvBox.PackEnd(tvActionBar, false, false, 0)
	v.Paned.Pack1(tvBox, false, false)
	v.Paned.Pack2(resultBox, true, false)

	v.disableAll()
	v.Paned.ShowAll()
	v.conditions.Hide()

	v.submitCallback = ctrl.OnExecQuery
	v.btnRsh.Connect("clicked", ctrl.OnRefresh)
	v.btnPrev.Connect("clicked", ctrl.OnRefresh)
	v.btnNext.Connect("clicked", ctrl.OnRefresh)
	v.btnDeleteRow.Connect("clicked", ctrl.OnDelete)
	v.btnCreateRow.Connect("clicked", ctrl.OnCreate)

	return &v, nil
}

func (v *ResultGrid) SetQuery(query string) {
	buff, err := v.textView.GetBuffer()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	iter := buff.GetIterAtMark(buff.GetInsert())
	offset := iter.GetOffset()

	query, err = ChromaHighlight(query)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
	buff.InsertMarkup(buff.GetStartIter(), query)

	buff.PlaceCursor(buff.GetIterAtOffset(offset))
}

func (v *ResultGrid) onSaveQuery(
	openDialog func(string, func(string, string)),
	onSaveQuery func(string, string),
) func() {
	return func() {
		buff, err := v.textView.GetBuffer()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		openDialog(txt, onSaveQuery)
	}
}

func (v *ResultGrid) PageSize() int64 {
	return 100
	//s, err := v.perPage.GetText()
	//if err != nil {
	//return int64(config.Env.GUI.PageSize)
	//}

	//size, err := strconv.ParseInt(s, 10, 64)
	//if err != nil {
	//return int64(config.Env.GUI.PageSize)
	//}

	//return size
}
func (v *ResultGrid) Offset() int64 {
	return v.offset
}

func (v *ResultGrid) Conditions() ([]driver.CondStmt, error) {
	return v.conditions.Statements()
}

func (v *ResultGrid) UpdateColumns(cols []driver.ColDef) error {
	v.colFilter.SetText("")
	v.offset = 0
	v.conditions.Update(cols)
	return v.result.UpdateColumns(cols)
}

func (v *ResultGrid) UpdateData(data [][]interface{}) error {
	v.pagerEnable(true)
	v.btnAddRow.SetSensitive(true)

	return v.result.UpdateData(data)
}

func (v *ResultGrid) UpdateRawData(cols []string, data [][]interface{}) error {
	v.pagerEnable(false)
	v.colFilter.SetText("")
	v.offset = 0
	return v.result.UpdateRawData(cols, data)
}

func (v *ResultGrid) SelectedIsNewRecord() (bool, error) {
	return v.result.SelectedIsNewRecord()
}

func (v *ResultGrid) RemoveSelected() error {
	err := v.result.RemoveSelected()
	if err != nil {
		return err
	}

	v.newRecordEnable(false)

	return nil
}

func (v *ResultGrid) GetRowID() ([]driver.ColDef, []interface{}, error) {
	return v.result.GetRowID()
}

func (u *ResultGrid) GetRow() ([]driver.ColDef, []interface{}, error) {
	return u.result.GetRow()
}

func (u *ResultGrid) UpdateRow(values []interface{}) error {
	err := u.result.UpdateRow(values)
	if err == nil {
		u.newRecordEnable(false)
	}
	return err
}

func (u *ResultGrid) SortOptions() []driver.SortOption {
	return u.result.SortOptions()
}

func (v *ResultGrid) actionbar() (*gtk.ActionBar, error) {
	actionbar, err := gtk.ActionBarNew()
	if err != nil {
		return nil, err
	}

	// new-add-delete
	{
		v.btnAddRow, err = gtk.ButtonNewFromIconName("gtk-add", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnDeleteRow, err = gtk.ButtonNewFromIconName("gtk-delete", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnCreateRow, err = gtk.ButtonNewFromIconName("gtk-apply", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnShowFilters, err = gtk.ButtonNewFromIconName("gtk-find", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.newRecordEnable(false)

		actionbar.Add(v.btnAddRow)
		actionbar.Add(v.btnDeleteRow)
		actionbar.Add(v.btnCreateRow)
		actionbar.Add(v.btnShowFilters)
	}

	// column filter
	{
		v.colFilter, err = gtk.SearchEntryNew()
		if err != nil {
			return nil, err
		}
		v.colFilter.SetPlaceholderText("Column filter: .*")
		actionbar.PackEnd(v.colFilter)
	}

	// menu
	//{
	//v.pagerMenu, err = gtk.MenuButtonNew()
	//if err != nil {
	//return nil, err
	//}
	//actionbar.PackEnd(v.pagerMenu)
	//}

	// pagination
	{
		v.btnPrev, err = gtk.ButtonNewFromIconName("gtk-go-back", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnNext, err = gtk.ButtonNewFromIconName("gtk-go-forward", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnRsh, err = gtk.ButtonNewFromIconName("gtk-refresh", gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		actionbar.PackEnd(v.btnRsh)
		actionbar.PackEnd(v.btnNext)
		actionbar.PackEnd(v.btnPrev)
	}

	return actionbar, nil
}

func (v *ResultGrid) pagerEnable(b bool) {
	v.btnPrev.SetSensitive(b)
	v.btnNext.SetSensitive(b)
	v.btnRsh.SetSensitive(b)
	//v.perPage.SetSensitive(b)
	//v.offset.SetSensitive(b)
	//v.perPage.SetSensitive(b)
	//v.offset.SetSensitive(b)
}

func (v *ResultGrid) disableAll() {
	v.btnPrev.SetSensitive(false)
	v.btnNext.SetSensitive(false)
	v.btnRsh.SetSensitive(false)
	v.btnAddRow.SetSensitive(false)
	v.btnDeleteRow.SetSensitive(false)
	v.btnCreateRow.SetSensitive(false)
}

func (v *ResultGrid) newRecordEnable(b bool) {
	v.btnCreateRow.SetSensitive(b)
}

func (v *ResultGrid) onTextViewKeyRelease(_ *gtk.TextView, e *gdk.Event) {
	buff, err := v.textView.GetBuffer()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	// hacky. easiest way to selection reset and bad behavior on non-printable key strokes
	if txt == v.prevText {
		return
	}
	v.prevText = txt

	iter := buff.GetIterAtMark(buff.GetInsert())
	offset := iter.GetOffset()

	txt, err = ChromaHighlight(txt)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
	buff.InsertMarkup(buff.GetStartIter(), txt)

	buff.PlaceCursor(buff.GetIterAtOffset(offset))
}

func (v *ResultGrid) onTextViewKeyPress(_ *gtk.TextView, e *gdk.Event) bool {
	keyEvent := gdk.EventKeyNewFromEvent(e)

	if keyEvent.KeyVal() == gdk.KEY_Return && keyEvent.State()&gdk.CONTROL_MASK > 0 {
		buff, err := v.textView.GetBuffer()
		if err != nil {
			config.Env.Log.Error(err)
			return true
		}
		txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			config.Env.Log.Error(err)
		}
		if v.submitCallback != nil {
			v.submitCallback(txt)
		}
		return true
	}

	return false
}

func (v *ResultGrid) onColFilterSearchChanged() {
	txt, err := v.colFilter.GetText()
	if err != nil {
		config.Env.Log.Error(err, "colFilter.GetText")
		return
	}

	rg, err := regexp.Compile(txt)
	if err != nil {
		rg = regexp.MustCompile(fmt.Sprintf(".*%s.*", regexp.QuoteMeta(txt)))
	}

	v.result.GetColumns().Foreach(func(i interface{}) {
		c := i.(*gtk.TreeViewColumn)

		c.SetVisible(rg.MatchString(strings.Replace(c.GetTitle(), "__", "_", -1)))
	})
}

func (v *ResultGrid) onRowActivated(_ *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	if v.result.mode == MODE_RAW {
		return
	}

	iter, err := v.result.store.GetIter(path)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	s, err := v.result.store.GetValue(iter, len(v.result.cols))
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	status, err := s.GoValue()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	v.newRecordEnable(status == STATUS_NEW)
	v.btnDeleteRow.SetSensitive(true)
}
