package server

import "github.com/rivo/tview"

// SidePanel ...
type SidePanel struct {
	*tview.Flex

	app *tview.Application

	databaseList *DatabaseList
	tableList    *TableList

	OnSelectDatabase func(string)
	OnSelectTable    func(string)

	width int
}

// NewSidePanel ...
func NewSidePanel(width int) *SidePanel {
	if width == 0 {
		width = 30
	}

	s := &SidePanel{
		tableList:    NewTableList(),
		databaseList: NewDatabaseList(),
		width:        width,
	}

	s.databaseList.OnSelectDatabase = s.onSelectDatabase
	s.tableList.OnSelectTable = s.onSelectTable

	s.Flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(s.databaseList, 3, 1, false).
		AddItem(s.tableList, 0, 1, false)

	return s
}

// SetDatabases ...
func (s *SidePanel) SetDatabases(dbs []string) {
	s.databaseList.SetDatabases(dbs)
	s.tableList.SetTables([]string{})
}

// SetTables ...
func (s *SidePanel) SetTables(tables []string) {
	s.tableList.SetTables(tables)
}

func (s *SidePanel) onSelectDatabase(db string) {
	if s.OnSelectDatabase != nil {
		s.OnSelectDatabase(db)
	}
}

func (s *SidePanel) onSelectTable(table string) {
	if s.OnSelectTable != nil {
		s.OnSelectTable(table)
	}
}

// SetRect ...
func (s *SidePanel) SetRect(x, y, width, height int) {
	s.Box.SetRect(x, y, width, height)
	s.databaseList.SetRect(x, y, s.width, height)
	s.tableList.SetRect(x, y, s.width, height)
}

// Focus ...
func (s *SidePanel) Focus(delegate func(tview.Primitive)) {
	delegate(s.databaseList)
}
