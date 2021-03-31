public class Benchwell.Database.MysqlForm : Gtk.Box {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.Database.MysqlTCPForm tcp_form;
	public Benchwell.Database.MysqlSocketForm socket_form;
	public Benchwell.ConnectionInfo? connection;

	public signal void changed(Benchwell.ConnectionInfo c);

	public MysqlForm(Benchwell.ApplicationWindow w) {
		Object(
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5,
			window: w
		);
		set_size_request (300, 200);
		set_valign (Gtk.Align.CENTER);
		set_halign (Gtk.Align.CENTER);

		var notebook = new Gtk.Notebook ();
		notebook.show ();
		notebook.set_can_focus (true);

		tcp_form = new Benchwell.Database.MysqlTCPForm ();
		tcp_form.show ();
		notebook.append_page (tcp_form, tcp_form.tab_label);

		socket_form = new Benchwell.Database.MysqlSocketForm ();
		socket_form.show ();
		notebook.append_page (socket_form, socket_form.tab_label);

		tcp_form.changed.connect ((entry) =>{
			try {
				if ( connection == null ) {
					connection = Config.add_connection ();
				}
				connection.name = tcp_form.name_entry.get_text ();
				connection.host = tcp_form.host_entry.get_text ();
				connection.port = int.parse (tcp_form.port_entry.get_text ());
				connection.user = tcp_form.user_entry.get_text ();
				connection.password = tcp_form.password_entry.get_text ();
				connection.database = tcp_form.database_entry.get_text ();

				changed (connection);
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
			}
		});

		socket_form.changed.connect ((entry) =>{
			try {
				if ( connection == null ) {
					connection = Config.add_connection ();
				}
				connection.name = socket_form.name_entry.get_text ();
				connection.socket = socket_form.socket_entry.get_text ();
				connection.user = socket_form.user_entry.get_text ();
				connection.password = socket_form.password_entry.get_text ();
				connection.database = socket_form.database_entry.get_text ();

				changed (connection);
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
			}
		});

		add (notebook);
	}

	public void set_connection(ref Benchwell.ConnectionInfo conn) {
		connection = conn;

		switch (conn.ttype) {
			case "tcp":
				tcp_form.set_connection (conn);
				break;
			case "socket":
				socket_form.set_connection (conn);
				break;
		}
	}

	public Benchwell.ConnectionInfo? get_connection() {
		return connection;
	}
}

public class Benchwell.Database.MysqlTCPForm : Gtk.Grid {
	public Gtk.Label tab_label;

	public Gtk.Entry name_entry;
	public Gtk.Entry host_entry;
	public Gtk.Entry port_entry;
	public Gtk.Entry user_entry;
	public SecretEntry password_entry;
	public Gtk.Entry database_entry;

	public Gtk.Label name_lbl;
	public Gtk.Label host_lbl;
	public Gtk.Label port_lbl;
	public Gtk.Label user_lbl;
	public Gtk.Label password_lbl;
	public Gtk.Label database_lbl;
	private bool _setting_connection = false;

	public signal void changed(Gtk.Entry entry);

	public MysqlTCPForm(){
		Object();
		set_name ("form");
		set_column_homogeneous (true);
		set_row_spacing (5);

		tab_label = new Gtk.Label ("TCP/IP");
		tab_label.show ();

		name_entry = new Gtk.Entry();
		name_entry.show ();

		host_entry = new Gtk.Entry();
		host_entry.show ();

		port_entry = new Gtk.Entry();
		port_entry.set_property("input-purpose", Gtk.InputPurpose.NUMBER);
		port_entry.show ();

		user_entry = new Gtk.Entry();
		user_entry.show ();

		password_entry = new SecretEntry(true);
		password_entry.show ();

		database_entry = new Gtk.Entry();
		database_entry.show ();

		name_lbl = new Gtk.Label("Name");
		name_lbl.set_halign (Gtk.Align.START);
		name_lbl.show ();

		host_lbl = new Gtk.Label("Host");
		host_lbl.set_halign (Gtk.Align.START);
		host_lbl.show ();

		port_lbl = new Gtk.Label("Port");
		port_lbl.set_halign (Gtk.Align.START);
		port_lbl.show ();

		user_lbl = new Gtk.Label("User");
		user_lbl.set_halign (Gtk.Align.START);
		user_lbl.show ();

		password_lbl = new Gtk.Label("Password");
		password_lbl.set_halign (Gtk.Align.START);
		password_lbl.show ();

		database_lbl = new Gtk.Label("Database");
		database_lbl.set_halign (Gtk.Align.START);
		database_lbl.show ();

		attach(name_lbl, 0, 0, 1, 1);
		attach(name_entry, 1, 0, 2, 1);

		attach(host_lbl, 0, 1, 1, 1);
		attach(host_entry, 1, 1, 2, 1);

		attach(port_lbl, 0, 2, 1, 1);
		attach(port_entry, 1, 2, 2, 1);

		attach(user_lbl, 0, 3, 1, 1);
		attach(user_entry, 1, 3, 2, 1);

		attach(password_lbl, 0, 4, 1, 1);
		attach(password_entry, 1, 4, 2, 1);

		attach(database_lbl, 0, 5, 1, 1);
		attach(database_entry, 1, 5, 2, 1);

		name_entry.changed.connect ((entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (name_entry);
		});
		host_entry.changed.connect ((entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (host_entry);
		});
		port_entry.changed.connect ((entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (port_entry);
		});
		user_entry.changed.connect ((entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (user_entry);
		});
		password_entry.changed.connect ((entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (password_entry);
		});
		database_entry.changed.connect ((entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (database_entry);
		});
	}

	public void clear() {
		name_entry.set_text ("");
		host_entry.set_text ("");
		port_entry.set_text ("");
		user_entry.set_text ("");
		password_entry.set_text ("");
		database_entry.set_text ("");
	}

	public void set_connection (Benchwell.ConnectionInfo conn) {
		_setting_connection = true;
		name_entry.set_text (conn.name);
		host_entry.set_text (conn.host);
		port_entry.set_text (conn.port.to_string ());
		user_entry.set_text (conn.user);
		//password_entry.set_text (conn.password);
		database_entry.set_text (conn.database);
		_setting_connection = false;
		password_entry.open = false;
	}
}

public class Benchwell.Database.MysqlSocketForm : Gtk.Grid {
	public Gtk.Label tab_label;

	public Gtk.Entry name_entry;
	public Gtk.Entry socket_entry;
	public Gtk.Entry user_entry;
	public SecretEntry password_entry;
	public Gtk.Entry database_entry;

	public Gtk.Label name_lbl;
	public Gtk.Label socket_lbl;
	public Gtk.Label user_lbl;
	public Gtk.Label password_lbl;
	public Gtk.Label database_lbl;
	private bool _setting_connection = false;

	public signal void changed(Gtk.Entry entry);

	public MysqlSocketForm(){
		Object();
		set_name ("form");
		set_column_homogeneous (true);
		set_row_spacing (5);

		tab_label = new Gtk.Label ("Socket");
		tab_label.show ();

		name_entry = new Gtk.Entry();
		name_entry.show ();

		socket_entry = new Gtk.Entry();
		socket_entry.show ();

		user_entry = new Gtk.Entry();
		user_entry.show ();

		password_entry = new SecretEntry(true);
		password_entry.show ();

		database_entry = new Gtk.Entry();
		database_entry.show ();


		name_lbl = new Gtk.Label("Name");
		name_lbl.set_halign (Gtk.Align.START);
		name_lbl.show ();

		socket_lbl = new Gtk.Label("Socket");
		socket_lbl.set_halign (Gtk.Align.START);
		socket_lbl.show ();

		user_lbl = new Gtk.Label("User");
		user_lbl.set_halign (Gtk.Align.START);
		user_lbl.show ();

		password_lbl = new Gtk.Label("Password");
		password_lbl.set_halign (Gtk.Align.START);
		password_lbl.show ();

		database_lbl = new Gtk.Label("Database");
		database_lbl.set_halign (Gtk.Align.START);
		database_lbl.show ();

		attach(name_lbl, 0, 1, 1, 1);
		attach(name_entry, 1, 1, 2, 1);

		attach(socket_lbl, 0, 2, 1, 1);
		attach(socket_entry, 1, 2, 2, 1);

		attach(user_lbl, 0, 3, 1, 1);
		attach(user_entry, 1, 3, 2, 1);

		attach(password_lbl, 0, 4, 1, 1);
		attach(password_entry, 1, 4, 2, 1);

		attach(database_lbl, 0, 5, 1, 1);
		attach(database_entry, 1, 5, 2, 1);


		name_entry.changed.connect ( (entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (name_entry);
		});
		socket_entry.changed.connect ( (entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (socket_entry);
		});
		user_entry.changed.connect ( (entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (user_entry);
		});
		password_entry.changed.connect ( (entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (password_entry);
		});
		database_entry.changed.connect ( (entry) => {
			if ( _setting_connection ) {
				return;
			}
			changed (database_entry);
		});
	}

	public void clear() {
		name_entry.set_text ("");
		socket_entry.set_text ("");
		user_entry.set_text ("");
		password_entry.set_text ("");
		database_entry.set_text ("");
	}

	public void set_connection (Benchwell.ConnectionInfo conn) {
		_setting_connection = true;
		name_entry.set_text (conn.name);
		socket_entry.set_text (conn.socket);
		user_entry.set_text (conn.user);
		//password_entry.set_text (conn.password);
		database_entry.set_text (conn.database);
		_setting_connection = false;
		password_entry.open = false;
	}
}

public class SecretEntry : Gtk.Entry {
	private bool _open;
	public bool open {
		get { return _open; }
		set {
			_open = value;
			if ( _open ) {
				placeholder_text = "";
			} else {
				set_text ("");
				placeholder_text = "Stored";
			}
		}
	}

	public SecretEntry (bool open = false) {
		Object(
			caps_lock_warning: true,
			placeholder_text: open ? "" : "Stored",
			input_purpose: Gtk.InputPurpose.PASSWORD,
			visibility: false
		);
		this.open = open;
	}
}
