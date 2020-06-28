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
func (e *Engine) Connect(ctx context.Context, cfg config.Connection) (*Context, error) {
	timeout := 2 * time.Second
	tmctx, cancel := context.WithTimeout(ctx, timeout)
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

	dbNames, err := conn.Databases(c.Context())
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

	dbs, err := conn.Databases(c.Context())
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

	db, err := conn.UseDatabase(c.Context(), dbName)
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

	tables, err := db.Tables(c.Context())
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

	return db.FetchTable(c.Context(), table.Name, opts)
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

	return db.DeleteRecord(c.Context(), tableName, defs, values)
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

	return db.UpdateFields(c.Context(), tableName, defs, values, keycount)
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

	return db.UpdateField(c.Context(), tableName, defs, values)
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

	return db.UpdateRecord(c.Context(), tableName, defs, values, oldValues)
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

	return db.InsertRecord(c.Context(), tableName, defs, values)
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

	return db.Query(c.Context(), query)
}

func (e *Engine) Execute(c *Context, query string) (string, int64, error) {
	db := c.Database()
	if db == nil {
		return "", 0, ErrNoDatabase
	}

	return db.Execute(c.Context(), query)
}

func (e *Engine) GetCreateTable(c *Context, tableName string) (string, error) {
	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	return db.GetCreateTable(c.Context(), tableName)
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

	return db.GetInsertStatement(c.Context(), tableName, cols, values)
}

func (e *Engine) GetSelectStatement(
	c *Context,
	table driver.TableDef,
) (string, error) {
	db := c.Database()
	if db == nil {
		return "", ErrNoDatabase
	}

	return db.GetSelectStatement(c.Context(), table)
}

func (e *Engine) TruncateTable(
	c *Context,
	table driver.TableDef,
) error {
	db := c.Database()
	if db == nil {
		return ErrNoDatabase
	}

	return db.TruncateTable(c.Context(), table)
}

func (e *Engine) DeleteTable(
	c *Context,
	table driver.TableDef,
) error {
	db := c.Database()
	if db == nil {
		return ErrNoDatabase
	}

	return db.DeleteTable(c.Context(), table)
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

func (e *Engine) timeoutErr(err error, timeout time.Duration) error {
	if err == nil {
		return nil
	}

	if err.Error() == "context timeout" {
		return errors.New("timeout after " + timeout.String())
	}

	return err
}
