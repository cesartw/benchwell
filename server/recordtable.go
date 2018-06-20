package server

import (
	"fmt"

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
	OnDeleteRecord func(string, []driver.ColDef, []*string, []*string) bool
	OnReload       func(string)
}

// NewRecordTable ...
func NewRecordTable() *RecordTable {
	t := &RecordTable{
		Box:          tview.NewBox().SetBorder(true).SetTitleAlign(tview.AlignLeft),
		Table:        tview.NewTable().SetBorders(true),
		editingField: tview.NewInputField(),
	}
	t.editingField.SetBorder(true)
	t.Table.SetSelectable(true, true)

	/*
		t.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
			switch key.Rune() {
			case 'e': // edit field
				row, col := t.GetSelection()
				if row == 0 {
					return nil
				}

				t.ctrl.OnEditRecord(t.tableDef, t.workingData[row-1], col, func(values []string) {
					if len(values) != len(t.tableDef) {
						// TODO: report error
						return
					}

					for col := range t.tableDef {
						cell := t.GetCell(row, col)
						cell.SetText(values[col])
					}
				})

				return nil
			case 's': // save record
				row, _ := t.GetSelection()

				values := make([]string, len(t.tableDef))
				for i := 0; i < len(t.tableDef); i++ {
					cell := t.GetCell(row, i)
					values[i] = cell.Text
				}

				if row < len(t.workingData) {
					t.ctrl.OnSaveRecord(t.tableDef, values, t.workingData[row-1])
				} else {
					t.ctrl.OnInsertRecord(t.tableDef, values)
				}

			case 'd': // delete
				row, _ := t.GetSelection()
				t.ctrl.OnDeleteRecord(t.tableDef, t.workingData[row-1])
			case 'r': // reload
				break
			case 'S': // save all records
			case 'a':
				for i := range t.tableDef {
					t.SetCell(len(t.workingData)+1, i,
						tview.NewTableCell("").
							SetTextColor(tcell.ColorWhite).
							SetAlign(tview.AlignLeft))
				}

				t.Select(len(t.workingData)+1, 0)
				return nil
			default:
				return key
			}

			t.ctrl.OnTableSelected(t.tableName)

			return nil
		})
	*/

	return t
}

// SetData ...
func (t *RecordTable) SetData(tableName string, colDef []driver.ColDef, rows [][]*string) {
	t.workingData = rows
	t.sourceData = make([][]*string, len(rows))
	copy(t.sourceData, t.workingData)

	t.tableDef = colDef
	t.tableName = tableName
	t.Box.SetTitle(fmt.Sprintf("Table: %s", tableName))

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

// SetRect ...
func (t *RecordTable) SetRect(x, y, width, height int) {
	t.Box.SetRect(x, y, width, height)
	t.Table.SetRect(x+1, y+1, width-2, height-2)
	t.editingField.SetRect(x+width/2-36/2, y+height/2-1, 36, 3)
}

// InputHandler ...
func (t *RecordTable) InputHandler() func(e *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(e *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := e.Key()

		switch key {
		case tcell.KeyRune:
			switch e.Rune() {
			case 'e':
				row, col := t.Table.GetSelection()
				if row == 0 {
					return
				}

				value := t.workingData[row-1][col]
				t.editingField.SetText(*value)
				t.editingField.SetFinishedFunc(func(e tcell.Key) {
					if e == tcell.KeyEnter {
						*t.workingData[row-1][col] = t.editingField.GetText()
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

				if t.OnSaveRecord != nil {
					_ = t.workingData[row-1]
					_ = t.sourceData[row-1]
					t.OnSaveRecord(t.tableName, t.tableDef, t.workingData[row-1], t.sourceData[row-1])
				}
				t.sourceData[row-1] = t.workingData[row-1]
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
	})
}
