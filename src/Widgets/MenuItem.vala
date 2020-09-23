public class Benchwell.MenuItem : Gtk.MenuItem {
	public MenuItem(owned string label, string asset) {
		var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		box.show ();

		//var pixbuf = Gdk.Pixbuf.from_file (@"../assets/iconset-temp/$asset.png");
		//pixbuf.scale_simple (size, size, Gdk.InterpType.BILINEAR);

		//var image = Gtk.Image.from_pixbuf (pixbuf);
		Gtk.Image image;
		if (asset.has_prefix ("gtk-")){
			image =	new Gtk.Image.from_icon_name (asset, Gtk.IconSize.MENU);
		} else {
			image =	new Image (asset, "orange", 16);
		}

		image.show ();

		var lbl = new Gtk.Label ((owned) label);
		lbl.show ();
		lbl.set_halign (Gtk.Align.START);

		lbl.set_use_underline (true);
		//lbl.set_xalign ((double) 0.0);

		box.pack_start (image, false, false, 0);
		box.pack_start (lbl, false, false, 5);

		add(box);
	}
}
