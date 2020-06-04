package gtk

import (
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
	"github.com/gotk3/gotk3/gtk"
)

type Conditions struct {
	*gtk.Frame
	grid       *gtk.Grid
	btnAdd     *gtk.Button
	conditions []*Condition
	cols       []driver.ColDef
}

type Condition struct {
	fieldCb    *gtk.ComboBoxText
	opCb       *gtk.ComboBoxText
	valueEntry *gtk.Entry
	btnRm      *gtk.Button
}

func NewConditions(cols []driver.ColDef) (c *Conditions, err error) {
	c = &Conditions{cols: cols}
	c.Frame, err = gtk.FrameNew("Filter:")
	if err != nil {
		return nil, err
	}
	c.Frame.SetProperty("shadow-type", gtk.SHADOW_ETCHED_IN)

	c.grid, err = gtk.GridNew()
	if err != nil {
		return nil, err
	}
	c.grid.SetRowSpacing(5)
	c.grid.SetColumnSpacing(5)

	c.btnAdd, err = gtk.ButtonNewFromIconName("gtk-add", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	c.btnAdd.Connect("clicked", c.Add)

	c.grid.Attach(c.btnAdd, 6, 0, 1, 1)

	c.grid.Show()

	c.Frame.Add(c.grid)

	return c, c.Add()
}

func (c *Conditions) Add() error {
	cond, err := newCondition(c.cols)
	if err != nil {
		return err
	}
	c.grid.Remove(c.btnAdd)

	y := len(c.conditions)
	c.grid.Attach(cond.fieldCb, 0, y, 2, 1)
	c.grid.Attach(cond.opCb, 2, y, 1, 1)
	c.grid.Attach(cond.valueEntry, 4, y, 2, 1)
	c.grid.Attach(c.btnAdd, 6, y, 1, 1)
	c.conditions = append(c.conditions, cond)

	if y >= 1 {
		c.grid.Attach(c.conditions[y-1].btnRm, 6, y-1, 1, 1)
	}

	cond.btnRm.Connect("clicked", func() {
		for i, con := range c.conditions {
			if con == cond {
				c.grid.RemoveRow(i)
				c.conditions = append(c.conditions[:i], c.conditions[i+1:]...)
				break
			}
		}

	})

	c.btnAdd.Show()

	if len(c.conditions) >= 2 {
		c.conditions[len(c.conditions)-2].btnRm.Show()
	}

	return nil
}

func (c *Conditions) Statements() ([]driver.CondStmt, error) {
	stmts := []driver.CondStmt{}
	for _, cond := range c.conditions {
		var field driver.ColDef
		textField := cond.fieldCb.GetActiveText()
		for _, col := range c.cols {
			if col.Name == textField {
				field = col
				break
			}
		}

		op := driver.Operator(cond.opCb.GetActiveText())
		value, err := cond.valueEntry.GetText()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, driver.CondStmt{Op: op, Value: value, Field: field})
	}

	return stmts, nil
}

func (c *Conditions) Update(cols []driver.ColDef) error {
	c.cols = cols

	for _, cond := range c.conditions {
		c.grid.Remove(cond.fieldCb)
		c.grid.Remove(cond.opCb)
		c.grid.Remove(cond.valueEntry)
	}

	c.conditions = nil

	return c.Add()
}

func newCondition(cols []driver.ColDef) (c *Condition, err error) {
	c = &Condition{}

	c.fieldCb, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	c.fieldCb.Append("", "")
	for _, col := range cols {
		c.fieldCb.Append(col.Name, col.Name)
	}

	c.opCb, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}

	for _, op := range driver.Operators {
		c.opCb.Append(string(op), string(op))
	}

	c.valueEntry, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	c.btnRm, err = gtk.ButtonNewFromIconName("gtk-remove", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	c.fieldCb.Show()
	c.opCb.Show()
	c.valueEntry.Show()
	c.btnRm.Show()

	return c, nil
}
