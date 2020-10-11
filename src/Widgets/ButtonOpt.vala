public class Benchwell.OptButton : Gtk.Grid {
	public Gtk.MenuButton menu_btn;
	public GLib.Menu menu;
	public Gtk.Button btn;

	public OptButton (string label, ...) {
		Object();
		get_style_context ().add_class ("linked");

		menu_btn = new Gtk.MenuButton ();
		menu_btn.show ();

		menu = new GLib.Menu ();

		var l = va_list ();
		while (true) {
			string? key = l.arg ();
			if (key == null) {
				break;
			}

			string? action = l.arg ();
			if (action == null) {
				break;
			}

			menu.append (key, action);
		}
		menu_btn.set_menu_model (menu);

		btn = new Gtk.Button.with_label (label);
		btn.show ();

		attach(btn, 0, 0, 2, 1);
		attach(menu_btn, 2, 0, 1, 1);
	}
}
