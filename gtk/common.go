package gtk

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"bitbucket.org/goreorto/benchwell/assets"
	"bitbucket.org/goreorto/benchwell/config"
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

	if size > 0 {
		pixbuf, err = pixbuf.ScaleSimple(size, size, gdk.INTERP_BILINEAR)
		if err != nil {
			return nil, err
		}
	} else {

		pixbuf, err = pixbuf.ScaleSimple(52, 9, gdk.INTERP_BILINEAR)
		if err != nil {
			return nil, err
		}
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

func BWToggleButtonNewFromIconName(asset, color string, size int) (*gtk.ToggleButton, error) {
	btn, err := gtk.ToggleButtonNew()
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
	ob.menubtn.Show()

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
	ob.btn.Show()

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
	c.box.Show()
	c.spinner.Start()
	c.AddOverlay(c.box)
	c.onCancel = onCancel
}

func (c *CancelOverlay) Stop() {
	c.Remove(c.box)
	c.spinner.Stop()
}

type KeyValues struct {
	*gtk.Box
	kvs       reflect.Value
	keyvalues []*KeyValue
	onChange  func()
}

type kv interface {
	Name() string
	Val() string
	IsEnabled() bool
	SetName(string)
	SetVal(string)
	SetEnabled(bool)
}

func (c KeyValues) Init(values interface{}, onChange func()) (*KeyValues, error) {
	defer config.LogStart("KeyValues.Init", nil)()

	c.kvs = reflect.ValueOf(values)
	if c.kvs.Kind() != reflect.Ptr || c.kvs.Elem().Kind() != reflect.Slice {
		return nil, errors.New("values must be pointer to slice")
	}

	var err error
	c.onChange = onChange
	c.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}

	return &c, c.AddEmpty()
}

func (c *KeyValues) AddEmpty() error {
	defer config.LogStart("KeyValues.AddEmpty", nil)()

	kv, err := c.add(nil)
	if err != nil {
		return err
	}
	kv.Show()

	return nil
}

func (c *KeyValues) Add(v kv) error {
	defer config.LogStart("KeyValues.AddWithValues", nil)()

	// if there's an empty one already
	//for _, kv := range c.keyvalues {
	//s, _ := kv.key.GetText()
	//if s != "" {
	//continue
	//}
	//kv.key.SetText(v.Key())
	//kv.value.SetText(v.Value())
	//return nil
	//}

	kv, err := c.add(v)
	if err != nil {
		return err
	}
	kv.enabled.SetActive(v.IsEnabled())

	kv.key.SetText(v.Name())
	kv.value.SetText(v.Val())

	return nil
}

func (c *KeyValues) add(v kv) (*KeyValue, error) {
	defer config.LogStart("KeyValues.add", nil)()

	if v == nil {
		v = reflect.New(
			c.kvs.Type().Elem().Elem().Elem(),
		).Interface().(kv)
	}

	kV, err := KeyValue{}.Init(v)
	if err != nil {
		return nil, err
	}
	kV.Show()
	kV.remove.Connect("clicked", c.onRemove(kV))

	focused := func() {
		if c.keyvalues[len(c.keyvalues)-1] != kV {
			return
		}

		c.AddEmpty()
	}
	kV.key.Connect("grab-focus", focused)
	kV.value.Connect("grab-focus", focused)
	kV.key.Connect("key-release-event", c.onChange)
	kV.value.Connect("key-release-event", c.onChange)
	kV.enabled.Connect("toggled", c.onChange)
	kV.remove.Connect("clicked", c.onChange)

	c.Box.PackStart(kV, false, false, 0)

	c.keyvalues = append(c.keyvalues, kV)
	c.kvs.Elem().Set(
		reflect.Append(
			c.kvs.Elem(),
			reflect.ValueOf(v),
		),
	)

	return kV, nil
}

func (c *KeyValues) onRemove(kv *KeyValue) func() {
	defer config.LogStart("KeyValues.onRemove", nil)()

	return func() {
		c.Box.Remove(kv)

		for i, v := range c.keyvalues {
			if v != kv {
				continue
			}
			c.keyvalues = append(c.keyvalues[:i], c.keyvalues[i+1:]...)

			// TODO: splice c.kvs
			kvsStart := c.kvs.Elem().Slice(0, i)
			kvsEnd := c.kvs.Elem().Slice(i+1, c.kvs.Elem().Len())

			for i := 0; i < kvsEnd.Len(); i++ {
				kvsStart = reflect.Append(kvsStart, kvsEnd.Index(i))
			}

			c.kvs.Elem().Set(kvsStart)
		}

		if len(c.keyvalues) == 0 {
			c.AddEmpty()
		}
	}
}

func (c *KeyValues) Clear() {
	defer config.LogStart("KeyValues.Clear", nil)()

	for _, kv := range c.keyvalues {
		c.Remove(kv)
	}
	c.keyvalues = nil
	c.AddEmpty()
}

func (c KeyValues) Collect() (map[string][]string, error) {
	defer config.LogStart("KeyValues.Collect", nil)()

	keyvalues := map[string][]string{}

	for _, kv := range c.keyvalues {
		if !kv.enabled.GetActive() {
			continue
		}

		key, value, err := kv.Get()
		if err != nil {
			return nil, err
		}

		if key == "" {
			continue
		}

		keyvalues[key] = append(keyvalues[key], value)
	}

	return keyvalues, nil
}

type KeyValue struct {
	*gtk.Box
	enabled *gtk.CheckButton
	key     *gtk.Entry
	value   *gtk.Entry
	remove  *gtk.Button
	kv      kv
}

func (c KeyValue) Init(v kv) (*KeyValue, error) {
	defer config.LogStart("KeyValue.Init", nil)()
	c.kv = v

	var err error
	c.key, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	c.key.Show()
	c.key.SetPlaceholderText("Name")

	c.value, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	c.value.Show()
	c.value.SetPlaceholderText("Value")

	c.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}

	c.remove, err = BWButtonNewFromIconName("close", "orange", ICON_SIZE_TAB)
	if err != nil {
		return nil, err
	}
	c.remove.Show()

	c.enabled, err = gtk.CheckButtonNew()
	if err != nil {
		return nil, err
	}
	c.enabled.Show()
	c.enabled.SetActive(true)

	c.enabled.Connect("toggled", c.onChange)
	c.key.Connect("key-release-event", c.onChange)
	c.value.Connect("key-release-event", c.onChange)

	c.Box.PackStart(c.enabled, false, false, 5)
	c.Box.PackStart(c.key, true, true, 0)
	c.Box.PackStart(c.value, true, true, 0)
	c.Box.PackEnd(c.remove, false, false, 5)

	return &c, nil
}

func (c *KeyValue) Get() (string, string, error) {
	defer config.LogStart("KeyValue.Get", nil)()

	key, err := c.key.GetText()
	if err != nil {
		return "", "", err
	}
	value, err := c.value.GetText()
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func (c *KeyValue) onChange() {
	val, _ := c.value.GetText()
	c.kv.SetVal(val)

	key, _ := c.key.GetText()
	c.kv.SetName(key)

	c.kv.SetEnabled(c.enabled.GetActive())
}
