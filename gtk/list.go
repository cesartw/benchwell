package gtk

import (
	"fmt"
	"regexp"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ListOptions struct {
	Names              []fmt.Stringer
	SelectOnRightClick bool
	FilterRegex        *regexp.Regexp
	IconFunc           func(fmt.Stringer) *gdk.Pixbuf
	StockIcon          string
}

type List struct {
	*gtk.ListBox

	options ListOptions

	activeItem        MVar
	activeItemIndex   MVar
	selectedItem      MVar
	selectedItemIndex MVar
}

func (list List) Init(opts ListOptions) (*List, error) {
	var err error
	list.options = opts

	list.ListBox, err = gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}

	list.SetProperty("activate-on-single-click", false)
	list.Connect("row-activated", list.onRowActivated)
	list.Connect("row-selected", list.onRowSelected)
	if list.options.SelectOnRightClick {
		list.Connect("button-press-event", list.onRightClick)
	}

	list.UpdateItems(list.options.Names)

	list.SetFilterFunc(func(row *gtk.ListBoxRow, userData ...interface{}) bool {
		if list.options.FilterRegex == nil {
			return true
		}
		if row.GetIndex() >= len(list.options.Names) {
			config.Env.Log.Debug("mmh, list is larger than the model(fake model)")
			return true
		}

		name := list.options.Names[row.GetIndex()]
		return list.options.FilterRegex.Match([]byte(name.String()))
	})

	return &list, nil
}

func (u *List) Clear() {
	u.UnselectAll()
	u.activeItem.Set(nil)
	u.activeItemIndex.Set(nil)
	u.selectedItem.Set(nil)
	u.selectedItemIndex.Set(nil)
}

func (u *List) onRightClick(_ *gtk.ListBox, e *gdk.Event) {
	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	row := u.GetRowAtY(int(keyEvent.Y()))
	u.GrabFocus()
	u.SelectRow(row)
}

func (u *List) OnButtonPress(f interface{}) {
	u.ListBox.Connect("button-press-event", f)
}

func (u *List) UpdateItems(names []fmt.Stringer) error {
	u.ListBox.GetChildren().Foreach(func(item interface{}) {
		u.ListBox.Remove(item.(gtk.IWidget))
	})
	u.options.Names = names

	for _, name := range names {
		_, err := u.addItem(name, false)
		if err != nil {
			return err
		}
	}

	u.ShowAll()

	return nil
}

func (u *List) AddItem(name fmt.Stringer) (*gtk.ListBoxRow, error) {
	return u.addItem(name, true)
}

func (u *List) addItem(name fmt.Stringer, appendToStore bool) (*gtk.ListBoxRow, error) {
	label, err := gtk.LabelNew(name.String())
	if err != nil {
		return nil, err
	}
	label.SetHAlign(gtk.ALIGN_START)

	row, err := gtk.ListBoxRowNew()
	if err != nil {
		return nil, err
	}

	var widget gtk.IWidget = label

	switch {
	case u.options.IconFunc != nil:
		image, err := gtk.ImageNewFromPixbuf(u.options.IconFunc(name))
		if err != nil {
			return nil, err
		}

		box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
		if err != nil {
			return nil, err
		}

		box.Add(image)
		box.Add(label)
		widget = box
	case u.options.StockIcon != "":
		image, err := gtk.ImageNewFromIconName(u.options.StockIcon, gtk.ICON_SIZE_MENU)
		if err != nil {
			return nil, err
		}

		box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
		if err != nil {
			return nil, err
		}

		box.Add(image)
		box.Add(label)
		widget = box
	}

	row.Add(widget)
	row.ShowAll()
	u.Add(row)
	if appendToStore {
		u.options.Names = append(u.options.Names, name)
	}

	return row, nil
}

func (u *List) onRowActivated(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	u.activeItem.Set(u.options.Names[row.GetIndex()])
	u.activeItemIndex.Set(row.GetIndex())
}

func (u *List) onRowSelected(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	if row != nil && row.GetIndex() >= 0 {
		u.selectedItem.Set(u.options.Names[row.GetIndex()])
		u.selectedItemIndex.Set(row.GetIndex())
	}
}

func (u *List) ActiveItem() (string, bool) {
	v := u.activeItem.Get()
	if v == nil {
		return "", false
	}
	table, ok := u.activeItem.Get().(string)

	return table, ok
}

func (u *List) ActiveItemIndex() (int, bool) {
	v := u.activeItemIndex.Get()
	if v == nil {
		return 0, false
	}
	i, ok := u.activeItemIndex.Get().(int)

	return i, ok
}
func (u *List) SelectedItem() (fmt.Stringer, bool) {
	v := u.selectedItem.Get()
	if v == nil {
		return nil, false
	}
	table, ok := u.selectedItem.Get().(fmt.Stringer)
	return table, ok
}

func (u *List) SelectedItemIndex() (int, bool) {
	v := u.selectedItemIndex.Get()
	if v == nil {
		return 0, false
	}
	i, ok := u.selectedItemIndex.Get().(int)

	return i, ok
}

func (u *List) SetFilterRegex(rg *regexp.Regexp) {
	u.options.FilterRegex = rg
	u.InvalidateFilter()
}
