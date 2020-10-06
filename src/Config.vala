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
	public static List<SQL.Query> queries;
	public static Secret.Schema schema;

	public Config () {
		settings = new GLib.Settings ("io.benchwell");
		string dbpath = Environment.get_user_config_dir () + "/benchwell/config.db";
		int ec = Sqlite.Database.open_v2 (dbpath, out db, Sqlite.OPEN_READWRITE);
		if (ec != Sqlite.OK) {
			stderr.printf ("could not open config database: %d: %s\n", db.errcode (), db.errmsg ());
		}

		schema = new Secret.Schema (Constants.PROJECT_NAME, Secret.SchemaFlags.NONE,
                                 "id", Secret.SchemaAttributeType.INTEGER,
                                 "schema", Secret.SchemaAttributeType.STRING);

		stdout.printf ("Using config db: %s\n", dbpath);
		load_connections ();
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

	public static void save_query (ref SQL.Query query) throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (query.id > 0) {
			 prepared_query_str = """
				UPDATE db_queries
					SET name = $NAME, query = $query
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO db_queries(name, query, connections_id)
				VALUES($NAME, $QUERY, $CONNECTION_ID)
			""";
		}

		var ec = db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", db.errcode (), db.errmsg ());
			return;
		}

		int param_position;
		if (query.id > 0) {
			param_position = stmt.bind_parameter_index ("$ID");
			assert (param_position > 0);
			stmt.bind_int64 (param_position, query.id);
		} else {
			param_position = stmt.bind_parameter_index ("$CONNECTION_ID");
			stmt.bind_int64 (param_position, query.connection_id);
		}

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, query.name);

		param_position = stmt.bind_parameter_index ("$QUERY");
		stmt.bind_text (param_position, query.query);

		string errmsg = "";
		ec = db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_CONNECTION(errmsg);
		}

		if (query.id == 0) {
			query.id = db.last_insert_rowid ();
			queries.append (query);
		}
	}

	public static void delete_query (SQL.Query query) throws ConfigError {
		if (query.id == 0) {
			return;
		}

		queries.foreach ( (q) => {
			if ( query.id != q.id ){
				return;
			}

			queries.remove (q);

			string errmsg = "";
			var ec = db.exec (@"DELETE FROM db_queries WHERE ID = $(q.id)", null, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.DELETE_CONNECTION (errmsg);
			}
		});
	}

	public static void save_connection (ref SQL.ConnectionInfo conn) throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (conn.id > 0) {
			 prepared_query_str = """
				UPDATE db_connections
					SET adapter = $ADAPTER, type = $TYPE, name = $NAME, socket = $SOCKET, file = $FILE, host = $HOST, port = $PORT,
					user = $USER, database = $DATABASE, options = $OPT, encrypted = $ENC
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO db_connections(adapter, type, name, socket, file, host, port,
					user, database, options, encrypted)
				VALUES($ADAPTER, $TYPE, $NAME, $SOCKET, $FILE, $HOST, $PORT,
					$USER, $DATABASE, $OPT, $ENC)
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

		encrypt (conn);
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

	private static void load_connections () throws ConfigError {
		string errmsg;
		var ec = db.exec ("SELECT * FROM db_connections",
						  connections_cb,
						  out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}
		load_queries ();
	}

	private static void load_queries () throws ConfigError {
		string errmsg;
		var ec = db.exec ("SELECT * FROM db_queries",
						  queries_cb,
						  out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		connections.foreach ((connection) => {
			Benchwell.SQL.Query[] qq = {};
			queries.foreach ((query) => {
				if (query.connection_id == connection.id) {
					qq += query;
				}
			});
			connection.queries = qq;
		});
	}

	private static int connections_cb(int n_columns, string[] values, string[] column_names){
		var info = new Benchwell.SQL.ConnectionInfo ();
		info.id = int.parse (values[0]);
		info.name = values[1];
		info.adapter = values[2];
		info.ttype = values[3];
		info.database = values[4];
		info.host = values[5];
		info.user = values[7];
		info.port = int.parse(values[8]);
		info.encrypted = bool.parse(values[9]);
		info.socket = values[10];
		info.file = values[11];

		connections.append (info);
		return 0;
	}

	private static int queries_cb(int n_columns, string[] values, string[] column_names){
		var query = new Benchwell.SQL.Query ();
		query.id = int.parse (values[0]);
		query.name = values[1];
		query.query = values[2];
		query.connection_id = int64.parse (values[3]);

		queries.append (query);
		return 0;
	}

	public static async void encrypt (Benchwell.SQL.ConnectionInfo info) {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
        attributes["id"] = info.id.to_string ();
        attributes["schema"] = Constants.PROJECT_NAME;

        var key_name = Constants.PROJECT_NAME + "." + info.id.to_string ();

        bool result = yield Secret.password_storev (schema, attributes, Secret.COLLECTION_DEFAULT, key_name, info.password, null);

        if (! result) {
            debug ("Unable to store password for \"%s\" in libsecret keyring", key_name);
        }
	}

	public static async string? decrypt (Benchwell.SQL.ConnectionInfo info) throws Error {
        var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
        attributes["id"] = info.id.to_string ();
        attributes["schema"] = Constants.PROJECT_NAME;

        var key_name = Constants.PROJECT_NAME + "." + info.id.to_string ();

        string? password = yield Secret.password_lookupv (schema, attributes, null);

        if (password == null) {
            debug ("Unable to fetch password in libsecret keyring for %s", key_name);
        }

        return password;
    }
}

