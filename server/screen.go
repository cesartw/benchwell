package server

import (
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Screen ...
type Screen struct {
	*tview.Flex
	sidePanel *SidePanel
	tabber    *Tabber

	OnSelectDatabase func(db string)
	OnSelectTable    func(tableName string)
	OnSaveRecord     func(string, []driver.ColDef, []*string, []*string) bool
	OnInsertRecord   func(string, []driver.ColDef, []*string) []*string
	OnDeleteRecord   func(string, []driver.ColDef, []*string, []*string) bool
	OnReload         func(string)

	ctx sqlengine.Context
}

// New ...
func New(app *tview.Application) *Screen {
	s := &Screen{}

	sidePanelWidth := 30
	s.sidePanel = NewSidePanel(sidePanelWidth)
	s.tabber = NewTabber()

	s.Flex = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(s.sidePanel, sidePanelWidth, 1, false).
		AddItem(s.tabber, 0, 1, false)

	s.sidePanel.OnSelectDatabase = s.onSelectDatabase
	s.sidePanel.OnSelectTable = s.onSelectTable

	return s
}

// AddTableTab ...
func (s *Screen) AddTableTab(tableName string) {
	recordTable := NewRecordTable()
	recordTable.OnDeleteRecord = s.onDeleteRecord
	recordTable.OnReload = s.onReload
	recordTable.OnSaveRecord = s.onSaveRecord
	recordTable.OnInsertRecord = s.onInsertRecord

	s.tabber.AddTab(tableName, recordTable, true)
}

// RecordTable ...
func (s *Screen) RecordTable() tview.Primitive {
	return s.tabber
}

// DatabaseList ...
func (s *Screen) DatabaseList() tview.Primitive {
	return s.sidePanel.databaseList
}

// TableList ...
func (s *Screen) TableList() tview.Primitive {
	return s.sidePanel.tableList
}

// SetDatabases ...
func (s *Screen) SetDatabases(dbs []string) {
	s.sidePanel.SetDatabases(dbs)
}

// SetTables ...
func (s *Screen) SetTables(dbs []string) {
	s.sidePanel.SetTables(dbs)
}

// SetData ...
func (s *Screen) SetData(tableName string, def []driver.ColDef, rows [][]*string) {
	var table *RecordTable

	s.AddTableTab(tableName)
	table = s.tabber.focusedTab.content.(*RecordTable)
	table.SetData(tableName, def, rows)
}

func (s *Screen) onSelectDatabase(db string) {
	if s.OnSelectDatabase != nil {
		s.OnSelectDatabase(db)
	}
}

func (s *Screen) onSelectTable(tableName string) {
	if s.OnSelectTable != nil {
		s.OnSelectTable(tableName)
	}
}

func (s *Screen) onDeleteRecord(tableName string, def []driver.ColDef, row, oldRow []*string) bool {
	if s.OnDeleteRecord != nil {
		return s.OnDeleteRecord(tableName, def, row, oldRow)
	}

	return false
}

func (s *Screen) onReload(tableName string) {
	if s.OnReload != nil {
		s.OnReload(tableName)
	}
}

func (s *Screen) onSaveRecord(tableName string, def []driver.ColDef, row, oldRow []*string) {
	if s.OnSaveRecord != nil {
		s.OnSaveRecord(tableName, def, row, oldRow)
	}
}

func (s *Screen) onInsertRecord(tableName string, def []driver.ColDef, row []*string) []*string {
	if s.OnInsertRecord != nil {
		return s.OnInsertRecord(tableName, def, row)
	}
	return nil
}

// Keybinds ...
func (s *Screen) Keybinds() map[tcell.Key]func() tview.Primitive {
	return map[tcell.Key]func() tview.Primitive{
		tcell.KeyCtrlD: func() tview.Primitive { return s.sidePanel.databaseList },
		tcell.KeyCtrlT: func() tview.Primitive { return s.sidePanel.tableList },
		tcell.KeyCtrlL: func() tview.Primitive { return s.tabber },
	}
}

// Focus ...
func (s *Screen) Focus(delegate func(tview.Primitive)) {
	delegate(s.sidePanel)
}

// SetContext ...
func (s *Screen) SetContext(ctx sqlengine.Context) {
	s.ctx = ctx
}

// Context ...
func (s *Screen) Context() sqlengine.Context {
	return s.ctx
}
