package driver

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
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
func Connect(dsn string) (Connection, error) {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	driverMU.Lock()
	d, ok := drivers[dsnURL.Scheme]
	driverMU.Unlock()

	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", dsnURL.Scheme)
	}

	// kind hacky
	return d.Connect(strings.TrimPrefix(dsn, dsnURL.Scheme+"://"))
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
	Name   string
	Type   string
	PK, FK bool
}

// Database ...
type Database interface {
	Tables() ([]string, error)
	TableDefinition(tableName string) ([]ColDef, error)
	FetchTable(tableName string, page, pageSize int64) ([]ColDef, [][]*string, error)
	DeleteRecord(tableName string, defs []ColDef, values []*string) error
	UpdateRecord(tableName string, cols []ColDef, values, oldValues []*string) (string, error)
	InsertRecord(tableName string, cols []ColDef, values []*string) (string, error)
	// Query(string, interface{}) ([]string,[]string,error)
	// Execute(string, interface{}) (in,error)
	Name() string
}
