package textarea

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	cursorStyle = tcell.StyleDefault.
			Blink(true).
			Background(tcell.ColorYellow).
			Foreground(tcell.ColorBlack)

	textStyle = tcell.StyleDefault.Background(tview.Styles.PrimitiveBackgroundColor)
)

type line struct {
	start, end int
}

func (l line) length() int {
	return l.end - l.start
}

// Mode typing mode
type Mode int

// Typing modes
const (
	ModeNormal = iota
	ModeInsert
)

// TextArea is multi-line input with vim-like bindings
type TextArea struct {
	*tview.Box

	buffer *textBuffer

	mode      Mode
	appending bool
}

// New TextArea
func New() *TextArea {
	t := &TextArea{
		Box:    tview.NewBox(),
		buffer: newtextbuffer(100),
	}

	return t
}

// InputHandler none
func (t *TextArea) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	if t.mode == ModeInsert {
		return t.WrapInputHandler(t.inputHandlerInsert())
	}
	return t.WrapInputHandler(t.inputHandlerNormal())
}

func (t *TextArea) inputHandlerInsert() func(*tcell.EventKey, func(tview.Primitive)) {
	return t.WrapInputHandler(func(e *tcell.EventKey, setFocus func(tview.Primitive)) {
		switch e.Key() {
		case tcell.KeyBackspace2:
			if t.buffer.CursorMoveLeft() {
				t.buffer.Remove()
			}
		case tcell.KeyDelete:
			t.buffer.CursorMoveRight()
			if t.buffer.Remove() {
				t.buffer.CursorMoveLeft()
			}
		case tcell.KeyCR:
			t.buffer.Append(10)
			t.buffer.CursorMoveDown()
		case tcell.KeyRune:
			t.buffer.Append(e.Rune())
		case tcell.KeyLeft:
			t.buffer.CursorMoveLeft()
		case tcell.KeyRight:
			t.buffer.CursorMoveRight()
		case tcell.KeyUp:
			t.buffer.CursorMoveUp()
		case tcell.KeyDown:
			t.buffer.CursorMoveDown()
		case tcell.KeyEsc:
			t.actionSwitchMode(ModeNormal)
			if t.appending {
				t.buffer.CursorMoveLeft()
			}
		}
	})
}

func (t *TextArea) inputHandlerNormal() func(*tcell.EventKey, func(tview.Primitive)) {
	return t.WrapInputHandler(func(e *tcell.EventKey, setFocus func(tview.Primitive)) {
		switch e.Key() {
		case tcell.KeyDelete:
			t.buffer.Remove()
		case tcell.KeyRune:
			switch e.Rune() {
			case 'i':
				t.actionSwitchMode(ModeInsert)
			case 'I':
				t.actionSwitchMode(ModeInsert)
				t.buffer.CursorMoveBOL()
			case 'j':
				t.buffer.CursorMoveDown()
			case 'k':
				t.buffer.CursorMoveUp()
			case 'h':
				t.buffer.CursorMoveLeft()
			case 'l':
				if t.buffer.CursorMoveRight() && t.buffer.RuneAtCursor() == 10 {
					t.buffer.CursorMoveLeft()
				}
			case 'x':
				t.buffer.Remove()
			case 'A':
				t.buffer.CursorMoveEOL()
				t.actionSwitchMode(ModeInsert)
				t.appending = t.buffer.CursorMoveRight()
			case 'a':
				t.actionSwitchMode(ModeInsert)
				t.appending = t.buffer.CursorMoveRight()
			case 'o':
				t.buffer.CursorMoveEOL()
				t.buffer.Append(10)
				t.actionSwitchMode(ModeInsert)
			}
		case tcell.KeyLeft:
			t.buffer.CursorMoveLeft()
		case tcell.KeyRight:
			if t.buffer.CursorMoveRight() && t.buffer.RuneAtCursor() == 10 {
				t.buffer.CursorMoveLeft()
			}
		case tcell.KeyUp:
			t.buffer.CursorMoveUp()
		case tcell.KeyDown:
			t.buffer.CursorMoveDown()
		case tcell.KeyCtrlE:
			t.buffer.CursorMoveEOL()
		case tcell.KeyCtrlA:
			t.buffer.CursorMoveBOL()
		}
	})
}

// Draw none
func (t *TextArea) Draw(screen tcell.Screen) {
	x, y, w, h := t.Box.GetInnerRect()
	t.Box.Draw(screen)
	t.buffer.SetEditingWidth(100)

	runes := t.buffer.Runes()
	lines := t.buffer.Lines()

	curX, curLine := t.buffer.Position()
	for yOffset, line := range lines {
		for i := line.Start; i <= line.End; i++ {
			xOffset := i - line.Start

			style := textStyle
			if curLine == yOffset && curX == xOffset {
				style = cursorStyle
			}

			if runes[i] == 10 {
				screen.SetContent(x+xOffset, y+yOffset, ' ', nil, style)
				continue
			}

			screen.SetContent(x+xOffset, y+yOffset, runes[i], nil, style)
		}
	}

	// mode indicator
	labelStyle := tcell.StyleDefault.Background(tview.Styles.PrimitiveBackgroundColor)
	screen.SetContent(x+w-4, y+h, ' ', nil, labelStyle)
	screen.SetContent(x+w-6, y+h, ' ', nil, labelStyle)
	if t.mode == ModeInsert {
		screen.SetContent(x+w-5, y+h, 'i', nil, labelStyle)
	} else {
		screen.SetContent(x+w-5, y+h, 'n', nil, labelStyle)
	}
}

// SetText none
func (t *TextArea) SetText(s string) {
	t.buffer.SetText(s)
}

func (t *TextArea) actionBreak() {
	t.buffer.Append(10)
}

func (t *TextArea) actionSwitchMode(m Mode) {
	t.mode = m
}
