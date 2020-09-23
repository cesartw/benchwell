public class Benchwell.Views.DBDatabase : Gtk.Box {
	public string title { get; set; }
	public Benchwell.ApplicationWindow window { get; construct; }
	private Benchwell.Views.DBConnect connect_view;
	private Benchwell.Views.DBData data_view;
	private Benchwell.SQL.Engine engine;

	public DBDatabase (Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);
		engine = new Benchwell.SQL.Engine ();

		connect_view = new Benchwell.Views.DBConnect (window);
		//data = new Benchwell.Views.DBData (window);
		connect_view.show ();

		show_connect ();

		connect_view.dbconnect.connect ((c) => {
			Benchwell.SQL.Connection connection;
			try {
				connection = engine.connect (c);
			} catch (Benchwell.SQL.ErrorConnection e) {
				return;
			}

			data_view = new Benchwell.Views.DBData(window, connection, c);
			show_data ();
		});
	}

	public void show_data () {
		remove (connect_view);
		add (data_view);
		data_view.show ();
	}

	public void show_connect () {
		remove (data_view);
		add (connect_view);
		connect_view.show ();
		//add (connect);
	}
}
