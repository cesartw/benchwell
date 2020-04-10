package gtk

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type App struct {
	*gtk.Application
	windows []*Window
	Menu    struct {
		Application struct {
			NewWindow   *glib.SimpleAction
			Preferences *glib.SimpleAction
		}
	}
}

func New(appid string) (*App, error) {
	var err error
	f := &App{}

	f.Application, err = gtk.ApplicationNew(appid, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}

	f.Connect("startup", func() {
		f.Menu.Application.NewWindow = glib.SimpleActionNew("new", nil)
		f.Menu.Application.Preferences = glib.SimpleActionNew("prefs", nil)

		f.AddAction(f.Menu.Application.NewWindow)
		f.AddAction(f.Menu.Application.Preferences)
	})

	f.Application.SetAccelsForAction("app.new", []string{"<control>N"})
	f.Application.SetAccelsForAction("win.new", []string{"<control>T"})
	f.Application.SetAccelsForAction("win.close", []string{"<control>W"})

	return f, nil
}
