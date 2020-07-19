// Package test is testing driver
package test

import (
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
)

type testDriver struct{}

type testConn struct {
	dsn string
}

type testDb struct {
	name string
}

func init() {
	driver.RegisterDriver("test", &testDriver{})
}

func (d *testDriver) Connect(dsn string) (driver.Connection, error) {
	return &testConn{dsn: dsn}, nil
}

func (c *testConn) UseDatabase(s string) error {
	return nil
}

func (c *testConn) Reconnect() error {
	return nil
}

// Disconnect ...
func (c *testConn) Disconnect() error {
	return nil
}

// LastError ...
func (c *testConn) LastError() error {
	return nil
}

// Databases ...
func (c *testConn) Databases() ([]driver.Database, error) {
	dbs := make([]driver.Database, 3)
	dbs = append(dbs, driver.Database(&testDb{"database1"}))
	dbs = append(dbs, driver.Database(&testDb{"database2"}))
	dbs = append(dbs, driver.Database(&testDb{"database3"}))

	return dbs, nil
}

func (d *testDb) Name() string {
	return d.name
}

func (d *testDb) Tables() (tables []string, err error) {
	return tables, err
}

func (d *testDb) FetchTable(page, pageSize int64) (columns []string, rows [][]string, err error) {
	return columns, rows, err
}

func (d *testDb) DeleteRecord(ID string) error {
	return nil
}

func (d *testDb) UpdateRecord(ID string, cols, values []string) (int, error) {
	return 0, nil
}
