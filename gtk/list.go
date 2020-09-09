package gtk

import (
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/goreorto/benchwell/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ListOptions struct {
	Names              []fmt.Stringer
	SelectOnRightClick bool
	FilterRegex        *regexp.Regexp
	IconFunc           func(fmt.Stringer) (string, int)
	StockIcon          string
}

type listrow struct {
	*gtk.ListBoxRow
	imageOff *gtk.Image
	imageOn  *gtk.Image
	label    *gtk.Label
	box      *gtk.Box
	name     fmt.Stringer
}

type List struct {
	*gtk.ListBox

	options *ListOptions
	rows    []*listrow
	ctrlMod bool

	// TODO: these may not be needed
	activeItem        MVar
	activeItemIndex   MVar
	selectedItem      MVar
	selectedItemIndex MVar
	ctrl              listCtrl
}

type listCtrl interface {
}

func (list List) Init(_ *Window, opts *ListOptions, ctrl listCtrl) (*List, error) {
	defer config.LogStart("List.Init", nil)()

	var err error
	list.options = opts
	list.ctrl = ctrl

	list.ListBox, err = gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}

	list.Connect("button-press-event", func(_ *gtk.ListBox, e *gdk.Event) bool {
		keyEvent := gdk.EventButtonNewFromEvent(e)
		list.ctrlMod = keyEvent.State()&gdk.CONTROL_MASK > 0

		return false
	})

	list.Connect("button-release-event", func(_ *gtk.ListBox, e *gdk.Event) bool {
		list.ctrlMod = false

		return false
	})

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
			config.Debug("mmh, list is larger than the model(fake model)")
			return true
		}

		name := list.options.Names[row.GetIndex()]
		return list.options.FilterRegex.Match([]byte(strings.ToLower(name.String())))
	})

	return &list, nil
}

func (u *List) CtrlMod() bool {
	defer config.LogStart("List.CtrlMod", nil)()

	return u.ctrlMod
}

func (u *List) ClearSelection() {
	defer config.LogStart("List.ClearSelection", nil)()

	u.UnselectAll()
	u.activeItem.Set(nil)
	u.activeItemIndex.Set(nil)
	u.selectedItem.Set(nil)
	u.selectedItemIndex.Set(nil)
}

func (u *List) Clear() {
	defer config.LogStart("List.Clear", nil)()

	u.ClearSelection()
	for _, row := range u.rows {
		u.Remove(row)
	}

	u.rows = nil
	u.options.Names = nil
}

func (u *List) onRightClick(_ *gtk.ListBox, e *gdk.Event) {
	defer config.LogStart("List.onRightClick", nil)()

	keyEvent := gdk.EventButtonNewFromEvent(e)

	if keyEvent.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	row := u.GetRowAtY(int(keyEvent.Y()))
	u.GrabFocus()
	u.ListBox.SelectRow(row)
}

func (u *List) OnButtonPress(f interface{}) {
	defer config.LogStart("List.OnButtonPress", nil)()

	u.ListBox.Connect("button-press-event", f)
}

func (u *List) UpdateItems(names []fmt.Stringer) error {
	defer config.LogStart("List.UpdateItems", nil)()

	u.ListBox.GetChildren().Foreach(func(item interface{}) {
		u.ListBox.Remove(item.(gtk.IWidget))
	})
	u.options.Names = names

	for _, name := range names {
		_, err := u.appendItem(name, false)
		if err != nil {
			//return err
		}
	}

	u.Show()

	return nil
}

func (u *List) AppendItem(name fmt.Stringer) (*gtk.ListBoxRow, error) {
	//defer config.LogStart("List.AppendItem", nil)()
	row, err := u.appendItem(name, true)
	if err != nil {
		return nil, err
	}

	return row.ListBoxRow, nil
}

func (u *List) appendItem(name fmt.Stringer, addToStore bool) (*listrow, error) {
	//defer config.LogStart("List.appendItem", nil)()

	row, err := u.buildItem(name, addToStore)
	if err != nil {
		return nil, err
	}

	u.Add(row)
	if addToStore {
		u.options.Names = append(u.options.Names, name)
	}
	u.rows = append(u.rows, row)

	return row, nil
}

func (u *List) PrependItem(name fmt.Stringer) (*gtk.ListBoxRow, error) {
	//defer config.LogStart("List.PrependItem", nil)()
	row, err := u.prependItem(name, true)
	if err != nil {
		return nil, err
	}

	return row.ListBoxRow, nil
}

func (u *List) prependItem(name fmt.Stringer, addToStore bool) (*listrow, error) {
	//defer config.LogStart("List.prependItem", nil)()

	row, err := u.buildItem(name, addToStore)
	if err != nil {
		return nil, err
	}

	u.Prepend(row)
	if addToStore {
		u.options.Names = append([]fmt.Stringer{name}, u.options.Names...)
	}
	u.rows = append(u.rows, row)

	return row, nil
}

func (u *List) buildItem(name fmt.Stringer, appendToStore bool) (*listrow, error) {
	row := &listrow{name: name}

	var err error
	row.label, err = gtk.LabelNew(name.String())
	if err != nil {
		return nil, err
	}
	row.label.SetHAlign(gtk.ALIGN_START)

	row.ListBoxRow, err = gtk.ListBoxRowNew()
	if err != nil {
		return nil, err
	}

	var widget gtk.IWidget = row.label

	switch {
	case u.options.IconFunc != nil:
		fileName, size := u.options.IconFunc(name)
		row.imageOff, err = BWImageNewFromFile(fileName, "orange", size)
		if err != nil {
			return nil, err
		}
		row.imageOn, err = BWImageNewFromFile(fileName, "white", size)
		if err != nil {
			return nil, err
		}

		row.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		if err != nil {
			return nil, err
		}

		row.box.PackStart(row.imageOff, false, false, 5)
		row.box.PackStart(row.imageOn, false, false, 5)
		row.box.PackStart(row.label, false, false, 0)
		widget = row.box
	case u.options.StockIcon != "":
		row.imageOff, err = gtk.ImageNewFromIconName(u.options.StockIcon, gtk.ICON_SIZE_MENU)
		if err != nil {
			return nil, err
		}

		row.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		if err != nil {
			return nil, err
		}

		row.box.PackStart(row.imageOff, false, false, 5)
		row.box.PackStart(row.label, false, false, 5)
		widget = row.box
	}

	row.Add(widget)
	if row.box != nil {
		row.box.Show()
		row.imageOff.Show()
	}
	row.label.Show()
	row.ListBoxRow.Show()

	return row, nil
}

func (u *List) onRowActivated(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	defer config.LogStart("List.onRowActivated", nil)()

	u.activeItem.Set(u.options.Names[row.GetIndex()])
	u.activeItemIndex.Set(row.GetIndex())
}

func (u *List) onRowSelected(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	defer config.LogStart("List.onRowSelected", nil)()

	if row != nil && row.GetIndex() >= 0 {
		u.selectedItem.Set(u.options.Names[row.GetIndex()])
		u.selectedItemIndex.Set(row.GetIndex())

		if u.options.IconFunc != nil {
			lr := u.rows[row.GetIndex()]
			if lr.imageOn == nil {
				return
			}

			for _, row := range u.rows {
				row.imageOff.Show()
				row.imageOn.Hide()
			}

			lr.imageOff.Hide()
			lr.imageOn.Show()
		}
	}
}

func (u *List) ActiveItem() (string, bool) {
	defer config.LogStart("List.ActiveItem", nil)()

	v := u.activeItem.Get()
	if v == nil {
		return "", false
	}
	table, ok := u.activeItem.Get().(string)

	return table, ok
}

func (u *List) ActiveItemIndex() (int, bool) {
	defer config.LogStart("List.ActiveItemIndex", nil)()

	v := u.activeItemIndex.Get()
	if v == nil {
		return 0, false
	}
	i, ok := u.activeItemIndex.Get().(int)

	return i, ok
}

func (u *List) SelectedItem() (fmt.Stringer, bool) {
	defer config.LogStart("List.SelectedItem", nil)()

	v := u.selectedItem.Get()
	if v == nil {
		return nil, false
	}
	table, ok := u.selectedItem.Get().(fmt.Stringer)
	return table, ok
}

func (u *List) SelectedItemIndex() (int, bool) {
	defer config.LogStart("List.SelectedItemIndex", nil)()

	v := u.selectedItemIndex.Get()
	if v == nil {
		return 0, false
	}
	i, ok := u.selectedItemIndex.Get().(int)

	return i, ok
}

func (u *List) SetFilterRegex(rg *regexp.Regexp) {
	defer config.LogStart("List.SetFilterRegex", nil)()

	u.options.FilterRegex = rg
	u.InvalidateFilter()
}
