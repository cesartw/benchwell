public class Benchwell.Views.DBConnect : Gtk.Paned {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.Views.DBConnectionList connections;
	public string title;

	// form buttons
	public Gtk.Button btn_connect;
	public Gtk.Button btn_test;
	public Gtk.Button btn_save;

	private Gtk.Stack stack;
	private MysqlForm mysql;
	private SQLiteForm sqlite;

	public signal void dbconnect(Benchwell.SQL.ConnectionInfo c);

	public DBConnect (Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL
		);

		title = "Connect";

		set_vexpand (true);
		set_hexpand (true);
		set_wide_handle (true);

		var leftsw = new Gtk.ScrolledWindow (null, null);
		leftsw.show ();

		connections = new Benchwell.Views.DBConnectionList ();
		connections.set_hexpand (true);
		connections.set_vexpand (true);
		connections.show ();
		leftsw.add (connections);

		connections.update_items ();

		var btn_box = new Gtk.ButtonBox (Gtk.Orientation.HORIZONTAL);
		btn_box.show ();
		btn_box.set_layout (Gtk.ButtonBoxStyle.EDGE);

		btn_connect = new Gtk.Button.with_label ("Connect");
		btn_connect.show ();
		btn_test = new Gtk.Button.with_label ("Test");
		btn_test.show ();
		btn_save = new Gtk.Button.with_label ("Save");
		btn_save.show ();

		btn_box.add (btn_connect);
		btn_box.add (btn_test);
		btn_box.add (btn_save);

		var adapter_stack = new Gtk.StackSwitcher ();
		adapter_stack.show ();
		adapter_stack.set_vexpand (true);
		adapter_stack.set_hexpand (true);

		stack = new Gtk.Stack ();
		stack.show();
		stack.set_homogeneous (true);
		stack.set_vexpand (true);
		stack.set_hexpand (true);
		adapter_stack.set_stack (stack);

		mysql = new MysqlForm();
		mysql.show ();
		sqlite = new SQLiteForm();
		sqlite.show ();

		stack.add_titled (mysql, "mysql", "Mysql");
		stack.set_visible_child_name ("mysql");
		stack.add_titled (sqlite, "sqlite", "SQLite");

		var formbox = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		formbox.show ();
		formbox.set_size_request (300, 200);
		formbox.set_valign (Gtk.Align.CENTER);
		formbox.set_halign (Gtk.Align.CENTER);

		formbox.pack_start (adapter_stack, false, true, 0);
		formbox.pack_start (stack, true, false, 0);
		formbox.pack_end (btn_box, false, true, 0);


		pack1 (leftsw, true, true);
		pack2 (formbox, true, false);

		connections.row_selected.connect (on_row_selected);

		btn_save.clicked.connect (on_save);

		enable_buttons (false);

		mysql.changed.connect ((conn) => {
			var ok = SQL.MysqlDB.validate_connection (conn);
			enable_buttons (ok);
		});

		btn_connect.clicked.connect (() => {
			var c = get_connection ();
			dbconnect (c);
		});

		var engine = new SQL.Engine ();
		btn_test.clicked.connect (() => {
			var c = get_connection ();
			if (c == null) {
				return;
			}

			try {
				engine.connect (c);
			} catch (SQL.ErrorConnection e) {
				//info_message (e.message, Gtk.MessageType.ERROR);
			}
		});

		connections.row_activated.connect (() => {
			var c = get_connection ();
			dbconnect (c);
		});
	}

	public SQL.ConnectionInfo? get_connection () {
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

		var conn = Config.connections.nth_data (index);

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

		var ok = SQL.MysqlDB.validate_connection (conn);
		enable_buttons (ok);
	}

	private void on_save () {
		var c = get_connection ();
		if ( c == null ) {
			return;
		}

		Config.save_connection (ref c);
		connections.update_items (c.id);
	}

	private void enable_buttons (bool ok) {
		btn_save.set_sensitive (ok);
		btn_test.set_sensitive (ok);
		btn_connect.set_sensitive (ok);
	}
}

public class Benchwell.Views.DBConnectionList : Gtk.ListBox {
	// connection menu
	private Gtk.Menu menu;
	private Gtk.MenuItem menu_new;
	private Gtk.MenuItem menu_connect;
	private Gtk.MenuItem menu_test;
	private Gtk.MenuItem menu_del;

	//public Regex? filter { get; set; }

	public DBConnectionList () {
		Object();

		set_property ("activate-on-single-click", false);

		button_press_event.connect ( (list, event) => {
			if (event.button == Gdk.BUTTON_SECONDARY) {
				grab_focus ();
				select_row (get_row_at_y ((int)event.y));
			}

			return false;
		});

		// connection menu
		menu = new Gtk.Menu ();
		menu_new = new Benchwell.MenuItem ("New", "connection");
		menu_new.show ();
		menu_connect = new Benchwell.MenuItem ("Connect", "next");
		menu_connect.show ();
		menu_test = new Benchwell.MenuItem ("Test", "refresh");
		menu_test.show ();
		menu_del = new Benchwell.MenuItem ("Delete", "close");
		menu_del.show ();

		menu.add (menu_new);
		menu.add (menu_connect);
		menu.add (menu_test);
		menu.add (menu_del);

		button_press_event.connect((list, event) => {
			if ( event.button != Gdk.BUTTON_SECONDARY){
				return false;
			}

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

		Config.connections.foreach ( (item) => {
			var row = build_row (item);
			add (row);
			if (item.id == selected_id) {
				select_row (row);
			}
		});
	}

	private Gtk.ListBoxRow build_row (SQL.ConnectionInfo item) {
		var row = new Gtk.ListBoxRow ();
		row.show ();

		var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		box.show ();

		var label = new Gtk.Label (item.to_string());
		label.set_halign (Gtk.Align.START);
		label.show ();

		var image = new Benchwell.Image ("connection", "orange", 16);
		image.show ();

		box.pack_start (image, false, false, 5);
		box.pack_start (label, false, false, 0);

		row.add (box);

		return row;
	}

	private void on_new () {
		var c = new Benchwell.SQL.ConnectionInfo();
		c.name = "New connection";
		c.adapter = "mysql";
		c.ttype = "tcp";
		c.port = 3306;

		Config.save_connection (ref c);
		update_items (c.id);
	}

	private void on_delete () {
		var row = get_selected_row ();
		var conn = Benchwell.Config.connections.nth_data (row.get_index ());
		Benchwell.Config.delete_connection (conn);
		update_items ();
	}
}
