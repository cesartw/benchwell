package ctrl

import (
	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
)

// tableCtrl manages table result screen
type TableCtrl struct {
	parent    *ConnectionCtrl
	ctx       sqlengine.Context
	tableName string

	// ui
	resultView *gtk.ResultView
}

func (c TableCtrl) init(ctx sqlengine.Context, parent *ConnectionCtrl, tableName string) (*TableCtrl, error) {
	var err error

	c.ctx = ctx
	c.parent = parent
	c.tableName = tableName

	c.resultView, err = gtk.NewResultView(nil, nil)
	if err != nil {
		return nil, err
	}
	c.resultView.ShowAll()

	return &c, nil
}

func (tc *TableCtrl) OnConnect() {
	def, data, err := tc.parent.engine.FetchTable(tc.ctx, tc.tableName, 0, 40)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	err = tc.resultView.UpdateData(def, data)
	if err != nil {
		config.Env.Log.Error(err)
	}
}

func (tc *TableCtrl) Screen() interface{} {
	return tc.resultView
}
