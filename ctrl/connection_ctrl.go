package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type ConnectionCtrl struct {
	*ConnectionTabCtrl

	ctx  sqlengine.Context
	scr  *gtk.ConnectionScreen
	conn *config.Connection

	tabs []*TableCtrl
}

func (c ConnectionCtrl) Init(
	ctx sqlengine.Context,
	p *ConnectionTabCtrl,
	conn *config.Connection,
) (*ConnectionCtrl, error) {
	c.ConnectionTabCtrl = p
	c.ctx = ctx
	c.conn = conn

	dbNames, err := c.engine.Databases(c.ctx)
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
		c.onDatabaseSelected()
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
	tab, err := TableCtrl{}.init(c.ctx, TableCtrlOpts{
		Parent:       c,
		TableDef:     tableDef,
		OnTabRemoved: c.onTabRemove,
	})
	if err != nil {
		return err
	}

	c.tabs = append(c.tabs, tab)
	return c.scr.AddTab(tab.connectionTab, true)
}

func (c *ConnectionCtrl) onTabRemove(ctrl *TableCtrl) {
	for i, tabCtrl := range c.tabs {
		if tabCtrl == ctrl {
			c.tabs = append(c.tabs[:i], c.tabs[i+1:]...)
			break
		}
	}
}

func (c *ConnectionCtrl) UpdateOrAddTab(tableDef driver.TableDef) error {
	tab, err := TableCtrl{}.init(c.ctx, TableCtrlOpts{Parent: c, TableDef: tableDef})
	if err != nil {
		return err
	}

	return c.scr.AddTab(tab.connectionTab, true)
}

func (c *ConnectionCtrl) onDatabaseSelected() {
	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.window.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}
	c.ctx, err = c.engine.UseDatabase(c.ctx, dbName)
	if err != nil {
		c.window.PushStatus("Error selecting database: `%s`", err.Error())
		return
	}

	tables, err := c.engine.Tables(c.ctx)
	if err != nil {
		c.window.PushStatus("Error getting tables: `%s`", err.Error())
		return
	}

	c.scr.SetTables(tables)
}

func (c *ConnectionCtrl) onTableSelected() {
	tableDef, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}

	c.AddTab(tableDef)
}

func (c *ConnectionCtrl) Screen() interface{} {
	return c.scr
}

func (c *ConnectionCtrl) onSchemaMenu() {
	tableName, ok := c.scr.SelectedTable()
	if !ok {
		return
	}

	schema, err := c.engine.GetCreateTable(c.ctx, tableName)
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
