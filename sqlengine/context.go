package sqlengine

import (
	"context"

	// drivers
	_ "bitbucket.org/goreorto/sqlhero/sqlengine/mysql"
	//_ "bitbucket.org/goreorto/sqlhero/sqlengine/test"

	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
)

type contextkey struct {
	string
}

// Context represent an engine session
type Context context.Context

var (
	ckConnection = contextkey{"connection"}
	ckDatabase   = contextkey{"database"}
)

var _ Context = (*sqlctx)(nil)

type sqlctx struct {
	context.Context
}

// NewContext returns a new engine session
func NewContext(c context.Context, conn driver.Connection) Context {
	if c == nil {
		c = context.Background()
	}

	return &sqlctx{
		Context: context.WithValue(c, ckConnection, conn),
	}
}
