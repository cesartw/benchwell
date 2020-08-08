package gtk

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io"
	"sync"

	"bitbucket.org/goreorto/benchwell/assets"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/quick"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// GTK_* consts
const (
	ICON_SIZE_TAB           = 10
	ICON_SIZE_MENU          = 16
	ICON_SIZE_SMALL_TOOLBAR = 16
	ICON_SIZE_LARGE_TOOLBAR = 24
	ICON_SIZE_BUTTON        = 16
	ICON_SIZE_DND           = 32
	ICON_SIZE_DIALOG        = 48
)

func init() {
	// Registrering pango formatter
	formatters.Register("pango", chroma.FormatterFunc(pangoFormatter))
}

type MVar struct {
	value interface{}
	sync.RWMutex
}

func (mv *MVar) Set(v interface{}) {
	mv.Lock()
	defer mv.Unlock()

	mv.value = v
}

func (mv *MVar) Get() interface{} {
	mv.RLock()
	defer mv.RUnlock()

	return mv.value
}

func BWAddClass(i interface {
	GetStyleContext() (*gtk.StyleContext, error)
}, class string) error {
	style, err := i.GetStyleContext()
	if err != nil {
		return err
	}
	style.AddClass(class)
	return nil
}

func BWRemoveClass(i interface {
	GetStyleContext() (*gtk.StyleContext, error)
}, class string) error {
	style, err := i.GetStyleContext()
	if err != nil {
		return err
	}
	style.RemoveClass(class)
	return nil
}

func menuItemWithImage(txt string, stockImage string) (*gtk.MenuItem, error) {
	item, err := gtk.MenuItemNew()
	if err != nil {
		return nil, err
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	icon, err := gtk.ImageNewFromIconName(stockImage, gtk.ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}

	label, err := gtk.LabelNew(txt)
	if err != nil {
		return nil, err
	}
	label.SetUseUnderline(true)
	label.SetXAlign(0.0)

	box.PackStart(icon, false, false, 0)
	box.PackEnd(label, true, true, 5)
	item.Add(box)

	return item, nil
}

func BWLabelNewWithClass(title, class string) (*gtk.Label, error) {
	label, err := gtk.LabelNew(title)
	if err != nil {
		return nil, err
	}

	style, err := label.GetStyleContext()
	if err != nil {
		return nil, err
	}
	style.AddClass(class)

	return label, nil
}

func BWMenuItemWithImage(txt string, asset string) (*gtk.MenuItem, error) {
	item, err := gtk.MenuItemNew()
	if err != nil {
		return nil, err
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	icon, err := BWImageNewFromFile(asset, ICON_SIZE_MENU)
	if err != nil {
		return nil, err
	}

	label, err := gtk.LabelNew(txt)
	if err != nil {
		return nil, err
	}
	label.SetUseUnderline(true)
	label.SetXAlign(0.0)

	box.PackStart(icon, false, false, 0)
	box.PackEnd(label, true, true, 5)
	item.Add(box)

	return item, nil
}

var loader *gdk.PixbufLoader

func BWImageNewFromFile(asset string, size int) (*gtk.Image, error) {
	pixbuf, err := BWPixbufFromFile(asset, size)
	if err != nil {
		return nil, err
	}

	img, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func BWPixbufFromFile(asset string, size int) (*gdk.Pixbuf, error) {
	loader, err := gdk.PixbufLoaderNewWithType("png")
	if err != nil {
		return nil, err
	}

	data, ok := assets.Iconset48[asset]
	if !ok {
		return nil, fmt.Errorf("`%s` icon not found", asset)
	}

	pixbuf, err := loader.WriteAndReturnPixbuf(data)
	if err != nil {
		return nil, err
	}

	pixbuf, err = pixbuf.ScaleSimple(size, size, gdk.INTERP_BILINEAR)
	if err != nil {
		return nil, err
	}

	return pixbuf, nil
}

func BWButtonNewFromIconName(asset string, size int) (*gtk.Button, error) {
	btn, err := gtk.ButtonNew()
	if err != nil {
		return nil, err
	}

	img, err := BWImageNewFromFile(asset, size)
	if err != nil {
		return nil, err
	}
	btn.SetImage(img)

	return btn, nil
}

func BWRadioButtonNew(label string, l *glib.SList) (*gtk.RadioButton, *glib.SList, error) {
	if l == nil {
		l = &glib.SList{}
	}

	radio, err := gtk.RadioButtonNewWithLabel(l, label)
	if err != nil {
		return nil, nil, err
	}
	l, err = radio.GetGroup()
	if err != nil {
		return nil, nil, err
	}

	return radio, l, nil
}

func ChromaHighlight(theme string, inputString string) (out string, err error) {
	buff := new(bytes.Buffer)
	writer := bufio.NewWriter(buff)

	// Doing the job (io.Writer, SourceText, language(go), Lexer(pango), style(pygments))
	if err = quick.Highlight(writer, inputString, "sql", "pango", theme); err != nil {
		return
	}
	writer.Flush()
	return string(buff.Bytes()), err
}

func pangoFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
	var r, g, b uint8
	var closer, out string

	var getColour = func(color chroma.Colour) string {
		r, g, b = color.Red(), color.Green(), color.Blue()
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}

	for tkn := it(); tkn != chroma.EOF; tkn = it() {

		entry := style.Get(tkn.Type)
		if !entry.IsZero() {
			if entry.Bold == chroma.Yes {
				out = `<b>`
				closer = `</b>`
			}
			if entry.Underline == chroma.Yes {
				out += `<u>`
				closer = `</u>` + closer
			}
			if entry.Italic == chroma.Yes {
				out += `<i>`
				closer = `</i>` + closer
			}
			if entry.Colour.IsSet() {
				out += `<span foreground="` + getColour(entry.Colour) + `">`
				closer = `</span>` + closer
			}
			if entry.Background.IsSet() {
				out += `<span background="` + getColour(entry.Background) + `">`
				closer = `</span>` + closer
			}
			if entry.Border.IsSet() {
				out += `<span background="` + getColour(entry.Border) + `">`
				closer = `</span>` + closer
			}
			fmt.Fprint(w, out)
		}
		fmt.Fprint(w, html.EscapeString(tkn.Value))
		if !entry.IsZero() {
			fmt.Fprint(w, closer)
		}
		closer, out = "", ""
	}
	return nil
}

type CancelOverlay struct {
	*gtk.Overlay
	btnCancel *gtk.Button
	spinner   *gtk.Spinner
	box       *gtk.Box

	onCancel func()
}

func (c CancelOverlay) Init(widget gtk.IWidget) (*CancelOverlay, error) {
	var err error
	c.Overlay, err = gtk.OverlayNew()
	if err != nil {
		return nil, err
	}

	c.btnCancel, err = gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		return nil, err
	}
	c.btnCancel.SetSizeRequest(100, 30)

	c.spinner, err = gtk.SpinnerNew()
	if err != nil {
		return nil, err
	}

	c.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	box.SetSizeRequest(100, 150)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetVExpand(true)
	box.SetHExpand(true)

	c.box.Add(box)
	box.PackStart(c.spinner, true, true, 0)
	box.PackStart(c.btnCancel, false, false, 0)

	c.Add(widget)

	c.btnCancel.Connect("clicked", func() {
		c.Stop()
		c.onCancel()
	})

	return &c, nil
}

func (c *CancelOverlay) Run(onCancel func()) {
	c.box.ShowAll()
	c.spinner.Start()
	c.AddOverlay(c.box)
	c.onCancel = onCancel
}

func (c *CancelOverlay) Stop() {
	c.Remove(c.box)
	c.spinner.Stop()
}
