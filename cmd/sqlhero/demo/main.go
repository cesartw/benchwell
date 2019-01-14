package main

import (
	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/connect"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	t := NewTabber()

	conf, err := config.New("")
	if err != nil {
		panic(err)
	}
	screen := connect.New(conf)

	t.AddTab("tab 1", nil, true)
	t.AddTab("tab 2", screen, false)
	t.AddTab("tab 3", nil, false)
	t.AddTab("tab 4", nil, false)

	app.SetRoot(t, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

// NewTabber ...
func NewTabber() *Tabber {
	t := &Tabber{
		Box:         tview.NewBox(),
		tabs:        make([]*Tab, 0),
		minTabWidth: 5,
		maxTabWidth: 10,
	}

	t.Box.SetBorder(true)

	return t
}

// Tabber ...
type Tabber struct {
	*tview.Box
	tabs []*Tab

	minTabWidth int
	maxTabWidth int
}

// Tab ...
type Tab struct {
	tview.Primitive
	tabber *Tabber

	header *tview.TextView

	title    string
	hasFocus bool

	focused tview.Primitive
}

// AddTab ...
func (t *Tabber) AddTab(title string, content tview.Primitive, focus bool) {
	if content == nil {
		content = tview.NewBox()
	}
	item := &Tab{
		header:    tview.NewTextView(),
		Primitive: content,
		tabber:    t,
		hasFocus:  focus,
	}

	item.header.SetText(title)
	item.header.SetBorderPadding(0, 0, 1, 1)
	t.tabs = append(t.tabs, item)
}

// SetRect ...
func (t *Tabber) SetRect(x, y, w, h int) {
	headerHeight := 1
	t.Box.SetRect(x, y+headerHeight, w, h-headerHeight)

	tabWidth := w / len(t.tabs)
	if tabWidth < t.minTabWidth {
		tabWidth = t.minTabWidth
	}
	if tabWidth > t.maxTabWidth {
		tabWidth = t.maxTabWidth
	}

	for i, tab := range t.tabs {
		tab.header.SetRect(x+tabWidth*i+i, y, tabWidth, headerHeight)
		tab.Primitive.SetRect(x+1, y+headerHeight+1, w-2, h-headerHeight-2)
	}
}

// Focus ...
func (t *Tabber) Focus(delegate func(tview.Primitive)) {
	for _, tab := range t.tabs {
		if tab.hasFocus {
			delegate(t.tabs[0])
			return
		}
	}

	if len(t.tabs) > 0 {
		delegate(t.tabs[0])
		return
	}
}

// Draw ...
func (t *Tabber) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)

	for _, tab := range t.tabs {
		tab.Draw(screen)
	}
}

// HasFocus ...
func (t *Tabber) HasFocus() bool {
	for _, tab := range t.tabs {
		if tab.GetFocusable().HasFocus() {
			return true
		}
	}

	return false
}

// InputHandler ...
func (t *Tab) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, setFocus func(tview.Primitive)) {
		// if ALT isn't not, let the current focus tab item handle the event
		if e.Modifiers()^tcell.ModAlt != 0 {
			t.focused.InputHandler()(e, t.focusdelegate)
			return
		}

		shortcut := map[rune]int{
			'1': 0,
			'2': 1,
			'3': 2,
			'4': 3,
			'5': 4,
			'6': 5,
			'7': 6,
			'8': 7,
			'9': 8,
			'0': 9,
		}
		// Switch tab
		switch e.Rune() {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if shortcut[e.Rune()] >= len(t.tabber.tabs) {
				return
			}

			for _, tab := range t.tabber.tabs {
				tab.Blur()
			}

			setFocus(t.tabber.tabs[shortcut[e.Rune()]])
		}
	}
}

// Blur ...
func (t *Tab) Blur() {
	t.hasFocus = false
	t.Primitive.Blur()
	if t.focused != nil {
		t.focused.Blur()
	}
}

// Focus ...
func (t *Tab) Focus(delegate func(tview.Primitive)) {
	t.focused = t.Primitive
	t.Primitive.Focus(t.focusdelegate)
}

func (t *Tab) focusdelegate(p tview.Primitive) {
	if t.focused != nil {
		t.focused.Blur()
	}
	t.focused = p
	t.focused.Focus(t.focusdelegate)
}

// Draw ...
func (t *Tab) Draw(screen tcell.Screen) {
	t.header.SetBackgroundColor(tcell.ColorDimGray)
	if t.GetFocusable().HasFocus() {
		t.Primitive.Draw(screen)
		t.header.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	}
	t.header.Draw(screen)
}
