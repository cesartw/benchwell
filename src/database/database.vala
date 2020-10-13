public class Benchwell.Database.Database : Gtk.Box {
	public string title { get; set; }
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.Services.Database service { get; construct; }
	private Benchwell.Database.Connect? connect_view;
	private Benchwell.Database.Data? data_view;
	private Benchwell.Backend.Sql.Engine engine;

	public Database (Benchwell.ApplicationWindow window, Benchwell.Services.Database service) {
		Object(
			window: window,
			service: service,
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);
		engine = new Benchwell.Backend.Sql.Engine ();
		title = _("Connect");

		connect_view = new Benchwell.Database.Connect (window);
		//data = new Benchwell.Database.Data (window);
		connect_view.show ();

		show_connect ();

		connect_view.dbconnect.connect ((c) => {
			launch_connection (c);
			//if (c.password != "") {
				//launch_connection (c);
			//} else {
				//Config.decrypt.begin (c, (obj, res) => {
					//c.password = Config.decrypt.end (res);
					//launch_connection (c);
				//});
			//}
		});
	}

	public void launch_connection (owned Benchwell.Backend.Sql.ConnectionInfo c, Benchwell.Backend.Sql.TableDef? selected_tabledef = null) {
		//Benchwell.Backend.Sql.Connection connection;

		try {
			//connection = engine.connect (c);
			service.connect (c);
		} catch (Benchwell.Backend.Sql.Error err) {
			show_error_dialog (err.message);
			return;
		}

		data_view = new Benchwell.Database.Data(window, service);
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