package main

import (
	"flag"
	"log"
	"os"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/logger"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	errorCheck(err)

	// Connect function to application startup event, this is not required.
	application.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	application.Connect("activate", func() {
		log.Println("application activate")
		win := &Window{}
		win.init()

		connect := &ConnectScreen{}
		err := connect.init()
		errorCheck(err)
		win.Add(connect)

		connect.SetConnections(conf.Connection)

		win.Show()
		application.AddWindow(win.Window)
	})

	// Connect function to application shutdown event, this is not required.
	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(application.Run(os.Args))
}

type Window struct {
	*gtk.Window
	builder *gtk.Builder
	box     *gtk.Box
}

func (w *Window) init() {
	var err error
	w.builder, err = gtk.BuilderNewFromFile("ui/main.glade")
	errorCheck(err)

	signals := map[string]interface{}{
		"on_main_window_destroy": w.onMainWindowDestroy,
	}

	w.builder.ConnectSignals(signals)
	obj, err := w.builder.GetObject("MainWindow")
	errorCheck(err)

	w.Window = obj.(*gtk.Window)

	obj, err = w.builder.GetObject("MainWindowBox")
	errorCheck(err)

	w.box = obj.(*gtk.Box)
}

func (w *Window) Add(wd gtk.IWidget) {
	w.box.Add(wd)
	//w.box.ReorderChild(wd, 1)
}

func (w Window) onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
}

func errorCheck(e error) {
	if e != nil {
		// panic for any errors.
		log.Panic(e)
	}
}
