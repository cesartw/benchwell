package gtk

import (
	"bitbucket.org/goreorto/benchwell/config"
	"github.com/google/uuid"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

type tabCtrl interface {
	Close()               //tab
	Removed()             //tab
	Content() gtk.IWidget //tab
	Title() string        //tab
	OnCloseTab(id string)
	SetWindowCtrl(interface{}) // tab
}

type ToolTab struct {
	tabCtrl

	id        string
	w         *Window
	label     *gtk.Label
	btn       *gtk.Button
	content   *gtk.Box
	mainW     gtk.IWidget
	header    *gtk.Box
	btnHandle glib.SignalHandle
}

type ToolTabOptions struct {
	Content gtk.IWidget
	Ctrl    tabCtrl
}

func (t ToolTab) Init(w *Window) (*ToolTab, error) {
	defer config.LogStart("ToolTab.Init", nil)()
	t.id = uuid.New().String()

	var err error
	t.w = w

	t.content, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	t.content.Show()
	t.content.SetVExpand(true)
	t.content.SetHExpand(true)

	t.header, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	t.header.Show()

	t.label, err = gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	t.label.Show()
	t.label.SetEllipsize(pango.ELLIPSIZE_END)
	t.label.SetWidthChars(20)

	t.btn, err = BWButtonNewFromIconName("close", "orange", ICON_SIZE_TAB)
	if err != nil {
		return nil, err
	}
	t.btn.Show()
	t.btn.SetRelief(gtk.RELIEF_NONE)
	t.btn.Connect("clicked", t.onCloseTab)

	t.header.PackStart(t.label, true, true, 0)
	t.header.PackEnd(t.btn, false, false, 0)

	return &t, nil
}

func (t *ToolTab) onCloseTab() {
	defer config.LogStart("ToolTab.onCloseTab", nil)()
	t.OnCloseTab(t.id)
}

func (t *ToolTab) SetWindowCtrl(
	ctrl interface {
		OnCloseTab(string)
	},
) {
	defer config.LogStart("ToolTab.SetWindowCtrl", nil)()

	t.tabCtrl.SetWindowCtrl(ctrl)
}

func (t *ToolTab) SetContent(opts ToolTabOptions) {
	defer config.LogStart("ToolTab.SetContent", nil)()

	t.tabCtrl = opts.Ctrl
	t.content.PackStart(opts.Content, true, true, 0)
	t.mainW = opts.Content
	t.content.Show()
}

func (t *ToolTab) SetTitle(title string) {
	defer config.LogStart("ToolTab.SetTitle", nil)()

	t.w.nb.SetMenuLabelText(t.Content(), title)
	t.label.SetText(title)
}

func (t *ToolTab) Header() *gtk.Box {
	defer config.LogStart("ToolTab.Header", nil)()

	return t.header
}

func (t *ToolTab) Content() *gtk.Box {
	defer config.LogStart("ToolTab.Content", nil)()

	return t.content
}
