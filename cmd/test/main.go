package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appID = "com.iodone.sqlhero"

func main() {
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

		win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		errorCheck(err)

		paned, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
		errorCheck(err)

		paned.SetHExpand(true)
		paned.SetVExpand(true)

		frame1, err := gtk.FrameNew("")
		errorCheck(err)
		frame2, err := gtk.FrameNew("")
		errorCheck(err)

		frame1.SetShadowType(gtk.SHADOW_IN)
		frame2.SetShadowType(gtk.SHADOW_IN)
		frame1.SetSizeRequest(50, -1)
		frame2.SetSizeRequest(50, -1)

		paned.Pack1(frame1, false, true)
		paned.Pack2(frame2, true, false)

		list, err := gtk.ListBoxNew()
		errorCheck(err)
		label, err := gtk.LabelNew("a")
		errorCheck(err)
		list.Add(label)

		frame1.Add(list)
		win.Add(paned)

		win.ShowAll()
		application.AddWindow(win)
	})

	// Connect function to application shutdown event, this is not required.
	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(application.Run(os.Args))
}

func errorCheck(e error) {
	if e != nil {
		// panic for any errors.
		log.Panic(e)
	}
}
