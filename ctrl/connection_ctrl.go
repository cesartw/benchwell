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
	mainCtx sqlengine.Context

	dbCtx  map[string]sqlengine.Context
	scr    *gtk.ConnectionScreen
	conn   *config.Connection
	dbName string

	tabs []*TableCtrl
}

func (c ConnectionCtrl) Init(
	ctx sqlengine.Context,
	p *WindowTabCtrl,
	conn *config.Connection,
) (*ConnectionCtrl, error) {
	c.dbCtx = map[string]sqlengine.Context{}
	c.WindowTabCtrl = p
	c.mainCtx = ctx
	c.conn = conn

	dbNames, err := c.engine.Databases(c.mainCtx)
	if err != nil {
		return nil, err
	}

	c.scr, err = c.app.NewConnectionScreen()
	if err != nil {
		return nil, err
	}

	c.scr.SetDatabases(dbNames)

	c.scr.ShowAll()

	c.scr.OnDatabaseSelected(c.onDatabaseSelected)
	c.scr.OnTableSelected(c.onTableSelected)

	if c.conn.Database != "" {
		c.scr.SetActiveDatabase(c.conn.Database)
	}

	c.scr.OnSchemaMenu(c.onSchemaMenu)
	c.scr.OnNewTabMenu(c.onNewTabMenu)

	return &c, nil
}

func (c *ConnectionCtrl) Close() bool {
	return c.scr.Close()
}

func (c *ConnectionCtrl) AddEmptyTab() error {
	return c.AddTab(driver.TableDef{})
}

func (c *ConnectionCtrl) AddTab(tableDef driver.TableDef) error {
	tab, err := TableCtrl{}.init(c.dbCtx[c.dbName], TableCtrlOpts{
		Parent:       c,
		TableDef:     tableDef,
		OnTabRemoved: c.onTabRemove,
	})
	if err != nil {
		return err
	}

	if tableDef.IsZero() {
		tab.connectionTab.SetTitle("New")
	}

	c.tabs = append(c.tabs, tab)
	return c.scr.AddTab(tab.connectionTab, true)
}

func (c *ConnectionCtrl) onTabRemove(ctrl *TableCtrl) {
	defer c.disconnect()

	for i, tabCtrl := range c.tabs {
		if tabCtrl == ctrl {
			c.tabs = append(c.tabs[:i], c.tabs[i+1:]...)
			break
		}
	}
}

func (c *ConnectionCtrl) UpdateOrAddTab(tableDef driver.TableDef) error {
	if len(c.tabs) == 0 || c.tabs[c.scr.CurrentTabIndex()].ctx != c.dbCtx[c.dbName] {
		return c.AddTab(tableDef)
	}

	c.tabs[c.scr.CurrentTabIndex()].SetTableDef(c.dbCtx[c.dbName], tableDef)
	return nil
}

func (c *ConnectionCtrl) onDatabaseSelected() {
	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.window.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}

	if c.dbCtx[dbName] == nil {
		c.dbCtx[dbName], err = c.engine.UseDatabase(c.mainCtx, dbName)
		if err != nil {
			c.window.PushStatus("Error selecting database: `%s`", err.Error())
			return
		}
	}

	tables, err := c.engine.Tables(c.dbCtx[dbName])
	if err != nil {
		c.window.PushStatus("Error getting tables: `%s`", err.Error())
		return
	}

	c.scr.SetTables(tables)
	c.dbName = dbName
}

func (c *ConnectionCtrl) onTableSelected() {
	defer c.disconnect()

	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	c.UpdateOrAddTab(tableDef)
}

func (c *ConnectionCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectionCtrl) onSchemaMenu() {
	tableName, ok := c.scr.SelectedTable()
	if !ok {
		return
	}

	schema, err := c.engine.GetCreateTable(c.dbCtx[c.dbName], tableName)
	if err != nil {
		config.Env.Log.Error(err, "getting table schema")
	}

	c.scr.ShowTableSchemaModal(tableName, schema)
}

func (c *ConnectionCtrl) onNewTabMenu() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	c.AddTab(tableDef)
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

		c.engine.Disconnect(ctx)
		delete(c.dbCtx, dbName)
	}
}
