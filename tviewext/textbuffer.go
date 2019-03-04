package tviewext

type textBuffer struct {
	buffer []rune
	cursor int
	atEnd  bool
}

func newtextbuffer(w int) *textBuffer {
	return &textBuffer{
		buffer: make([]rune, 0),
	}
}

func (t *textBuffer) Append(r rune) {
	t.atEnd = false

	if t.cursor == len(t.buffer) && r == 10 {
		t.atEnd = true
	}

	if len(t.buffer) == 0 {
		t.buffer = append(t.buffer, r)
		t.cursor++
		return
	}

	nbuffer := []rune{}
	nbuffer = append(nbuffer, t.buffer[:t.cursor]...)
	nbuffer = append(nbuffer, r)
	nbuffer = append(nbuffer, t.buffer[t.cursor:]...)

	t.buffer = nbuffer
	t.cursor++
}

func (t *textBuffer) SetText(s string) {
	t.buffer = []rune(s)
	t.CursorMoveEOF()
}

func (t *textBuffer) Remove() bool {
	if len(t.buffer) == 0 {
		return false
	}

	cursor := t.cursor // t.Cursor()
	if cursor > len(t.buffer)-1 {
		return false
	}

	t.buffer = append(t.buffer[:cursor], t.buffer[cursor+1:]...)

	if t.cursor > len(t.buffer) {
		t.cursor = len(t.buffer)
	}

	if t.cursor == len(t.buffer) && t.HasEOFLF() {
		t.atEnd = true
	}

	return true
}

func (t *textBuffer) RemoveCurrentLine() bool {
	if len(t.buffer) == 0 {
		return false
	}

	_, lindex := t.Position()
	line := t.Lines()[lindex]

	t.buffer = append(t.buffer[:line.Start], t.buffer[line.End+1:]...)

	if t.cursor > len(t.buffer) {
		t.cursor = len(t.buffer)
	}

	if t.cursor == len(t.buffer) && t.HasEOFLF() {
		t.atEnd = true
	}

	return true
}

func (t textBuffer) RuneAtCursor() rune {
	if t.cursor >= len(t.buffer) {
		return 0
	}

	return t.buffer[t.cursor]
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
	if t.cursor == 0 {
		return false
	}

	if t.atEnd {
		return false
	}

	t.cursor--
	// Prevent from abandoning line
	if t.buffer[t.cursor] == 10 {
		t.cursor++
		return false
	}

	return true
}

func (t *textBuffer) CursorMoveRight() bool {
	if t.RuneAtCursor() == 10 {
		return false
	}

	t.cursor++
	if t.cursor > len(t.buffer) {
		t.cursor--
		return false
	}

	return true
}

func (t *textBuffer) CursorMoveDown() bool {
	lines := t.Lines()

	_, curL := t.Position()
	if curL == len(lines)-1 {
		return false
	}

	nextL := lines[curL+1]

	// last line and it a LF
	if nextL.Start == nextL.End && curL+1 == len(lines)-1 {
		t.atEnd = true
		t.cursor = len(t.buffer)
		return true
	}

	t.cursor = nextL.Start + t.cursor - lines[curL].Start
	if t.cursor > nextL.End {
		t.cursor = nextL.End
	}

	return true
}

func (t *textBuffer) CursorMoveUp() bool {
	if t.atEnd {
		t.atEnd = false
		lines := t.Lines()
		t.cursor = lines[len(lines)-2].Start
		return true
	}

	_, curL := t.Position()

	if curL == 0 {
		return false
	}

	lines := t.Lines()
	nextL := lines[curL-1]

	t.cursor = nextL.Start + t.cursor - lines[curL].Start
	if t.cursor > nextL.End {
		t.cursor = nextL.End
	}

	return true
}

func (t *textBuffer) CursorMoveEOL() {
	lines := t.Lines()
	_, curL := t.Position()

	t.cursor = lines[curL].End
}

func (t *textBuffer) CursorMoveBOL() {
	lines := t.Lines()
	_, curL := t.Position()

	t.cursor = lines[curL].Start
}

func (t *textBuffer) CursorMoveEOF() {
	t.cursor = len(t.buffer)
	if t.HasEOFLF() {
		t.atEnd = true
	}
}

func (t *textBuffer) CursorMoveBOF() {
	t.cursor = 0
	t.atEnd = false
}

func (t textBuffer) Cursor() int {
	return t.cursor
}

func (t textBuffer) Position() (cursor int, line int) {
	lines := t.Lines()
	if t.atEnd {
		return t.cursor, len(lines) - 1
	}

	for i, line := range lines {
		if line.Start <= t.cursor && line.End >= t.cursor {
			return t.cursor, i
		}
	}

	if len(lines) > 0 {
		return t.cursor, len(lines) - 1
	}

	return t.cursor, 0
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
		return []Line{{}}
	}

	l := Line{}
	for i := 0; i < len(t.buffer); i++ {
		l.End = i

		if t.buffer[i] == 10 {
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
