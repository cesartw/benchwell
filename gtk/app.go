package gtk

import (
	"bitbucket.org/goreorto/sqlaid/assets"
	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gdk"
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
			DarkMode    *glib.SimpleAction
		}
	}

	DarkMode bool
}

func New(appid string) (*App, error) {
	var err error
	f := &App{}

	f.Application, err = gtk.ApplicationNew(appid, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}
	f.DarkMode = config.Env.GUI.DarkMode

	f.Connect("startup", func() {
		f.Menu.Application.NewWindow = glib.SimpleActionNew("new", nil)
		f.Menu.Application.Preferences = glib.SimpleActionNew("prefs", nil)
		f.Menu.Application.DarkMode = glib.SimpleActionNew("darkmode", nil)

		f.AddAction(f.Menu.Application.NewWindow)
		f.AddAction(f.Menu.Application.Preferences)
		f.AddAction(f.Menu.Application.DarkMode)
		f.SetTheme()
		f.loadSettingsCSS()
	})

	//f.Application.SetAccelsForAction("app.new", []string{"<control>N"})
	// main tab
	f.Application.SetAccelsForAction("win.new", []string{"<control>N"})
	// sub tab
	f.Application.SetAccelsForAction("win.tabnew", []string{"<control>T"})
	// close sub tab, and main tab when there's no sub tabs left
	f.Application.SetAccelsForAction("win.close", []string{"<control>W"})

	return f, nil
}

func (a *App) ToggleMode() {
	a.DarkMode = !a.DarkMode
	a.SetTheme()
}

func (a *App) SetTheme() {
	stylePath := assets.THEME_DARK + assets.BRAND_DARK
	if !a.DarkMode {
		stylePath = assets.THEME_LIGHT
	}
	a.loadCSS(stylePath)
}

func (a *App) loadCSS(path string) {
	css, err := gtk.CssProviderNew()
	if err != nil {
		panic(err)
	}

	err = css.LoadFromData(path)
	if err != nil {
		panic(err)
	}

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		panic(err)
	}

	gtk.AddProviderForScreen(screen, css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (a *App) loadSettingsCSS() {
	css, err := gtk.CssProviderNew()
	if err != nil {
		panic(err)
	}

	// TODO: works, need to check vendoring
	//css.LoadFromResource("/org/gtk/libgtk/theme/Adwaita/gtk-contained-dark.css")

	err = css.LoadFromData(config.Env.CSS())
	if err != nil {
		panic(err)
	}

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		panic(err)
	}

	gtk.AddProviderForScreen(screen, css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
