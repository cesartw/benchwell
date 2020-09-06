package gtk

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"bitbucket.org/goreorto/benchwell/assets"
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

	var iconOff, iconOn *gtk.Image
	if strings.HasPrefix(asset, "gtk-") {
		iconOff, err = gtk.ImageNewFromIconName(asset, gtk.ICON_SIZE_MENU)
		if err != nil {
			return nil, err
		}
	} else {
		iconOff, err = BWImageNewFromFile(asset, "orange", ICON_SIZE_MENU)
		if err != nil {
			return nil, err
		}
		iconOn, err = BWImageNewFromFile(asset, "white", ICON_SIZE_MENU)
		if err != nil {
			return nil, err
		}
		item.Connect("select", func() {
			iconOff.Hide()
			iconOn.Show()
		})
		item.Connect("deselect", func() {
			iconOn.Hide()
			iconOff.Show()
		})
	}

	label, err := gtk.LabelNew(txt)
	if err != nil {
		return nil, err
	}
	label.SetUseUnderline(true)
	label.SetXAlign(0.0)

	box.PackStart(iconOff, false, false, 0)
	if iconOn != nil {
		box.PackStart(iconOn, false, false, 0)
		iconOn.Hide()
	}
	box.PackEnd(label, true, true, 5)
	item.Add(box)

	label.Show()
	iconOff.Show()
	box.Show()
	item.Show()

	return item, nil
}

func BWImageNewFromFile(asset, color string, size int) (*gtk.Image, error) {
	pixbuf, err := BWPixbufFromFile(asset, color, size)
	if err != nil {
		return nil, err
	}

	img, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func BWPixbufFromFile(asset, color string, size int) (*gdk.Pixbuf, error) {
	loader, err := gdk.PixbufLoaderNewWithType("png")
	if err != nil {
		return nil, err
	}
	colors := assets.IconsetOrange48
	if color == "white" {
		colors = assets.IconsetWhite48
	}

	data, ok := colors[asset]
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

func BWButtonNewFromIconName(asset, color string, size int) (*gtk.Button, error) {
	btn, err := gtk.ButtonNew()
	if err != nil {
		return nil, err
	}

	img, err := BWImageNewFromFile(asset, color, size)
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

type OptionButton struct {
	*gtk.Grid
	menubtn *gtk.MenuButton
	menu    *glib.Menu
	btn     *gtk.Button
	actions map[string]*glib.SimpleAction
}

func (o *OptionButton) ConnectAction(action string, v interface{}) {
	o.actions[action].Connect("activate", v)
}

func BWOptionButtonNew(label string, w *Window, options []string) (*OptionButton, error) {
	var err error

	if len(options)%2 != 0 {
		return nil, errors.New("options must be pair slice")
	}

	ob := &OptionButton{actions: map[string]*glib.SimpleAction{}}
	ob.menubtn, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}

	ob.menu = glib.MenuNew()
	for i := 0; i < len(options)-1; i = i + 2 {
		actionName := options[i+1]
		actionLabel := options[i]

		ob.actions[actionName] = glib.SimpleActionNew(strings.TrimPrefix(actionName, "win."), nil)
		w.AddAction(ob.actions[actionName])
		ob.menu.Append(actionLabel, actionName)
	}

	ob.menubtn.SetMenuModel(&ob.menu.MenuModel)

	ob.btn, err = gtk.ButtonNewWithLabel(label)
	if err != nil {
		return nil, err
	}

	ob.Grid, err = gtk.GridNew()
	if err != nil {
		return nil, err
	}
	ob.Grid.Attach(ob.btn, 0, 0, 2, 1)
	ob.Grid.Attach(ob.menubtn, 2, 0, 1, 1)
	BWAddClass(ob.Grid, "linked")

	return ob, nil
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
