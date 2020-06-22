package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type ConnectionCtrl struct {
	*WindowTabCtrl

	// db-less connection
	mainCtx *sqlengine.Context

	dbCtx  map[string]*sqlengine.Context
	scr    *gtk.ConnectionScreen
	conn   *config.Connection
	dbName string

	tabs []*TableCtrl
}

func (c ConnectionCtrl) Init(
	ctx *sqlengine.Context,
	p *WindowTabCtrl,
	conn *config.Connection,
) (*ConnectionCtrl, error) {
	c.dbCtx = map[string]*sqlengine.Context{}
	c.WindowTabCtrl = p
	c.mainCtx = ctx
	c.conn = conn

	dbNames, err := c.Engine.Databases(c.mainCtx)
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
		return &c, c.AddEmptyTab()
	}

	return &c, nil
}

func (c *ConnectionCtrl) OnCopyLog() {
}

func (c *ConnectionCtrl) Close() bool {
	return c.scr.Close()
}

func (c *ConnectionCtrl) AddEmptyTab() error {
	if _, ok := c.scr.ActiveDatabase(); ok {
		return c.AddTab(driver.TableDef{})
	}
	return nil
}

func (c *ConnectionCtrl) SetFileText(s string) {
	if len(c.tabs) == 0 {
		return
	}
	tabCtrl := c.tabs[c.scr.CurrentTabIndex()]
	tabCtrl.SetQuery(tabCtrl.ctx, s)
}

func (c *ConnectionCtrl) AddTab(tableDef driver.TableDef) error {
	// TODO: control doesn't know it's a tab. good or bad?
	tab, err := TableCtrl{}.Init(c.dbCtx[c.dbName], TableCtrlOpts{
		Parent:   c,
		TableDef: tableDef,
		//OnTabRemoved: c.onTabRemove, // change to c.OnTabRemoved
	})
	if err != nil {
		return err
	}

	c.tabs = append(c.tabs, tab)
	return c.scr.AddTab(tab.connectionTab, tab, true)
}

func (c *ConnectionCtrl) UpdateOrAddTab(tableDef driver.TableDef) error {
	ok, err := c.scr.SetTableDef(c.dbCtx[c.dbName], tableDef)
	if err != nil {
		return err
	}

	if !ok {
		return c.AddTab(tableDef)
	}

	return nil
}

func (c *ConnectionCtrl) OnDatabaseSelected() {
	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.window.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}

	if c.dbCtx[dbName] == nil {
		c.dbCtx[dbName], err = c.Engine.UseDatabase(c.mainCtx, dbName)
		if err != nil {
			c.window.PushStatus("Error selecting database: `%s`", err.Error())
			return
		}
		c.dbCtx[dbName].Logger = c.scr.Log
	}

	tables, err := c.Engine.Tables(c.dbCtx[dbName])
	if err != nil {
		c.window.PushStatus("Error getting tables: `%s`", err.Error())
		return
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
}

func (c *ConnectionCtrl) OnTableSelected() {
	defer c.disconnect()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	if c.scr.CtrlMod() {
		c.AddTab(tableDef)
	} else {
		c.UpdateOrAddTab(tableDef)
	}
}

func (c *ConnectionCtrl) OnEditTable() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		return
	}

	if tableDef.Type == driver.TableTypeDummy {
		ok, err := c.scr.SetQuery(c.dbCtx[c.dbName], tableDef.Query)
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		if ok {
			return
		}

		// add tab for connection and try again
		err = c.AddTab(tableDef)
		if err != nil {
			config.Env.Log.Error(err)
			return
		}

		c.scr.SetQuery(c.dbCtx[c.dbName], tableDef.Query)
	}
}

func (c *ConnectionCtrl) OnTruncateTable() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	c.Engine.TruncateTable(c.dbCtx[c.dbName], tableDef)
}

func (c *ConnectionCtrl) OnDeleteTable() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	c.Engine.DeleteTable(c.dbCtx[c.dbName], tableDef)
}

func (c *ConnectionCtrl) OnCopySelect() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	sql, err := c.Engine.GetSelectStatement(c.dbCtx[c.dbName], tableDef)
	if err != nil {
		c.window.PushStatus(err.Error())
		return
	}

	gtk.ClipboardCopy(sql)
	config.Env.Log.Debugf("select copied: %s", sql)
}

func (c *ConnectionCtrl) OnSchemaMenu() {
	tableName, ok := c.scr.SelectedTable()
	if !ok {
		return
	}

	schema, err := c.Engine.GetCreateTable(c.dbCtx[c.dbName], tableName)
	if err != nil {
		config.Env.Log.Error(err, "getting table schema")
	}

	c.scr.ShowTableSchemaModal(tableName, schema)
}

func (c *ConnectionCtrl) OnRefreshMenu() {
	if dbName, ok := c.scr.ActiveDatabase(); ok {
		c.dbCtx[dbName].CacheTable = nil
	}

	c.OnDatabaseSelected()
}

func (c *ConnectionCtrl) OnNewTabMenu() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	c.AddTab(tableDef)
}

func (c *ConnectionCtrl) OnSaveFav(name, query string) {
	c.conn.Queries = append(c.conn.Queries, config.Query{
		Name:  name,
		Query: query,
	})
	config.Env.Save(c.window.ApplicationWindow)
}

func (c *ConnectionCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectionCtrl) OnTabRemove(ctrl *TableCtrl) {
	defer c.disconnect()

	for i, tabCtrl := range c.tabs {
		if tabCtrl == ctrl {
			c.tabs = append(c.tabs[:i], c.tabs[i+1:]...)
			break
		}
	}
}

func (c *ConnectionCtrl) disconnect() {
	if len(c.tabs) == 0 {
		return
	}
NEXT:
	for dbName, ctx := range c.dbCtx {
		for _, tab := range c.tabs {
			if tab.ctx == ctx {
				continue NEXT
			}
		}

		// db dropdown is showing the tables
		selectedDB, _ := c.scr.ActiveDatabase()
		if dbName == selectedDB {
			continue
		}

		config.Env.Log.Debug("disconnecting: ", c.Engine.Database(ctx).Name())

		c.Engine.Disconnect(ctx)
		delete(c.dbCtx, dbName)
	}
}
