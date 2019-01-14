package connect

import (
	"strconv"

	"bitbucket.org/goreorto/sqlhero/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ConnectionList ...
type ConnectionList struct {
	*tview.List
	connections []*config.Connection

	OnSelectConnection func(*config.Connection)
	OnNewConnection    func()
	OnDeleteConnection func(*config.Connection)
}

// NewConnectionList ...
func NewConnectionList() *ConnectionList {
	list := &ConnectionList{}
	list.List = tview.NewList()
	list.List.ShowSecondaryText(false)

	list.List.SetTitle("Favorities")
	list.List.SetTitleAlign(tview.AlignLeft)
	list.List.SetBorder(true)

	list.List.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		if i == 0 {
			list.onNewConnection()
			return
		}
		list.onSelectConnection(list.connections[i-1])
	})

	list.List.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Rune() == 'd' {
			list.onDeleteConnection()
			return nil
		}

		return e
	})

	return list
}

// SetConnections ...
func (c *ConnectionList) SetConnections(cs []*config.Connection) {
	c.List.Clear()
	c.connections = cs

	c.List.AddItem("New connection", "", rune('n'), nil)
	for i := 0; i < len(c.connections); i++ {
		c.List.AddItem(c.connections[i].Name, "", rune(strconv.Itoa(i)[0]), nil)
	}
}

// AddConnection ...
func (c *ConnectionList) AddConnection(conn *config.Connection) {
	c.connections = append(c.connections, conn)
	c.List.AddItem(conn.Name, "", rune(strconv.Itoa(len(c.connections) - 1)[0]), nil)
}

func (c *ConnectionList) onSelectConnection(conn *config.Connection) {
	if c.OnSelectConnection != nil {
		c.OnSelectConnection(conn)
	}
}

func (c *ConnectionList) onNewConnection() {
	if c.OnNewConnection != nil {
		c.OnNewConnection()
	}
}

func (c *ConnectionList) onDeleteConnection() {
	i := c.List.GetCurrentItem()
	if i == 0 || c.OnDeleteConnection == nil {
		return
	}

	c.OnDeleteConnection(c.connections[i-1])
}

/*
func (c *ConnectionList) Focus(delegate func(tview.Primitive)) {
	panic("--")
	delegate(c.List)
}

func (c *ConnectionList) Blur() {
	c.List.Blur()
}
func (c *ConnectionList) Draw(screen tcell.Screen) {
	c.List.Draw(screen)
}
func (c *ConnectionList) GetFocusable() tview.Focusable {
	return c.List
}
func (c *ConnectionList) GetRect() (int, int, int, int) {
	return c.List.GetRect()
}

func (c *ConnectionList) SetRect(x, y, w, h int) {
	c.List.SetRect(x, y, w, h)
}
func (c *ConnectionList) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return c.List.InputHandler()
}
*/
