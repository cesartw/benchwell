package gtk

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type tabCtrl interface {
	Close()               //tab
	Removed()             //tab
	Title() string        //tab
	Content() gtk.IWidget //tab
	SetFileText(string)   //tab
	OnCloseTab()
	SetWindowCtrl(interface{}) // tab
}

type ToolTab struct {
	tabCtrl

	w         *Window
	label     *gtk.Label
	btn       *gtk.Button
	content   *gtk.Box
	mainW     gtk.IWidget
	header    *gtk.Box
	btnHandle glib.SignalHandle
}

type ToolTabOptions struct {
	Title   string
	Content gtk.IWidget
	Ctrl    tabCtrl
}

func (t ToolTab) Init(w *Window) (*ToolTab, error) {
	var err error
	t.w = w

	t.content, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	//t.content.PackStart(opts.Content, true, true, 0)
	t.content.SetVExpand(true)
	t.content.SetHExpand(true)
	t.content.Show()

	t.header, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	t.label, err = gtk.LabelNew("")
	if err != nil {
		return nil, err
	}

	t.btn, err = BWButtonNewFromIconName("close", ICON_SIZE_TAB)
	if err != nil {
		return nil, err
	}
	t.btn.SetRelief(gtk.RELIEF_NONE)

	t.header.PackStart(t.label, true, true, 0)
	t.header.PackEnd(t.btn, false, false, 0)
	t.header.ShowAll()

	return &t, nil
}

func (t *ToolTab) SetWindowCtrl(
	ctrl interface {
		OnCloseTab()
	},
) {
	t.tabCtrl.SetWindowCtrl(ctrl)
}

func (t *ToolTab) SetContent(opts ToolTabOptions) {
	if opts.Ctrl != nil {
		t.tabCtrl = opts.Ctrl
	}

	if opts.Title != "" {
		t.SetTitle(opts.Title)
	}

	if opts.Content != nil {
		if t.mainW != nil {
			t.content.Remove(t.mainW)
		}

		if t.btnHandle > 0 {
			t.btn.HandlerDisconnect(t.btnHandle)
		}

		t.content.PackStart(opts.Content, true, true, 0)
		t.mainW = opts.Content
		t.content.Show()
	}

	if opts.Ctrl != nil {
		t.btnHandle, _ = t.btn.Connect("clicked", t.OnCloseTab)
	}
}

func (t *ToolTab) SetTitle(title string) {
	t.label.SetText(title)
}

func (t *ToolTab) Label() *gtk.Box {
	return t.header
}

func (t *ToolTab) Content() *gtk.Box {
	return t.content
}
