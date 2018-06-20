package connect

import (
	"bitbucket.org/goreorto/sqlhero/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Screen ...
type Screen struct {
	*tview.Flex

	list *ConnectionList
	form *Form

	OnTest    func(config.Connection)
	OnSave    func(*config.Connection)
	OnConnect func(config.Connection)
}

// New ...
func New(app *tview.Application, conf *config.Config) *Screen {
	s := &Screen{}

	// FORM
	s.form = NewForm()

	// FAV LIST
	s.list = NewConnectionList()
	s.list.SetConnections(conf.Connection)

	s.list.OnSelectConnection = func(con *config.Connection) {
		s.form.SetConnection(con)
		app.SetFocus(s.form)
	}

	s.list.OnNewConnection = func() {
		s.form.SetConnection(nil)
		app.SetFocus(s.form)
	}

	s.list.OnDeleteConnection = func(con *config.Connection) {
		for i, c := range conf.Connection {
			if c == con {
				conf.Connection = append(conf.Connection[:i], conf.Connection[i+1:]...)
				conf.Save()
				s.list.SetConnections(conf.Connection)

				return
			}
		}
	}

	s.form.OnSave = func(c *config.Connection) {
		var selectedItem int

		defer func() {
			conf.Save()
			s.list.SetConnections(conf.Connection)
			s.list.List.SetCurrentItem(selectedItem)
		}()

		for i, cc := range conf.Connection {
			if c == cc {
				selectedItem = i + 1
				conf.Save()
				return
			}
		}

		s.list.AddConnection(c)
		conf.Connection = append(conf.Connection, c)
		selectedItem = len(conf.Connection) + 1
		s.onSave(c)
	}

	s.form.OnConnect = s.onConnect
	s.form.OnTest = s.onTest

	s.Flex = tview.NewFlex().
		AddItem(s.list, 30, 1, true).
		AddItem(s.form, 0, 2, false)

	return s
}

func (s *Screen) onTest(c config.Connection) {
	if s.OnTest != nil {
		s.OnTest(c)
	}
}

func (s *Screen) onSave(c *config.Connection) {
	if s.OnSave != nil {
		s.OnSave(c)
	}
}

func (s *Screen) onConnect(c config.Connection) {
	if s.OnConnect != nil {
		s.OnConnect(c)
	}
}

// Keybinds ...
func (s *Screen) Keybinds() map[tcell.Key]tview.Primitive {
	return map[tcell.Key]tview.Primitive{
		tcell.KeyCtrlL: s.list,
		tcell.KeyCtrlN: s.form,
	}
}
