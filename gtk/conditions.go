package gtk

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

// TODO: look into gtk.EntryCompletion for combobox
type Conditions struct {
	w *Window
	*gtk.Frame
	grid       *gtk.Grid
	btnAdd     *gtk.Button
	conditions []*Condition
	cols       []driver.ColDef

	ctrl conditionsCtrl
}

type conditionsCtrl interface {
	OnApplyConditions()
	Config() *config.Config
}
type Condition struct {
	cols       []driver.ColDef
	activeCb   *gtk.CheckButton
	store      *gtk.ListStore
	fieldCb    *gtk.ComboBox
	opCb       *gtk.ComboBoxText
	valueEntry *gtk.Entry
	btnRm      *gtk.Button
	ctrl       conditionsCtrl
}

func (c Conditions) Init(w *Window, ctrl conditionsCtrl) (*Conditions, error) {
	var err error
	c.w = w
	c.ctrl = ctrl

	c.Frame, err = gtk.FrameNew("Filter:")
	if err != nil {
		return nil, err
	}
	c.Frame.SetProperty("shadow-type", gtk.SHADOW_NONE)
	c.Frame.SetName("conditions")

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

	return &c, nil //c.Add()
}

func (c *Conditions) Add() error {
	cond, err := Condition{}.Init(c.w, c.cols, c.ctrl)
	if err != nil {
		return err
	}
	cond.valueEntry.Connect("activate", c.ctrl.OnApplyConditions)

	c.grid.Remove(c.btnAdd)

	y := len(c.conditions)
	c.grid.Attach(cond.activeCb, 0, y, 2, 1)
	c.grid.Attach(cond.fieldCb, 2, y, 2, 1)
	c.grid.Attach(cond.fieldCb, 3, y, 2, 1)
	c.grid.Attach(cond.opCb, 4, y, 1, 1)
	c.grid.Attach(cond.valueEntry, 5, y, 2, 1)
	c.grid.Attach(c.btnAdd, 8, y, 1, 1)
	c.conditions = append(c.conditions, cond)

	if y >= 1 {
		c.grid.Attach(c.conditions[y-1].btnRm, 8, y-1, 1, 1)
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
		if !cond.activeCb.GetActive() {
			continue
		}
		var field driver.ColDef
		// ffs
		textField, err := cond.Field()
		if err != nil {
			return nil, err
		}

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
	conditions := []*Condition{}

	// remove empty
	if len(c.conditions) > 1 {
		for i, cond := range c.conditions {
			field, err := cond.Field()
			if err != nil {
				return err
			}

			if field == "" {
				c.grid.RemoveRow(i)
				continue
			}
			conditions = append(conditions, cond)
		}
	} else {
		conditions = c.conditions
	}

	//update columns in remaining conditions
	for _, cond := range c.conditions {
		cond.cols = cols
		field, err := cond.Field()
		if err != nil {
			return err
		}

		foundAt := -1
		for i, col := range cols {
			if col.Name != field {
				continue
			}
			foundAt = i
			break
		}

		// field is not part of the new table disable widget and move on
		if foundAt == -1 && field != "" {
			cond.fieldCb.SetSensitive(false)
			cond.opCb.SetSensitive(false)
			cond.valueEntry.SetSensitive(false)
			cond.activeCb.SetSensitive(false)
			cond.activeCb.SetActive(false)
			break
		}

		cond.fieldCb.SetSensitive(true)
		cond.opCb.SetSensitive(true)
		cond.valueEntry.SetSensitive(true)
		cond.activeCb.SetSensitive(true)
		cond.activeCb.SetActive(true)

		cond.store.Clear()
		cond.store.SetValue(cond.store.Append(), 0, "")
		for _, col := range c.cols {
			cond.store.SetValue(cond.store.Append(), 0, col.Name)
		}
		cond.fieldCb.SetActiveID(field)
	}

	if len(c.conditions) == 0 {
		return c.Add()
	}

	c.conditions = conditions

	return nil
}

func (c Condition) Init(_ *Window, cols []driver.ColDef, ctrl conditionsCtrl) (*Condition, error) {
	c.cols = cols
	c.ctrl = ctrl
	var err error

	c.store, _ = gtk.ListStoreNew(glib.TYPE_STRING)
	c.store.SetValue(c.store.Append(), 0, "")
	for _, col := range cols {
		c.store.SetValue(c.store.Append(), 0, col.Name)
	}

	c.activeCb, err = gtk.CheckButtonNew()
	if err != nil {
		return nil, err
	}
	c.activeCb.SetActive(true)

	c.fieldCb, err = gtk.ComboBoxNewWithModelAndEntry(c.store.ToTreeModel())
	if err != nil {
		return nil, err
	}
	c.fieldCb.SetEntryTextColumn(0)
	c.fieldCb.SetProperty("id-column", 0)
	completion, err := gtk.EntryCompletionNew()
	if err != nil {
		return nil, err
	}
	completion.SetProperty("text-column", 0)
	completion.SetProperty("inline-completion", true)
	completion.SetProperty("inline-selection", true)
	completion.SetMinimumKeyLength(2)
	completion.SetModel(c.store)
	entry, err := c.fieldCb.GetEntry()
	if err != nil {
		return nil, err
	}
	entry.SetCompletion(completion)
	entry.Connect("focus-out-event", c.onFocusOut)

	c.opCb, err = gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	for _, op := range driver.Operators {
		c.opCb.Append(string(op), string(op))
	}
	c.opCb.SetActive(0)
	c.opCb.Connect("changed", func() {
		enable := c.opCb.GetActiveText() != string(driver.IsNull) && c.opCb.GetActiveText() != string(driver.IsNotNull)
		c.valueEntry.SetSensitive(enable)
	})

	c.valueEntry, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	c.btnRm, err = gtk.ButtonNewFromIconName("gtk-remove", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	c.activeCb.ShowAll()
	c.fieldCb.ShowAll()
	c.opCb.ShowAll()
	c.valueEntry.ShowAll()
	c.btnRm.ShowAll()

	return &c, nil
}

func (c *Condition) Field() (string, error) {
	if c.fieldCb.GetActiveID() == "" {
		return "", nil
	}

	iter, err := c.fieldCb.GetActiveIter()
	if err != nil {
		return "", err
	}

	gvalue, err := c.store.GetValue(iter, 0)
	if err != nil {
		return "", err
	}

	textField, err := gvalue.GetString()
	if err != nil {
		return "", err
	}

	return textField, nil
}

func (c *Condition) onFocusOut() {
	entry, err := c.fieldCb.GetEntry()
	if err != nil {
		c.ctrl.Config().Error(err)
		return
	}

	field, err := entry.GetText()
	if err != nil {
		c.ctrl.Config().Error(err)
		return
	}

	selectedText, err := c.Field()
	if err != nil {
		c.ctrl.Config().Error(err)
		return

	}

	if field == selectedText {
		return
	}

	for i, col := range c.cols {
		if col.Name != field {
			continue
		}

		c.fieldCb.SetActive(i + 1) // +1 because of the blank row
	}
}
