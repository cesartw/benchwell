package ctrl

import (
	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/gtk"
	"bitbucket.org/goreorto/benchwell/sqlengine"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type ConnectionCtrl struct {
	*DbTabCtrl

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
	p *DbTabCtrl,
	conn *config.Connection,
) (*ConnectionCtrl, error) {
	defer config.LogStart("ConnectionCtrl.Init", nil)()

	c.dbCtx = map[string]*sqlengine.Context{}
	c.DbTabCtrl = p
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

func (c *ConnectionCtrl) Title() string {
	defer config.LogStart("ConnectionCtrl.Title", nil)()

	if c.conn.Name != "" {
		return c.conn.Name
	}

	return c.conn.Host
}

func (c *ConnectionCtrl) Content() ggtk.IWidget {
	defer config.LogStart("ConnectionCtrl.Content", nil)()

	return c.scr
}

func (c *ConnectionCtrl) OnCopyLog() {
	defer config.LogStart("ConnectionCtrl.OnCopyLog", nil)()
}

func (c *ConnectionCtrl) Close() bool {
	defer config.LogStart("ConnectionCtrl.Close", nil)()

	return c.scr.Close()
}

func (c *ConnectionCtrl) FullClose() {
	defer config.LogStart("ConnectionCtrl.FullClose", nil)()

	c.scr.CloseAll()
}

func (c *ConnectionCtrl) AddEmptyTab() error {
	defer config.LogStart("ConnectionCtrl.AddEmptyTab", nil)()

	if _, ok := c.scr.ActiveDatabase(); ok {
		return c.AddTab(driver.TableDef{})
	}
	return nil
}

func (c *ConnectionCtrl) SetFileText(s string) {
	defer config.LogStart("ConnectionCtrl.SetFileText", nil)()

	if len(c.tabs) == 0 {
		return
	}
	tabCtrl := c.tabs[c.scr.CurrentTabIndex()]
	tabCtrl.SetQuery(tabCtrl.ctx, s)
}

func (c *ConnectionCtrl) AddTab(tableDef driver.TableDef) error {
	defer config.LogStart("ConnectionCtrl.AddTab", nil)()

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
	defer config.LogStart("ConnectionCtrl.UpdateOrAddTab", nil)()

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
	defer config.LogStart("ConnectionCtrl.OnDatabaseSelected", nil)()

	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.window.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}
	config.Debug("dbName", dbName)
	config.Debug("dbs", c.dbCtx)

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
}

func (c *ConnectionCtrl) OnTableSelected() {
	defer config.LogStart("ConnectionCtrl.OnTableSelected", nil)()

	defer c.disconnect()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	if c.scr.CtrlMod() {
		c.AddTab(tableDef)
	} else {
		c.UpdateOrAddTab(tableDef)
	}
}

func (c *ConnectionCtrl) OnEditTable() {
	defer config.LogStart("ConnectionCtrl.OnEditTable", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		return
	}

	if tableDef.Type == driver.TableTypeDummy {
		ok, err := c.scr.SetQuery(c.dbCtx[c.dbName], tableDef.Query)
		if err != nil {
			config.Error(err)
			return
		}

		if ok {
			return
		}

		// add tab for connection and try again
		err = c.AddTab(tableDef)
		if err != nil {
			config.Error(err)
			return
		}

		c.scr.SetQuery(c.dbCtx[c.dbName], tableDef.Query)
	}
}

func (c *ConnectionCtrl) OnTruncateTable() {
	defer config.LogStart("ConnectionCtrl.OnTruncateTable", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	c.Engine.TruncateTable(c.dbCtx[c.dbName], tableDef)
}

func (c *ConnectionCtrl) OnDeleteTable() {
	defer config.LogStart("ConnectionCtrl.OnDeleteTable", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	err := c.Engine.DeleteTable(c.dbCtx[c.dbName], tableDef)
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

	sql, err := c.Engine.GetSelectStatement(c.dbCtx[c.dbName], tableDef)
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

	schema, err := c.Engine.GetCreateTable(c.dbCtx[c.dbName], tableName)
	if err != nil {
		config.Error(err, "getting table schema")
	}

	c.scr.ShowTableSchemaModal(tableName, schema)
}

func (c *ConnectionCtrl) OnRefreshMenu() {
	defer config.LogStart("ConnectionCtrl.OnRefreshMenu", nil)()

	if dbName, ok := c.scr.ActiveDatabase(); ok {
		c.dbCtx[dbName].CacheTable = nil
	}

	c.OnDatabaseSelected()
}

func (c *ConnectionCtrl) OnNewTabMenu() {
	defer config.LogStart("ConnectionCtrl.OnNewTabMenu", nil)()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Debug("no table selected. odd!")
		return
	}

	c.AddTab(tableDef)
}

func (c *ConnectionCtrl) OnSaveFav(name, query string) {
	defer config.LogStart("ConnectionCtrl.OnSaveFav", nil)()

	config.SaveQuery(&config.Query{
		Name:         name,
		Query:        query,
		ConnectionID: c.conn.ID,
	})
}

func (c *ConnectionCtrl) Screen() interface{} {
	defer config.LogStart("ConnectionCtrl.Screen", nil)()

	return c.scr
}

func (c *ConnectionCtrl) OnTabRemove(ctrl *TableCtrl) {
	defer config.LogStart("ConnectionCtrl.OnTabRemove", nil)()

	defer c.disconnect()

	for i, tabCtrl := range c.tabs {
		if tabCtrl == ctrl {
			c.tabs = append(c.tabs[:i], c.tabs[i+1:]...)
			break
		}
	}
}

func (c *ConnectionCtrl) disconnect() {
	return
	defer config.LogStart("ConnectionCtrl.disconnect", nil)()

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

		config.Debug("disconnecting: ", c.Engine.Database(ctx).Name())

		c.Engine.Disconnect(ctx)
		delete(c.dbCtx, dbName)
		config.Debug("dbCtx ", c.dbCtx)
	}
}
