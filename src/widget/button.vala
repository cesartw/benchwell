public class Benchwell.Button : Gtk.Button {
	public Button (string asset, Gtk.IconSize? size = null) {
		Object();

		var img = new Benchwell.Image(asset, size);
		set_image(img);
	}
}

public class Benchwell.ToggleButton : Gtk.ToggleButton {
	public ToggleButton (string asset, Gtk.IconSize? size = null) {
		Object ();

		var img = new Benchwell.Image(asset, size);
		set_image (img);
	}
}

public class Benchwell.LoadButton : Gtk.Button {
	private Gtk.Spinner spinner;
	private bool _loading = false;
	private Gtk.Label lbl;
	public bool loading {
		get {
			return _loading;
		}
		set {
			_loading = value;
			if (_loading) {
				sensitive = false;
				spinner.start ();
				spinner.show ();
				remove (lbl);
				add (spinner);
			} else {
				sensitive = true;
				spinner.stop ();
				spinner.hide ();
				remove (spinner);
				add (lbl);
			}
		}
	}

	public LoadButton (string l) {
		Object(
			label: l
		);

		lbl = get_child () as Gtk.Label;
		spinner = new Gtk.Spinner ();
		remove (lbl);
	}
}
