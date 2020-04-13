package gtk

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io"
	"regexp"
	"strconv"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/quick"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// ResultGrid is a table result tab content
type ResultGrid struct {
	*gtk.Paned

	textView *gtk.TextView

	result  *Result
	btnPrev *gtk.Button
	btnNext *gtk.Button
	btnRsh  *gtk.Button
	perPage *gtk.Entry
	offset  *gtk.Entry

	btnAddRow    *gtk.Button
	btnDeleteRow *gtk.Button
	btnCreateRow *gtk.Button

	colFilter *gtk.SearchEntry

	submitCallbacks []func(string)
}

func NewResultGrid(
	cols []driver.ColDef,
	data [][]interface{},
	parser parser,
) (v *ResultGrid, err error) {
	v = &ResultGrid{}

	v.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}

	v.textView, err = gtk.TextViewNew()
	if err != nil {
		return nil, err
	}

	v.textView.Connect("key-release-event", v.onTextViewKeyPress)

	var resultSW, textViewSW *gtk.ScrolledWindow

	resultSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	// buttonbox for add/remove rows
	resultBtnBox, err := v.resultButtonBox()
	if err != nil {
		return nil, err
	}

	// buttonbox for pagination control
	pagerBtnBox, err := v.paginationButtonBox()
	if err != nil {
		return nil, err
	}

	v.btnAddRow.Connect("clicked", func() {
		v.result.AddEmptyRow()
	})

	v.btnDeleteRow.Connect("clicked", func() {
	})

	v.btnNext.Connect("clicked", func() {
		p := v.Offset() + v.PageSize()
		v.offset.SetText(fmt.Sprintf("%d", p))
	})
	v.btnPrev.Connect("clicked", func() {
		p := v.Offset() - v.PageSize()
		if p < 0 {
			p = 0
		}
		v.offset.SetText(fmt.Sprintf("%d", p))
	})
	v.colFilter.Connect("search-changed", v.onColFilterSearchChanged)

	btnGridBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	btnGridBox.PackStart(resultBtnBox, false, false, 0)
	btnGridBox.PackEnd(resultSW, true, true, 0)

	resultBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	resultBox.PackStart(btnGridBox, true, true, 0)
	resultBox.PackEnd(pagerBtnBox, false, false, 0)

	textViewSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	textViewSW.SetSizeRequest(-1, 200)

	v.result, err = NewResult(cols, data, parser)
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

	v.Paned.Pack1(textViewSW, false, false)
	v.Paned.Pack2(resultBox, true, false)

	return v, nil
}

func (v *ResultGrid) Offset() int64 {
	s, _ := v.offset.GetText()
	p, _ := strconv.ParseInt(s, 10, 64)
	return p
}

func (v *ResultGrid) PageSize() int64 {
	s, err := v.perPage.GetText()
	if err != nil {
		return int64(config.Env.GUI.PageSize)
	}

	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(config.Env.GUI.PageSize)
	}

	return size
}

func (v *ResultGrid) UpdateData(cols []driver.ColDef, data [][]interface{}) error {
	v.pagerEnable(true)
	return v.result.UpdateData(cols, data)
}

func (v *ResultGrid) UpdateRawData(cols []string, data [][]interface{}) error {
	v.pagerEnable(false)
	return v.result.UpdateRawData(cols, data)
}

func (v *ResultGrid) SetUpdateRecordFunc(fn func([]driver.ColDef, []interface{}) error) *ResultGrid {
	v.result.SetUpdateRecordFunc(fn)
	return v
}

func (v *ResultGrid) SetCreateRecordFunc(fn func([]driver.ColDef, []interface{}) ([]interface{}, error)) *ResultGrid {
	v.result.SetCreateRecordFunc(fn)
	return v
}

func (v *ResultGrid) OnSubmit(fn func(value string)) *ResultGrid {
	v.submitCallbacks = append(v.submitCallbacks, fn)
	return v
}

func (v *ResultGrid) OnRefresh(fn interface{}) *ResultGrid {
	v.btnRsh.Connect("clicked", fn)
	return v
}

func (v *ResultGrid) OnBack(fn interface{}) *ResultGrid {
	v.btnPrev.Connect("clicked", fn)
	return v
}

func (v *ResultGrid) OnForward(fn interface{}) *ResultGrid {
	v.btnNext.Connect("clicked", fn)
	return v
}

func (v *ResultGrid) OnCreate(fn interface{}) *ResultGrid {
	v.btnCreateRow.Connect("clicked", fn)
	return v
}

func (v *ResultGrid) OnDelete(fn interface{}) *ResultGrid {
	v.btnDeleteRow.Connect("clicked", fn)
	return v
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

func (v *ResultGrid) paginationButtonBox() (*gtk.ButtonBox, error) {
	pagerBtnBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	pagerBtnBox.SetLayout(gtk.BUTTONBOX_CENTER)
	pagerBtnBox.SetProperty("spacing", 5)

	perPageLabel, err := gtk.LabelNew("Size")
	if err != nil {
		return nil, err
	}

	v.perPage, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	v.perPage.SetText(fmt.Sprintf("%d", config.Env.GUI.PageSize))
	v.perPage.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	v.offset, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	v.offset.SetText("0")
	v.offset.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

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

	offsetLabel, err := gtk.LabelNew("Offset")
	if err != nil {
		return nil, err
	}

	v.offset, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	v.offset.SetText("0")
	v.offset.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	v.colFilter, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	v.colFilter.SetPlaceholderText("Column filter: .*")

	pagerBtnBox.Add(v.btnPrev)
	pagerBtnBox.Add(perPageLabel)
	pagerBtnBox.Add(v.perPage)
	pagerBtnBox.Add(offsetLabel)
	pagerBtnBox.Add(v.offset)
	pagerBtnBox.Add(v.btnNext)
	pagerBtnBox.Add(v.btnRsh)

	pagerBtnBox.PackEnd(v.colFilter, false, false, 0)

	return pagerBtnBox, nil
}

func (v *ResultGrid) resultButtonBox() (*gtk.ButtonBox, error) {
	btnbox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	btnbox.SetLayout(gtk.BUTTONBOX_START)
	btnbox.SetProperty("spacing", 5)

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
	v.newRecordEnable(false)

	btnbox.Add(v.btnAddRow)
	btnbox.Add(v.btnDeleteRow)
	btnbox.Add(v.btnCreateRow)

	return btnbox, nil
}

func (v *ResultGrid) pagerEnable(b bool) {
	v.btnPrev.SetSensitive(b)
	v.btnNext.SetSensitive(b)
	v.btnRsh.SetSensitive(b)
	v.perPage.SetSensitive(b)
	v.offset.SetSensitive(b)
}

func (v *ResultGrid) newRecordEnable(b bool) {
	v.btnCreateRow.SetSensitive(b)
}

func (v *ResultGrid) onTextViewKeyPress(_ *gtk.TextView, e *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(e)
	if keyEvent.KeyVal() >= gdk.KEY_Home && keyEvent.KeyVal() <= gdk.KEY_End {
		return
	}
	if keyEvent.KeyVal() == gdk.KEY_Shift_R || keyEvent.KeyVal() == gdk.KEY_Shift_L {
		return
	}

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

	if keyEvent.KeyVal() == gdk.KEY_Return && keyEvent.State()&gdk.GDK_CONTROL_MASK > 0 {
		for _, fn := range v.submitCallbacks {
			fn(txt)
		}
		return
	}

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
		c.SetVisible(rg.MatchString(c.GetTitle()))
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
}

func init() {
	// Registrering pango formatter
	formatters.Register("pango", chroma.FormatterFunc(pangoFormatter))
}

func ChromaHighlight(inputString string) (out string, err error) {
	buff := new(bytes.Buffer)
	writer := bufio.NewWriter(buff)

	// Doing the job (io.Writer, SourceText, language(go), Lexer(pango), style(pygments))
	if err = quick.Highlight(writer, inputString, "sql", "pango", "pygments"); err != nil {
		return
	}
	writer.Flush()
	return string(buff.Bytes()), err
}

func pangoFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
	var r, g, b uint8
	var closer, out string

	var getColour = func(color chroma.Colour) string {
		r, g, b = color.Red(), color.Green(), color.Blue()
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}

	for tkn := it(); tkn != chroma.EOF; tkn = it() {

		entry := style.Get(tkn.Type)
		if !entry.IsZero() {
			if entry.Bold == chroma.Yes {
				out = `<b>`
				closer = `</b>`
			}
			if entry.Underline == chroma.Yes {
				out += `<u>`
				closer = `</u>` + closer
			}
			if entry.Italic == chroma.Yes {
				out += `<i>`
				closer = `</i>` + closer
			}
			if entry.Colour.IsSet() {
				out += `<span foreground="` + getColour(entry.Colour) + `">`
				closer = `</span>` + closer
			}
			if entry.Background.IsSet() {
				out += `<span background="` + getColour(entry.Background) + `">`
				closer = `</span>` + closer
			}
			if entry.Border.IsSet() {
				out += `<span background="` + getColour(entry.Border) + `">`
				closer = `</span>` + closer
			}
			fmt.Fprint(w, out)
		}
		fmt.Fprint(w, html.EscapeString(tkn.Value))
		if !entry.IsZero() {
			fmt.Fprint(w, closer)
		}
		closer, out = "", ""
	}
	return nil
}
