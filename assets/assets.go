package assets

import (
	"github.com/gotk3/gotk3/gdk"
)

var (
	Table       *gdk.Pixbuf
	TableCustom *gdk.Pixbuf
)

// icon size matching GTK's
const (
	SizeMenu         = 16
	SizeButton       = 16
	SizeSmallToolbar = 16
	SizeLargeToolbar = 24
)

func Load() (err error) {
	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return err
	}

	Table, err = loader.WriteAndReturnPixbuf(tableBytes)
	if err != nil {
		return err
	}

	Table, err = Table.ScaleSimple(SizeMenu, SizeMenu, gdk.INTERP_NEAREST)
	if err != nil {
		return err
	}

	loader, err = gdk.PixbufLoaderNew()
	if err != nil {
		return err
	}

	TableCustom, err = loader.WriteAndReturnPixbuf(tableCustomBytes)
	if err != nil {
		return err
	}

	TableCustom, err = TableCustom.ScaleSimple(SizeMenu, SizeMenu, gdk.INTERP_NEAREST)
	if err != nil {
		return err
	}

	return nil
}