package driver

import (
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
func Connect(dsn string) (Connection, error) {
	colonS := strings.Split(dsn, ":")

	driverMU.Lock()
	d, ok := drivers[colonS[0]]
	driverMU.Unlock()

	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", colonS[0])
	}

	// kind hacky
	return d.Connect(strings.TrimPrefix(dsn, colonS[0]+"://"))
}

// Driver is a database implementation for SQLHero
type Driver interface {
	Connect(string) (Connection, error)
	// Sanitize(string) string
}

// Connection ...
type Connection interface {
	UseDatabase(string) error
	Databases() ([]Database, error)
	Disconnect() error
	Reconnect() error
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

// Database ...
type Database interface {
	Tables() ([]string, error)
	TableDefinition(tableName string) ([]ColDef, error)
	FetchTable(tableName string, page, pageSize int64) ([]ColDef, [][]interface{}, error)
	DeleteRecord(tableName string, defs []ColDef, values []*string) error
	UpdateRecord(tableName string, cols []ColDef, values, oldValues []interface{}) (string, error)
	// UpdateField updates a single field. cols[-1] is the changed values, cols[:-1] are primary keys
	UpdateField(tableName string, cols []ColDef, values []interface{}) (string, error)
	InsertRecord(tableName string, cols []ColDef, values []*string) ([]*string, error)
	ParseValue(def ColDef, value string) interface{}
	Query(string) ([]string, [][]interface{}, error)
	// Execute(string, interface{}) (in,error)
	Name() string
	// DDL
	// GetCreateTable(string) string
}
