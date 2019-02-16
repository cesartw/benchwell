package textarea

import (
	"fmt"
	"reflect"
	"testing"
)

func TestComplex(t *testing.T) {
	subject := newtextbuffer(10)

	subject.Append('a') // a_ _
	subject.Append('c') // ac_ _
	if string(subject.Runes()) != "ac" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ac")
	}

	subject.Remove() // ac_ _
	if string(subject.Runes()) != "ac" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ac")
	}

	subject.CursorMoveLeft()
	subject.Remove() // _a_
	if string(subject.Runes()) != "a" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "a")
	}

	subject.Append('d') // ad_ _
	subject.Append('e') // ade_ _
	subject.Append('f') // adef_ _
	if string(subject.Runes()) != "adef" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "adef")
	}

	subject.CursorMoveLeft() // ade_f_
	subject.CursorMoveLeft() // ad_e_f
	subject.Remove()         // ad_f_
	if string(subject.Runes()) != "adf" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "adf")
	}

	subject.Append('g') // adg_f_
	subject.Append('h') // adgh_f_
	if string(subject.Runes()) != "adghf" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "adghf")
	}

	subject.CursorMoveLeft() // adg_h_f
	subject.CursorMoveLeft() // ad_g_hf
	subject.CursorMoveLeft() // a_d_ghf
	subject.Remove()         // a_g_hf
	if string(subject.Runes()) != "aghf" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "aghf")
	}

	subject.Remove() // a_h_f
	if string(subject.Runes()) != "ahf" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ahf")
	}

	subject.CursorMoveRight() // ah_f_
	subject.Remove()          // ah_ _
	if string(subject.Runes()) != "ah" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ah")
	}

	subject.Remove() // ah_ _
	if string(subject.Runes()) != "ah" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ah")
	}

	subject.CursorMoveLeft() // a_h_
	subject.Remove()         // a_ _
	if string(subject.Runes()) != "a" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "a")
	}

	subject.Append('b') // ab_ _
	if string(subject.Runes()) != "ab" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ab")
	}

	subject.Append(10)  // ab\n_ _
	subject.Append('c') // ab\nc_ _
	if string(subject.Runes()) != "ab\nc" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "ab\nc")
	}
}

func TestBasic(t *testing.T) {
	subject := newtextbuffer(10)

	subject.Append('c')
	subject.Append('e')
	subject.Append('s')
	subject.Append('a')
	subject.Append('r')
	subject.Append(10)
	subject.Append('c')
	subject.Append('e')
	subject.Append('j')
	subject.Append('u')
	subject.Append('d')
	subject.Append('o')
	if string(subject.Runes()) != "cesar\ncejudo" {
		t.Errorf("got: `%s`, expected %s", string(subject.Runes()), "cesar\ncejudo")
	}
}

func TestAppend(t *testing.T) {
	testcases := []struct {
		name     string
		text     []rune
		cursor   int
		expected []rune
	}{
		{
			text:     []rune{},
			cursor:   0,
			expected: []rune{'a', 'b', 'c'},
		},
		{
			text:     []rune{'f', 'f', 'f'},
			cursor:   0,
			expected: []rune{'a', 'b', 'c', 'f', 'f', 'f'},
		},
		{
			text:     []rune{'f', 'f', 'f'},
			cursor:   1,
			expected: []rune{'f', 'a', 'b', 'c', 'f', 'f'},
		},
		{
			text:     []rune{'f', 'f', 'f'},
			cursor:   2,
			expected: []rune{'f', 'f', 'a', 'b', 'c', 'f'},
		},
		{
			text:     []rune{'f', 'f', 'f'},
			cursor:   3,
			expected: []rune{'f', 'f', 'f', 'a', 'b', 'c'},
		},
	}

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			subject := newtextbuffer(6)
			subject.buffer = tc.text
			subject.cursor = tc.cursor

			subject.Append('a')
			subject.Append('b')
			subject.Append('c')

			if !reflect.DeepEqual(subject.buffer, tc.expected) {
				t.Fatalf("\nexpected\t%s\ngot\t\t%s", string(tc.expected), string(subject.buffer))
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testcases := []struct {
		name     string
		text     []rune
		cursor   int
		expected []rune
	}{
		{
			name:     "from empty",
			text:     []rune{},
			cursor:   0,
			expected: []rune{},
		},
		{
			name:     "first rune",
			text:     []rune{'a'},
			cursor:   0,
			expected: []rune{},
		},
		{
			name:     "in between",
			text:     []rune{'a', 'b', 'c', 'd'},
			cursor:   3,
			expected: []rune{'a', 'b', 'c'},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			subject := newtextbuffer(10)
			subject.buffer = tc.text
			subject.cursor = tc.cursor

			subject.Remove()
			subject.Remove()

			if !reflect.DeepEqual(subject.buffer, tc.expected) {
				t.Fatalf("\nexpected\t%s\ngot\t\t%s", string(tc.expected), string(subject.buffer))
			}
		})
	}
}

func TestLines(t *testing.T) {
	testcases := []struct {
		name     string
		text     string
		expected []Line
	}{
		{
			name:     "blank",
			text:     "",
			expected: []Line{},
		},
		{
			name:     "single line",
			text:     "adc",
			expected: []Line{{Start: 0, End: 2}},
		},
		{
			name:     "with lf",
			text:     "adc\n",
			expected: []Line{{Start: 0, End: 3}, {Start: 3, End: 3}},
		},
		{
			name:     "mix lf and widthbr",
			text:     "01234567890123456789\n",
			expected: []Line{{Start: 0, End: 9}, {Start: 10, End: 19}, {Start: 20, End: 20}},
		},
		{
			name: "width breaks",
			text: "01234567890123456789\nadb",
			expected: []Line{
				{Start: 0, End: 9},
				{Start: 10, End: 19},
				{Start: 20, End: 20},
				{Start: 21, End: 23},
			},
		},
		{
			name: "LF matches W",
			text: "01234567890123456789\nadb\n\nrgrg\n",
			expected: []Line{
				{Start: 0, End: 9},
				{Start: 10, End: 19},
				{Start: 20, End: 20},
				{Start: 21, End: 24},
				{Start: 25, End: 25},
				{Start: 26, End: 30},
				{Start: 30, End: 30},
			},
		},
		{
			name: "blank lines",
			text: "0123456789012345678\nadb\n\nrgrg\n",
			expected: []Line{
				{Start: 0, End: 9},
				{Start: 10, End: 19},
				{Start: 20, End: 23},
				{Start: 24, End: 24},
				{Start: 25, End: 29},
				{Start: 29, End: 29},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			subject := newtextbuffer(10)
			subject.buffer = []rune(tc.text)

			lines := subject.Lines()

			if !reflect.DeepEqual(lines, tc.expected) {
				t.Fatalf("\nexpected\t%+v\ngot\t\t%+v", tc.expected, lines)
			}
		})
	}
}

func TestNav(t *testing.T) {
	subject := newtextbuffer(10)
	subject.SetText(
		"123456789labc\ndef\nghijk\n\nlmnopqr\n")

	/*
			123456789l  0-9
			abc        10-13
			def        14-17
			ghijk      18-23
		               24-24
			lmnopqr    25-32
			           32-32
	*/

	subject.CursorMoveUp()
	if subject.RuneAtCursor() != 'l' {
		t.Fatalf("expected l got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveDown()
	subject.CursorMoveDown()
	if subject.RuneAtCursor() != 10 {
		t.Fatalf("expected 10 got %s", string(subject.RuneAtCursor()))
	}
	subject.CursorMoveLeft()
	if subject.RuneAtCursor() != 10 {
		t.Fatalf("expected 10 got %s", string(subject.RuneAtCursor()))
	}
	subject.CursorMoveRight()
	if subject.RuneAtCursor() != 10 {
		t.Fatalf("expected 10 got %s", string(subject.RuneAtCursor()))
	}
	subject.CursorMoveEOL()
	if subject.RuneAtCursor() != 10 {
		t.Fatalf("expected 10 got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveUp()
	subject.CursorMoveEOL()
	subject.CursorMoveLeft()
	subject.CursorMoveLeft()
	subject.CursorMoveLeft()
	subject.CursorMoveLeft()
	subject.CursorMoveLeft()
	subject.CursorMoveLeft()
	subject.CursorMoveLeft()
	if subject.RuneAtCursor() != 'l' {
		t.Fatalf("expected l, got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveLeft()
	if subject.RuneAtCursor() != 'l' {
		t.Fatalf("expected l, got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveRight()
	subject.CursorMoveRight()
	subject.CursorMoveRight()
	if subject.RuneAtCursor() != 'o' {
		t.Fatalf("expected o, got %s", string(subject.RuneAtCursor()))
	}

	fmt.Println(subject.atEnd, subject.cursor)
	subject.CursorMoveUp()
	if subject.RuneAtCursor() != 10 {
		fmt.Println(subject.atEnd, subject.cursor)
		t.Fatalf("expected LF(10), got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveUp()
	if subject.RuneAtCursor() != 'g' {
		t.Fatalf("expected g, got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveRight()
	subject.CursorMoveRight()
	subject.CursorMoveUp()
	if subject.RuneAtCursor() != 'f' {
		t.Fatalf("expected f, got %s", string(subject.RuneAtCursor()))
	}

	/*
			123456789l  0-9
			abc        10-13
			def        14-17
			ghijk      18-23
		               24-24
			lmnopqr    25-31
			           32-32
	*/

	subject.CursorMoveDown()
	if subject.RuneAtCursor() != 'i' {
		t.Fatalf("expected i, got %s", string(subject.RuneAtCursor()))
	}
	subject.CursorMoveUp()
	subject.CursorMoveUp()
	subject.CursorMoveUp()
	if subject.RuneAtCursor() != '3' {
		t.Fatalf("expected 3, got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveRight()
	subject.CursorMoveRight()
	if subject.RuneAtCursor() != '5' {
		t.Fatalf("expected 5, got %s", string(subject.RuneAtCursor()))
	}

	subject.CursorMoveDown()
	if subject.RuneAtCursor() != 10 {
		t.Fatalf("expected LF(10), got %s", string(subject.RuneAtCursor()))
	}
}
