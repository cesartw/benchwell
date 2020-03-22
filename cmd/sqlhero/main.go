package main

import (
	"flag"
	"log"
	"os"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/ctrl"
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

	ctl, err := ctrl.MainCtrl{}.Init(ctrl.Options{
		Engine:  eng,
		Config:  conf,
		Factory: app,
	})
	errorCheck(err)

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", ctl.OnActivate)

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
