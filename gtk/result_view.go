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
	ParseValue(driver.ColDef, string) (interface{}, error)
}

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
	btnShowFilters *gtk.ToggleButton
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

func (v ResultView) Init(
	w *Window,
	ctrl resultViewCtrl,
) (*ResultView, error) {
	defer config.LogStart("ResultView.Init", nil)()

	v.w = w
	v.ctrl = ctrl
	var err error

	v.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	v.Paned.Show()
	v.CancelOverlay, err = CancelOverlay{}.Init(v.Paned)
	if err != nil {
		return nil, err
	}

	v.sourceView, err = SourceView{}.Init(v.w, SourceViewOptions{true, true, "sql"}, ctrl)
	if err != nil {
		return nil, err
	}
	v.sourceView.Show()
	v.sourceView.SetShowLineNumbers(false)
	v.sourceView.SetShowRightMargin()
	v.sourceView.SetHExpand(true)
	v.sourceView.SetVExpand(true)

	buff, err := v.sourceView.GetBuffer()
	if err != nil {
		return nil, err
	}
	buff.Connect("insert-text", func(_ *sourceview.SourceBuffer, iter *gtk.TextIter, txt string, _ int) {
		// TODO: attemp to autocomplete
		if iter.GetOffset() == 0 {
			return
		}

		cursorAt := iter.GetOffset() + len(txt)

		start := buff.GetStartIter()
		end := buff.GetEndIter()
		query, err := buff.GetText(start, end, false)
		if err != nil {
			config.Error(err)
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
	resultSW.Show()

	actionbar, err := v.actionbar()
	if err != nil {
		return nil, err
	}
	actionbar.Show()

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
	v.btnShowFilters.Connect("toggled", func() {
		if v.btnShowFilters.GetActive() {
			v.conditions.Show()
		} else {
			v.conditions.Hide()
		}
	})
	v.colFilter.Connect("search-changed", v.onColFilterSearchChanged)

	btnGridBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	btnGridBox.Show()
	btnGridBox.PackStart(actionbar, false, false, 0)
	btnGridBox.PackStart(v.conditions, false, false, 5)
	btnGridBox.PackEnd(resultSW, true, true, 0)

	resultBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	resultBox.Show()
	resultBox.PackStart(btnGridBox, true, true, 0)

	tvActionBar, err := gtk.ActionBarNew()
	if err != nil {
		return nil, err
	}
	tvActionBar.Show()

	img, err := BWImageNewFromFile("save", "orange", ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	img.Show()

	v.btnSaveMenu, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	v.btnSaveMenu.Show()
	v.btnSaveMenu.SetImage(img)

	menu := glib.MenuNew()
	menu.Append("Save As", "win.save.file")
	menu.Append("Save fav", "win.save.fav")

	v.actionSaveFile = glib.SimpleActionNew("save.file", nil)
	v.actionSaveFav = glib.SimpleActionNew("save.fav", nil)

	v.w.AddAction(v.actionSaveFile)
	v.w.AddAction(v.actionSaveFav)

	v.actionSaveFile.Connect("activate", v.onSaveQuery)
	v.actionSaveFav.Connect("activate", v.onSaveFav)
	v.btnSaveMenu.SetMenuModel(&menu.MenuModel)

	v.btnLoadQuery, err = BWButtonNewFromIconName("open", "orange", ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	v.btnLoadQuery.Show()
	v.btnLoadQuery.Connect("clicked", v.onOpenFile)

	tvActionBar.PackEnd(v.btnSaveMenu)
	tvActionBar.PackEnd(v.btnLoadQuery)
	tvActionBar.SetName("queryactionbar")

	textViewSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	textViewSW.Show()

	v.result, err = Result{}.Init(v.w, ctrl)
	if err != nil {
		return nil, err
	}
	v.result.Show()
	v.result.TreeView.Connect("row-activated", v.onRowActivated)

	v.sourceView.SetProperty("highlight-current-line", false)
	v.sourceView.SetProperty("show-line-numbers", true)
	v.sourceView.SetProperty("show-right-margin", true)
	//v.sourceView.SetProperty("show-left-margin", true)
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

	tvBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	tvBox.Show()

	tvBox.PackStart(textViewSW, false, true, 0)
	tvBox.PackEnd(tvActionBar, false, false, 0)
	v.Paned.Pack1(tvBox, false, false)
	v.Paned.Pack2(resultBox, true, false)

	v.disableAll()
	v.Paned.Show()
	v.conditions.Hide()

	v.btnRsh.Connect("clicked", ctrl.OnRefresh)
	v.btnPrev.Connect("clicked", ctrl.OnRefresh)
	v.btnNext.Connect("clicked", ctrl.OnRefresh)
	v.btnDeleteRow.Connect("clicked", ctrl.OnDelete)
	v.btnCreateRow.Connect("clicked", ctrl.OnCreate)

	return &v, nil
}

func (v *ResultView) onOpenFile() {
	defer config.LogStart("ResultView.onOpenFile", nil)()

	openfileDialog, err := gtk.FileChooserDialogNewWith2Buttons("Select file", v.w, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Open", gtk.RESPONSE_OK,
		"Cancel", gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		config.Error("open file dialog", err)
		return
	}
	defer openfileDialog.Destroy()

	response := openfileDialog.Run()
	if response == gtk.RESPONSE_OK && openfileDialog.GetFilename() != "" {
		v.ctrl.OnFileSelected(openfileDialog.GetFilename())
	}
}

func (v *ResultView) ShowAutoComplete(words []string) {
	defer config.LogStart("ResultView.ShowAutoComplete", nil)()

	v.sourceView.ShowAutoComplete(words)
}

func (v *ResultView) Block(cancel func()) {
	defer config.LogStart("ResultView.Block", nil)()

	v.Run(cancel)
}

func (v *ResultView) SetQuery(query string) {
	defer config.LogStart("ResultView.SetQuery", nil)()

	buff, err := v.sourceView.GetBuffer()
	if err != nil {
		config.Error(err)
		return
	}

	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
	buff.Insert(buff.GetStartIter(), query)
}

func (v *ResultView) onSaveQuery() {
	defer config.LogStart("ResultView.onSaveQuery", nil)()

	openfileDialog, err := gtk.FileChooserDialogNewWith2Buttons("Save file", v.w, gtk.FILE_CHOOSER_ACTION_SAVE,
		"Save", gtk.RESPONSE_OK,
		"Cancel", gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		config.Error("save file dialog", err)
		return
	}
	defer openfileDialog.Destroy()

	response := openfileDialog.Run()
	if response == gtk.RESPONSE_CANCEL {
		return
	}

	buff, err := v.sourceView.GetBuffer()
	if err != nil {
		config.Error(err)
		return
	}

	txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
	if err != nil {
		config.Error(err)
		return
	}

	v.ctrl.OnSaveQuery(txt, openfileDialog.GetFilename())
}

func (v *ResultView) onSaveFav() {
	defer config.LogStart("ResultView.onSaveFac", nil)()

	name, err := v.askFavName()
	if err != nil {
		config.Error(err)
		return
	}

	if name == "" {
		return
	}

	buff, err := v.sourceView.GetBuffer()
	if err != nil {
		config.Error(err)
		return
	}

	query, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
	if err != nil {
		config.Error(err)
		return
	}

	v.ctrl.OnSaveFav(name, query)
}

func (v *ResultView) PageSize() int64 {
	defer config.LogStart("ResultView.PageSize", nil)()

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
	defer config.LogStart("ResultView.Offset", nil)()

	return v.offset
}

func (v *ResultView) Conditions() ([]driver.CondStmt, error) {
	defer config.LogStart("ResultView.Conditions", nil)()

	return v.conditions.Statements()
}

func (v *ResultView) UpdateColumns(cols []driver.ColDef) error {
	defer config.LogStart("ResultView.UpdateColumns", nil)()

	defer v.Stop()

	v.colFilter.SetText("")
	v.offset = 0
	v.conditions.Update(cols)
	return v.result.UpdateColumns(cols)
}

func (v *ResultView) UpdateData(data [][]interface{}) error {
	defer config.LogStart("ResultView.UpdateData", nil)()

	defer v.Stop()

	v.pagerEnable(true)
	v.btnAddRow.SetSensitive(true)
	v.btnShowFilters.SetSensitive(true)

	return v.result.UpdateData(data)
}

func (v *ResultView) UpdateRawData(cols []string, data [][]interface{}) error {
	defer config.LogStart("ResultView.UpdateRawData", nil)()

	defer v.Stop()

	v.pagerEnable(false)
	v.colFilter.SetText("")
	v.offset = 0
	return v.result.UpdateRawData(cols, data)
}

func (v *ResultView) SelectedIsNewRecord() (bool, error) {
	defer config.LogStart("ResultView.SelectedIsNewRecord", nil)()

	return v.result.SelectedIsNewRecord()
}

func (v *ResultView) RemoveSelected() error {
	defer config.LogStart("ResultView.RemoveSelected", nil)()

	err := v.result.RemoveSelected()
	if err != nil {
		return err
	}

	v.newRecordEnable(false)

	return nil
}

func (v *ResultView) ForEachSelected(f func([]driver.ColDef, []interface{})) error {
	defer config.LogStart("ResultView.ForEachSelected", nil)()

	return v.result.ForEachSelected(f)
}

func (v *ResultView) GetRowID() ([]driver.ColDef, []interface{}, error) {
	defer config.LogStart("ResultView.GetRowID", nil)()

	return v.result.GetRowID()
}

func (u *ResultView) GetRow() ([]driver.ColDef, []interface{}, error) {
	defer config.LogStart("ResultView.GetRow", nil)()

	return u.result.GetRow()
}

func (u *ResultView) UpdateRow(values []interface{}) error {
	defer config.LogStart("ResultView.UpdateRow", nil)()

	err := u.result.UpdateRow(values)
	if err == nil {
		u.newRecordEnable(false)
	}
	return err
}

func (u *ResultView) SortOptions() []driver.SortOption {
	defer config.LogStart("ResultView.SortOptions", nil)()

	return u.result.SortOptions()
}

func (v *ResultView) actionbar() (*gtk.ActionBar, error) {
	defer config.LogStart("ResultView.actionbar", nil)()

	actionbar, err := gtk.ActionBarNew()
	if err != nil {
		return nil, err
	}

	// new-add-delete
	{
		v.btnAddRow, err = BWButtonNewFromIconName("add-record", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnAddRow.Show()

		v.btnDeleteRow, err = BWButtonNewFromIconName("delete-record", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnDeleteRow.Show()

		v.btnCreateRow, err = BWButtonNewFromIconName("save-record", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnCreateRow.Show()

		v.btnShowFilters, err = BWToggleButtonNewFromIconName("filter", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnShowFilters.Show()
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
		v.colFilter.Show()
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
		v.btnPrev, err = BWButtonNewFromIconName("back", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnPrev.Show()

		v.btnNext, err = BWButtonNewFromIconName("next", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnNext.Show()

		v.btnRsh, err = BWButtonNewFromIconName("refresh", "orange", ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
		v.btnRsh.Show()

		actionbar.PackEnd(v.btnRsh)
		actionbar.PackEnd(v.btnNext)
		actionbar.PackEnd(v.btnPrev)
	}

	return actionbar, nil
}

func (v *ResultView) pagerEnable(b bool) {
	defer config.LogStart("ResultView.pagerEnable", nil)()

	v.btnPrev.SetSensitive(b)
	v.btnNext.SetSensitive(b)
	v.btnRsh.SetSensitive(b)
	//v.perPage.SetSensitive(b)
	//v.offset.SetSensitive(b)
	//v.perPage.SetSensitive(b)
	//v.offset.SetSensitive(b)
}

func (v *ResultView) disableAll() {
	defer config.LogStart("ResultView.disableAll", nil)()

	v.btnPrev.SetSensitive(false)
	v.btnNext.SetSensitive(false)
	v.btnRsh.SetSensitive(false)
	v.btnAddRow.SetSensitive(false)
	v.btnShowFilters.SetSensitive(false)
	v.btnDeleteRow.SetSensitive(false)
	v.btnCreateRow.SetSensitive(false)
}

func (v *ResultView) newRecordEnable(b bool) {
	defer config.LogStart("ResultView.newRecordEnable", nil)()

	v.btnCreateRow.SetSensitive(b)
}

func (v *ResultView) onTextViewKeyPress(_ *sourceview.SourceView, e *gdk.Event) bool {
	defer config.LogStart("ResultView.onTextViewKeyPress", nil)()

	keyEvent := gdk.EventKeyNewFromEvent(e)

	if keyEvent.KeyVal() == gdk.KEY_Return && keyEvent.State()&gdk.CONTROL_MASK > 0 {
		buff, err := v.sourceView.GetBuffer()
		if err != nil {
			config.Error(err)
			return true
		}
		txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			config.Error(err)
		}

		v.ctrl.OnExecQuery(txt)
		return true
	}

	return false
}

func (v *ResultView) onColFilterSearchChanged() {
	defer config.LogStart("ResultView.onColFilterSearchChanged", nil)()

	txt, err := v.colFilter.GetText()
	if err != nil {
		config.Error(err, "colFilter.GetText")
		return
	}

	rg, err := regexp.Compile(strings.ToLower(txt))
	if err != nil {
		rg = regexp.MustCompile(fmt.Sprintf(".*%s.*", regexp.QuoteMeta(txt)))
	}

	v.result.GetColumns().Foreach(func(i interface{}) {
		c := i.(*gtk.TreeViewColumn)

		c.SetVisible(rg.MatchString(strings.Replace(strings.ToLower(c.GetTitle()), "__", "_", -1)))
	})
}

func (v *ResultView) onRowActivated(_ *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	defer config.LogStart("ResultView.onRowActivated", nil)()

	if v.result.mode == MODE_RAW {
		return
	}

	iter, err := v.result.store.GetIter(path)
	if err != nil {
		config.Error(err)
		return
	}

	s, err := v.result.store.GetValue(iter, len(v.result.cols))
	if err != nil {
		config.Error(err)
		return
	}
	status, err := s.GoValue()
	if err != nil {
		config.Error(err)
		return
	}

	v.newRecordEnable(status == STATUS_NEW)
	v.btnDeleteRow.SetSensitive(true)
}

func (v *ResultView) askFavName() (string, error) {
	defer config.LogStart("ResultView.askFavName", nil)()

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
	content.Show()

	label, err := gtk.LabelNew("Enter favorite name")
	if err != nil {
		return "", err
	}
	label.Show()

	entry, err := gtk.EntryNew()
	if err != nil {
		return "", err
	}
	entry.Show()

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return "", err
	}
	box.Show()

	box.PackStart(label, true, true, 0)
	box.PackStart(entry, true, true, 0)
	content.Add(box)

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
