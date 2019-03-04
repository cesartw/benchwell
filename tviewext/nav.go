package tviewext

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Navigable struct {
	focused   tview.Primitive
	keybinds  map[tcell.Key]func() tview.Primitive
	runebinds map[rune]map[tcell.ModMask]func() tview.Primitive

	Fallback tview.Primitive

	tview.Primitive
}

// Focus ...
func (n *Navigable) Focus(_ func(tview.Primitive)) {
	n.Blur()
	n.focused = n.Fallback
	n.Fallback.Focus(n.DelegateFocus)
}

func (n *Navigable) DelegateFocus(p tview.Primitive) {
	n.Blur()
	n.focused = p
	p.Focus(n.DelegateFocus)
}

func (n *Navigable) Blur() {
	if n.focused != nil {
		n.focused.Blur()
		n.focused = nil
	}
}

func (n *Navigable) GetFocusable() tview.Focusable {
	return n
}

func (n Navigable) HasFocus() bool {
	return n.focused != nil
}

func (n *Navigable) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, _ func(tview.Primitive)) {
		switch e.Key() {
		case tcell.KeyRune:
			if f, ok := n.runebinds[e.Rune()][e.Modifiers()]; ok {
				n.Blur()
				n.focused = f()
				n.focused.Focus(n.DelegateFocus)

				return
			}
		default:
			if f, ok := n.keybinds[e.Key()]; ok {
				n.Blur()
				n.focused = f()
				n.focused.Focus(n.DelegateFocus)

				return
			}
		}

		if n.focused != nil {
			n.focused.InputHandler()(e, n.DelegateFocus)
		}
	}
}

func (n *Navigable) Keybind(k tcell.Key, f func() tview.Primitive) {
	if n.keybinds == nil {
		n.keybinds = map[tcell.Key]func() tview.Primitive{}
	}

	n.keybinds[k] = f
}

func (n *Navigable) Runebind(r rune, mods tcell.ModMask, f func() tview.Primitive) {
	if n.runebinds == nil {
		n.runebinds = map[rune]map[tcell.ModMask]func() tview.Primitive{}
	}

	if n.runebinds[r] == nil {
		n.runebinds[r] = map[tcell.ModMask]func() tview.Primitive{}
	}

	n.runebinds[r][mods] = f
}

func (n Navigable) GetFocused() tview.Primitive {
	return n.focused
}
