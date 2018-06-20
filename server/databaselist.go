package server

import (
	"github.com/rivo/tview"
)

// DatabaseList ...
type DatabaseList struct {
	*tview.DropDown

	OnSelectDatabase func(string)
}

// NewDatabaseList ...
func NewDatabaseList() *DatabaseList {
	list := &DatabaseList{}

	list.DropDown = tview.NewDropDown()
	list.DropDown.
		SetBorder(true).
		SetBorderPadding(0, 0, 0, 0).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Databases")

	return list
}

// SetDatabases ...
func (d *DatabaseList) SetDatabases(dbs []string) {
	d.DropDown.
		SetOptions(dbs, d.onSelectDatabase).
		SetCurrentOption(-1)
}

func (d *DatabaseList) onSelectDatabase(db string, _ int) {
	if d.OnSelectDatabase != nil {
		d.OnSelectDatabase(db)
	}
}

// SetRect ...
func (d *DatabaseList) SetRect(x, y, width, height int) {
	d.DropDown.SetRect(x, y, width, 3)
}
