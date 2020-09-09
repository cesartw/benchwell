package ctrl

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
	"bitbucket.org/goreorto/benchwell/sqlengine"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type ConnectionCtrl struct {
	*DbTabCtrl

	ctx      *sqlengine.Context
	tableDef driver.TableDef

	scr    *gtk.ConnectionScreen
	conn   *config.Connection
	dbName string
}

func (c ConnectionCtrl) Init(
	ctx *sqlengine.Context,
	p *DbTabCtrl,
	conn *config.Connection,
) (*ConnectionCtrl, error) {
	defer config.LogStart("ConnectionCtrl.Init", nil)()

	c.DbTabCtrl = p
	c.ctx = ctx
	c.conn = conn

	dbNames, err := c.Engine.Databases(c.ctx)
	if err != nil {
		return nil, err
	}

	c.scr, err = gtk.ConnectionScreen{}.Init(c.window, &c)
	if err != nil {
		return nil, err
	}

	c.scr.SetDatabases(dbNames)

	c.scr.ShowAll()

	if c.conn.Database != "" {
		c.scr.SetActiveDatabase(c.conn.Database)
		c.dbName = c.conn.Database
	}
	c.updateTitle()

	return &c, nil
}

func (c *ConnectionCtrl) Content() ggtk.IWidget {
	defer config.LogStart("ConnectionCtrl.Content", nil)()

	return c.scr
}

func (c *ConnectionCtrl) OnCopyLog() {
	defer config.LogStart("ConnectionCtrl.OnCopyLog", nil)()
}

func (c *ConnectionCtrl) OnDatabaseSelected() {
	defer config.LogStart("ConnectionCtrl.OnDatabaseSelected", nil)()

	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.window.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}

	c.ctx, err = c.Engine.UseDatabase(c.ctx, dbName)
	if err != nil {
		c.window.PushStatus("Error selecting database: `%s`", err.Error())
		return
	}
	c.ctx.Logger = c.scr.Log

	tables, err := c.Engine.Tables(c.ctx)
	if err != nil {
		c.window.PushStatus("Error getting tables: `%s`", err.Error())
		return
	}
	err = c.conn.LoadQueries()
	if err != nil {
		c.window.PushStatus("Error loading fake tables: `%s`", err.Error())
	}

	for _, q := range c.conn.Queries {
		tables = append(tables, driver.TableDef{
			Name:  q.Name,
			Type:  driver.TableTypeDummy,
			Query: q.Query,
		})
	}

	c.scr.SetTables(tables)
	c.dbName = dbName
	c.tableDef = driver.TableDef{}
	c.updateTitle()
}

func (c *ConnectionCtrl) OnTableSelected() {
	defer config.LogStart("ConnectionCtrl.OnTableSelected", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}
	c.tableDef = tableDef

	c.OnLoadTable()
}

func (c *ConnectionCtrl) OnEditTable() {
	defer config.LogStart("ConnectionCtrl.OnEditTable", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		return
	}

	if tableDef.Type == driver.TableTypeDummy {
		c.scr.SetQuery(tableDef.Query)
	}
}

func (c *ConnectionCtrl) OnTruncateTable() {
	defer config.LogStart("ConnectionCtrl.OnTruncateTable", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	c.Engine.TruncateTable(c.ctx, tableDef)
}

func (c *ConnectionCtrl) OnDeleteTable() {
	defer config.LogStart("ConnectionCtrl.OnDeleteTable", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	err := c.Engine.DeleteTable(c.ctx, tableDef)
	if err != nil {
		c.window.PushStatus(err.Error())
	}
	c.OnDatabaseSelected()
}

func (c *ConnectionCtrl) OnCopySelect() {
	defer config.LogStart("ConnectionCtrl.OnCopySelect", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	sql, err := c.Engine.GetSelectStatement(c.ctx, tableDef)
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	gtk.ClipboardCopy(sql)
	config.Debugf("select copied: %s", sql)
}

func (c *ConnectionCtrl) OnSchemaMenu() {
	defer config.LogStart("ConnectionCtrl.OnSchemaMenu", nil)()

	tableName, ok := c.scr.SelectedTable()
	if !ok {
		return
	}

	schema, err := c.Engine.GetCreateTable(c.ctx, tableName)
	if err != nil {
		config.Error(err, "getting table schema")
	}

	c.scr.ShowTableSchemaModal(tableName, schema)
}

func (c *ConnectionCtrl) OnNewTabMenu() {
	defer config.LogStart("ConnectionCtrl.OnNewTabMenu", nil)()

	table, ok := c.scr.ActiveTable()
	if !ok {
		return
	}

	ctrl, err := c.AddTab(TAB_TYPE_DB)
	if err != nil {
		config.Error(err)
		return
	}

	conn := *c.conn
	dbCtrl := ctrl.(*DbTabCtrl)
	connectCtrl := dbCtrl.currentCtrl.(*ConnectCtrl)
	connectCtrl.scr.SetConnection(&conn)
	dbCtrl.onConnect(func() {
		connectionCtrl := dbCtrl.currentCtrl.(*ConnectionCtrl)
		connectionCtrl.scr.SetActiveTable(table)
		connectionCtrl.OnTableSelected()
	})
}

func (c *ConnectionCtrl) OnRefreshMenu() {
	defer config.LogStart("ConnectionCtrl.OnRefreshMenu", nil)()

	c.ctx.CacheTable = nil

	c.OnDatabaseSelected()
}

func (c *ConnectionCtrl) OnSaveFav(name, query string) {
	defer config.LogStart("ConnectionCtrl.OnSaveFav", nil)()

	config.SaveQuery(&config.Query{
		Name:         name,
		Query:        query,
		ConnectionID: c.conn.ID,
	})

	c.OnDatabaseSelected()
}

func (c *ConnectionCtrl) Screen() interface{} {
	defer config.LogStart("ConnectionCtrl.Screen", nil)()

	return c.scr
}

func (c *ConnectionCtrl) OnTextChange(query string, cursorAt int) {
	defer config.LogStart("ConnectionCtrl.OnTextChange", nil)()

	return
	// TODO: need to implement sourceview completion
	//columnMachines := driver.CompleteColumnMachines(*c.conn)
	//tableMachines := driver.CompleteTableMachines(*c.conn)
	//l := lexers.Get("sql")
	//it, _ := l.Tokenise(nil, string(query[:cursorAt]))
	//tokens := it.Tokens()

	//for _, m := range tableMachines {
	//_, ok := m.Match(tokens)
	//if !ok {
	//continue
	//}

	//tables, err := c.Engine.Tables(c.ctx)
	//if err != nil {
	//config.Error(err)
	//return
	//}

	//words := make([]string, len(tables))
	//for i, t := range tables {
	//words[i] = t.String()
	//}

	//c.scr.ShowAutoComplete(words)
	//}
}

func (c *ConnectionCtrl) String() string {
	defer config.LogStart("ConnectionCtrl.String", nil)()

	return c.ctx.Database().Name() + "." + c.tableDef.Name
}

func (c *ConnectionCtrl) OnCopyInsert(cols []driver.ColDef, values []interface{}) {
	defer config.LogStart("ConnectionCtrl.OnCopyInsert", nil)()

	sql, err := c.Engine.GetInsertStatement(c.ctx, c.tableDef.Name, cols, values)
	if err != nil {
		c.window.PushStatus(err.Error())
	}

	gtk.ClipboardCopy(sql)
	config.Debugf("insert copied: %s", sql)
}

func (c *ConnectionCtrl) OnUpdateRecord(cols []driver.ColDef, values []interface{}) error {
	defer config.LogStart("ConnectionCtrl.OnUpdateRecord", nil)()

	_, err := c.Engine.UpdateField(c.ctx, c.tableDef.Name, cols, values)
	if err != nil {
		return err
	}

	c.window.PushStatus("Saved")
	return nil
}

func (c *ConnectionCtrl) OnCreateRecord(cols []driver.ColDef, values []interface{}) ([]interface{}, error) {
	defer config.LogStart("ConnectionCtrl.OnCreateRecord", nil)()

	data, err := c.Engine.InsertRecord(c.ctx, c.tableDef.Name, cols, values)
	if err != nil {
		return nil, err
	} else {
		c.window.PushStatus("Inserted")
	}

	return data, nil
}

func (c *ConnectionCtrl) OnExecQuery(value string) {
	defer config.LogStart("ConnectionCtrl.OnExecQuery", nil)()

	columns, data, err := c.Engine.Query(c.ctx, value)
	if err != nil {
		return
	}
	c.scr.UpdateRawData(columns, data)
	c.window.PushStatus("%d rows loaded", len(data))

	/*dml, ddl := c.parseQuery(value)

	for _, query := range dml {
		columns, data, err := c.Engine.Query(c.ctx, query)
		if err != nil {
			config.Error(err)
			c.window.PushStatus("Error: %s", err.Error())
			return
		}
		c.scr.UpdateRawData(columns, data)
		c.window.PushStatus("%d rows loaded", len(data))
	}

	for _, query := range ddl {
		id, affected, err := c.Engine.Execute(c.ctx, query)
		if err != nil {
			config.Error(err)
			c.window.PushStatus("Error: %s", err.Error())
			return
		}
		c.window.PushStatus("Last inserted id: %s Affected rows: %d", id, affected)
	}
	*/
}

func (c *ConnectionCtrl) OnDelete() {
	defer config.LogStart("ConnectionCtrl.OnDelete", nil)()

	newRecord, err := c.scr.SelectedIsNewRecord()
	if err != nil {
		return
	}

	if newRecord {
		c.scr.RemoveSelected()
		c.window.PushStatus("Record removed")
	} else {
		c.scr.ForEachSelected(func(cols []driver.ColDef, values []interface{}) {
			err = c.Engine.DeleteRecord(c.ctx, c.tableDef.Name, cols, values)
			if err != nil {
				config.Error(err, "deleting record")
				return
			}
			c.scr.RemoveSelected()
			c.window.PushStatus("Record deleted")
		})
	}
}

func (c *ConnectionCtrl) ParseValue(col driver.ColDef, value string) (interface{}, error) {
	return c.Engine.ParseValue(c.ctx, col, value)
}

func (c *ConnectionCtrl) OnLoadTable() {
	defer config.LogStart("ConnectionCtrl.OnLoadTable", nil)()

	cancel := c.window.Go(func(ctx context.Context) func() {
		switch c.tableDef.Type {
		case driver.TableTypeDummy:
			// TODO: query is done on the background
			return func() {
				c.OnExecQuery(c.tableDef.Query)
			}
		default:
			def, data, err := c.Engine.FetchTable(
				c.ctx.WithContext(ctx), c.tableDef,
				driver.FetchTableOptions{
					Offset: c.scr.Offset(),
					Limit:  c.scr.PageSize(),
				},
			)
			if err != nil {
				return func() {}
			}

			if c.tableDef.Type == driver.TableTypeRegular {
				return func() {
					err = c.scr.UpdateColumns(def)
					if err != nil {
						config.Error(err)
						return
					}

					err = c.scr.UpdateData(data)
					if err != nil {
						config.Error(err)
					}
				}
			} else {
				columns := []string{}
				for _, d := range def {
					columns = append(columns, d.Name)
				}

				return func() {
					c.scr.UpdateRawData(columns, data)
					if err != nil {
						config.Error(err)
					}
				}
			}
		}
	})

	c.scr.Block(cancel)
	c.updateTitle()
}

func (c *ConnectionCtrl) OnRefresh() {
	defer config.LogStart("ConnectionCtrl.OnRefresh", nil)()

	conditions, err := c.scr.Conditions()
	if err != nil {
		config.Error(err)
		return
	}

	_, data, err := c.Engine.FetchTable(
		c.ctx, c.tableDef,
		driver.FetchTableOptions{
			Offset:     c.scr.Offset(),
			Limit:      c.scr.PageSize(),
			Sort:       c.scr.SortOptions(),
			Conditions: conditions,
		},
	)
	if err != nil {
		return
	}

	err = c.scr.UpdateData(data)
	if err != nil {
		return
	}

	c.window.PushStatus("Table reloaded")
}

func (c *ConnectionCtrl) OnApplyConditions() {
	defer config.LogStart("ConnectionCtrl.OnApplyConditions", nil)()

	c.OnRefresh()
}

func (c *ConnectionCtrl) OnCreate() {
	defer config.LogStart("ConnectionCtrl.OnCreate", nil)()

	newRecord, err := c.scr.SelectedIsNewRecord()
	if err != nil {
		return
	}
	if !newRecord {
		return
	}

	cols, values, err := c.scr.GetRow()
	if err != nil {
		return
	}

	values, err = c.Engine.InsertRecord(c.ctx, c.tableDef.Name, cols, values)
	if err != nil {
		return
	}

	err = c.scr.UpdateRow(values)
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	c.window.PushStatus("Record saved")
}

func (c *ConnectionCtrl) SetTableDef(ctx *sqlengine.Context, tableDef driver.TableDef) (bool, error) {
	defer config.LogStart("ConnectionCtrl.SetTableDef", nil)()

	if c.ctx != nil && c.ctx != ctx {
		return false, nil
	}

	c.tableDef = tableDef
	//c.SetTitle(fmt.Sprintf("%s.%s", c.dbName, tableDef.Name))
	c.OnLoadTable()

	return true, nil
}

func (c *ConnectionCtrl) SetQuery(query string) {
	defer config.LogStart("ConnectionCtrl.SetQuery", nil)()

	c.scr.SetQuery(query)
}

func (c *ConnectionCtrl) Close() {
	c.Engine.Disconnect(c.ctx)
}

func (c *ConnectionCtrl) updateTitle() {
	c.ChangeTitle(c.Title())
}

func (c *ConnectionCtrl) Title() string {
	defer config.LogStart("ConnectionCtrl.Title", nil)()

	switch {
	case c.tableDef.Name != "":
		return fmt.Sprintf("%s.%s.%s", c.conn.Name, c.dbName, c.tableDef.Name)
	case c.dbName != "":
		return fmt.Sprintf("%s.%s", c.conn.Name, c.dbName)
	default:
		return fmt.Sprintf("%s", c.conn.Name)
	}
}

func (c *ConnectionCtrl) OnFileSelected(filepath string) {
	defer config.LogStart("ConnectionCtrl.OnFileSelected", nil)()

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		config.Error("reading file", err)
		return
	}

	c.scr.SetQuery(string(bytes))
}

func (c *ConnectionCtrl) OnSaveQuery(query, path string) {
	defer config.LogStart("ConnectionCtrl.OnSaveQuery", nil)()

	err := ioutil.WriteFile(path, []byte(query), os.FileMode(0666))
	if err != nil {
		c.window.PushStatus("failed to save file: %#v", err)
	}
}

/*
func (c *ConnectionCtrl) parseQuery(src string) (dml []string, ddl []string) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(src, "", "")
	if err != nil {
		config.Error(err)
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
