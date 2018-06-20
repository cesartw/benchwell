package sqlhero

import (
	"fmt"

	"github.com/rivo/tview"
)

// layout
type layout struct {
	*tview.Flex

	screen     tview.Primitive
	statusline *statusline
}

func newLayout() *layout {
	l := &layout{}
	l.statusline = newStatusLine()

	l.Flex = tview.NewFlex().
		SetDirection(tview.FlexRow)

	return l
}

// SetText update the displayed text
func (l *layout) SetStatus(format string, args ...interface{}) {
	l.statusline.TextView.Clear()
	fmt.Fprintf(l.statusline.TextView, format, args...)
}

func (l *layout) SetScreen(screen tview.Primitive) {
	if l.screen != nil {
		l.Flex.RemoveItem(l.screen)
	}
	l.Flex.RemoveItem(l.statusline)

	l.screen = screen
	l.Flex.AddItem(l.screen, 0, 1, false)
	l.Flex.AddItem(l.statusline, 3, 1, false)
}

type statusline struct {
	*tview.TextView
}

func newStatusLine() *statusline {
	s := &statusline{}

	s.TextView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignRight).
		SetText("Ready")

	s.TextView.
		SetBorderPadding(0, 0, 0, 0).
		SetBorder(true)

	return s
}
