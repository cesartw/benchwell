package server

import (
	"github.com/rivo/tview"
)

// TableList ...
type TableList struct {
	*tview.List

	OnSelectTable func(string)
}

// NewTableList ...
func NewTableList() *TableList {
	list := &TableList{}
	list.List = tview.NewList()
	list.ShowSecondaryText(false)

	list.SetTitle("Tables")
	list.SetTitleAlign(tview.AlignLeft)
	list.SetBorder(true)

	return list
}

// Selected
//func (t *TableList) SelectedTable(table *tview.Table) func() {
//return func() {
//tl.GetItemText(tl.GetCurrentItem())
//lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
//cols, rows := 10, 40
//word := 0
//for r := 0; r < rows; r++ {
//for c := 0; c < cols; c++ {
//color := tcell.ColorWhite
//if c < 1 || r < 1 {
//color = tcell.ColorYellow
//}
//table.SetCell(r, c,
//tview.NewTableCell(lorem[word]).
//SetTextColor(color).
//SetAlign(tview.AlignCenter))
//word = (word + 1) % len(lorem)
//}
//}
//}
//}

// SetTables ...
func (t *TableList) SetTables(tables []string) {
	t.Clear()
	for _, table := range tables {
		t.AddItem(table, "", 0, func() {
			tableName, _ := t.GetItemText(t.GetCurrentItem())
			t.onSelectTable(tableName)
		})
	}
	t.SetCurrentItem(-1)
}

func (t *TableList) onSelectTable(table string) {
	if t.OnSelectTable != nil {
		t.OnSelectTable(table)
	}
}
