public class Benchwell.ApplicationWindow : Gtk.ApplicationWindow {
	public Gtk.Notebook notebook;
	public Gtk.ListStore env_store;
	public Gtk.Button btn_env;
	public Gtk.ComboBox env_combo;

	public SimpleAction new_connection_action;
	public SimpleAction new_database_tab_action;
	public SimpleAction new_http_tab_action;
	public SimpleAction new_tab_action;
	public SimpleAction close_action;
	public SimpleAction copy_curl_action;

	public Gtk.InfoBar infobar;
	public Gtk.Label infobar_label;

	public class ApplicationWindow(Gtk.Application app) {
		Object (
			application: app
		);

		icon_name = "io.benchwell";
		set_title ("Benchwell");

		new_database_tab_action = new SimpleAction("new.db", null);
		new_http_tab_action = new SimpleAction("new.http", null);
		new_tab_action = new SimpleAction("new.tab", null);
		close_action = new SimpleAction("close", null);
		copy_curl_action = new SimpleAction("copycurl", null);

		add_action(new_database_tab_action);
		add_action(new_http_tab_action);
		add_action(new_tab_action);
		add_action(close_action);
		add_action(copy_curl_action);

		notebook = new Gtk.Notebook ();
		notebook.set_name ("MainNotebook");
		notebook.scrollable = true;
		notebook.group_name = "mainwindow";
		notebook.tab_pos = Config.tab_position ();
		notebook.show_border = true;
		notebook.key_press_event.connect ( (e) => {
			var page = 0;
			switch (e.keyval) {
				case Gdk.Key.@1:
					page = 1;
					break;
				case Gdk.Key.@2:
					page = 2;
					break;
				case Gdk.Key.@3:
					page = 3;
					break;
				case Gdk.Key.@4:
					page = 4;
					break;
				case Gdk.Key.@5:
					page = 5;
					break;
				case Gdk.Key.@6:
					page = 6;
					break;
				case Gdk.Key.@7:
					page = 7;
					break;
				case Gdk.Key.@8:
					page = 8;
					break;
				case Gdk.Key.@9:
					page = 9;
					break;
				default:
					return false;
			}

			if (e.state != Gdk.ModifierType.MOD1_MASK) {
				return false;
			}

			if (page <= notebook.get_n_pages ()) {
				notebook.set_current_page (page - 1);
				return true;
			}

			return false;
		});
		notebook.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		box.show ();

		// header bar
		var logo_box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		logo_box.show ();

		var header_logo = new Gtk.Image.from_icon_name ("io.benchwell", Gtk.IconSize.LARGE_TOOLBAR);
		header_logo.show ();

		var titles_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		titles_box.show ();

		var header_title = new Gtk.Label ("Benchwell");
		header_title.get_style_context ().add_class ("title");
		header_title.show ();

		var header_subtitle = new Gtk.Label (Benchwell.Constants.VERSION);
		header_subtitle.get_style_context ().add_class ("subtitle");
		header_subtitle.show ();

		titles_box.pack_start (header_title, false, false, 0);
		titles_box.pack_start (header_subtitle, false, false, 0);

		logo_box.pack_start (header_logo, false, false, 0);
		logo_box.pack_start (titles_box, false, false, 0);

		var header = new Gtk.HeaderBar ();
		header.custom_title = logo_box;
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

		var pomodoro = new Benchwell.Pomodoro ();
		pomodoro.show ();

		header.pack_start (window_btn_menu);
		header.pack_end (env);
		header.pack_end (pomodoro);

		box.pack_start (notebook, true, true, 0);

		infobar = new Gtk.InfoBar ();
		var infobar_label_sw = new Gtk.ScrolledWindow (null, null);
		infobar_label_sw.show ();
		infobar_label_sw.hexpand = true;
		infobar.no_show_all = true;
		infobar.add_button (_("Ok"), Gtk.ResponseType.OK);

		infobar_label = new Gtk.Label("");
		infobar_label.wrap = true;
		infobar_label.show ();
		infobar_label_sw.add (infobar_label);
		infobar.get_content_area ().add (infobar_label_sw);
		infobar.response.connect (infobar.hide);
		box.pack_end (infobar, false, false);

		add(box);

		set_default_size (Config.settings.window_size_w,
							Config.settings.window_size_h);
		move (Config.settings.window_pos_x,
			Config.settings.window_pos_y);

		new_database_tab_action.activate.connect (() => {
			add_database_tab ();
		});

		new_http_tab_action.activate.connect (() => {
			add_http_tab ();
		});

		close_action.activate.connect (() => {
			notebook.remove_page (notebook.get_current_page ());
		});

		var css_provider = new Gtk.CssProvider ();
        css_provider.load_from_resource ("/io/benchwell/stylesheet.css");
		Gtk.StyleContext.add_provider_for_screen (
			Gdk.Screen.get_default (), css_provider, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
		);

		delete_event.connect (before_destroy);
	}

	private bool before_destroy () {
		int width, height, x, y;
		get_size (out width, out height);
		get_position (out x, out y);

		Config.settings.window_size_w = width;
		Config.settings.window_size_h = height;
		Config.settings.window_pos_x =  x;
		Config.settings.window_pos_y =  y;
		Config.save_http_tree_state ();

		return false;
	}

	public void add_database_tab (Benchwell.ConnectionInfo? connection_info=null, Benchwell.TableDef? tabledef = null) {
		var tab  = new Benchwell.Tab ();
		tab.show ();

		var database = new Benchwell.Database.Database (this, new Benchwell.DatabaseService ());
		database.notify["title"].connect ((s, p) => {
			tab.label.set_text (database.title);
			tab.label.tooltip_text = database.title;
		});
		database.show ();

		tab.label.set_text ( database.title );
		tab.pack_start (database, true, true, 0);

		notebook.append_page (tab, tab.header);
		notebook.set_tab_reorderable (tab, true);
		tab.btn.clicked.connect( () => {
			notebook.remove_page (notebook.page_num (tab));
		});
		notebook.set_current_page (notebook.get_n_pages () - 1);

		if (connection_info != null) {
			database.launch_connection.begin (connection_info, tabledef);
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
		notebook.set_tab_reorderable (tab, true);
		tab.btn.clicked.connect( () => {
			notebook.remove_page (notebook.page_num (tab));
		});
		notebook.set_current_page (notebook.get_n_pages () - 1);
	}

	private Gtk.Grid env_selector () {
		env_store = new Gtk.ListStore (2, GLib.Type.INT64, GLib.Type.STRING);

		var env_cell = new Gtk.CellRendererText ();
		var env_cell_box = new Gtk.CellAreaBox ();
		env_cell_box.pack_start (env_cell, true);
		env_cell_box.add_attribute (env_cell, "text", 1);

		env_combo = new Gtk.ComboBox.with_area (env_cell_box);
		env_combo.set_model (env_store);
		env_combo.set_id_column (0);
		env_combo.set_entry_text_column (1);
		env_combo.show ();

		btn_env = new Benchwell.Button ("config", Gtk.IconSize.BUTTON);
		btn_env.show ();

		var popover = new Gtk.Popover (btn_env);
		btn_env.clicked.connect (() => {
			popover.show ();
		});

		var settings_panel = new Benchwell.SettingsPanel (this);
		settings_panel.show ();
		popover.add (settings_panel);

		var grid = new Gtk.Grid ();
		grid.attach (env_combo, 0, 0, 4, 1);
		grid.attach (btn_env, 4, 0, 1, 1);
		grid.get_style_context ().add_class ("link");

		var selected_environment_id = Config.settings.environment_id;
		for (var i = 0; i < Config.environments.length; i++) {
			var env = Config.environments[i];
			Gtk.TreeIter iter;
			env_store.append (out iter);

			env_store.set_value (iter, 0, env.id);
			env_store.set_value (iter, 1, env.name);
			if (selected_environment_id == env.id) {
				env_combo.set_active_iter (iter);
				Config.environment = env;
				settings_panel.select_env (env);
			}
		}

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
					settings_panel.select_env (env);
					break;
				}
			}

			Config.settings.environment_id = Config.environment.id;
		});

		Config.environment_added.connect ((env) => {
			Gtk.TreeIter iter;
			env_store.append (out iter);
			env_store.set_value (iter, 0, env.id);
			env_store.set_value (iter, 1, env.name);
		});

		Config.environment_removed.connect ((env) => {
			Gtk.TreeIter? iter = null;
			env_store.foreach ((model, path, i) => {
				iter = i;
				GLib.Value val;
				env_store.get_value (iter, 0, out val);
				return val.get_int64 () == env.id;
			});

			if (iter != null) {
				env_store.remove (ref iter);
			}
		});

		return grid;
	}

	public void show_alert (string message, Gtk.MessageType type = Gtk.MessageType.ERROR, bool autohide = false, int timeout = 0) {
		infobar_label.set_text (message);
		infobar.message_type = type;
		infobar.show ();

		if (autohide) {
			if (timeout == 0) {
				switch (message.split (" ").length) {
					case 1:
						timeout = 1000;
						break;
					case 2, 3:
						timeout = 1500;
						break;
					default:
						timeout = 3000;
						break;
				}
			}

			Timeout.add (timeout, () => {
				infobar.hide ();
				return false;
			});
		}
	}
}
