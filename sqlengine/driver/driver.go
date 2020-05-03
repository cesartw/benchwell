package driver

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

const NULL_PATTERN = "<NULL>"

type TYPE uint

const (
	TYPE_BOOLEAN TYPE = iota
	TYPE_STRING
	TYPE_FLOAT
	TYPE_INT
	TYPE_DATE
	TYPE_LIST
)

var (
	drivers  map[string]Driver
	driverMU sync.Mutex
)

func init() {
	drivers = make(map[string]Driver)
}

// RegisterDriver registers a Driver
func RegisterDriver(name string, d Driver) {
	driverMU.Lock()
	defer driverMU.Unlock()

	if _, ok := drivers[name]; ok {
		panic(fmt.Sprintf("%s already registrated", name))
	}

	drivers[name] = d
}

// Connect returns a Connection
// TODO: dsn parsing too optimistic
func Connect(ctx context.Context, dsn string) (Connection, error) {
	colonS := strings.Split(dsn, ":")

	driverMU.Lock()
	d, ok := drivers[colonS[0]]
	driverMU.Unlock()

	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", colonS[0])
	}

	// kind hacky
	return d.Connect(ctx, strings.TrimPrefix(dsn, colonS[0]+"://"))
}

// Driver is a database implementation for SQLHero
type Driver interface {
	Connect(context.Context, string) (Connection, error)
	// Sanitize(string) string
}

// Connection ...
type Connection interface {
	UseDatabase(context.Context, string) error
	Databases(context.Context) ([]Database, error)
	Disconnect(context.Context) error
	Reconnect(context.Context) error
	LastError() error
}

// ColDef describe a column
type ColDef struct {
	Name               string
	PK, FK             bool
	Precision          int
	Unsigned, Nullable bool
	Type               TYPE
	Values             []string
}

func (c ColDef) String() string {
	return c.Name
}

type SortDirection int

const (
	SortDirectionAsc SortDirection = iota
	SortDirectionDesc
)

type SortOption struct {
	Column    ColDef
	Direction SortDirection
}
type FetchTableOptions struct {
	Offset, Limit int64
	Sort          []SortOption
}

// Database ...
type Database interface {
	Tables(context.Context) ([]string, error)
	TableDefinition(ctx context.Context, tableName string) ([]ColDef, error)
	FetchTable(ctx context.Context, tableName string, opts FetchTableOptions) ([]ColDef, [][]interface{}, error)
	DeleteRecord(ctx context.Context, tableName string, defs []ColDef, values []interface{}) error
	UpdateRecord(ctx context.Context, tableName string, cols []ColDef, values, oldValues []interface{}) (string, error)
	UpdateField(ctx context.Context, tableName string, cols []ColDef, values []interface{}) (string, error)
	UpdateFields(ctx context.Context, tableName string, cols []ColDef, values []interface{}, keycount int) (string, error)
	InsertRecord(ctx context.Context, tableName string, cols []ColDef, values []interface{}) ([]interface{}, error)
	ParseValue(def ColDef, value string) interface{}
	Query(context.Context, string) ([]string, [][]interface{}, error)
	Execute(context.Context, string) (string, int64, error)
	Name() string
	// DDL
	GetCreateTable(context.Context, string) (string, error)
	GetInsertStatement(context.Context, string, []ColDef, []interface{}) (string, error)
}
