package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/clipboard"
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

func (tc TableCtrl) init(ctx sqlengine.Context, parent *ConnectionCtrl, tableName string) (*TableCtrl, error) {
	var err error

	tc.ctx = ctx
	tc.ConnectionCtrl = parent
	tc.tableName = tableName

	tc.grid, err = gtk.NewResultGrid(nil, nil,
		func(col driver.ColDef, value string) (interface{}, error) {
			return tc.engine.ParseValue(tc.ctx, col, value)
		})
	if err != nil {
		return nil, err
	}
	tc.grid.ShowAll()
	tc.grid.SetUpdateRecordFunc(func(
		cols []driver.ColDef,
		values []interface{},
	) error {
		_, err := tc.engine.UpdateField(ctx, tc.tableName, cols, values)
		if err != nil {
			tc.window.PushStatus(err.Error())
			return err
		}

		tc.window.PushStatus("Saved")
		return nil
	}).SetCreateRecordFunc(func(
		cols []driver.ColDef,
		values []interface{},
	) ([]interface{}, error) {
		data, err := tc.engine.InsertRecord(ctx, tc.tableName, cols, values)
		if err != nil {
			tc.window.PushStatus(err.Error())
			return nil, err
		} else {
			tc.window.PushStatus("Inserted")
		}

		return data, nil
	}).OnSubmit(func(value string) {
		columns, data, err := tc.engine.Query(tc.ctx, value)
		if err != nil {
			config.Env.Log.Error(err)
			tc.window.PushStatus("Error: %s", err.Error())
			return
		}
		tc.grid.UpdateRawData(columns, data)
		tc.window.PushStatus("%d rows loaded", len(data))

		/*dml, ddl := tc.parseQuery(value)

		for _, query := range dml {
			columns, data, err := tc.engine.Query(tc.ctx, query)
			if err != nil {
				config.Env.Log.Error(err)
				tc.window.PushStatus("Error: %s", err.Error())
				return
			}
			tc.grid.UpdateRawData(columns, data)
			tc.window.PushStatus("%d rows loaded", len(data))
		}

		for _, query := range ddl {
			id, affected, err := tc.engine.Execute(tc.ctx, query)
			if err != nil {
				config.Env.Log.Error(err)
				tc.window.PushStatus("Error: %s", err.Error())
				return
			}
			tc.window.PushStatus("Last inserted id: %s Affected rows: %d", id, affected)
		}
		*/
	}).OnRefresh(func() {
		tc.OnConnect()
	}).OnBack(func() {
		tc.OnConnect()
	}).OnForward(func() {
		tc.OnConnect()
	}).OnDelete(func() {
		newRecord, err := tc.grid.SelectedIsNewRecord()
		if err != nil {
			return
		}

		if newRecord {
			tc.grid.RemoveSelected()
		} else {
			cols, values, err := tc.grid.GetRowID()
			if err != nil {
				config.Env.Log.Error(err)
				return
			}

			tc.engine.DeleteRecord(tc.ctx, tc.tableName, cols, values)
			tc.grid.RemoveSelected()
		}
	}).OnCreate(func() {
		newRecord, err := tc.grid.SelectedIsNewRecord()
		if err != nil {
			return
		}
		if !newRecord {
			return
		}

		cols, values, err := tc.grid.GetRow()
		if err != nil {
			return
		}

		values, err = tc.engine.InsertRecord(tc.ctx, tc.tableName, cols, values)
		if err != nil {
			return
		}

		err = tc.grid.UpdateRow(values)
		if err != nil {
			return
		}
	}).OnCopyInsert(func(cols []driver.ColDef, values []interface{}) {
		sql, err := tc.engine.GetInsertStatement(tc.ctx, tc.tableName, cols, values)
		if err != nil {
			tc.window.PushStatus(err.Error())
		}

		clipboard.Copy(sql)
		config.Env.Log.Debugf("insert copied: %s", sql)
	})

	return &tc, nil
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

/*
func (tc *TableCtrl) parseQuery(src string) (dml []string, ddl []string) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(src, "", "")
	if err != nil {
		config.Env.Log.Error(err)
	}

	for _, node := range stmtNodes {
		_, isSelect := node.(*ast.SelectStmt)
		_, isShow := node.(*ast.ShowStmt)
		_, isExplain := node.(*ast.ExplainStmt)

		if isShow || isSelect || isExplain {
			dml = append(dml, node.Text())
		} else {
			ddl = append(ddl, node.Text())
		}
	}

	return
}
*/

func (tc *TableCtrl) Screen() interface{} {
	return tc.grid
}
