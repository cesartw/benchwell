package ctrl

import (
	"context"
	"fmt"

	"bitbucket.org/goreorto/sqlaid/clipboard"
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// tableCtrl manages table result screen
type TableCtrl struct {
	*ConnectionCtrl
	ctx      sqlengine.Context
	tableDef driver.TableDef

	// ui
	connectionTab *gtk.ConnectionTab
	grid          *gtk.ResultGrid
}

type TableCtrlOpts struct {
	Parent       *ConnectionCtrl
	TableDef     driver.TableDef
	OnTabRemoved func(*TableCtrl)
	Log          func(string)
}

func (tc TableCtrl) Init(
	ctx sqlengine.Context,
	opts TableCtrlOpts,
) (*TableCtrl, error) {
	var err error

	tc.ctx = driver.SetLogger(ctx, opts.Log)
	tc.ConnectionCtrl = opts.Parent
	tc.tableDef = opts.TableDef

	tc.grid, err = gtk.ResultGrid{}.Init(
		&tc,
		func(col driver.ColDef, value string) (interface{}, error) {
			return tc.Engine.ParseValue(tc.ctx, col, value)
		})
	if err != nil {
		return nil, err
	}

	tc.grid.Show()

	tc.connectionTab, err = gtk.ConnectionTab{}.Init(gtk.ConnectionTabOpts{
		Title:   fmt.Sprintf("%s.%s", tc.dbName, opts.TableDef.Name),
		Content: tc.grid,
		OnRemove: func() {
			opts.OnTabRemoved(&tc)
		},
	})
	if err != nil {
		return nil, err
	}

	if !tc.tableDef.IsZero() {
		tc.OnConnect()
	}

	return &tc, nil
}

func (tc *TableCtrl) OnCopyInsert(cols []driver.ColDef, values []interface{}) {
	sql, err := tc.Engine.GetInsertStatement(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		tc.window.PushStatus(err.Error())
	}

	clipboard.Copy(sql)
	config.Env.Log.Debugf("insert copied: %s", sql)
}

func (tc *TableCtrl) OnUpdateRecord(cols []driver.ColDef, values []interface{}) error {
	_, err := tc.Engine.UpdateField(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		tc.window.PushStatus(err.Error())
		return err
	}

	tc.window.PushStatus("Saved")
	return nil
}

func (tc *TableCtrl) OnCreateRecord(cols []driver.ColDef, values []interface{}) ([]interface{}, error) {
	data, err := tc.Engine.InsertRecord(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		tc.window.PushStatus(err.Error())
		return nil, err
	} else {
		tc.window.PushStatus("Inserted")
	}

	return data, nil
}

func (tc *TableCtrl) OnExecQuery(value string) {
	columns, data, err := tc.Engine.Query(tc.ctx, value)
	if err != nil {
		config.Env.Log.Error(err)
		tc.window.PushStatus("Error: %s", err.Error())
		return
	}
	tc.grid.UpdateRawData(columns, data)
	tc.window.PushStatus("%d rows loaded", len(data))

	/*dml, ddl := tc.parseQuery(value)

	for _, query := range dml {
		columns, data, err := tc.Engine.Query(tc.ctx, query)
		if err != nil {
			config.Env.Log.Error(err)
			tc.window.PushStatus("Error: %s", err.Error())
			return
		}
		tc.grid.UpdateRawData(columns, data)
		tc.window.PushStatus("%d rows loaded", len(data))
	}

	for _, query := range ddl {
		id, affected, err := tc.Engine.Execute(tc.ctx, query)
		if err != nil {
			config.Env.Log.Error(err)
			tc.window.PushStatus("Error: %s", err.Error())
			return
		}
		tc.window.PushStatus("Last inserted id: %s Affected rows: %d", id, affected)
	}
	*/
}

func (tc *TableCtrl) OnDelete() {
	newRecord, err := tc.grid.SelectedIsNewRecord()
	if err != nil {
		return
	}

	if newRecord {
		tc.grid.RemoveSelected()
		tc.window.PushStatus("Record removed")
	} else {
		cols, values, err := tc.grid.GetRowID()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		tc.Engine.DeleteRecord(tc.ctx, tc.tableDef.Name, cols, values)
		tc.grid.RemoveSelected()
		tc.window.PushStatus("Record deleted")
	}
}

func (tc *TableCtrl) OnConnect() {
	def, data, err := tc.Engine.FetchTable(
		tc.ctx, tc.tableDef.Name,
		driver.FetchTableOptions{
			Offset: tc.grid.Offset(),
			Limit:  tc.grid.PageSize(),
			Sort:   tc.grid.SortOptions(),
		},
	)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	if tc.tableDef.Type == driver.TableTypeRegular {
		err = tc.grid.UpdateColumns(def)
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		err = tc.grid.UpdateData(data)
		if err != nil {
			config.Env.Log.Error(err)
		}
	} else {
		columns := []string{}
		for _, d := range def {
			columns = append(columns, d.Name)
		}
		tc.grid.UpdateRawData(columns, data)
	}
}

func (tc *TableCtrl) OnRefresh() {
	conditions, err := tc.grid.Conditions()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	_, data, err := tc.Engine.FetchTable(
		tc.ctx, tc.tableDef.Name,
		driver.FetchTableOptions{
			Offset:     tc.grid.Offset(),
			Limit:      tc.grid.PageSize(),
			Sort:       tc.grid.SortOptions(),
			Conditions: conditions,
		},
	)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	err = tc.grid.UpdateData(data)
	if err != nil {
		config.Env.Log.Error(err)
	}

	tc.window.PushStatus("Table reloaded")
}

func (tc *TableCtrl) OnCreate() {
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

	values, err = tc.Engine.InsertRecord(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		tc.window.PushStatus(err.Error())
		return
	}

	err = tc.grid.UpdateRow(values)
	if err != nil {
		tc.window.PushStatus(err.Error())
		return
	}

	tc.window.PushStatus("Record saved")
}

func (tc *TableCtrl) SetTableDef(ctx context.Context, tableDef driver.TableDef) {
	tc.tableDef = tableDef
	tc.ctx = ctx
	tc.connectionTab.SetTitle(fmt.Sprintf("%s.%s", tc.dbName, tableDef.Name))
	tc.OnConnect()
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
