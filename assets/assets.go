package assets

import "github.com/gotk3/gotk3/gdk"

var (
	Table       *gdk.Pixbuf
	TableCustom *gdk.Pixbuf
)

func Load() (err error) {
	Table, err = gdk.PixbufNewFromFile("assets/table.png")
	if err != nil {
		return err
	}

	Table, err = Table.ScaleSimple(15, 15, gdk.INTERP_NEAREST)
	if err != nil {
		return err
	}

	TableCustom, err = gdk.PixbufNewFromFile("assets/table-custom.png")
	if err != nil {
		return err
	}

	TableCustom, err = TableCustom.ScaleSimple(15, 15, gdk.INTERP_NEAREST)
	if err != nil {
		return err
	}

	return nil
}
