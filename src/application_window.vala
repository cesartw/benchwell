public class Benchwell.ApplicationWindow : Gtk.ApplicationWindow {
	public Gtk.Notebook notebook;
	public Gtk.ListStore env_store;
	public Gtk.Button btn_env;
	public Gtk.ComboBox env_combo;

	public SimpleAction new_connection_menu;
	public SimpleAction new_database_tab_menu;
	public SimpleAction new_http_tab_menu;
	public SimpleAction new_tab_menu;
	public SimpleAction close_menu;


	public class ApplicationWindow(Gtk.Application app) {
		Object (
			application: app
		);

		set_title ("Benchwell");

		new_database_tab_menu = new SimpleAction("new.db", null);
		new_http_tab_menu = new SimpleAction("new.http", null);
		new_tab_menu = new SimpleAction("new.tab", null);
		close_menu = new SimpleAction("close", null);

		add_action(new_database_tab_menu);
		add_action(new_http_tab_menu);
		add_action(new_tab_menu);
		add_action(close_menu);

		notebook = new Gtk.Notebook ();
		notebook.set_name ("MainNotebook");
		notebook.set_property ("scrollable", true);
		notebook.set_group_name ("mainwindow");
		notebook.popup_enable ();

		notebook.set_property ("tab-pos", Config.tab_position());
		notebook.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		box.show ();

		// header bar
		var header = new Gtk.HeaderBar ();
		header.title ="Benchwell";
		header.subtitle ="version";
		header.show_close_button =true;
		header.show ();

		var window_btn_menu = new Gtk.MenuButton();
		window_btn_menu.show ();

		var add_img = new Benchwell.Image ("add-tab", Gtk.IconSize.BUTTON);
		add_img.show ();
		window_btn_menu.set_image (add_img);

		var window_menu = new GLib.Menu ();
		window_btn_menu.set_menu_model (window_menu);

		window_menu.append (_("Window"), "app.new");
		window_menu.append (_("Database"), "win.new.db");
		window_menu.append (_("HTTP"), "win.new.http");

		var env = env_selector ();
		env.show ();
		/////////////

		set_titlebar (header);

		header.pack_start (window_btn_menu);
		//header.pack_end (app_btn_menu);
		header.pack_end (env);

		box.pack_start (notebook, true, true, 0);

		add(box);

		set_default_size (Config.window_width (), Config.window_height ());
		move (Config.window_position_x (), Config.window_position_y ());

		new_database_tab_menu.activate.connect (() => {
			add_database_tab ();
		});

		new_http_tab_menu.activate.connect (() => {
			add_http_tab ();
		});

		close_menu.activate.connect (() => {
			notebook.remove_page (notebook.get_current_page ());
		});

		var css_provider = new Gtk.CssProvider ();
        css_provider.load_from_resource ("/io/benchwell/stylesheet.css");
		Gtk.StyleContext.add_provider_for_screen (
			Gdk.Screen.get_default (), css_provider, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
		);
	}

	public void add_database_tab (Benchwell.Backend.Sql.ConnectionInfo? connection_info=null, Benchwell.Backend.Sql.TableDef? tabledef = null) {
		var tab  = new Benchwell.Tab ();
		tab.show ();

		var database = new Benchwell.Database.Database (this, new Benchwell.Services.Database ());
		database.notify["title"].connect ((s, p) => {
			tab.label.set_text (database.title);
			tab.label.tooltip_text = database.title;
		});
		database.show ();

		tab.label.set_text ( database.title );
		tab.pack_start (database, true, true, 0);

		notebook.append_page (tab, tab.header);
		tab.btn.clicked.connect( () => {
			notebook.remove_page (notebook.page_num (tab));
		});
		notebook.set_current_page (notebook.get_n_pages () - 1);

		if (connection_info != null) {
			database.launch_connection (connection_info, tabledef);
		}
	}

	public void add_http_tab () {
		var tab  = new Benchwell.Tab ();
		tab.show ();

		var http = new Benchwell.Http.Http (this);
		http.notify["title"].connect ((s, p) => {
			tab.label.set_text (http.title);
			tab.label.tooltip_text = http.title;
		});
		http.show ();

		tab.label.set_text ( http.title );
		tab.pack_start (http, true, true, 0);

		notebook.append_page (tab, tab.header);
		tab.btn.clicked.connect( () => {
			notebook.remove_page (notebook.page_num (tab));
		});
		notebook.set_current_page (notebook.get_n_pages () - 1);
	}

	private Gtk.Grid env_selector () {
		env_store = new Gtk.ListStore (2, GLib.Type.INT64, GLib.Type.STRING);

		env_combo = new Gtk.ComboBox.with_model_and_entry (env_store);
		env_combo.set_id_column (0);
		env_combo.set_entry_text_column (1);
		env_combo.show ();

		Config.environments.foreach( (env) => {
			Gtk.TreeIter iter;
			env_store.append (out iter);

			env_store.set_value (iter, 0, env.id);
			env_store.set_value (iter, 1, env.name);
		});

		btn_env = new Benchwell.Button ("config", Gtk.IconSize.BUTTON);
		btn_env.show ();

		var popover = new Gtk.Popover (btn_env);
		btn_env.clicked.connect (() => {
			popover.show ();
		});

		var editor = new Benchwell.EnvironmentEditor ();
		editor.show ();
		popover.add (editor);

		var grid = new Gtk.Grid ();
		grid.attach (env_combo, 0, 0, 4, 1);
		grid.attach (btn_env, 4, 0, 1, 1);
		grid.get_style_context ().add_class ("link");

		env_combo.changed.connect (() => {
			Gtk.TreeIter? iter = null;
			env_combo.get_active_iter (out iter);
			if ( iter == null ) {
				return;
			}

			GLib.Value val;
			env_store.get_value (iter, 0, out val);
			var id = val.get_int64 ();
			foreach (var env in Config.environments) {
				if (env.id == id) {
					Config.environment = env;
					break;
				}
			}
		});

		return grid;
	}
}
