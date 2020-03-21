package gtk

import (
	"bitbucket.org/goreorto/sqlhero/gtk/controls"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gotk3/gotk3/gtk"
)

// ResultView is a table result tab content
type ResultView struct {
	*gtk.ScrolledWindow
	result *controls.Result
}

func NewResultView() (rv *ResultView, err error) {
	rv = &ResultView{}

	rv.ScrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	rv.result, err = controls.NewResult(nil, nil)
	if err != nil {
		return nil, err
	}
	rv.ScrolledWindow.Add(rv.result)

	return rv, nil
}

func (v *ResultView) UpdateData(cols []driver.ColDef, data [][]interface{}) error {
	return v.result.UpdateData(cols, data)
}
