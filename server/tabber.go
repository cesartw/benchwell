package server

import (
	"fmt"

	"bitbucket.org/goreorto/sqlhero/tviewext"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var tabshortcut = map[rune]int{
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

// Tabber ...
type Tabber struct {
	Box            *tview.Box
	tabs           []*Tab
	selectTabIndex int

	tviewext.Navigable
}

// NewTabber ...
func NewTabber() *Tabber {
	t := &Tabber{
		Box:  tview.NewBox(),
		tabs: make([]*Tab, 0),
	}

	t.Box.SetBorder(false)

	return t
}

// CurrentTab ...
func (t Tabber) CurrentTab() *Tab {
	for _, t := range t.tabs {
		if t.HasFocus() {
			return t
		}
	}

	return nil
}

// AddTab ...
func (t *Tabber) AddTab(title string, content tview.Primitive, focus bool) {
	if content == nil {
		content = tview.NewBox()
	}

	item := &Tab{
		content: content,
		header:  tview.NewTextView(),
		title:   title,
	}
	item.header.SetBorderPadding(0, 0, 1, 1)
	t.tabs = append(t.tabs, item)

	if focus {
		t.selectTabIndex = len(t.tabs) - 1
		t.Navigable.DelegateFocus(item.content)
	}
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
		tabsWidth += len(tab.prefixedTitle(i)) + headerPadding
		if tab.HasFocus() {
			selectedRightEdge = tabsWidth
		}
	}

	var requiredSpace int
	if selectedRightEdge > w {
		requiredSpace = selectedRightEdge - w
	}

	// arrows
	headerx := x + 1
	w -= 2

	tabOffset := headerx
	var skipped int
	for i, tab := range t.tabs {
		tabWidth := len(tab.prefixedTitle(i)) + headerPadding
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
		tab.header.SetText(tab.prefixedTitle(i))
		tab.header.SetBackgroundColor(tcell.ColorDimGray)
		if i == t.selectTabIndex {
			tab.header.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
			tab.content.Draw(screen)
		}
		tab.header.Draw(screen)
	}

	var tabsWidth int
	for i, tab := range t.tabs {
		tabsWidth += len(tab.prefixedTitle(i)) + 2
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
func (t *Tabber) Focus(_ func(tview.Primitive)) {
	if len(t.tabs) == 0 {
		return
	}

	t.DelegateFocus(t.tabs[0].content)
}

// HasFocus ...
func (t *Tabber) HasFocus() bool {
	for _, tab := range t.tabs {
		if tab.HasFocus() {
			return true
		}
	}

	return false
}

// InputHandler ...
func (t *Tabber) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, _ func(tview.Primitive)) {
		// Switch tab
		if e.Modifiers()^tcell.ModAlt == 0 {
			switch e.Rune() {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if tabshortcut[e.Rune()] >= len(t.tabs) {
					return
				}

				t.DelegateFocus(t.tabs[tabshortcut[e.Rune()]].content)
				t.selectTabIndex = tabshortcut[e.Rune()]

				return
			}

			if e.Key() == tcell.KeyLeft {
				t.selectTabIndex--
				if t.selectTabIndex < 0 {
					t.selectTabIndex = len(t.tabs) - 1
				}
				t.DelegateFocus(t.tabs[t.selectTabIndex].content)
				return
			}

			if e.Key() == tcell.KeyRight {
				t.selectTabIndex++
				if t.selectTabIndex >= len(t.tabs) {
					t.selectTabIndex = 0
				}
				t.DelegateFocus(t.tabs[t.selectTabIndex].content)
				return
			}
		}

		// Close tab
		if e.Key() == tcell.KeyCtrlW {
			t.tabs = append(t.tabs[:t.selectTabIndex], t.tabs[t.selectTabIndex+1:]...)

			if t.selectTabIndex >= len(t.tabs) {
				t.selectTabIndex = len(t.tabs) - 1
			}

			if len(t.tabs) > 0 {
				t.DelegateFocus(t.tabs[t.selectTabIndex].content)
			}

			return
		}

		if len(t.tabs) > 0 {
			t.tabs[t.selectTabIndex].content.InputHandler()(e, t.DelegateFocus)
		}
	}
}

// Tab ...
type Tab struct {
	header  *tview.TextView
	content tview.Primitive
	title   string
}

// Focus ...
func (t *Tab) Focus(delegate func(tview.Primitive)) {
	delegate(t.content)
}

// Blur ...
func (t *Tab) Blur() {
	t.content.Blur()
}

func (t Tab) prefixedTitle(index int) string {
	return fmt.Sprintf("%d. %s", index+1, t.title)
}

// HasFocus ...
func (t Tab) HasFocus() bool {
	return t.content.GetFocusable().HasFocus()
}
