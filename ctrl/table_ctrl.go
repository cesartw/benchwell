package ctrl

import (
	"context"
	"fmt"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// tableCtrl manages table result screen
type TableCtrl struct {
	*ConnectionCtrl
	ctx      *sqlengine.Context
	tableDef driver.TableDef

	// ui
	connectionTab *gtk.ConnectionTab
	resultView    *gtk.ResultView
}

type TableCtrlOpts struct {
	Parent   *ConnectionCtrl
	TableDef driver.TableDef
	//OnTabRemoved func(*TableCtrl)
}

func (tc TableCtrl) Init(
	ctx *sqlengine.Context,
	opts TableCtrlOpts,
) (*TableCtrl, error) {
	var err error
	tc.ctx = ctx
	tc.ConnectionCtrl = opts.Parent
	tc.tableDef = opts.TableDef

	tc.resultView, err = gtk.ResultView{}.Init(
		tc.window,
		&tc,
		func(col driver.ColDef, value string) (interface{}, error) {
			return tc.Engine.ParseValue(tc.ctx, col, value)
		})
	if err != nil {
		return nil, err
	}

	tc.resultView.Show()
	tabName := tc.dbName
	if opts.TableDef.Name != "" {
		tabName += "." + opts.TableDef.Name
	}

	tc.connectionTab, err = gtk.ConnectionTab{}.Init(gtk.ConnectionTabOpts{
		Database: tc.dbName,
		Title:    tabName,
		Content:  tc.resultView,
	})
	if err != nil {
		return nil, err
	}

	if !tc.tableDef.IsZero() {
		tc.OnLoadTable()
	}

	return &tc, nil
}

func (tc *TableCtrl) String() string {
	return tc.ctx.Database().Name() + "." + tc.tableDef.Name
}

func (tc *TableCtrl) OnCopyInsert(cols []driver.ColDef, values []interface{}) {
	sql, err := tc.Engine.GetInsertStatement(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		tc.window.PushStatus(err.Error())
	}

	gtk.ClipboardCopy(sql)
	config.Env.Log.Debugf("insert copied: %s", sql)
}

func (tc *TableCtrl) OnUpdateRecord(cols []driver.ColDef, values []interface{}) error {
	_, err := tc.Engine.UpdateField(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		return err
	}

	tc.window.PushStatus("Saved")
	return nil
}

func (tc *TableCtrl) OnCreateRecord(cols []driver.ColDef, values []interface{}) ([]interface{}, error) {
	data, err := tc.Engine.InsertRecord(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		return nil, err
	} else {
		tc.window.PushStatus("Inserted")
	}

	return data, nil
}

func (tc *TableCtrl) OnExecQuery(value string) {
	columns, data, err := tc.Engine.Query(tc.ctx, value)
	if err != nil {
		return
	}
	tc.resultView.UpdateRawData(columns, data)
	tc.window.PushStatus("%d rows loaded", len(data))

	/*dml, ddl := tc.parseQuery(value)

	for _, query := range dml {
		columns, data, err := tc.Engine.Query(tc.ctx, query)
		if err != nil {
			config.Env.Log.Error(err)
			tc.window.PushStatus("Error: %s", err.Error())
			return
		}
		tc.resultView.UpdateRawData(columns, data)
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
	newRecord, err := tc.resultView.SelectedIsNewRecord()
	if err != nil {
		return
	}

	if newRecord {
		tc.resultView.RemoveSelected()
		tc.window.PushStatus("Record removed")
	} else {
		cols, values, err := tc.resultView.GetRowID()
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		err = tc.Engine.DeleteRecord(tc.ctx, tc.tableDef.Name, cols, values)
		if err != nil {
			return
		}
		tc.resultView.RemoveSelected()
		tc.window.PushStatus("Record deleted")
	}
}

func (tc *TableCtrl) OnLoadTable() {
	cancel := tc.window.Go(func(ctx context.Context) func() {
		switch tc.tableDef.Type {
		case driver.TableTypeDummy:
			// TODO: query is done on the background
			return func() {
				tc.OnExecQuery(tc.tableDef.Query)
			}
		default:
			def, data, err := tc.Engine.FetchTable(
				tc.ctx.WithContext(ctx), tc.tableDef,
				driver.FetchTableOptions{
					Offset: tc.resultView.Offset(),
					Limit:  tc.resultView.PageSize(),
				},
			)
			if err != nil {
				return func() {}
			}

			if tc.tableDef.Type == driver.TableTypeRegular {
				return func() {
					err = tc.resultView.UpdateColumns(def)
					if err != nil {
						config.Env.Log.Error(err)
						return
					}

					err = tc.resultView.UpdateData(data)
					if err != nil {
						config.Env.Log.Error(err)
					}
				}
			} else {
				columns := []string{}
				for _, d := range def {
					columns = append(columns, d.Name)
				}

				return func() {
					tc.resultView.UpdateRawData(columns, data)
					if err != nil {
						config.Env.Log.Error(err)
					}
				}
			}
		}
	})

	tc.resultView.Block(cancel)
}

func (tc *TableCtrl) OnRefresh() {
	conditions, err := tc.resultView.Conditions()
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	_, data, err := tc.Engine.FetchTable(
		tc.ctx, tc.tableDef,
		driver.FetchTableOptions{
			Offset:     tc.resultView.Offset(),
			Limit:      tc.resultView.PageSize(),
			Sort:       tc.resultView.SortOptions(),
			Conditions: conditions,
		},
	)
	if err != nil {
		return
	}

	err = tc.resultView.UpdateData(data)
	if err != nil {
		return
	}

	tc.window.PushStatus("Table reloaded")
}

func (tc *TableCtrl) OnApplyConditions() {
	tc.OnRefresh()
}

func (tc *TableCtrl) OnCreate() {
	newRecord, err := tc.resultView.SelectedIsNewRecord()
	if err != nil {
		return
	}
	if !newRecord {
		return
	}

	cols, values, err := tc.resultView.GetRow()
	if err != nil {
		return
	}

	values, err = tc.Engine.InsertRecord(tc.ctx, tc.tableDef.Name, cols, values)
	if err != nil {
		return
	}

	err = tc.resultView.UpdateRow(values)
	if err != nil {
		tc.window.PushStatus(err.Error())
		return
	}

	tc.window.PushStatus("Record saved")
}

func (tc *TableCtrl) SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error) {
	if tc.ctx != nil && tc.ctx != ctx {
		return false, nil
	}

	tc.tableDef = tableDef
	tc.connectionTab.SetTitle(fmt.Sprintf("%s.%s", tc.dbName, tableDef.Name))
	tc.OnLoadTable()

	return true, nil
}

func (tc *TableCtrl) SetQuery(ctx *sqlengine.Context, query string) (bool, error) {
	if tc.ctx != nil && tc.ctx != ctx {
		return false, nil
	}

	tc.resultView.SetQuery(query)
	return true, nil
}

func (tc *TableCtrl) OnTabRemove() {
	tc.ConnectionCtrl.OnTabRemove(tc)
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
