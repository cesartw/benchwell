package gtk

import (
	"fmt"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var WHITESPACE = [2]string{" ", "\t"}

func isWhitespace(s string) bool {
	for _, w := range WHITESPACE {
		if w == s {
			return true
		}
	}
	return false
}

type TextView struct {
	*gtk.TextView
	buffer *TextBuffer
}

type TextViewOptions struct {
	Highlight bool
	Undoable  bool
}

func (t TextView) Init(opts TextViewOptions) (*TextView, error) {
	var err error
	t.TextView, err = gtk.TextViewNew()
	if err != nil {
		return nil, err
	}

	t.buffer, err = TextBuffer{}.Init(opts.Undoable, opts.Highlight)
	if err != nil {
		return nil, err
	}

	t.TextView.SetBuffer(t.buffer.TextBuffer)

	t.Connect("key-press-event", t.onTextViewKeyPress) // ctrl+enter exec query

	return &t, nil
}

func (t *TextView) onTextViewKeyPress(_ *gtk.TextView, e *gdk.Event) bool {
	keyEvent := gdk.EventKeyNewFromEvent(e)

	if keyEvent.KeyVal() == gdk.KEY_z && keyEvent.State()&gdk.CONTROL_MASK > 0 {
		t.buffer.Undo()
		return true
	}
	if keyEvent.KeyVal() == gdk.KEY_Z && keyEvent.State()&gdk.CONTROL_MASK > 0 {
		t.buffer.Redo()
		return true
	}

	return false
}

type stack []interface{}

func (s *stack) Pop() (interface{}, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	i := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return i, true
}

func (s *stack) Append(i interface{}) {
	*s = append(*s, i)
}

func (s stack) String() string {
	out := ""
	for _, i := range s {
		switch v := i.(type) {
		case undoableInsert:
			out = out + v.text
		case undoableDelete:
			out = out + "<"
		}

	}
	return out
}

type undoableInsert struct {
	offset    int
	text      string
	length    int
	mergeable bool
}

func (u undoableInsert) Init(iter *gtk.TextIter, text string, length int) undoableInsert {
	u.offset = iter.GetOffset()
	u.text = text
	u.length = length

	// TODO: it's broken with newlines
	//u.mergeable = length <= 1
	u.mergeable = false

	return u
}

type undoableDelete struct {
	text          string
	start         int
	end           int
	deleteKeyUsed bool
	mergeable     bool
}

func (u undoableDelete) Init(buffer *gtk.TextBuffer, start, end *gtk.TextIter) (undoableDelete, error) {
	var err error
	u.text, err = buffer.GetText(start, end, false)
	if err != nil {
		return u, err
	}

	u.start = start.GetOffset()
	u.end = end.GetOffset()

	u.deleteKeyUsed = buffer.GetIterAtMark(buffer.GetInsert()).GetOffset() <= u.start
	u.mergeable = !(u.end-u.start > 1 || u.text == "\r" || u.text == "\n" || u.text == " ")
	u.mergeable = false

	return u, nil
}

type TextBuffer struct {
	*gtk.TextBuffer

	undostack         stack
	redostack         stack
	inProgress        bool
	notUndoableAction bool
}

func (t TextBuffer) Init(undoable, highlight bool) (*TextBuffer, error) {
	var (
		err      error
		tagTable *gtk.TextTagTable
	)
	tagTable, err = gtk.TextTagTableNew()
	if err != nil {
		return nil, err
	}

	t.TextBuffer, err = gtk.TextBufferNew(tagTable)
	if err != nil {
		return nil, err
	}

	if undoable {
		t.Connect("insert-text", t.onInsertText)
		t.Connect("delete-range", t.onDeleteRange)
		t.undostack = stack{}
		t.redostack = stack{}
	}

	if highlight {
		t.Connect("end-user-action", t.onChanged)
	}

	return &t, nil
}

func (t *TextBuffer) onChanged() {
	if t.notUndoableAction {
		return
	}
	t.notUndoableAction = true
	t.highlight()
	t.notUndoableAction = false
}

func (t *TextBuffer) highlight() {
	txt, err := t.GetText(t.GetStartIter(), t.GetEndIter(), false)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	iter := t.GetIterAtMark(t.GetInsert())
	offset := iter.GetOffset()

	txt, err = ChromaHighlight(txt)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}
	t.Delete(t.GetStartIter(), t.GetEndIter())
	t.InsertMarkup(t.GetStartIter(), txt)
	t.PlaceCursor(t.GetIterAtOffset(offset))
}

func (t *TextBuffer) onInsertText(_ *gtk.TextBuffer, location *gtk.TextIter, text string, length int) {
	if !t.inProgress {
		t.redostack = stack{}
	}
	if t.notUndoableAction {
		return
	}

	undoAction := undoableInsert{}.Init(location, text, length)
	prevAction, ok := t.undostack.Pop()
	if !ok {
		t.undostack.Append(undoAction)
		return
	}

	prevInsert, ok := prevAction.(undoableInsert)
	if !ok {
		t.undostack.Append(prevInsert)
		t.undostack.Append(undoAction)
		return
	}

	if t.canMergeInserts(prevInsert, undoAction) {
		prevInsert.length += undoAction.length
		prevInsert.text += undoAction.text
		t.undostack.Append(prevInsert)
		return
	}

	t.undostack.Append(prevInsert)
	t.undostack.Append(undoAction)
}

func (t *TextBuffer) onDeleteRange(_ *gtk.TextBuffer, start, end *gtk.TextIter) {
	if !t.inProgress {
		t.redostack = stack{}
	}
	if t.notUndoableAction {
		return
	}

	undoAction, err := undoableDelete{}.Init(t.TextBuffer, start, end)
	if err != nil {
		config.Env.Log.Error(err)
		return
	}

	prevAction, ok := t.undostack.Pop()
	if !ok {
		t.undostack.Append(undoAction)
		return
	}

	prevDelete, ok := prevAction.(undoableDelete)
	if !ok {
		t.undostack.Append(prevAction)
		t.undostack.Append(undoAction)
		return
	}

	if t.canMergeDeletes(prevDelete, undoAction) {
		if prevDelete.start == undoAction.start {
			prevDelete.text += undoAction.text
			prevDelete.end += (undoAction.end - undoAction.start)
		} else {
			prevDelete.text = fmt.Sprintf("%s%s", undoAction.text, prevDelete.text)
			prevDelete.start = undoAction.start
		}
		t.undostack.Append(prevDelete)
		return
	}

	t.undostack.Append(prevDelete)
	t.undostack.Append(undoAction)
}

func (t *TextBuffer) canUndo() bool {
	return len(t.undostack) > 0
}

func (t *TextBuffer) canRedo() bool {
	return len(t.redostack) > 0
}

func (t *TextBuffer) canMergeInserts(prev, cur undoableInsert) bool {
	if !prev.mergeable || !cur.mergeable {
		return false
	}

	if cur.offset != (prev.offset + prev.length) {
		return false
	}

	if isWhitespace(cur.text) && !isWhitespace(prev.text) {
		return false
	}

	if isWhitespace(prev.text) && !isWhitespace(cur.text) {
		return false
	}

	return true
}

func (t *TextBuffer) canMergeDeletes(prev, cur undoableDelete) bool {
	if !cur.mergeable || !prev.mergeable {
		return false
	}

	if prev.deleteKeyUsed != cur.deleteKeyUsed {
		return false
	}

	if prev.start != cur.start && prev.start != cur.end {
		return false
	}

	if !isWhitespace(cur.text) && isWhitespace(prev.text) {
		return false
	}
	if isWhitespace(cur.text) && !isWhitespace(prev.text) {
		return false
	}

	return true
}

func (t *TextBuffer) Undo() {
	if len(t.undostack) == 0 {
		return
	}

	t.notUndoableAction = true
	t.inProgress = true
	undoAction, _ := t.undostack.Pop()
	t.redostack.Append(undoAction)

	undoInsert, ok := undoAction.(undoableInsert)
	if ok {
		start := t.GetIterAtOffset(undoInsert.offset)
		end := t.GetIterAtOffset(undoInsert.offset + undoInsert.length)
		t.Delete(start, end)
		t.PlaceCursor(start)
	} else {
		undoDelete, _ := undoAction.(undoableDelete)
		start := t.GetIterAtOffset(undoDelete.start)
		t.Insert(start, undoDelete.text)
		stop := t.GetIterAtOffset(undoDelete.end)

		if undoDelete.deleteKeyUsed {
			t.PlaceCursor(start)
		} else {
			t.PlaceCursor(stop)
		}
	}
	t.notUndoableAction = false
	t.inProgress = false
}

func (t *TextBuffer) Redo() {
	if len(t.redostack) == 0 {
		return
	}

	t.notUndoableAction = true
	t.inProgress = true
	redoAction, _ := t.redostack.Pop()

	t.undostack.Append(redoAction)
	redoInsert, ok := redoAction.(undoableInsert)
	if ok {
		start := t.GetIterAtOffset(redoInsert.offset)
		t.Insert(start, redoInsert.text)
		t.PlaceCursor(t.GetIterAtOffset(redoInsert.offset + redoInsert.length))
	} else {
		redoDelete := redoAction.(undoableDelete)
		start := t.GetIterAtOffset(redoDelete.start)
		end := t.GetIterAtOffset(redoDelete.end)
		t.Delete(start, end)
		t.PlaceCursor(start)
	}

	t.notUndoableAction = false
	t.inProgress = false
}
