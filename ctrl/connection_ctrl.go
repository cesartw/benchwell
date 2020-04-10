package ctrl

import (
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	ggtk "github.com/gotk3/gotk3/gtk"
)

type ConnectionCtrl struct {
	*TabCtrl

	ctx  sqlengine.Context
	scr  *gtk.ConnectionScreen
	conn *config.Connection

	tabs []*TableCtrl
}

func (c ConnectionCtrl) init(ctx sqlengine.Context, p *TabCtrl, conn *config.Connection) (*ConnectionCtrl, error) {
	c.TabCtrl = p
	c.ctx = ctx
	c.conn = conn

	dbNames, err := c.engine.Databases(c.ctx)
	if err != nil {
		return nil, err
	}

	c.scr, err = c.factory.NewConnectionScreen()
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

	return &c, nil
}

func (c *ConnectionCtrl) onDatabaseSelected() {
	var err error
	dbName, ok := c.scr.ActiveDatabase()
	if !ok {
		c.factory.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}
	c.ctx, err = c.engine.UseDatabase(c.ctx, dbName)

	tables, err := c.engine.Tables(c.ctx)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.scr.SetTables(tables)
}

func (c *ConnectionCtrl) onTableSelected() {
	tableName, ok := c.scr.ActiveTable()
	if !ok {
		config.Env.Log.Debug("no table selected. odd!")
		return
	}
	tab, err := TableCtrl{}.init(c.ctx, c, tableName)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	c.tabs = append(c.tabs, tab)
	c.scr.AddTab(tableName, tab.Screen().(ggtk.IWidget), true)
	tab.OnConnect()
}

func (c *ConnectionCtrl) Screen() interface{} {
	return c.scr
}
