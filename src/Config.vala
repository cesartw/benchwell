public errordomain Benchwell.ConfigError {
	GET_CONNECTIONS,
	SAVE_CONNECTION,
	DELETE_CONNECTION,
	ENVIRONMENTS
}

public class Benchwell.Config : Object {
	private static GLib.Settings settings;
	private static Sqlite.Database db;
	public static List<SQL.ConnectionInfo> connections;

	public Config () {
		settings = new GLib.Settings ("io.benchwell");
		string dbpath = Environment.get_user_config_dir () + "/benchwell/config.db";
		int ec = Sqlite.Database.open_v2 (dbpath, out db, Sqlite.OPEN_READWRITE);
		if (ec != Sqlite.OK) {
			stderr.printf ("could not open config database: %d: %s\n", db.errcode (), db.errmsg ());
		}

		stdout.printf ("Using config db: %s\n", dbpath);
		loadConnections ();
	}

	public static int window_width() {
		return settings.get_int ("window-size-w");
	}

	public static int window_height() {
		return settings.get_int ("window-size-h");
	}

	public static int window_position_x() {
		return settings.get_int ("window-pos-x");
	}

	public static int window_position_y() {
		return settings.get_int ("window-pos-y");
	}

	public static Gtk.PositionType tab_position() {
		Gtk.PositionType v;

		switch (settings.get_string ("tab-position")) {
			case "TOP":
				v = Gtk.PositionType.TOP;
				break;
			case "BOTTOM":
				v = Gtk.PositionType.BOTTOM;
				break;
			default:
				v = Gtk.PositionType.TOP;
				break;
		}

		return v;
	}

	public static void save_connection (ref SQL.ConnectionInfo conn) throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (conn.id > 0) {
			 prepared_query_str = """
				UPDATE db_connections
					SET adapter = $ADAPTER, type = $TYPE, name = $NAME, socket = $SOCKET, file = $FILE, host = $HOST, port = $PORT,
					user = $USER, password = $PASSWORD, database = $DATABASE, options = $OPT, encrypted = $ENC
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO db_connections(adapter, type, name, socket, file, host, port,
					user, password, database, options, encrypted)
				VALUES($ADAPTER, $TYPE, $NAME, $SOCKET, $FILE, $HOST, $PORT,
					$USER, $PASSWORD, $DATABASE, $OPT, $ENC)
			""";
		}

		var ec = db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", db.errcode (), db.errmsg ());
			return;
		}

		if (conn.id > 0) {
			int param_position = stmt.bind_parameter_index ("$ID");
			assert (param_position > 0);
			stmt.bind_int64 (param_position, conn.id);
		}

		int param_position = stmt.bind_parameter_index ("$ADAPTER");
		stmt.bind_text (param_position, conn.adapter);

		param_position = stmt.bind_parameter_index ("$TYPE");
		stmt.bind_text (param_position, conn.ttype);

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, conn.name);

		param_position = stmt.bind_parameter_index ("$SOCKET");
		stmt.bind_text (param_position, conn.socket);

		param_position = stmt.bind_parameter_index ("$FILE");
		stmt.bind_text (param_position, conn.file);

		param_position = stmt.bind_parameter_index ("$HOST");
		stmt.bind_text (param_position, conn.host);

		param_position = stmt.bind_parameter_index ("$PORT");
		stmt.bind_int (param_position, conn.port);

		param_position = stmt.bind_parameter_index ("$USER");
		stmt.bind_text (param_position, conn.user);

		param_position = stmt.bind_parameter_index ("$PASSWORD");
		stmt.bind_text (param_position, conn.password);

		param_position = stmt.bind_parameter_index ("$DATABASE");
		stmt.bind_text (param_position, conn.database);

		param_position = stmt.bind_parameter_index ("$ENC");
		stmt.bind_int (param_position, conn.encrypted ? 1 : 0);

		string errmsg = "";
		ec = db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_CONNECTION(errmsg);
		}

		if (conn.id == 0) {
			conn.id = db.last_insert_rowid ();
			connections.append (conn);
		}

		//loadConnections ();
	}

	public static void delete_connection (SQL.ConnectionInfo c) throws ConfigError {
		if ( c.id == 0 ) {
			return;
		}

		connections.foreach ( (conn) => {
			if ( c.id != conn.id ){
				return;
			}

			connections.remove (conn);

			string errmsg = "";
			var ec = db.exec (@"DELETE FROM db_connections WHERE ID = $(c.id)", null, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.DELETE_CONNECTION (errmsg);
			}
		});
	}

	private static void loadConnections () throws ConfigError {
		string errmsg;
		var ec = db.exec ("SELECT * FROM db_connections",
						  connections_cb,
						  out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}
	}

	private static int connections_cb(int n_columns, string[] values, string[] column_names){
		var info = new SQL.ConnectionInfo ();
		info.id = int.parse (values[0]);
		info.name = values[1];
		info.adapter = values[2];
		info.ttype = values[3];
		info.database = values[4];
		info.host = values[5];
		info.user = values[7];
		info.password = values[8];
		info.port = int.parse(values[9]);
		info.encrypted = bool.parse(values[10]);
		info.socket = values[11];
		info.file = values[12];

		connections.append (info);
		return 0;
	}
}