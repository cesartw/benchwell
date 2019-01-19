package server

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// NewTabber ...
func NewTabber() *Tabber {
	t := &Tabber{
		Box:  tview.NewBox(),
		tabs: make([]*Tab, 0),
	}

	t.Box.SetBorder(false)

	return t
}

// Tabber ...
type Tabber struct {
	*tview.Box
	tabs []*Tab

	focused    tview.Primitive
	focusedTab *Tab
}

// AddTab ...
func (t *Tabber) AddTab(title string, content tview.Primitive, focus bool) {
	if content == nil {
		content = tview.NewBox()
	}

	item := &Tab{
		content:  content,
		header:   tview.NewTextView(),
		title:    title,
		hasFocus: focus,
	}
	item.header.SetBorderPadding(0, 0, 1, 1)

	if focus {
		t.blurFocusedTab()
		t.focusedTab = item
		t.focusedTab.Focus(t.focusdelegate)
	}

	t.tabs = append(t.tabs, item)
}

// SetRect ...
func (t *Tabber) SetRect(x, y, w, h int) {
	headerHeight := 1
	headerPadding := 2
	t.Box.SetRect(x, y+headerHeight, w, h-headerHeight)

	if len(t.tabs) == 0 {
		return
	}

	var (
		tabsWidth    int
		tabsPosition = make([]int, len(t.tabs))
	)

	var selectedRightEdge int
	for i, tab := range t.tabs {
		tabsPosition[i] = tabsWidth
		tabsWidth += len(tab.PrefixedTitle(i)) + headerPadding
		if tab == t.focusedTab {
			selectedRightEdge = tabsWidth
		}
	}

	var requiredSpace int
	if selectedRightEdge > w {
		requiredSpace = selectedRightEdge - w
	}

	// arrows
	x = x + 1
	w -= 2

	tabOffset := x
	var skipped int
	for i, tab := range t.tabs {
		tabWidth := len(tab.PrefixedTitle(i)) + headerPadding
		if requiredSpace > 0 {
			requiredSpace -= tabWidth
			skipped++
			continue
		}
		i -= skipped

		tab.header.SetRect(tabOffset+i+1, y, tabWidth, headerHeight)
		tab.content.SetRect(x, y+headerHeight, w, h-headerHeight-2)
		tabOffset += tabWidth
	}
}

// Draw ...
func (t *Tabber) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)
	x, y, w, _ := t.Box.GetRect()

	// clear header
	for i := x; i < x+w; i++ {
		screen.SetContent(i, y-1, ' ', nil, tcell.StyleDefault)
	}

	for i, tab := range t.tabs {
		tab.header.SetText(tab.PrefixedTitle(i))
		tab.header.SetBackgroundColor(tcell.ColorDimGray)
		if tab.hasFocus {
			tab.header.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
			tab.content.Draw(screen)
		}
		tab.header.Draw(screen)
	}

	var tabsWidth int
	for i, tab := range t.tabs {
		tabsWidth += len(tab.PrefixedTitle(i)) + 2
	}

	if tabsWidth > w {
		style := tcell.StyleDefault.
			Background(tview.Styles.PrimitiveBackgroundColor).
			Foreground(tview.Styles.TitleColor)

		screen.SetContent(x, y-1, '<', nil, style)
		screen.SetContent(x+w-1, y-1, '>', nil, style)
	}
}

// Focus ...
func (t *Tabber) Focus(delegate func(tview.Primitive)) {
	if len(t.tabs) == 0 {
		return
	}

	if t.focusedTab != nil {
		return
	}

	t.focusedTab = t.tabs[0]
	t.focusedTab.Focus(t.focusdelegate)
}

func (t *Tabber) blurFocusedTab() {
	if t.focusedTab != nil {
		t.focusedTab.Blur()
		t.focusedTab = nil
	}
	if t.focused != nil {
		t.focused.Blur()
		t.focused = nil
	}
}

func (t *Tabber) focusdelegate(p tview.Primitive) {
	if t.focused != nil {
		t.focused.Blur()
	}

	t.focused = p
	t.focused.Focus(t.focusdelegate)
}

// HasFocus ...
func (t *Tabber) HasFocus() bool {
	for _, tab := range t.tabs {
		if tab.hasFocus {
			return true
		}
	}

	return false
}

// InputHandler ...
func (t *Tabber) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, _ func(tview.Primitive)) {
		// if ALT isn't on, let the current focus tab item handle the event
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
			if shortcut[e.Rune()] >= len(t.tabs) {
				return
			}

			// unfocus selected tab
			t.blurFocusedTab()

			// select tab at
			t.focusedTab = t.tabs[shortcut[e.Rune()]]
			t.focusedTab.Focus(t.focusdelegate)

			return
		}

		if e.Key() == tcell.KeyLeft {
		}

		if e.Key() == tcell.KeyRight {
		}
	}
}

// Tab ...
type Tab struct {
	header   *tview.TextView
	content  tview.Primitive
	title    string
	hasFocus bool
}

// Focus ...
func (t *Tab) Focus(delegate func(tview.Primitive)) {
	t.hasFocus = true
	delegate(t.content)
}

// Blur ...
func (t *Tab) Blur() {
	t.hasFocus = false
	t.content.Blur()
}

// PrefixedTitle ...
func (t *Tab) PrefixedTitle(index int) string {
	return fmt.Sprintf("%d. %s", index+1, t.title)
}
