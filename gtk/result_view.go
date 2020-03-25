package gtk

import (
	"bitbucket.org/goreorto/sqlhero/gtk/controls"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gotk3/gotk3/gtk"
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
	// TODO: move to config
	//rv.textView.SetProperty("wrap-mode", gtk.WRAP_WORD)

	rv.resultSW.Add(rv.result)
	rv.textViewSW.Add(rv.textView)

	rv.Paned.Pack1(rv.textViewSW, false, false)
	rv.Paned.Pack2(rv.resultSW, true, false)

	return rv, nil
}

func (v *ResultView) UpdateData(cols []driver.ColDef, data [][]interface{}) error {
	return v.result.UpdateData(cols, data)
}
