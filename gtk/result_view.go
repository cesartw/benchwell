package gtk

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/sourceview"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
)

// ResultView is a table result tab content
type ResultView struct {
	w *Window
	*CancelOverlay
	Paned *gtk.Paned

	//textView   *TextView
	sourceView *SourceView
	prevText   string
	offset     int64

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

	btnLoadQuery *gtk.Button
	btnSaveMenu  *gtk.MenuButton

	actionSaveFile *glib.SimpleAction
	actionSaveFav  *glib.SimpleAction

	colFilter *gtk.SearchEntry

	ctrl resultViewCtrl
}

type resultViewCtrl interface {
	OnUpdateRecord([]driver.ColDef, []interface{}) error
	OnCreateRecord([]driver.ColDef, []interface{}) ([]interface{}, error)
	OnExecQuery(string)
	OnTextChange(string, int) //query, cursor position
	OnRefresh()
	OnDelete()
	OnCreate()
	OnCopyInsert([]driver.ColDef, []interface{})
	OnFileSelected(string)
	OnSaveQuery(string, string)
	OnSaveFav(string, string)
	OnApplyConditions()
	Config() *config.Config
}

func (v ResultView) Init(
	w *Window,
	ctrl resultViewCtrl,
	parser parser,
) (*ResultView, error) {
	v.w = w
	v.ctrl = ctrl
	var err error

	v.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	v.CancelOverlay, err = CancelOverlay{}.Init(v.Paned)
	if err != nil {
		return nil, err
	}

	v.sourceView, err = SourceView{}.Init(v.w, SourceViewOptions{true, true, "sql"}, ctrl)
	if err != nil {
		return nil, err
	}
	v.sourceView.SetShowLineNumbers(true)
	v.sourceView.SetShowRightMargin()
	v.sourceView.SetHExpand(true)
	v.sourceView.SetVExpand(true)

	//v.textView, err = TextView{}.Init(v.w, TextViewOptions{true, true}, ctrl)
	//if err != nil {
	//return nil, err
	//}
	//v.textView.SetName("query")
	//v.textView.SetHExpand(true)
	//v.textView.SetVExpand(true)
	//v.textView.buffer.Connect("insert-text", func(_ *gtk.TextBuffer, iter *gtk.TextIter, txt string, _ int) {
	//if iter.GetOffset() == 0 {
	//return
	//}

	//cursorAt := iter.GetOffset() + len(txt)

	//start := v.textView.buffer.GetStartIter()
	//end := v.textView.buffer.GetEndIter()
	//query, err := v.textView.buffer.GetText(start, end, false)
	//if err != nil {
	//ctrl.Config().Error(err)
	//return
	//}

	//ctrl.OnTextChange(query+txt, cursorAt)
	//})
	buff, err := v.sourceView.GetBuffer()
	if err != nil {
		return nil, err
	}
	buff.Connect("insert-text", func(_ *sourceview.SourceBuffer, iter *gtk.TextIter, txt string, _ int) {
		if iter.GetOffset() == 0 {
			return
		}

		cursorAt := iter.GetOffset() + len(txt)

		start := buff.GetStartIter()
		end := buff.GetEndIter()
		query, err := buff.GetText(start, end, false)
		if err != nil {
			ctrl.Config().Error(err)
			return
		}

		ctrl.OnTextChange(query+txt, cursorAt)
	})

	v.sourceView.Connect("key-press-event", v.onTextViewKeyPress) // ctrl+enter exec query

	var resultSW, textViewSW *gtk.ScrolledWindow
	resultSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	actionbar, err := v.actionbar()
	if err != nil {
		return nil, err
	}

	v.conditions, err = Conditions{}.Init(v.w, ctrl)
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

	img, err := BWImageNewFromFile("save", ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	v.btnSaveMenu, err = gtk.MenuButtonNew()
	v.btnSaveMenu.SetImage(img)
	menu := glib.MenuNew()
	menu.Append("Save As", "win.save.file")
	menu.Append("Save fav", "win.save.fav")

	v.actionSaveFile = glib.SimpleActionNew("save.file", nil)
	v.actionSaveFav = glib.SimpleActionNew("save.fav", nil)

	v.w.AddAction(v.actionSaveFile)
	v.w.AddAction(v.actionSaveFav)

	v.actionSaveFile.Connect("activate", v.onSaveQuery(w.OnSaveQuery, ctrl.OnSaveQuery))
	v.actionSaveFav.Connect("activate", v.onSaveFav(ctrl.OnSaveFav))
	v.btnSaveMenu.SetMenuModel(&menu.MenuModel)

	v.btnLoadQuery, err = BWButtonNewFromIconName("open", ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	v.btnLoadQuery.Connect("clicked", w.OnOpenFile(ctrl.OnFileSelected))

	tvActionBar.PackEnd(v.btnSaveMenu)
	tvActionBar.PackEnd(v.btnLoadQuery)
	tvActionBar.SetName("queryactionbar")

	textViewSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	v.result, err = Result{}.Init(v.w, ctrl, parser)
	if err != nil {
		return nil, err
	}

	v.result.TreeView.Connect("row-activated", v.onRowActivated)

	v.sourceView.SetProperty("highlight-current-line", false)
	v.sourceView.SetProperty("show-line-numbers", true)
	v.sourceView.SetProperty("show-right-margin", true)
	v.sourceView.SetProperty("show-left-margin", true)
	//v.sourceView.SetProperty("top-margin", 10)

	//v.sourceView.SetProperty("wrap-mode", map[string]gtk.WrapMode{
	//"none":      gtk.WRAP_NONE,
	//"char":      gtk.WRAP_CHAR,
	//"word":      gtk.WRAP_WORD,
	//"word_char": gtk.WRAP_WORD_CHAR,
	//}[v.ctrl.Config().Editor.WordWrap.String()])

	v.Paned.SetProperty("wide-handle", true)

	resultSW.Add(v.result)
	textViewSW.Add(v.sourceView)

	tvBox.PackStart(textViewSW, false, true, 0)
	tvBox.PackEnd(tvActionBar, false, false, 0)
	v.Paned.Pack1(tvBox, false, false)
	v.Paned.Pack2(resultBox, true, false)

	v.disableAll()
	v.Paned.ShowAll()
	v.conditions.Hide()

	v.btnRsh.Connect("clicked", ctrl.OnRefresh)
	v.btnPrev.Connect("clicked", ctrl.OnRefresh)
	v.btnNext.Connect("clicked", ctrl.OnRefresh)
	v.btnDeleteRow.Connect("clicked", ctrl.OnDelete)
	v.btnCreateRow.Connect("clicked", ctrl.OnCreate)

	return &v, nil
}

func (v *ResultView) ShowAutoComplete(words []string) {
	v.sourceView.ShowAutoComplete(words)
}

func (v *ResultView) Block(cancel func()) {
	v.Run(cancel)
}

func (v *ResultView) SetQuery(query string) {
	buff, err := v.sourceView.GetBuffer()
	if err != nil {
		v.ctrl.Config().Error(err)
		return
	}

	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
	buff.InsertMarkup(buff.GetStartIter(), query)
}

func (v *ResultView) onSaveQuery(
	openDialog func(string, func(string, string)),
	onSaveQuery func(string, string),
) func() {
	return func() {
		buff, err := v.sourceView.GetBuffer()
		if err != nil {
			v.ctrl.Config().Error(err)
			return
		}

		txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			v.ctrl.Config().Error(err)
			return
		}

		openDialog(txt, onSaveQuery)
	}
}

func (v *ResultView) onSaveFav(
	onSaveQuery func(string, string),
) func() {
	return func() {
		name, err := v.askFavName()
		if err != nil {
			v.ctrl.Config().Error(err)
			return
		}

		if name == "" {
			return
		}

		buff, err := v.sourceView.GetBuffer()
		if err != nil {
			v.ctrl.Config().Error(err)
			return
		}

		query, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			v.ctrl.Config().Error(err)
			return
		}

		onSaveQuery(name, query)
	}
}

func (v *ResultView) PageSize() int64 {
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
func (v *ResultView) Offset() int64 {
	return v.offset
}

func (v *ResultView) Conditions() ([]driver.CondStmt, error) {
	return v.conditions.Statements()
}

func (v *ResultView) UpdateColumns(cols []driver.ColDef) error {
	defer v.Stop()

	v.colFilter.SetText("")
	v.offset = 0
	v.conditions.Update(cols)
	return v.result.UpdateColumns(cols)
}

func (v *ResultView) UpdateData(data [][]interface{}) error {
	defer v.Stop()

	v.pagerEnable(true)
	v.btnAddRow.SetSensitive(true)

	return v.result.UpdateData(data)
}

func (v *ResultView) UpdateRawData(cols []string, data [][]interface{}) error {
	defer v.Stop()

	v.pagerEnable(false)
	v.colFilter.SetText("")
	v.offset = 0
	return v.result.UpdateRawData(cols, data)
}

func (v *ResultView) SelectedIsNewRecord() (bool, error) {
	return v.result.SelectedIsNewRecord()
}

func (v *ResultView) RemoveSelected() error {
	err := v.result.RemoveSelected()
	if err != nil {
		return err
	}

	v.newRecordEnable(false)

	return nil
}

func (v *ResultView) ForEachSelected(f func([]driver.ColDef, []interface{})) error {
	return v.result.ForEachSelected(f)
}

func (v *ResultView) GetRowID() ([]driver.ColDef, []interface{}, error) {
	return v.result.GetRowID()
}

func (u *ResultView) GetRow() ([]driver.ColDef, []interface{}, error) {
	return u.result.GetRow()
}

func (u *ResultView) UpdateRow(values []interface{}) error {
	err := u.result.UpdateRow(values)
	if err == nil {
		u.newRecordEnable(false)
	}
	return err
}

func (u *ResultView) SortOptions() []driver.SortOption {
	return u.result.SortOptions()
}

func (v *ResultView) actionbar() (*gtk.ActionBar, error) {
	actionbar, err := gtk.ActionBarNew()
	if err != nil {
		return nil, err
	}

	// new-add-delete
	{
		v.btnAddRow, err = BWButtonNewFromIconName("add-record", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnDeleteRow, err = BWButtonNewFromIconName("delete-record", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnCreateRow, err = BWButtonNewFromIconName("save-record", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnShowFilters, err = BWButtonNewFromIconName("filter", ICON_SIZE_BUTTON)
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
		v.btnPrev, err = BWButtonNewFromIconName("back", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnNext, err = BWButtonNewFromIconName("next", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		v.btnRsh, err = BWButtonNewFromIconName("refresh", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}

		actionbar.PackEnd(v.btnRsh)
		actionbar.PackEnd(v.btnNext)
		actionbar.PackEnd(v.btnPrev)
	}

	return actionbar, nil
}

func (v *ResultView) pagerEnable(b bool) {
	v.btnPrev.SetSensitive(b)
	v.btnNext.SetSensitive(b)
	v.btnRsh.SetSensitive(b)
	//v.perPage.SetSensitive(b)
	//v.offset.SetSensitive(b)
	//v.perPage.SetSensitive(b)
	//v.offset.SetSensitive(b)
}

func (v *ResultView) disableAll() {
	v.btnPrev.SetSensitive(false)
	v.btnNext.SetSensitive(false)
	v.btnRsh.SetSensitive(false)
	v.btnAddRow.SetSensitive(false)
	v.btnDeleteRow.SetSensitive(false)
	v.btnCreateRow.SetSensitive(false)
}

func (v *ResultView) newRecordEnable(b bool) {
	v.btnCreateRow.SetSensitive(b)
}

func (v *ResultView) onTextViewKeyPress(_ *sourceview.SourceView, e *gdk.Event) bool {
	keyEvent := gdk.EventKeyNewFromEvent(e)

	if keyEvent.KeyVal() == gdk.KEY_Return && keyEvent.State()&gdk.CONTROL_MASK > 0 {
		buff, err := v.sourceView.GetBuffer()
		if err != nil {
			v.ctrl.Config().Error(err)
			return true
		}
		txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			v.ctrl.Config().Error(err)
		}

		v.ctrl.OnExecQuery(txt)
		return true
	}

	return false
}

func (v *ResultView) onColFilterSearchChanged() {
	txt, err := v.colFilter.GetText()
	if err != nil {
		v.ctrl.Config().Error(err, "colFilter.GetText")
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

func (v *ResultView) onRowActivated(_ *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	if v.result.mode == MODE_RAW {
		return
	}

	iter, err := v.result.store.GetIter(path)
	if err != nil {
		v.ctrl.Config().Error(err)
		return
	}

	s, err := v.result.store.GetValue(iter, len(v.result.cols))
	if err != nil {
		v.ctrl.Config().Error(err)
		return
	}
	status, err := s.GoValue()
	if err != nil {
		v.ctrl.Config().Error(err)
		return
	}

	v.newRecordEnable(status == STATUS_NEW)
	v.btnDeleteRow.SetSensitive(true)
}

func (v *ResultView) askFavName() (string, error) {
	modal, err := gtk.DialogNewWithButtons(
		"Favorite Name",
		v.w,
		gtk.DIALOG_DESTROY_WITH_PARENT|gtk.DIALOG_MODAL,
		[]interface{}{"Ok", gtk.RESPONSE_ACCEPT},
		[]interface{}{"Cancel", gtk.RESPONSE_CANCEL},
	)
	if err != nil {
		return "", err
	}
	modal.SetDefaultSize(250, 130)
	content, err := modal.GetContentArea()
	if err != nil {
		return "", err
	}

	label, err := gtk.LabelNew("Enter favorite name")
	if err != nil {
		return "", err
	}

	entry, err := gtk.EntryNew()
	if err != nil {
		return "", err
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return "", err
	}

	box.PackStart(label, true, true, 0)
	box.PackStart(entry, true, true, 0)
	content.Add(box)
	content.ShowAll()

	defer modal.Destroy()
	resp := modal.Run()
	if resp != gtk.RESPONSE_ACCEPT {
		return "", nil
	}

	name, err := entry.GetText()
	if err != nil {
		return "", err
	}

	return name, nil
}
