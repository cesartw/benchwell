enum Benchwell.Icons {
	AddRecord,
	AddTab,
	Add,
	Back,
	Close,
	Config,
	Connection,
	Copy,
	Cowboy,
	Db,
	DeleteRecord,
	DeleteTable,
	Directory,
	EditTable,
	Filter,
	Next,
	Open,
	Refresh,
	Regresh,
	SaveRecord,
	Save,
	TableV,
	Table,
	Truncate
}

public class Benchwell.Image : Gtk.Image {
	public Image (string asset, string color, int size) {
		var pixbuf = new Gdk.Pixbuf.from_file (@"/home/goreorto/gopath/src/bitbucket.org/goreorto/benchwell/assets/iconset-temp/$(asset)48.png");
		pixbuf = pixbuf.scale_simple (size, size, Gdk.InterpType.BILINEAR);

		set_from_pixbuf (pixbuf);
	}
}
