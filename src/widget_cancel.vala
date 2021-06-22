namespace Benchwell {
	public class CancelOverlay : Gtk.Overlay {
		public delegate void OnCancelFunc ();

		public Gtk.Button btn_cancel;
		public Gtk.Spinner spinner;
		public Gtk.Box controls;
		private OnCancelFunc* cancel;
		public Gtk.Widget overlayed { construct; }

		public CancelOverlay (Gtk.Widget overlayed) {
			Object(
				overlayed: overlayed
			);

			btn_cancel = new Gtk.Button.with_label (_("Cancel"));
			btn_cancel.set_size_request (100, 30);
			btn_cancel.show ();

			spinner = new Gtk.Spinner ();

			controls = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);

			var actions = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
			actions.set_size_request (100, 150);
			actions.valign = Gtk.Align.CENTER;
			actions.halign = Gtk.Align.CENTER;
			actions.vexpand = true;
			actions.hexpand = true;
			actions.show ();

			actions.pack_start (spinner, true, true, 0);
			actions.pack_start (btn_cancel, false, false, 0);
			controls.add (actions);

			add (overlayed);

			btn_cancel.clicked.connect ( () => {
				stop ();
				((OnCancelFunc)(*cancel)) ();
			});
		}

		public void run (OnCancelFunc c) {
			controls.show ();
			spinner.show ();
			add_overlay (controls);
			cancel = &c;
		}

		public void stop () {
			remove (controls);
			spinner.stop ();
		}
	}
}
