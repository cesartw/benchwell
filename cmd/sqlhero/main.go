package main

import (
	"context"
	"flag"
	"log"
	"os"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/logger"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
)

const appID = "com.iodone.sqlhero"

func main() {
	var phrase string
	flag.StringVar(&phrase, "phrase", "", "En/Decrypting phrase for connection passwords")
	flag.Parse()

	conf, err := config.New(phrase)
	if err != nil {
		panic(err)
	}

	eng := sqlengine.New(conf, logger.NewLogger(os.Stdout))
	defer eng.Dispose()

	// Create a new application.
	app, err := gtk.New(appID)
	errorCheck(err)

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		connect, err := app.NewConnectScreen()
		errorCheck(err)

		app.Add(connect)

		connect.SetConnections(conf.Connection)

		connect.OnConnect(func() {
			conn := connect.ActiveConnection()
			ctx, err := eng.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
			errorCheck(err)

			dbNames, err := eng.Databases(ctx)
			errorCheck(err)

			connectionScr, err := app.NewConnectionScreen(ctx)
			errorCheck(err)

			connectionScr.OnDatabaseSelected(func() {
				dbName := connectionScr.ActiveDatabase()
				ctx, err = eng.UseDatabase(ctx, dbName)

				tables, err := eng.Tables(ctx)
				errorCheck(err)

				connectionScr.SetTables(tables)
			})

			connectionScr.OnTableSelected(func() {
				tableName, ok := connectionScr.ActiveTable()
				if ok {
					def, data, err := eng.FetchTable(ctx, tableName, 0, 40)
					errorCheck(err)

					connectionScr.SetTableData(def, data)
				}
			})

			app.Remove(connect)
			app.Add(connectionScr)

			connectionScr.SetDatabases(dbNames)

			app.PushStatus("Connected to `%s`", conn.Host)
		})

		connect.OnTest(func() {
			conn := connect.ActiveConnection()
			ctx, err := eng.Connect(sqlengine.Context(context.TODO()), conn.GetDSN())
			if err != nil {
				errorCheck(err)
			}
			app.PushStatus("Connection to `%s` was successful", conn.Host)
			eng.Disconnect(ctx)
		})

		connect.OnSave(func() {
			//conn := connect.ActiveConnection()
			app.PushStatus("Saved")
		})

		app.Show()
		app.PushStatus("Ready")
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(os.Args))
}

func errorCheck(e error) {
	if e != nil {
		// panic for any errors.
		log.Panic(e)
	}
}
