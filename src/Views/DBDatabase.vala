public class Benchwell.Views.DBDatabase : Gtk.Box {
	public string title { get; set; }
	public Benchwell.ApplicationWindow window { get; construct; }
	private Benchwell.Views.DBConnect? connect_view;
	private Benchwell.Views.DBData? data_view;
	private Benchwell.SQL.Engine engine;

	public DBDatabase (Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);
		engine = new Benchwell.SQL.Engine ();
		title = _("Connect");

		connect_view = new Benchwell.Views.DBConnect (window);
		//data = new Benchwell.Views.DBData (window);
		connect_view.show ();

		show_connect ();

		connect_view.dbconnect.connect ((c) => {
			if (c.password != "") {
				launch_connection (c);
			} else {
				Config.decrypt.begin (c, (obj, res) => {
					c.password = Config.decrypt.end (res);
					launch_connection (c);
				});
			}
		});
	}

	public void launch_connection (Benchwell.SQL.ConnectionInfo c, Benchwell.SQL.TableDef? selected_tabledef = null) {
		Benchwell.SQL.Connection connection;

		try {
			connection = engine.connect (c);
		} catch (Benchwell.SQL.Error err) {
			show_error_dialog (err.message);
			return;
		}

		data_view = new Benchwell.Views.DBData(window, connection, c);
		show_data ();

		title = c.name;
		if (c.database != "") {
			title = @"$(c.name).$(c.database)";
		}
		data_view.database_selected.connect ((dbname) => {
			title = @"$(c.name).$(dbname)";
		});

		data_view.tables.new_tab_menu.activate.connect (() => {
			var tabledef = data_view.tables.selected_tabledef;
			if (tabledef == null) {
				return;
			}

			window.add_database_tab (c, tabledef);
		});

		data_view.tables.selected_tabledef = selected_tabledef;
	}

	public void show_error_dialog (string message) {
		var dialog = new Gtk.MessageDialog (window, Gtk.DialogFlags.DESTROY_WITH_PARENT, Gtk.MessageType.ERROR, Gtk.ButtonsType.CLOSE, message);
		dialog.response.connect (dialog.destroy);
		dialog.show ();
	}

	public void show_data () {
		remove (connect_view);
		add (data_view);
		data_view.show ();
	}

	public void show_connect () {
		if (data_view != null) {
			remove (data_view);
		}
		add (connect_view);
		connect_view.show ();
		//add (connect);
	}
}
