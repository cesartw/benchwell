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
}

// New return a new Engine
func New(conf *config.Config) *Engine {
	return &Engine{
		config:      conf,
		connections: make([]driver.Connection, 0),
	}
}

// Connect to a database
func (e *Engine) Connect(ctx Context, cfg config.Connection) (Context, error) {
	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := driver.Connect(tmctx, cfg)
	if err != nil {
		return nil, err
	}

	e.connections = append(e.connections, conn)

	c := NewContext(ctx, conn)
	return c, nil
}

// Databases returns a list databases
func (e *Engine) Databases(ctx Context) ([]string, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return nil, errors.New("no connection available")
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	dbNames, err := conn.Databases(tmctx)
	if err != nil {
		return nil, err
	}

	return dbNames, nil
}

// UseDatabase ...
func (e *Engine) UseDatabase(ctx Context, dbName string) (Context, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return ctx, ErrNoConnection
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	dbs, err := conn.Databases(tmctx)
	if err != nil {
		return ctx, err
	}

	var exists bool
	for _, db := range dbs {
		if db == dbName {
			exists = true
			break
		}
	}

	if !exists {
		return ctx, ErrDatabaseNotFound
	}

	db, err := conn.UseDatabase(tmctx, dbName)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(NewContext(nil, conn), ckDatabase, db), nil
}

// Tables ...
func (e *Engine) Tables(ctx Context) ([]driver.TableDef, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return nil, ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return nil, ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.Tables(tmctx)
}

// FetchTable returns table column definition and table data
func (e *Engine) FetchTable(
	ctx Context, tableName string, opts driver.FetchTableOptions,
) (
	[]driver.ColDef, [][]interface{}, error,
) {
	conn := e.connection(ctx)
	if conn == nil {
		return nil, nil, ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return nil, nil, ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.FetchTable(tmctx, tableName, opts)
}

// DeleteRecord ...
func (e *Engine) DeleteRecord(ctx Context, tableName string, defs []driver.ColDef, values []interface{}) error {
	conn := e.connection(ctx)
	if conn == nil {
		return ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.DeleteRecord(tmctx, tableName, defs, values)
}

// UpdateFields ...
func (e *Engine) UpdateFields(
	ctx Context,
	tableName string,
	defs []driver.ColDef,
	values []interface{},
	keycount int,
) (string, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return "", ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.UpdateFields(tmctx, tableName, defs, values, keycount)
}

// UpdateField ...
func (e *Engine) UpdateField(
	ctx Context,
	tableName string,
	defs []driver.ColDef,
	values []interface{},
) (string, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return "", ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.UpdateField(tmctx, tableName, defs, values)
}

// ParseValue ...
func (e *Engine) ParseValue(
	ctx Context,
	def driver.ColDef,
	v string,
) (interface{}, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return "", ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return "", ErrNoConnection
	}

	return db.ParseValue(def, v), nil
}

// UpdateRecord ...
func (e *Engine) UpdateRecord(
	ctx Context,
	tableName string,
	defs []driver.ColDef,
	values, oldValues []interface{},
) (string, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return "", ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.UpdateRecord(tmctx, tableName, defs, values, oldValues)
}

// InsertRecord ...
func (e *Engine) InsertRecord(
	ctx Context,
	tableName string,
	defs []driver.ColDef,
	values []interface{},
) ([]interface{}, error) {
	conn := e.connection(ctx)
	if conn == nil {
		return nil, ErrNoConnection
	}

	db := e.database(ctx)
	if db == nil {
		return nil, ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.InsertRecord(tmctx, tableName, defs, values)
}

// Disconnect ...
func (e *Engine) Disconnect(ctx Context) error {
	conn := e.connection(ctx)
	if conn == nil {
		return ErrNoConnection
	}

	return conn.Disconnect(ctx)
}

func (e *Engine) Query(ctx Context, query string) ([]string, [][]interface{}, error) {
	db := e.database(ctx)
	if db == nil {
		return nil, nil, ErrNoDatabase
	}

	return db.Query(ctx, query)
}

func (e *Engine) Execute(ctx Context, query string) (string, int64, error) {
	db := e.database(ctx)
	if db == nil {
		return "", 0, ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.Execute(tmctx, query)
}

func (e *Engine) GetCreateTable(ctx Context, tableName string) (string, error) {
	db := e.database(ctx)
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.GetCreateTable(tmctx, tableName)
}

func (e *Engine) GetInsertStatement(
	ctx Context,
	tableName string,
	cols []driver.ColDef,
	values []interface{},
) (string, error) {
	db := e.database(ctx)
	if db == nil {
		return "", ErrNoDatabase
	}

	tmctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return db.GetInsertStatement(tmctx, tableName, cols, values)
}

// Dispose ...
func (e *Engine) Dispose() {
	for _, c := range e.connections {
		c.Disconnect(context.Background())
	}
}

// GETTERS
func (e *Engine) connection(ctx Context) driver.Connection {
	connI := ctx.Value(ckConnection)
	if connI == nil {
		return nil
	}

	conn, ok := connI.(driver.Connection)
	if !ok {
		return nil
	}

	return conn
}

func (e *Engine) database(ctx Context) driver.Database {
	dbI := ctx.Value(ckDatabase)
	if dbI == nil {
		return nil
	}

	db, ok := dbI.(driver.Database)
	if !ok {
		return nil
	}

	return db
}

func (e *Engine) Database(ctx Context) driver.Database {
	return e.database(ctx)
}
