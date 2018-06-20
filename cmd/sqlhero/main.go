package main

import (
	"bitbucket.org/goreorto/sqlhero"
	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/logger"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	eng := sqlengine.New(conf, logger.NewLogger())
	defer eng.Dispose()

	app := sqlhero.New(conf, eng)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
