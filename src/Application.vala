public class Benchwell.Application : Gtk.Application {
	private SimpleAction new_window_action;
	private SimpleAction preference_action;
	private SimpleAction dark_mode_action;
	private Config config;

	public Application () {
		Object(
			application_id: "io.benchwell",
			flags: ApplicationFlags.FLAGS_NONE
		);

		// config initialization
		config = new Config ();

		new_window_action = new SimpleAction ("new", null);
		preference_action = new SimpleAction ("preference", null);
		dark_mode_action = new SimpleAction ("darkmode", null);

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
		Gtk.Settings.get_default ().gtk_application_prefer_dark_theme = true;
	}
}

public static int main(string[] args) {
	Gtk.init (ref args);

	var app = new Benchwell.Application ();

	return app.run (args);
}
