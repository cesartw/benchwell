package sqlengine

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// Errors
var (
	ErrNoConnection     = errors.New("no connection available")
	ErrNoDatabase       = errors.New("no database selected")
	ErrDatabaseNotFound = errors.New("database not found")
)

// Engine orchestrate multiple Connection
type Engine struct {
	config *config.Config

	connections []driver.Connection
	Logger      func(string)
}

// New return a new Engine
func New(conf *config.Config) *Engine {
	return &Engine{
		config:      conf,
		connections: make([]driver.Connection, 0),
	}
}

func (e *Engine) runWithTimeout(timeout time.Duration, f func(context.Context)) error {
	tmctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	go func() {
		f(tmctx)
		done <- struct{}{}
	}()

	select {
	case <-tmctx.Done():
		return errors.New("timeout after " + timeout.String())
	case <-done:
		return nil
	}
}

// Connect to a database
func (e *Engine) ConnectWithTimeout(cfg config.Connection) (*Context, error) {
	timeout := 2 * time.Second

	var (
		conn driver.Connection
		err  error
	)

	timeoutErr := e.runWithTimeout(timeout, func(ctx context.Context) {
		conn, err = driver.Connect(ctx, cfg)
	})
	if timeoutErr != nil {
		return nil, timeoutErr
	}
	if err != nil {
		return nil, err
	}

	e.connections = append(e.connections, conn)
	return NewContext(conn, nil), nil
}

// Connect to a database
func (e *Engine) Connect(cfg config.Connection) (*Context, error) {
	timeout := 2 * time.Second
	tmctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := driver.Connect(tmctx, cfg)
	if err != nil {
		return nil, e.timeoutErr(err, timeout)
	}
	e.connections = append(e.connections, conn)
	return NewContext(conn, nil), nil
}

// Databases returns a list databases
func (e *Engine) Databases(c *Context) ([]string, error) {
	conn := c.Connection()
	if conn == nil {
		return nil, errors.New("no connection available")
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	dbNames, err := conn.Databases(tmctx)
	if err != nil {
		return nil, err
	}

	return dbNames, nil
}

// UseDatabase ...
func (e *Engine) UseDatabase(c *Context, dbName string) (*Context, error) {
	conn := c.Connection()
	if conn == nil {
		return c, ErrNoConnection
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	dbs, err := conn.Databases(tmctx)
	if err != nil {
		return c, err
	}

	var exists bool
	for _, db := range dbs {
		if db == dbName {
			exists = true
			break
		}
	}

	if !exists {
		return c, ErrDatabaseNotFound
	}

	db, err := conn.UseDatabase(tmctx, dbName)
	if err != nil {
		return c, err
	}

	return NewContext(conn, db), nil
}

// Tables ...
func (e *Engine) Tables(c *Context) ([]driver.TableDef, error) {
	if c.CacheTable != nil {
		return c.CacheTable, nil
	}

	conn := c.Connection()
	if conn == nil {
		return nil, ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return nil, ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	tables, err := db.Tables(tmctx)
	if err != nil {
		return nil, err
	}
	c.CacheTable = tables

	return tables, err
}

// FetchTable returns table column definition and table data
func (e *Engine) FetchTable(
	c *Context,
	table driver.TableDef,
	opts driver.FetchTableOptions,
) (
	[]driver.ColDef, [][]interface{}, error,
) {
	if table.Type == driver.TableTypeDummy {
		return nil, nil, errors.New("not a table")
	}

	conn := c.Connection()
	if conn == nil {
		return nil, nil, ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return nil, nil, ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.FetchTable(tmctx, table.Name, opts)
}

// DeleteRecord ...
func (e *Engine) DeleteRecord(c *Context, tableName string, defs []driver.ColDef, values []interface{}) error {
	conn := c.Connection()
	if conn == nil {
		return ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.DeleteRecord(tmctx, tableName, defs, values)
}

// UpdateFields ...
func (e *Engine) UpdateFields(
	c *Context,
	tableName string,
	defs []driver.ColDef,
	values []interface{},
	keycount int,
) (string, error) {
	conn := c.Connection()
	if conn == nil {
		return "", ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.UpdateFields(tmctx, tableName, defs, values, keycount)
}

// UpdateField ...
func (e *Engine) UpdateField(
	c *Context,
	tableName string,
	defs []driver.ColDef,
	values []interface{},
) (string, error) {
	conn := c.Connection()
	if conn == nil {
		return "", ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.UpdateField(tmctx, tableName, defs, values)
}

// ParseValue ...
func (e *Engine) ParseValue(
	c *Context,
	def driver.ColDef,
	v string,
) (interface{}, error) {
	conn := c.Connection()
	if conn == nil {
		return "", ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return "", ErrNoConnection
	}

	return db.ParseValue(def, v), nil
}

// UpdateRecord ...
func (e *Engine) UpdateRecord(
	c *Context,
	tableName string,
	defs []driver.ColDef,
	values, oldValues []interface{},
) (string, error) {
	conn := c.Connection()
	if conn == nil {
		return "", ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.UpdateRecord(tmctx, tableName, defs, values, oldValues)
}

// InsertRecord ...
func (e *Engine) InsertRecord(
	c *Context,
	tableName string,
	defs []driver.ColDef,
	values []interface{},
) ([]interface{}, error) {
	conn := c.Connection()
	if conn == nil {
		return nil, ErrNoConnection
	}

	db := c.Database()
	if db == nil {
		return nil, ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.InsertRecord(tmctx, tableName, defs, values)
}

// Disconnect ...
func (e *Engine) Disconnect(c *Context) error {
	conn := c.Connection()
	if conn == nil {
		return ErrNoConnection
	}

	return conn.Disconnect(context.Background())
}

func (e *Engine) Query(c *Context, query string) ([]string, [][]interface{}, error) {
	db := c.Database()
	if db == nil {
		return nil, nil, ErrNoDatabase
	}
	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.Query(tmctx, query)
}

func (e *Engine) Execute(c *Context, query string) (string, int64, error) {
	db := c.Database()
	if db == nil {
		return "", 0, ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.Execute(tmctx, query)
}

func (e *Engine) GetCreateTable(c *Context, tableName string) (string, error) {
	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.GetCreateTable(tmctx, tableName)
}

func (e *Engine) GetInsertStatement(
	c *Context,
	tableName string,
	cols []driver.ColDef,
	values []interface{},
) (string, error) {
	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.GetInsertStatement(tmctx, tableName, cols, values)
}

func (e *Engine) GetSelectStatement(
	c *Context,
	table driver.TableDef,
) (string, error) {
	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.GetSelectStatement(tmctx, table)
}

func (e *Engine) TruncateTable(
	c *Context,
	table driver.TableDef,
) error {
	db := c.Database()
	if db == nil {
		return ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.TruncateTable(tmctx, table)
}

func (e *Engine) DeleteTable(
	c *Context,
	table driver.TableDef,
) error {
	db := c.Database()
	if db == nil {
		return ErrNoDatabase
	}

	tmctx, cancel := prepereCtx(c, time.Minute)
	defer cancel()

	return db.DeleteTable(tmctx, table)
}

// Dispose ...
func (e *Engine) Dispose() {
	for _, c := range e.connections {
		c.Disconnect(context.Background())
	}
}

// GETTERS

func (e *Engine) Database(c *Context) driver.Database {
	return c.Database()
}

func prepereCtx(c *Context, d time.Duration) (context.Context, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), d)

	return driver.SetLogger(ctx, c.Logger), cancel
}

func (e *Engine) timeoutErr(err error, timeout time.Duration) error {
	if err == nil {
		return nil
	}

	if err.Error() == "context timeout" {
		return errors.New("timeout after " + timeout.String())
	}

	return err
}
