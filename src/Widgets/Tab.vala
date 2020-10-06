public class Benchwell.Tab : Gtk.Box {
	public Gtk.Label label;
	public Benchwell.Button btn;
	public Gtk.Box header;

	public Tab() {
		Object(
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 0
		);

		set_vexpand (true);
		set_hexpand (true);

		btn = new Benchwell.Button("close", Gtk.IconSize.MENU);
		btn.show ();
		btn.set_relief (Gtk.ReliefStyle.NONE);

		label = new Gtk.Label ("");
		label.set_ellipsize (Pango.EllipsizeMode.END);
		label.set_width_chars(20);
		label.show ();

		header = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		header.pack_start (label, true, true, 0);
		header.pack_end (btn, false, false, 0);
		header.show ();
	}
}
