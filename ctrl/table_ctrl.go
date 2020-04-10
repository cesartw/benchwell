package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// tableCtrl manages table result screen
type TableCtrl struct {
	*ConnectionCtrl
	ctx       sqlengine.Context
	tableName string

	// ui
	grid *gtk.ResultGrid
}

func (c TableCtrl) init(ctx sqlengine.Context, parent *ConnectionCtrl, tableName string) (*TableCtrl, error) {
	var err error

	c.ctx = ctx
	c.ConnectionCtrl = parent
	c.tableName = tableName

	c.grid, err = gtk.NewResultGrid(nil, nil,
		func(cols driver.ColDef, values string) (interface{}, error) {
			return c.engine.ParseValue(c.ctx, cols, values)
		})
	if err != nil {
		return nil, err
	}
	c.grid.ShowAll()
	c.grid.OnEdited(func(
		cols []driver.ColDef,
		values []interface{},
	) {
		_, err := c.engine.UpdateField(ctx, c.tableName, cols, values)
		if err != nil {
			c.window.PushStatus(err.Error())
		} else {
			c.window.PushStatus("Saved")
		}
	}).OnSubmit(func(value string) {
		columns, data, err := c.engine.Query(c.ctx, value)
		if err != nil {
			config.Env.Log.Error(err)
		}

		c.grid.UpdateRawData(columns, data)
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
	def, data, err := tc.engine.FetchTable(tc.ctx, tc.tableName,
		tc.grid.Offset(), tc.grid.PageSize())
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	err = tc.grid.UpdateData(def, data)
	if err != nil {
		config.Env.Log.Error(err)
	}
}

func (tc *TableCtrl) Screen() interface{} {
	return tc.grid
}
