package textarea

import (
	"fmt"
	"os"

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

	viewPort struct {
		xScroll int
		yScroll int
	}
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
				t.moveViewPort(dirLeft, dirDown)
			}

		case tcell.KeyDelete:
			t.buffer.CursorMoveRight()
			if t.buffer.Remove() {
				t.buffer.CursorMoveLeft()
				t.moveViewPort(dirLeft, dirDown)
			}

		case tcell.KeyCR:
			t.buffer.Append(10)
			t.moveViewPort(dirLeft, dirDown)

		case tcell.KeyRune:
			t.buffer.Append(e.Rune())
			t.moveViewPort(dirRight, dirDown)

		case tcell.KeyLeft:
			t.buffer.CursorMoveLeft()
			t.moveViewPort(dirLeft, dirDown)

		case tcell.KeyRight:
			t.buffer.CursorMoveRight()
			t.moveViewPort(dirRight, dirDown)

		case tcell.KeyUp:
			t.buffer.CursorMoveUp()
			t.moveViewPort(dirRight, dirUp)

		case tcell.KeyDown:
			t.buffer.CursorMoveDown()
			t.moveViewPort(dirRight, dirDown)

		case tcell.KeyEsc:
			t.actionSwitchMode(ModeNormal)
			if t.appending {
				t.buffer.CursorMoveLeft()
				t.moveViewPort(dirLeft, dirDown)
				t.appending = false
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
				t.moveViewPort(dirLeft, dirUp)

			case 'j':
				t.buffer.CursorMoveDown()
				t.moveViewPort(dirRight, dirDown)

			case 'k':
				t.buffer.CursorMoveUp()
				t.moveViewPort(dirRight, dirUp)

			case 'g':
				t.buffer.CursorMoveBOF()
				t.moveViewPort(dirRight, dirUp)

			case 'G':
				t.buffer.CursorMoveEOF()
				t.moveViewPort(dirRight, dirDown)

			case 'D':
				t.buffer.RemoveCurrentLine()
				t.moveViewPort(dirRight, dirDown)

			case 'h':
				t.buffer.CursorMoveLeft()
				t.moveViewPort(dirLeft, dirDown)

			case 'l':
				t.buffer.CursorMoveRight()
				t.moveViewPort(dirRight, dirDown)

			case 'x':
				t.buffer.Remove()

			case 'A':
				t.buffer.CursorMoveEOL()
				t.moveViewPort(dirRight, dirDown)
				t.actionSwitchMode(ModeInsert)
				t.appending = t.buffer.CursorMoveRight()

			case 'a':
				t.actionSwitchMode(ModeInsert)
				t.appending = t.buffer.CursorMoveRight()

			case 'o':
				t.buffer.CursorMoveEOL()
				t.buffer.Append(10)
				t.actionSwitchMode(ModeInsert)
				t.moveViewPort(dirLeft, dirDown)

			case '$':
				t.buffer.CursorMoveEOL()
				t.moveViewPort(dirRight, dirDown)

			case '^', '0':
				t.buffer.CursorMoveBOL()
				t.moveViewPort(dirRight, dirDown)

			}
		case tcell.KeyLeft:
			t.buffer.CursorMoveLeft()
			t.moveViewPort(dirLeft, dirDown)

		case tcell.KeyRight:
			t.buffer.CursorMoveRight()
			t.moveViewPort(dirRight, dirDown)

		case tcell.KeyUp:
			t.buffer.CursorMoveUp()
			t.moveViewPort(dirRight, dirUp)

		case tcell.KeyDown:
			t.buffer.CursorMoveDown()
			t.moveViewPort(dirRight, dirDown)

		case tcell.KeyCtrlE:
			t.buffer.CursorMoveEOL()
			t.moveViewPort(dirRight, dirDown)

		case tcell.KeyCtrlA:
			t.buffer.CursorMoveBOL()
			t.moveViewPort(dirRight, dirDown)

		}
	})
}

// Draw none
func (t *TextArea) Draw(screen tcell.Screen) {
	x, y, w, h := t.Box.GetInnerRect()
	t.Box.Draw(screen)
	runes := t.buffer.Runes()
	lines := t.buffer.Lines()

	curX, curLine := t.buffer.Position()
	if t.buffer.atEnd {
		curX--
	}

	for yOffset := 0; yOffset < h; yOffset++ {
		var line *Line
		if yOffset+t.viewPort.yScroll < len(lines) {
			line = &lines[yOffset+t.viewPort.yScroll]
		}

		for xOffset := 0; xOffset < w; xOffset++ {
			r := rune(0)
			style := textStyle
			cursor := -1

			if line != nil {
				cursor = line.Start + xOffset + t.viewPort.xScroll
				if cursor < len(runes) && cursor <= line.End {
					r = runes[cursor]
				}
			}

			if curLine == yOffset+t.viewPort.yScroll && curX == cursor {
				style = cursorStyle
			}

			if r == 10 {
				screen.SetContent(x+xOffset, y+yOffset, 0, nil, style)
				continue
			}

			screen.SetContent(x+xOffset, y+yOffset, r, nil, style)
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
	t.moveViewPort(dirRight, dirDown)
}

func (t *TextArea) actionBreak() {
	t.buffer.Append(10)
}

func (t *TextArea) actionSwitchMode(m Mode) {
	t.mode = m
}

type xdir func(cursor, start, end, xScroll int) (bool, int)
type ydir func(lindex, yScroll, h int) (bool, int)

func dirLeft(cursor, start, end, xScroll int) (bool, int) {
	return cursor >= start+xScroll, cursor - start
}

func dirRight(cursor, start, end, xScroll int) (bool, int) {
	return cursor < end, cursor - end + 1
}

func dirUp(lindex, yScroll, h int) (bool, int) {
	return yScroll < lindex, lindex
}

func dirDown(lindex, yScroll, h int) (bool, int) {
	return lindex+1-yScroll < h, lindex - h + 1
}

func (t *TextArea) moveViewPort(x xdir, y ydir) {
	cursor, lindex := t.buffer.Position()
	lines := t.buffer.Lines()
	if len(lines) == 0 {
		return
	}

	line := lines[lindex]
	_, _, w, h := t.Box.GetInnerRect()

	start := line.Start
	end := start + w

	if ok, offset := x(cursor, start, end, t.viewPort.xScroll); !ok {
		t.viewPort.xScroll = offset
	}

	if ok, offset := y(lindex, t.viewPort.yScroll, h); !ok {
		t.viewPort.yScroll = offset
	}
}

func printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
