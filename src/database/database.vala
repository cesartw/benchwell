public class Benchwell.Database.Database : Gtk.Box {
	public string title { get; set; }
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.DatabaseService service { get; construct; }
	private Benchwell.Database.Connect? connect_view;
	private Benchwell.Database.Data? data_view;
	private Benchwell.Engine engine;

	public Database (Benchwell.ApplicationWindow window, Benchwell.DatabaseService service) {
		Object(
			window: window,
			service: service,
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);
		engine = new Benchwell.Engine ();
		title = _("Connect");

		connect_view = new Benchwell.Database.Connect (window);
		connect_view.show ();

		show_connect ();

		connect_view.dbconnect.connect ((c) => {
			launch_connection.begin (c);
		});
	}

	public async void launch_connection (owned Benchwell.ConnectionInfo c, Benchwell.TableDef? selected_tabledef = null) {
		try {
			yield service.dbconnect (c);
		} catch (Benchwell.Error err) {
			Config.show_alert (this, err.message);
			return;
		} catch (GLib.Error err) {
			Config.show_alert (this, err.message);
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

		data_view.disconnect_btn.clicked.connect (() => {
			service.dbdisconnect ();
			show_connect ();
			data_view.dispose ();
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
