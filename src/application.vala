public class Benchwell.Application : Gtk.Application {
	private SimpleAction new_window_action;
	private SimpleAction preference_action;

	public Application () {
		Object(
			application_id: "io.benchwell",
			flags: ApplicationFlags.FLAGS_NONE
		);

		new_window_action = new SimpleAction ("new", null);
		preference_action = new SimpleAction ("preference", null);

		set_accels_for_action ("win.new.db", {"<control>D"});
		set_accels_for_action ("win.new.http", {"<control>H"});
		set_accels_for_action ("win.new.tab", {"<control>T"});
		set_accels_for_action ("win.close", {"<control>W"});
	}

	protected override void activate () {
		var window = new Benchwell.ApplicationWindow (this);
		window.window_position = Gtk.WindowPosition.CENTER;
		add_window (window);
		window.show ();

		window.add_database_tab ();
		window.add_http_tab ();

		Gtk.Settings.get_default ().gtk_application_prefer_dark_theme = Config.settings.dark_mode;
	}

	public void get_cursor_position (out int x, out int y) {
		var pointer = Gdk.Display.get_default ().get_default_seat ().get_pointer ();
		Gdk.Screen screen;
		pointer.get_position (out screen, out x, out y);
	}
}

public static int main(string[] args) {
	Gtk.init (ref args);

	// config initialization
	try {
		Benchwell.Config = new Benchwell._Config ();
	} catch (Benchwell.ConfigError err) {
		stderr.printf (err.message);
		return 1;
	}


	var app = new Benchwell.Application ();

	return app.run (args);
}
