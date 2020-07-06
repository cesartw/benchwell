package gtk

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
)

type Application struct {
	*gtk.Application
	windows []*Window
	Menu    struct {
		Application struct {
			NewWindow   *glib.SimpleAction
			Preferences *glib.SimpleAction
			DarkMode    *glib.SimpleAction
		}
	}

	DarkMode bool
	ctrl     applicationCtrl
}
type applicationCtrl interface {
	AppID() string
	OnStartup()
	OnActivate()
	OnShutdown()
	OnNewWindow()
	OnPreferences()
	Config() *config.Config
}

func (a Application) Init(ctrl applicationCtrl) (*Application, error) {
	var err error
	a.ctrl = ctrl

	a.Application, err = gtk.ApplicationNew(ctrl.AppID(), glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}
	a.DarkMode = ctrl.Config().GUI.DarkMode.Bool()

	a.Connect("startup", func() {
		a.Menu.Application.NewWindow = glib.SimpleActionNew("new", nil)
		a.Menu.Application.Preferences = glib.SimpleActionNew("preferences", nil)
		a.Menu.Application.DarkMode = glib.SimpleActionNew("darkmode", nil)

		a.AddAction(a.Menu.Application.NewWindow)
		a.AddAction(a.Menu.Application.Preferences)
		a.AddAction(a.Menu.Application.DarkMode)
		a.loadSettingsCSS()
		a.Menu.Application.Preferences.Connect("activate", a.OnShowPreferences)
	})

	//a.Application.SetAccelsForAction("app.new", []string{"<control>N"})
	// main tab
	a.Application.SetAccelsForAction("win.new", []string{"<control>N"})
	// sub tab
	a.Application.SetAccelsForAction("win.tabnew", []string{"<control>T"})
	// close sub tab, and main tab when there's no sub tabs left
	a.Application.SetAccelsForAction("win.close", []string{"<control>W"})

	a.Application.Connect("activate", func() {
		a.Menu.Application.NewWindow.Connect("activate", ctrl.OnNewWindow)
		a.Menu.Application.Preferences.Connect("activate", ctrl.OnPreferences)

		a.Menu.Application.DarkMode.Connect("activate", func() {
			a.ToggleMode()
		})
	})
	a.Application.Connect("activate", ctrl.OnActivate)

	// Connect function to application shutdown event, this is not required.
	a.Application.Connect("shutdown", ctrl.OnShutdown)

	return &a, nil
}

func (a *Application) ToggleMode() {
	a.ctrl.Config().GUI.DarkMode.SetBool(!a.ctrl.Config().GUI.DarkMode.Bool())
	a.loadSettingsCSS()
}

func (a *Application) OnShowPreferences() {
	modal, err := gtk.DialogNewWithButtons("Preferences", nil,
		gtk.DIALOG_DESTROY_WITH_PARENT|gtk.DIALOG_MODAL,
		[]interface{}{"Done", gtk.RESPONSE_OK},
	)
	if err != nil {
		a.ctrl.Config().Error(err)
		return
	}

	//pref, err := Preferences{}.Init()
	//if err != nil {
	//config.Env.Log.Error(err)
	//return
	//}

	//b, err := modal.GetContentArea()
	//if err != nil {
	//config.Env.Log.Error(err)
	//return
	//}
	//b.Add(pref)

	modal.Run()
	defer modal.Destroy()
}

func (a *Application) loadSettingsCSS() {
	css, err := gtk.CssProviderNew()
	if err != nil {
		panic(err)
	}

	// TODO: works, need to check vendoring
	//css.LoadFromResource("/org/gtk/libgtk/theme/Adwaita/gtk-contained-dark.css")

	err = css.LoadFromData(a.ctrl.Config().CSS())
	if err != nil {
		panic(err)
	}

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		panic(err)
	}

	gtk.AddProviderForScreen(screen, css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
