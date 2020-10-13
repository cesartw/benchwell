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
	public Image (string asset, Gtk.IconSize size = Gtk.IconSize.BUTTON) {
		Object(
			icon_name: @"bw-$asset",
			icon_size: Gtk.IconSize.BUTTON
		);
	}
}
