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

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
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

	colFilter *gtk.SearchEntry

	submitCallbacks []func(string)
}

func NewResultGrid(
	cols []driver.ColDef,
	data [][]interface{},
	parser parser,
) (rv *ResultGrid, err error) {
	rv = &ResultGrid{}

	rv.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}

	rv.textView, err = gtk.TextViewNew()
	if err != nil {
		return nil, err
	}

	rv.textView.Connect("key-release-event", rv.onTextViewKeyPress)

	var resultSW, textViewSW *gtk.ScrolledWindow

	resultSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	btnbox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	btnbox.SetLayout(gtk.BUTTONBOX_CENTER)
	btnbox.SetProperty("spacing", 5)

	resultBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	perPageLabel, err := gtk.LabelNew("Size")
	if err != nil {
		return nil, err
	}

	rv.perPage, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	rv.perPage.SetText(fmt.Sprintf("%d", config.Env.GUI.PageSize))
	rv.perPage.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	rv.offset, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	rv.offset.SetText("0")
	rv.offset.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	rv.btnPrev, err = gtk.ButtonNewFromIconName("gtk-go-back", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	rv.btnPrev.Connect("clicked", func() {
		p := rv.Offset() - rv.PageSize()
		if p < 0 {
			p = 0
		}
		rv.offset.SetText(fmt.Sprintf("%d", p))
	})

	rv.btnNext, err = gtk.ButtonNewFromIconName("gtk-go-forward", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	rv.btnNext.Connect("clicked", func() {
		p := rv.Offset() + rv.PageSize()
		rv.offset.SetText(fmt.Sprintf("%d", p))
	})

	rv.btnRsh, err = gtk.ButtonNewFromIconName("gtk-refresh", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	offsetLabel, err := gtk.LabelNew("Offset")
	if err != nil {
		return nil, err
	}

	rv.offset, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	rv.offset.SetText("0")
	rv.offset.SetProperty("input_purpose", gtk.INPUT_PURPOSE_NUMBER)

	rv.colFilter, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	rv.colFilter.SetPlaceholderText("Column filter: .*")
	rv.colFilter.Connect("search-changed", rv.onColFilterSearchChanged)

	btnbox.Add(rv.btnPrev)
	btnbox.Add(perPageLabel)
	btnbox.Add(rv.perPage)
	btnbox.Add(offsetLabel)
	btnbox.Add(rv.offset)
	btnbox.Add(rv.btnNext)
	btnbox.Add(rv.btnRsh)

	btnbox.PackEnd(rv.colFilter, false, false, 0)

	resultBox.PackStart(resultSW, true, true, 0)
	resultBox.PackEnd(btnbox, false, false, 0)

	textViewSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	textViewSW.SetSizeRequest(-1, 200)

	rv.result, err = NewResult(cols, data, parser)
	if err != nil {
		return nil, err
	}

	rv.textView.SetProperty("accepts-tab", true)
	rv.Paned.SetProperty("wide-handle", true)
	rv.textView.SetLeftMargin(10)
	// this naming mess
	rv.textView.SetProperty("top-margin", 10)

	rv.textView.SetProperty("wrap-mode", map[string]gtk.WrapMode{
		"none":      gtk.WRAP_NONE,
		"char":      gtk.WRAP_CHAR,
		"word":      gtk.WRAP_WORD,
		"word_char": gtk.WRAP_WORD_CHAR,
	}[config.Env.GUI.Editor.WordWrap])

	resultSW.Add(rv.result)
	textViewSW.Add(rv.textView)

	rv.Paned.Pack1(textViewSW, false, false)
	rv.Paned.Pack2(resultBox, true, false)

	return rv, nil
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

func (v *ResultGrid) pagerEnable(b bool) {
	v.btnPrev.SetSensitive(b)
	v.btnNext.SetSensitive(b)
	v.btnRsh.SetSensitive(b)
	v.perPage.SetSensitive(b)
	v.offset.SetSensitive(b)
}

func (v *ResultGrid) OnEdited(fn func([]driver.ColDef, []interface{})) *ResultGrid {
	v.result.OnEdited(fn)
	return v
}

func (v *ResultGrid) OnSubmit(fn func(value string)) *ResultGrid {
	v.submitCallbacks = append(v.submitCallbacks, fn)
	return v
}

func (v *ResultGrid) onTextViewKeyPress(_ *gtk.TextView, e *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(e)

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

	if keyEvent.KeyVal() == 65293 && keyEvent.State()&gdk.GDK_CONTROL_MASK > 0 {
		for _, fn := range v.submitCallbacks {
			fn(txt)
		}
		return
	}

	txt, err = ChromaHighlight(txt)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
	buff.InsertMarkup(buff.GetStartIter(), txt)
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
