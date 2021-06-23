namespace Benchwell {
	public class ComboTab : Gtk.Box {
		public Gtk.ComboBoxText combo;
		public Gtk.Label label;
		public bool enabled  { get; set; }

		public ComboTab (string l, bool enabled = false) {
			Object (
				orientation: Gtk.Orientation.VERTICAL,
				spacing: 0
			);

			combo = new Gtk.ComboBoxText ();
			label = new Gtk.Label (l);

			pack_start (combo, true, true, 0);
			pack_start (label, true, true, 0);
			combo.changed.connect (() => {
				label.set_text (combo.get_active_text ());
			});

			this.enabled = enabled;
			on_toggle ();
			notify["enabled"].connect (on_toggle);
		}

		private void on_toggle () {
			if (enabled) {
				combo.show ();
				label.hide ();
				return;
			}
			combo.hide ();
			label.show ();
		}
	}
}
