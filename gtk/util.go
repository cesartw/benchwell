package gtk

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"bitbucket.org/goreorto/sqlaid/clipboard"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type Stringer string

func (s Stringer) String() string {
	return string(s)
}

func StringSliceToStringers(strings []string) []fmt.Stringer {
	result := make([]fmt.Stringer, len(strings))

	for i, s := range strings {
		result[i] = Stringer(s)
	}

	return result
}

func stringSliceToStringerSlice(sc []string) (r []fmt.Stringer) {
	for _, str := range sc {
		r = append(r, Stringer(str))
	}

	return r
}

func colDefSliceToStringerSlice(cols []driver.ColDef) (r []fmt.Stringer) {
	for _, col := range cols {
		r = append(r, col)
	}

	return r
}

func stringerSliceToColDefSlice(cols []fmt.Stringer) (r []driver.ColDef) {
	for _, col := range cols {
		r = append(r, col.(driver.ColDef))
	}

	return r
}

func ClipboardCopy(d string) {
	cb, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	if err != nil {
		clipboard.Copy(d)
	}

	cb.SetText(d)
	cb.Store()
}
