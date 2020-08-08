package gtk

import (
	"bitbucket.org/goreorto/benchwell/config"
	"github.com/gotk3/sourceview"
)

type SourceView struct {
	*sourceview.SourceView
}

type SourceViewOptions struct {
	Highlight bool
	Undoable  bool
	Language  string
}

func (t SourceView) Init(_ *Window, opts SourceViewOptions, ctrl interface{ Config() *config.Config }) (*SourceView, error) {
	var err error
	t.SourceView, err = sourceview.SourceViewNew()

	buff, err := t.SourceView.GetBuffer()
	if err != nil {
		return nil, err
	}

	err = t.SetLanguage(opts.Language)
	if err != nil {
		return nil, err
	}

	style, err := sourceview.SourceStyleSchemeManagerGetDefault()
	if err != nil {
		return nil, err
	}
	buff.SetStyleScheme(style.GetScheme("benchwell_dark"))

	return &t, err
}

func (t *SourceView) SetLanguage(lang string) error {
	buff, err := t.SourceView.GetBuffer()
	if err != nil {
		return err
	}
	if lang == "" {
		buff.SetLanguage(nil)
		return nil
	}

	language, err := sourceview.SourceLanguageManagerGetDefault()
	if err != nil {
		return err
	}

	l, err := language.GetLanguage(lang)
	if err != nil {
		return err
	}

	buff.SetLanguage(l)
	return nil
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
