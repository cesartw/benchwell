package controls

import (
	"github.com/gotk3/gotk3/gtk"
)

type List struct {
	activeItem        MVar
	activeItemIndex   MVar
	selectedItem      MVar
	selectedItemIndex MVar
	names             []string
	*gtk.ListBox
}

func NewList(names []string) (*List, error) {
	list := &List{
		names: names,
	}

	var err error
	list.ListBox, err = gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}

	list.SetProperty("activate-on-single-click", false)
	list.Connect("row-activated", list.onRowActivated)
	list.Connect("row-selected", list.onRowSelected)

	list.UpdateItems(names)

	return list, nil
}

func (u *List) UpdateItems(names []string) error {
	u.ListBox.GetChildren().Foreach(func(item interface{}) {
		u.ListBox.Remove(item.(gtk.IWidget))
	})
	u.names = names

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
	u.activeItem.Set(u.names[row.GetIndex()])
	u.activeItemIndex.Set(row.GetIndex())
}

func (u *List) onRowSelected(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	if row != nil {
		u.selectedItem.Set(u.names[row.GetIndex()])
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
