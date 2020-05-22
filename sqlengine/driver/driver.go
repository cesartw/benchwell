package driver

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"bitbucket.org/goreorto/sqlaid/config"
)

const NULL_PATTERN = "<NULL>"

type ColType uint
type TableType uint

// Column types
const (
	ColTypeBoolean ColType = iota
	ColTypeString
	ColTypeFloat
	ColTypeInt
	ColTypeDate
	ColTypeList
)

// Table types
const (
	TableTypeRegular TableType = iota
	TableTypeView
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
func Connect(ctx context.Context, cfg config.Connection) (Connection, error) {
	colonS := strings.Split(cfg.GetDSN(), ":")

	driverMU.Lock()
	d, ok := drivers[colonS[0]]
	driverMU.Unlock()

	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", colonS[0])
	}

	// kind hacky
	return d.Connect(ctx, cfg)
}

// Driver is a database implementation for SQLHero
type Driver interface {
	Connect(context.Context, config.Connection) (Connection, error)
}

// Connection ...
type Connection interface {
	Databases(context.Context) ([]string, error)
	UseDatabase(context.Context, string) (Database, error)
	Disconnect(context.Context) error
	Reconnect(context.Context) error
	LastError() error
}

// Database ...
type Database interface {
	Tables(context.Context) ([]TableDef, error)
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

// ColDef describe a column
type ColDef struct {
	Name               string
	PK, FK             bool
	Precision          int
	Unsigned, Nullable bool
	Type               ColType
	Values             []string
}

type TableDefs []TableDef
type TableDef struct {
	Name string
	Type TableType
}

func (t TableDef) String() string {
	return t.Name
}

func (t TableDef) IsZero() bool {
	if strings.TrimSpace(t.Name) == "" {
		return true
	}
	return false
}

func (t TableDefs) ToStringer() []fmt.Stringer {
	s := make([]fmt.Stringer, len(t))
	for i, def := range t {
		s[i] = def
	}
	return s
}

func (t *TableDefs) FromStringer(s []fmt.Stringer) {
	tt := make(TableDefs, len(s))
	for i, stringer := range s {
		tt[i], _ = stringer.(TableDef)
	}

	*t = tt
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
