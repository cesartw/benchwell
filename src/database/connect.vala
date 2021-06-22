namespace Benchwell {
	namespace Database {
		public class Connect : Gtk.Paned {
			public Benchwell.ApplicationWindow window { get; construct; }
			public Benchwell.Database.ConnectionList connections;

			public Gtk.Button btn_connect;
			public Gtk.Button btn_test;

			private Gtk.Stack stack;
			private MysqlForm mysql;
			private SQLiteForm sqlite;

			public signal void dbconnect (Benchwell.ConnectionInfo c);

			public Connect (Benchwell.ApplicationWindow window) {
				Object(
					window: window,
					orientation: Gtk.Orientation.HORIZONTAL,
					wide_handle: true
				);

				connections = new Benchwell.Database.ConnectionList ();
				connections.show ();

				var leftsw = new Gtk.ScrolledWindow (null, null);
				leftsw.show ();
				leftsw.add (connections);

				connections.update_items ();

				var btn_box = new Gtk.ButtonBox (Gtk.Orientation.HORIZONTAL);
				btn_box.set_layout (Gtk.ButtonBoxStyle.EDGE);
				btn_box.show ();

				btn_connect = new Gtk.Button.with_label (_("Connect"));
				btn_connect.get_style_context ().add_class ("suggested-action");
				btn_connect.show ();

				btn_test = new Gtk.Button.with_label (_("Test"));
				btn_test.show ();

				btn_box.add (btn_connect);
				btn_box.add (btn_test);

				var adapter_stack = new Gtk.StackSwitcher () {
					vexpand = true,
					hexpand = true,
					homogeneous = true
				};
				adapter_stack.icon_size = Gtk.IconSize.DIALOG;
				adapter_stack.show ();

				stack = new Gtk.Stack () {
					vexpand = true,
					hexpand = true
				};
				stack.show();
				adapter_stack.set_stack (stack);

				var size_group = new Gtk.SizeGroup (Gtk.SizeGroupMode.BOTH);
				mysql = new MysqlForm(window, size_group);
				mysql.show ();

				sqlite = new SQLiteForm(window, size_group);
				sqlite.show ();

				stack.add_titled (mysql, "mysql", "Mysql");
				stack.add_titled (sqlite, "sqlite", "SQLite");
				stack.set_visible_child_name ("mysql");

				stack.child_set_property (mysql, "icon-name", "bw-mariadb");
				stack.child_set_property (sqlite, "icon-name", "bw-sqlite");

				var formbox = new Gtk.Box (Gtk.Orientation.VERTICAL, 5) {
					valign = Gtk.Align.CENTER,
					halign = Gtk.Align.CENTER
				};
				formbox.show ();
				//formbox.set_size_request (500, 300);

				formbox.pack_start (adapter_stack, true, true, 0);
				formbox.pack_start (stack, true, false, 0);
				formbox.pack_end (btn_box, false, true, 0);

				pack1 (leftsw, true, true);
				pack2 (formbox, true, false);

				connections.row_selected.connect (on_row_selected);

				enable_buttons (false);

				mysql.changed.connect ((conn) => {
					var ok = Benchwell.MysqlDB.validate_connection (conn);
					enable_buttons (ok);
				});

				btn_connect.clicked.connect (() => {
					var c = get_connection ();
					dbconnect (c);
				});

				var engine = new Benchwell.Engine ();
				btn_test.clicked.connect (() => {
					var c = get_connection ();
					if (c == null) {
						return;
					}

					try {
						engine.connect (c);
					} catch (Benchwell.Error err) {
						Config.show_alert(this, err.message);
					}
				});

				connections.row_activated.connect (() => {
					var c = get_connection ();
					dbconnect (c);
				});
			}

			public Benchwell.ConnectionInfo? get_connection () {
				switch (stack.get_visible_child_name ()) {
					case "mysql":
						return mysql.get_connection ();
					case "sqlite":
						return sqlite.get_connection ();
				}

				return null;
			}

			private void on_row_selected (Gtk.ListBox list, Gtk.ListBoxRow? row) {
				enable_buttons (false);

				if ( row == null ) {
					return;
				}

				var index = row.get_index ();
				if ( index < 0 ){
					return;
				}

				var conn = Config.connections.at (index);

				if (conn == null) {
					return;
				}

				switch (conn.adapter){
					case "mysql":
						stack.set_visible_child_name ("mysql");
						mysql.set_connection (ref conn);
						break;
					case "sqlite":
						stack.set_visible_child_name ("sqlite");
						break;
				}

				var ok = Benchwell.MysqlDB.validate_connection (conn);
				enable_buttons (ok);
			}

			private void enable_buttons (bool ok) {
				btn_test.set_sensitive (ok);
				btn_connect.set_sensitive (ok);
			}
		}

		public class ConnectionList : Gtk.ListBox {
			private Gtk.Menu menu;
			private Gtk.MenuItem menu_new;
			private Gtk.MenuItem menu_connect;
			private Gtk.MenuItem menu_test;
			private Gtk.MenuItem menu_del;

			public ConnectionList () {
				Object(
					activate_on_single_click: false
				);

				get_style_context ().add_class ("bw-spacing");

				// connection menu
				menu = new Gtk.Menu ();
				menu_new = new Benchwell.MenuItem (_("New"), "connection");
				menu_new.show ();
				menu_connect = new Benchwell.MenuItem (_("Connect"), "next");
				menu_connect.show ();
				menu_test = new Benchwell.MenuItem (_("Test"), "refresh");
				menu_test.show ();
				menu_del = new Benchwell.MenuItem (_("Delete"), "close");
				menu_del.show ();

				menu.add (menu_new);
				menu.add (menu_connect);
				menu.add (menu_test);
				menu.add ( menu_del);

				button_press_event.connect((list, event) => {
					if (event.button != Gdk.BUTTON_SECONDARY){
						return false;
					}

					grab_focus ();
					select_row (get_row_at_y ((int)event.y));

					menu.show ();
					menu.popup_at_pointer (event);

					return true;
				});

				menu_new.activate.connect (on_new);
				menu_del.activate.connect (on_delete);
			}

			public void update_items (int64 selected_id = 0) {
				get_children().foreach( (row) => {
					remove (row);
				});

				Config.connections.for_each ((item) => {
					var conn = item as ConnectionInfo;
					var row = build_row (conn);
					add (row);
					if (conn.id == selected_id) {
						select_row (row);
					}
					return false;
				});
			}

			private Gtk.ListBoxRow build_row (Benchwell.ConnectionInfo item) {
				var row = new Gtk.ListBoxRow ();
				row.show ();

				var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
				box.show ();

				var label = new Gtk.Label (item.to_string());
				label.set_halign (Gtk.Align.START);
				label.show ();

				var image = new Benchwell.Image ("connection", Gtk.IconSize.BUTTON);
				image.show ();

				box.pack_start (image, false, false, 5);
				box.pack_start (label, false, false, 0);

				row.add (box);

				item.notify["name"].connect (() => {
					label.set_text (item.name);
				});

				return row;
			}

			private void on_new () {
				try {
					var c = Config.connections.add () as ConnectionInfo;

					c.adapter = "mysql";
					c.ttype = "tcp";
					c.port = 3306;

					update_items (c.id);
				} catch (Benchwell.ConfigError err) {
					Config.show_alert (this, err.message);
					return;
				}
			}

			private void on_delete () {
				var row = get_selected_row ();
				var conn = Config.connections.at (row.get_index ());
				try {
					conn.remove ();
					Config.connections.remove (conn);
				} catch (Benchwell.ConfigError err) {
					Config.show_alert (this, err.message);
					return;
				}
				update_items ();
			}
		}
	}
}
