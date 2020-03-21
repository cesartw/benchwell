package ctrl

import (
	"context"
	"errors"

	ggtk "github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
)

type Options struct {
	Config  *config.Config
	Log     *logrus.Logger
	Factory *gtk.Factory
	Engine  *sqlengine.Engine
}

func (o Options) Valid() error {
	if o.Config == nil {
		return errors.New("Fonfig is required")
	}

	if o.Factory == nil {
		return errors.New("Factory is required")
	}

	if o.Engine == nil {
		return errors.New("Engine is required")
	}

	return nil
}

// Controller controls a sqlhero tab
type Ctrl struct {
	ctx    sqlengine.Context
	config *config.Config
	log    *logrus.Logger
	conn   *config.Connection

	factory *gtk.Factory
	engine  *sqlengine.Engine

	// TODO: interface it up
	tabs []*tableCtrl

	// ui
	scrConnect    *gtk.ConnectScreen
	scrConnection *gtk.ConnectionScreen
}

func New(opts Options) (*Ctrl, error) {
	if err := opts.Valid(); err != nil {
		return nil, err
	}
	c := &Ctrl{
		engine:  opts.Engine,
		factory: opts.Factory,
		config:  opts.Config,
		log:     opts.Log,
	}

	if c.log == nil {
		c.log = logrus.New()
	}

	return c, nil
}

func (c *Ctrl) OnActivate() {
	var err error
	c.scrConnect, err = c.factory.NewConnectScreen()
	if err != nil {
		c.log.Error(err)
		return
	}

	c.factory.Add(c.scrConnect)

	c.scrConnect.SetConnections(c.config.Connection)

	c.scrConnect.OnConnect(c.onConnect)
	c.scrConnect.OnTest(c.onTest)

	c.scrConnect.OnSave(c.onSave)

	c.factory.Show()
	c.factory.PushStatus("Ready")
}

func (c *Ctrl) onConnect() {
	var err error

	c.conn = c.scrConnect.ActiveConnection()
	c.ctx, err = c.engine.Connect(sqlengine.Context(context.TODO()), c.conn.GetDSN())
	if err != nil {
		c.log.Error(err)
		return
	}

	c.factory.PushStatus("Connected to `%s`", c.conn.Host)

	c.scrConnection, err = c.factory.NewConnectionScreen(c.ctx)
	if err != nil {
		c.log.Error(err)
		return
	}

	dbNames, err := c.engine.Databases(c.ctx)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.scrConnection.SetDatabases(dbNames)

	c.scrConnection.ShowAll()

	c.scrConnection.OnDatabaseSelected(c.onDatabaseSelected)
	c.scrConnection.OnTableSelected(c.onTableSelected)

	c.factory.Remove(c.scrConnect)
	c.factory.Add(c.scrConnection)

	if c.conn.Database != "" {
		c.scrConnection.SetActiveDatabase(c.conn.Database)
		c.onDatabaseSelected()
	}
}

func (c *Ctrl) onDatabaseSelected() {
	var err error
	dbName, ok := c.scrConnection.ActiveDatabase()
	if !ok {
		c.factory.PushStatus("Database `%s` not found", c.conn.Database)
		return
	}
	c.ctx, err = c.engine.UseDatabase(c.ctx, dbName)

	tables, err := c.engine.Tables(c.ctx)
	if err != nil {
		c.log.Error(err)
		return
	}

	c.scrConnection.SetTables(tables)
}

func (c *Ctrl) onTableSelected() {
	rv, err := gtk.NewResultView()
	if err != nil {
		c.log.Error(err)
		return
	}
	rv.ShowAll()

	tableName, ok := c.scrConnection.ActiveTable()
	if !ok {
		c.log.Info("no table selected. odd!")
		return
	}

	ctl := newTableCtrl(c.ctx, c, tableName, rv)

	c.tabs = append(c.tabs, ctl)
	// TODO: meh. find better way to initialize
	ctl.OnConnect()

	c.AddTab(tableName, rv, true)
}

func (c *Ctrl) onTest() {
	conn := c.scrConnect.ActiveConnection()
	ctx, err := c.engine.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
	if err != nil {
		c.log.Error(err)
		return
	}

	c.factory.PushStatus("Connection to `%s` was successful", conn.Host)
	c.engine.Disconnect(ctx)
}

func (c *Ctrl) onSave() {
	c.factory.PushStatus("Saved")
}

func (c *Ctrl) AddTab(title string, w ggtk.IWidget, switchNow bool) {
	c.scrConnection.AddTab(title, w)
}

// tableCtrl manages table result screen
type tableCtrl struct {
	parent    *Ctrl
	ctx       sqlengine.Context
	tableName string

	// ui
	resultView *gtk.ResultView
}

func newTableCtrl(ctx sqlengine.Context, parent *Ctrl, tableName string, rv *gtk.ResultView) *tableCtrl {
	return &tableCtrl{ctx: ctx, parent: parent, tableName: tableName, resultView: rv}
}

func (tc *tableCtrl) OnConnect() {
	def, data, err := tc.parent.engine.FetchTable(tc.ctx, tc.tableName, 0, 40)
	if err != nil {
		tc.parent.log.Error(err)
		return
	}

	err = tc.resultView.UpdateData(def, data)
	if err != nil {
		tc.parent.log.Error(err)
	}
}
