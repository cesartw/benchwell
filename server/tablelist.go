package server

import (
	"regexp"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type matchString func(string) bool

func (f matchString) MatchString(txt string) bool {
	return f(txt)
}

type searchField struct {
	*tview.InputField
	list *list
}

type list struct {
	*tview.List
	searchField *searchField

	source []string
}

// TableList ...
type TableList struct {
	*tview.Flex
	list        *list
	searchField *searchField

	OnSelectTable func(string)
}

// NewTableList ...
func NewTableList() *TableList {
	list := &TableList{
		Flex:        tview.NewFlex(),
		list:        &list{List: tview.NewList()},
		searchField: &searchField{InputField: tview.NewInputField()},
	}

	list.searchField.list = list.list
	list.list.searchField = list.searchField

	list.searchField.SetBorder(false)

	list.list.ShowSecondaryText(false)
	list.list.SetTitleAlign(tview.AlignLeft)
	list.list.SetBorder(false)

	list.Flex.SetBorder(true)
	list.Flex.SetTitle("Tables")
	list.Flex.SetDirection(tview.FlexRow)
	list.Flex.AddItem(list.searchField, 1, 1, false)
	list.Flex.AddItem(list.list, 0, 1, false)
	list.Flex.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	list.searchField.SetChangedFunc(func(txt string) {
		list.Filter(txt)
	})

	return list
}

// SetTables ...
func (t *TableList) SetTables(tables []string) {
	t.list.source = tables

	t.list.Clear()
	for _, table := range tables {
		t.list.AddItem(table, "", 0, func() {
			tableName, _ := t.list.GetItemText(t.list.GetCurrentItem())
			t.onSelectTable(tableName)
		})
	}
	t.list.SetCurrentItem(-1)
}

func (t *TableList) onSelectTable(table string) {
	if t.OnSelectTable != nil {
		t.OnSelectTable(table)
	}
}

// Draw ...
func (t *TableList) Draw(screen tcell.Screen) {
	t.searchField.SetFieldBackgroundColor(tview.Styles.ContrastBackgroundColor)
	t.searchField.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	if t.searchField.HasFocus() {
		t.searchField.SetFieldBackgroundColor(tcell.ColorYellow)
		t.searchField.SetFieldTextColor(tcell.ColorBlack)
	}
	t.Flex.Draw(screen)
}

// Filter ...
func (t *TableList) Filter(txt string) {
	type matcher interface {
		MatchString(string) bool
	}

	var compare matcher = matchString(func(s string) bool { return strings.Contains(s, txt) })

	reg, err := regexp.Compile(txt)
	if err == nil {
		compare = matchString(reg.MatchString)
	}

	t.list.Clear()
	for _, item := range t.list.source {
		if compare.MatchString(item) {
			t.list.AddItem(item, "", 0, func() {
				tableName, _ := t.list.GetItemText(t.list.GetCurrentItem())
				t.onSelectTable(tableName)
			})
		}
	}
}

// Focus ...
func (t *TableList) Focus(delegate func(tview.Primitive)) {
	if t.list.GetCurrentItem() == -1 {
		delegate(t.searchField)
		return
	}
	delegate(t.list)
}

func (s *searchField) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, setFocus func(tview.Primitive)) {
		switch key := e.Key(); key {
		case tcell.KeyDown:
			s.list.SetCurrentItem(0)
			s.Blur()
			setFocus(s.list)
			return
		case tcell.KeyUp:
			s.list.SetCurrentItem(s.list.GetItemCount() - 1)
			s.Blur()
			setFocus(s.list)
			return
		case tcell.KeyEnter, tcell.KeyEsc, tcell.KeyTab:
			s.list.SetCurrentItem(0)
			s.Blur()
			setFocus(s.list)
			return
		}

		s.InputField.InputHandler()(e, setFocus)
	}
}

func (l *list) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(e *tcell.EventKey, setFocus func(tview.Primitive)) {
		switch key := e.Key(); key {
		case tcell.KeyDown:
			if l.GetCurrentItem() == l.GetItemCount()-1 {
				l.Blur()
				setFocus(l.searchField)
				return
			}
		case tcell.KeyUp:
			if l.GetCurrentItem() == 0 {
				l.Blur()
				setFocus(l.searchField)
				return
			}
		case tcell.KeyCtrlS:
			l.Blur()
			setFocus(l.searchField)
			return
		}

		l.List.InputHandler()(e, setFocus)
	}
}
