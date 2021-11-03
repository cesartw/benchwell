public class Benchwell.Http.Overlay : Gtk.Overlay {
	public Gtk.Button btn_cancel;
	public Gtk.Spinner spinner;
	public Gtk.Box box;

	public signal void cancel ();

	public Overlay () {
		Object(
			name: "HttpOverlay"
		);

		btn_cancel = new Gtk.Button.with_label (_("Cancel"));
		btn_cancel.set_size_request (100, 30);
		btn_cancel.show ();

		spinner = new Gtk.Spinner ();
		spinner.show ();

		box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		box.get_style_context ().add_class ("overlay-bg");

		var center_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		center_box.set_size_request (100, 150);
		center_box.valign = Gtk.Align.CENTER;
		center_box.halign = Gtk.Align.CENTER;
		center_box.vexpand = true;
		center_box.hexpand = true;
		center_box.show ();

		box.add (center_box);

		center_box.pack_start (spinner, true, true, 0);
		center_box.pack_start (btn_cancel, false, false, 0);
		add_overlay (box);

		btn_cancel.clicked.connect (on_cancel);
	}

	private void on_cancel () {
		spinner.stop ();
		box.hide ();
		cancel ();
	}

	public void start () {
		box.show ();
		spinner.start ();
	}

	public void stop () {
		spinner.stop ();
		box.hide ();
	}
}
