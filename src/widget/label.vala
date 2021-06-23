public class Benchwell.Label : Gtk.EventBox {
	public signal void clicked ();
	public Gtk.Label label {owned get; construct;}
	public int max_width_chars {
		get { return label.max_width_chars ; }
		set { label.max_width_chars = value; }
	}

	public Pango.EllipsizeMode ellipsize {
		get { return label.ellipsize; }
		set { label.ellipsize = value; }
	}

	public Label (string text) {
		Object (
			label: new Gtk.Label (text)
		);

		label.show ();

		add (label);

		button_press_event.connect (() => {
			clicked ();
			return true;
		});
	}

	public string get_text () {
		return label.get_text ();
	}

	public void set_text (string text) {
		label.set_text (text);
	}

	public new unowned Gtk.StyleContext get_style_context () {
		return label.get_style_context ();
	}
}
