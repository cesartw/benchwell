public class Benchwell.Button : Gtk.Button {
	public Button (string asset, string color, int size) {
		Object();

		var img = new Benchwell.Image(asset, color, size);
		set_image(img);
	}
}

public class Benchwell.ToggleButton : Gtk.ToggleButton {
	public ToggleButton (string asset, string color, int size) {
		Object ();

		var img = new Image(asset, color, size);
		set_image (img);
	}
}
