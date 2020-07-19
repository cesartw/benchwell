package gtk

import (
	"github.com/gotk3/sourceview"
)

type SourceView struct {
	*sourceview.SourceView
}

type SourceViewOptions struct {
	Highlight bool
	Undoable  bool
}

func (t SourceView) Init(_ *Window, opts SourceViewOptions, ctrl textViewCtrl) (*SourceView, error) {
	var err error
	t.SourceView, err = sourceview.SourceViewNew()

	return &t, err
}

func (t *SourceView) ShowAutoComplete(words []string) {
	return
	// TODO: implement SourceCompletion
	//completion, err := t.SourceView.GetCompletion()
	//if err != nil {
	//return
	//}

	//buff, err := gtk.TextBufferNew(nil)
	//if err != nil {
	//return
	//}
	//buff.Insert(buff.GetStartIter(), strings.Join(words, " "))

	//provider, err := sourceview.SourceCompletionWordsNew("Tables", nil)
	//if err != nil {
	//return
	//}
	//provider.Register(buff)

	//sourcebuff, err := t.SourceView.GetBuffer()
	//if err != nil {
	//return
	//}

	//context, err := completion.CreateContext(sourcebuff.GetEndIter())
	//if err != nil {
	//return
	//}

	//completion.Show([]sourceview.ISourceCompletionProvider{provider}, context)
}
