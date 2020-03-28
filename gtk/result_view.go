package gtk

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/quick"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk/controls"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
)

// ResultView is a table result tab content
type ResultView struct {
	*gtk.Paned

	textView   *gtk.TextView
	textViewSW *gtk.ScrolledWindow

	result   *controls.Result
	resultSW *gtk.ScrolledWindow
}

func NewResultView(cols []driver.ColDef, data [][]interface{}) (rv *ResultView, err error) {
	rv = &ResultView{}

	rv.Paned, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}

	rv.textView, err = gtk.TextViewNew()
	if err != nil {
		return nil, err
	}

	rv.textView.Connect("key-release-event", func() {
		buff, err := rv.textView.GetBuffer()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		txt, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		txt, err = ChromaHighlight(txt)
		if err != nil {
			config.Env.Log.Error(err)
			return
		}
		buff.Delete(buff.GetStartIter(), buff.GetEndIter())
		buff.InsertMarkup(buff.GetStartIter(), txt)
	})

	rv.resultSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	rv.textViewSW, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	rv.textViewSW.SetSizeRequest(-1, 200)

	rv.result, err = controls.NewResult(cols, data)
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

	rv.resultSW.Add(rv.result)
	rv.textViewSW.Add(rv.textView)

	rv.Paned.Pack1(rv.textViewSW, false, false)
	rv.Paned.Pack2(rv.resultSW, true, false)

	return rv, nil
}

func (v *ResultView) UpdateData(cols []driver.ColDef, data [][]interface{}) error {
	return v.result.UpdateData(cols, data)
}

func ChromaHighlight(inputString string) (out string, err error) {

	buff := new(bytes.Buffer)
	writer := bufio.NewWriter(buff)

	// Registrering pango formatter
	formatters.Register("pango", chroma.FormatterFunc(pangoFormatter))

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
