package gtk

import (
	"regexp"

	"bitbucket.org/goreorto/sqlhero/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ListOptions struct {
	Names              []string
	SelectOnRightClick bool
	FilterRegex        *regexp.Regexp
}

type List struct {
	activeItem        MVar
	activeItemIndex   MVar
	selectedItem      MVar
	selectedItemIndex MVar
	*gtk.ListBox
	options ListOptions
}

func NewList(opts ListOptions) (*List, error) {
	list := &List{
		options: opts,
	}

	var err error
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
		return list.options.FilterRegex.Match([]byte(name))
	})

	return list, nil
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

func (u *List) UpdateItems(names []string) error {
	u.ListBox.GetChildren().Foreach(func(item interface{}) {
		u.ListBox.Remove(item.(gtk.IWidget))
	})
	u.options.Names = names

	for _, name := range names {
		label, err := gtk.LabelNew(name)
		if err != nil {
			return err
		}
		label.SetHAlign(gtk.ALIGN_START)

		row, err := gtk.ListBoxRowNew()
		if err != nil {
			return err
		}

		row.Add(label)
		u.Add(row)
	}

	u.ShowAll()

	return nil
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
func (u *List) SelectedItem() (string, bool) {
	v := u.selectedItem.Get()
	if v == nil {
		return "", false
	}
	table, ok := u.selectedItem.Get().(string)

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
