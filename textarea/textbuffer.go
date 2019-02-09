package textarea

type textBuffer struct {
	buffer        []rune
	w             int
	line, xOffset int
}

func newtextbuffer(w int) *textBuffer {
	if w == 0 {
		panic("width must be greater than 0")
	}

	return &textBuffer{
		buffer: make([]rune, 0),
		w:      w,
	}
}

func (t *textBuffer) Append(r rune) {
	if len(t.buffer) == 0 {
		t.buffer = append(t.buffer, r)
		t.xOffset = 1
		return
	}

	lines := t.Lines()
	cursor := t.Cursor()

	nbuffer := []rune{}
	nbuffer = append(nbuffer, t.buffer[:cursor]...)
	nbuffer = append(nbuffer, r)
	nbuffer = append(nbuffer, t.buffer[cursor:]...)

	t.buffer = nbuffer

	t.xOffset++

	// We reach the end of the line(width based)
	if t.xOffset == lines[t.line].Len()+1 && lines[t.line].Len()+1 == t.w {
		t.xOffset = 0
		t.line++
	}

	if r == 10 {
		t.xOffset = 0
		t.line++
	}
}

func (t *textBuffer) SetText(s string) {
	t.buffer = []rune(s)

	t.xOffset = 0
	t.line = 0

	lines := t.Lines()
	t.line = len(lines) - 1

	if len(lines) > 0 {
		t.xOffset = 0
	}
}

func (t *textBuffer) SetEditingWidth(w int) {
	if w == 0 {
		panic("width must be greater than 0")
	}

	t.w = w
}

func (t *textBuffer) Remove() bool {
	if len(t.buffer) == 0 {
		return false
	}

	cursor := t.Cursor()
	if cursor > len(t.buffer)-1 {
		return false
	}

	if t.buffer[cursor] == 10 && t.line > 0 {
		t.line--
	}

	t.buffer = append(t.buffer[:cursor], t.buffer[cursor+1:]...)

	if t.xOffset > len(t.buffer) {
		t.xOffset = len(t.buffer)
	}

	return true
}

func (t textBuffer) RuneAtCursor() rune {
	return t.buffer[t.Cursor()]
}

func (t textBuffer) Runes() []rune {
	return t.buffer
}

func (t textBuffer) Text() string {
	return string(t.buffer)
}

func (t textBuffer) Empty() bool {
	return len(t.buffer) == 0
}

func (t *textBuffer) CursorMoveLeft() bool {
	if t.xOffset == 0 {
		return false
	}

	t.xOffset--
	return true
}

func (t *textBuffer) CursorMoveRight() bool {
	t.xOffset++

	lines := t.Lines()
	if len(lines) == 0 {
		return false
	}

	line := t.Lines()[t.line]

	if t.xOffset >= line.Len() {
		t.xOffset = line.Len() - 1
		return false
	}

	return true
}

func (t *textBuffer) CursorMoveDown() bool {
	lines := t.Lines()
	t.line++

	if t.line >= len(lines) {
		t.line = len(lines) - 1
		return false
	}

	if len(lines) == 0 {
		t.xOffset = 0
		return false
	}

	line := lines[t.line]
	if line.Start+t.xOffset > line.End {
		t.xOffset = line.End - line.Start
	}

	return true
}

func (t *textBuffer) CursorMoveUp() bool {
	t.line--

	if t.line < 0 {
		t.line = 0
		return false
	}

	lines := t.Lines()
	if len(lines) == 0 {
		t.xOffset = 0
		return true
	}

	line := lines[t.line]
	if line.Start+t.xOffset > line.End {
		t.xOffset = line.End - line.Start
	}

	return true
}

func (t *textBuffer) CursorMoveEOL() {
	lines := t.Lines()
	if len(lines) == 0 {
		t.xOffset = 0
		t.line = 0
		return
	}
	line := lines[t.line]
	t.xOffset = line.End - line.Start
}

func (t *textBuffer) CursorMoveBOL() {
	t.xOffset = 0
}

func (t textBuffer) Cursor() int {
	lines := t.Lines()
	if len(lines) > 0 {
		return lines[t.line].Start + t.xOffset
	}

	return t.xOffset

	//cursor := t.xOffset
	//if cursor == -1 {
	//cursor = 0
	//}

	//lines := t.Lines()
	//for i := 0; i < t.line; i++ {
	//cursor += (lines[i].End - lines[i].Start + 1)
	//}

	//if t.HasEOFLF() && cursor == len(t.buffer) {
	//return cursor - 1
	//}

	//return cursor
}

func (t *textBuffer) SetCursor(c int) {
	lines := t.Lines()
	if len(lines) == 0 {
		return
	}

	if c == 0 {
		t.xOffset = 0
		t.line = 0
	}

	for i, l := range lines {
		diff := l.End - l.Start
		if diff > c {
			t.xOffset = diff - c
			t.line = i
			return
		}
		c = c - diff
	}
}

func (t textBuffer) Position() (int, int) {
	return t.xOffset, t.line
}

// Line ..
type Line struct {
	Start, End int
}

// Len ...
func (l Line) Len() int {
	return l.End - l.Start + 1
}

func (t textBuffer) Lines() []Line {
	lines := []Line{}

	if len(t.buffer) == 0 {
		return []Line{}
	}

	l := Line{}
	for i := 0; i < len(t.buffer); i++ {
		l.End = i

		if t.buffer[i] == 10 {
			lines = append(lines, l)
			l = Line{Start: i + 1}

			continue
		}

		// runes since the last break
		if i-l.Start > 0 && (i-l.Start+1)%t.w == 0 {
			lines = append(lines, l)
			l = Line{Start: i + 1}

			continue
		}

		if i == len(t.buffer)-1 {
			lines = append(lines, l)
		}
	}

	if t.HasEOFLF() {
		l := lines[len(lines)-1]
		if l.End != l.Start {
			lines = append(lines, Line{Start: len(t.buffer) - 1, End: len(t.buffer) - 1})
		}
	}

	return lines
}

// HasEOFLF ...
func (t textBuffer) HasEOFLF() bool {
	if len(t.buffer) == 0 {
		return false
	}

	return t.buffer[len(t.buffer)-1] == 10
}
