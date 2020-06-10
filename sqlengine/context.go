package sqlengine

import (
	"context"

	// drivers
	_ "bitbucket.org/goreorto/sqlaid/sqlengine/mysql"
	//_ "bitbucket.org/goreorto/sqlaid/sqlengine/test"

	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type contextkey struct {
	string
}

var (
	ckConnection = contextkey{"connection"}
	ckDatabase   = contextkey{"database"}
)

type Context struct {
	context    context.Context
	connection driver.Connection
	database   driver.Database
	Logger     func(string)
}

func (c *Context) Connection() driver.Connection {
	return c.connection
}

func (c *Context) Database() driver.Database {
	return c.database
}

func (c *Context) Context() context.Context {
	if c == nil || c.context == nil {
		return context.Background()
	}

	return c.context
}

// NewContext returns a new engine session
func NewContext(conn driver.Connection, db driver.Database) *Context {
	return &Context{
		connection: conn,
		database:   db,
	}
}
