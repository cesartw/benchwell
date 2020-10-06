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
