package ctrl

import (
	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
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

	c.resultView, err = gtk.NewResultView(nil, nil,
		func(cols driver.ColDef, values string) (interface{}, error) {
			return c.parent.engine.ParseValue(c.ctx, cols, values)
		})
	if err != nil {
		return nil, err
	}
	c.resultView.ShowAll()
	c.resultView.OnEdited(func(
		cols []driver.ColDef,
		values []interface{},
	) {
		_, err := c.parent.engine.UpdateField(ctx, c.tableName, cols, values)
		if err != nil {
			c.parent.factory.PushStatus(err.Error())
		} else {
			c.parent.factory.PushStatus("Saved")
		}
	}).OnSubmit(func(value string) {
		columns, data, err := c.parent.engine.Query(c.ctx, value)
		if err != nil {
			config.Env.Log.Error(err)
		}

		c.resultView.UpdateRawData(columns, data)
	}).OnRefresh(func() {
		c.OnConnect()
	}).OnBack(func() {
		c.OnConnect()
	}).OnForward(func() {
		c.OnConnect()
	})

	return &c, nil
}

func (tc *TableCtrl) OnConnect() {
	def, data, err := tc.parent.engine.FetchTable(tc.ctx, tc.tableName,
		tc.resultView.Offset(), tc.resultView.PageSize())
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
