package server

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
)

// RecordTable ...
type RecordTable struct {
	*tview.Box
	Table        *tview.Table
	editingField *tview.InputField

	tableName   string
	tableDef    []driver.ColDef
	sourceData  [][]*string
	workingData [][]*string

	OnSaveRecord   func(string, []driver.ColDef, []*string, []*string)
	OnInsertRecord func(string, []driver.ColDef, []*string) []*string
	OnDeleteRecord func(string, []driver.ColDef, []*string, []*string) bool
	OnReload       func(string)
}

// NewRecordTable ...
func NewRecordTable() *RecordTable {
	t := &RecordTable{
		Box:          tview.NewBox().SetBorder(false),
		Table:        tview.NewTable().SetBorders(true),
		editingField: tview.NewInputField(),
	}

	t.editingField.SetBorder(true)
	t.Table.SetSelectable(true, true)

	return t
}

// SetData ...
func (t *RecordTable) SetData(tableName string, colDef []driver.ColDef, rows [][]*string) {
	t.workingData = rows
	t.sourceData = make([][]*string, len(rows))
	copy(t.sourceData, t.workingData)

	t.tableDef = colDef
	t.tableName = tableName

	t.Table.Clear()

	// add headers
	for i, def := range colDef {
		name := def.Name

		colorBG := tcell.ColorLime
		if def.PK {
			colorBG = tcell.ColorYellow
		}
		t.Table.SetCell(0, i,
			tview.NewTableCell(name).
				SetTextColor(tcell.ColorBlack).
				SetBackgroundColor(colorBG).
				SetAlign(tview.AlignLeft))
	}

	// add data
	for row := 0; row < len(rows); row++ {
		for col := 0; col < len(colDef); col++ {
			color := tcell.ColorWhite
			if colDef[col].PK {
				color = tcell.ColorYellow
			}

			var value string
			if rows[row][col] != nil {
				value = *rows[row][col]
			}

			t.Table.SetCell(row+1, col,
				tview.NewTableCell(value).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}

	t.Table.SetFixed(1, 0)
}

// Draw ...
func (t *RecordTable) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)
	t.Table.Draw(screen)

	if t.editingField.HasFocus() {
		t.editingField.Draw(screen)
	}
}

// HasFocus ...
func (t *RecordTable) HasFocus() bool {
	return t.Box.HasFocus() || t.Table.HasFocus() || t.editingField.HasFocus()
}

// GetFocusable ...
func (t *RecordTable) GetFocusable() tview.Focusable {
	return t
}

// SetRect ...
func (t *RecordTable) SetRect(x, y, width, height int) {
	t.Box.SetRect(x, y, width, height)
	t.Table.SetRect(x+1, y, width-2, height)
	t.editingField.SetRect(x+width/2-36/2, y+height/2-1, 36, 3)
}

// InputHandler ...
func (t *RecordTable) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, setFocus func(tview.Primitive)) {
		key := e.Key()

		switch key {
		case tcell.KeyRune:
			switch e.Rune() {
			case 'i':
				row := t.Table.GetRowCount()
				for col := range t.tableDef {
					color := tcell.ColorWhite
					if t.tableDef[col].PK {
						color = tcell.ColorYellow
					}

					t.Table.SetCell(row, col,
						tview.NewTableCell("").
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
				}

				t.Table.Select(row, 0)
				t.workingData = append(t.workingData, make([]*string, len(t.tableDef)))
			case 'e':
				row, col := t.Table.GetSelection()
				if row == 0 {
					return
				}

				value := t.workingData[row-1][col]
				if value == nil {
					t.workingData[row-1][col] = new(string)
					value = t.workingData[row-1][col]
				}

				t.editingField.SetText(*value)
				t.editingField.SetFinishedFunc(func(e tcell.Key) {
					if e == tcell.KeyEnter {
						*(t.workingData[row-1][col]) = t.editingField.GetText()
						t.Table.GetCell(row, col).SetText(t.editingField.GetText())
					}

					t.editingField.Blur()
					setFocus(t)
				})

				setFocus(t.editingField)

				return
			case 's':
				row, _ := t.Table.GetSelection()
				if row == 0 {
					return
				}

				if len(t.sourceData) < row-2 {
					if t.OnSaveRecord != nil {
						t.OnSaveRecord(t.tableName, t.tableDef, t.workingData[row-1], t.sourceData[row-1])
						t.sourceData[row-1] = t.workingData[row-1]
					}
				} else {
					if t.OnInsertRecord != nil {
						data := t.OnInsertRecord(t.tableName, t.tableDef, t.workingData[row-1])
						if data == nil {
							return
						}
						t.workingData[row-1] = data
						t.sourceData = append(t.sourceData, t.workingData[row-1])

						for col, d := range data {
							if d == nil {
								continue
							}
							t.Table.GetCell(row, col).SetText(*d)
						}

					}
				}

			case 'd':
				row, _ := t.Table.GetSelection()
				if row == 0 {
					return
				}

				if t.OnDeleteRecord != nil {
					if t.OnDeleteRecord(t.tableName, t.tableDef, t.workingData[row-1], t.sourceData[row-1]) {
						t.workingData = append(t.workingData[:row], t.workingData[row+1:]...)
						t.sourceData = append(t.sourceData[:row], t.sourceData[row+1:]...)
						t.Table.RemoveRow(row)
					}
				}
			case 'r':
				if t.OnReload != nil {
					t.OnReload(t.tableName)
				}
			default:
				t.Table.InputHandler()(e, setFocus)
			}
		default:
			t.Table.InputHandler()(e, setFocus)
		}
	}
}
