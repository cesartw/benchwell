public errordomain Benchwell.ConfigError {
	GET_CONNECTIONS,
	GET_ENVIRONMENTS,
	SAVE_CONNECTION,
	DELETE_CONNECTION,
	ENVIRONMENTS
}

public class Benchwell.Environment : Object {
	public int64  id;
	public string name;
	public Benchwell.EnvVar[] variables;

	public Regex regex;

	public Environment () {
		regex = /({{\s*([a-zA-Z0-9]+)\s*}})/;
	}

	public string interpolate (string s) {
		MatchInfo info;
		string result = s;

		for (regex.match (s, 0, out info); info.matches () ; info.next ()) {
			for (var i = info.get_match_count () - 1; i > 0; i-=2) {
				var var_name = info.fetch (i);
				var to_replace = info.fetch (i-1);

				foreach (var envar in variables) {
					if (envar._key == var_name) {
						result = result.replace (to_replace, envar._val);
					}
				}
			}
		}

		return result;
	}
}

public class Benchwell.EnvVar : Object, Benchwell.KeyValueI {
	public int64  id;
	public string _key;
	public string _val;
	public bool   _enabled;
	public string type; // header | param
	public int    sort;
	public int64  environment_id;

	public string key () {
		return _key;
	}

	public string val() {
		return _val;
	}

	public bool enabled() {
		return _enabled;
	}

	public void set_key(string n) {
		_key = n;
	}

	public void set_val(string v) {
		_val = v;
	}

	public void set_enabled(bool e) {
		_enabled = e;
	}
}

public interface Benchwell.KeyValueI {
	public abstract string key();
	public abstract string val();
	public abstract bool enabled();
	public abstract void set_key(string n);
	public abstract void set_val(string v);
	public abstract void set_enabled(bool e);
}

public class Benchwell.HttpCollection : Object {
	public int64      id;
	public string     name;
	public int        count;
	public HttpItem[] items;
}

public class Benchwell.HttpItem : Object {
	public int64  id;
	public int64  parent_id;
	public string name;
	public string description;
	public bool   is_folder;
	public int64  http_collection_id;
	public int    sort;
	public int64  count;

	public string             method;
	public string             url;
	public string             body;
	public string             mime;
	public Benchwell.HttpKv[] headers;
	public Benchwell.HttpKv[] query_params;

	public Benchwell.HttpItem[]?  items;

	internal bool                  loaded;
}

public class Benchwell.HttpKv : Object, Benchwell.KeyValueI {
	public int64  id;
	public string _key;
	public string _val;
	public bool   _enabled;
	public string type; // header | param
	public int    sort;
	public int64  http_item_id;

	public string key () {
		return _key;
	}

	public string val() {
		return _val;
	}

	public bool enabled() {
		return _enabled;
	}

	public void set_key(string n) {
		_key = n;
	}

	public void set_val(string v) {
		_val = v;
	}

	public void set_enabled(bool e) {
		_enabled = e;
	}
}

public class Benchwell.Config : Object {
	private static GLib.Settings settings;
	private static Sqlite.Database db;
	public static List<Benchwell.Environment> environments;
	public static List<Benchwell.Backend.Sql.ConnectionInfo> connections;
	public static List<Benchwell.Backend.Sql.Query> queries;
	public static List<Benchwell.HttpCollection> http_collections;
	public static Secret.Schema schema;

	public Config () {
		settings = new GLib.Settings ("io.benchwell");
		string dbpath = GLib.Environment.get_user_config_dir () + "/benchwell/config.db";
		int ec = Sqlite.Database.open_v2 (dbpath, out db, Sqlite.OPEN_READWRITE);
		if (ec != Sqlite.OK) {
			stderr.printf ("could not open config database: %d: %s\n", db.errcode (), db.errmsg ());
		}

		schema = new Secret.Schema (Constants.PROJECT_NAME, Secret.SchemaFlags.NONE,
                                 "id", Secret.SchemaAttributeType.INTEGER,
                                 "schema", Secret.SchemaAttributeType.STRING);

		stdout.printf ("Using config db: %s\n", dbpath);

		load_environments ();
		load_connections ();
		load_http_collections ();
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

	public static void save_query (ref Benchwell.Backend.Sql.Query query) throws ConfigError {
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

	public static void delete_query (Benchwell.Backend.Sql.Query query) throws ConfigError {
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

	public static void save_connection (ref Benchwell.Backend.Sql.ConnectionInfo conn) throws ConfigError {
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

	public static void delete_connection (Benchwell.Backend.Sql.ConnectionInfo c) throws ConfigError {
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
			Benchwell.Backend.Sql.Query[] qq = {};
			queries.foreach ((query) => {
				if (query.connection_id == connection.id) {
					qq += query;
				}
			});
			connection.queries = qq;
		});
	}

	private static int connections_cb(int n_columns, string[] values, string[] column_names){
		var info = new Benchwell.Backend.Sql.ConnectionInfo ();
		info.id = int.parse (values[0]);
		info.name = values[1];
		info.adapter = values[2];
		info.ttype = values[3];
		info.database = values[4];
		info.host = values[5];
		info.user = values[7];
		info.port = int.parse(values[8]);
		info.encrypted = values[9] == "1";
		info.socket = values[10];
		info.file = values[11];

		connections.append (info);
		return 0;
	}

	private static int queries_cb(int n_columns, string[] values, string[] column_names){
		var query = new Benchwell.Backend.Sql.Query ();
		query.id = int.parse (values[0]);
		query.name = values[1];
		query.query = values[2];
		query.connection_id = int64.parse (values[3]);

		queries.append (query);
		return 0;
	}

	public static async void encrypt (Benchwell.Backend.Sql.ConnectionInfo info) {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
        attributes["id"] = info.id.to_string ();
        attributes["schema"] = Constants.PROJECT_NAME;

        var key_name = Constants.PROJECT_NAME + "." + info.id.to_string ();

        bool result = yield Secret.password_storev (schema, attributes, Secret.COLLECTION_DEFAULT, key_name, info.password, null);

        if (! result) {
            debug ("Unable to store password for \"%s\" in libsecret keyring", key_name);
        }
	}

	public static async string? decrypt (Benchwell.Backend.Sql.ConnectionInfo info) throws Error {
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

	private static void load_http_collections () throws ConfigError {
		string errmsg;
		var query = """SELECT id, name, count
					 FROM http_collections
					 ORDER BY name
					""";
		var ec = db.exec (query, (n_columns, values, column_names) => {
			var collection = new Benchwell.HttpCollection ();
			collection.id = int.parse (values[0]);
			collection.name = values[1];
			collection.count = int.parse (values[2]);

			http_collections.append (collection);
			return 0;
		}, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}
	}

	public static void load_root_items (Benchwell.HttpCollection collection) {
		string errmsg;
		Benchwell.HttpItem[] items = {};
		var query = """SELECT id, name, is_folder, sort, http_collections_id, method
						FROM http_items
						WHERE http_collections_id = %lld AND (parent_id IS NULL OR parent_id = 0)
						ORDER BY sort ASC
						""".printf (collection.id);

		var ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.HttpItem ();

			item.id = int64.parse (values[0]);
			item.name = values[1];
			item.is_folder = values[2] == "1";
			item.sort = int.parse (values[3]);
			item.http_collection_id = int64.parse (values[4]);
			item.method = values[5];

			items += item;

			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		collection.items = items;
	}

	public static void load_full_item (Benchwell.HttpItem item) {
		if (item.loaded) {
			return;
		}

		string errmsg = "";

		Benchwell.HttpItem[] items = {};
		// folder
		if (item.is_folder) {
			var query = """SELECT id, name, parent_id, is_folder, sort,
										http_collections_id, method
							FROM http_items
							WHERE http_collections_id = %lld AND parent_id = %lld
							ORDER BY sort ASC
							""".printf (item.http_collection_id, item.id);
			var ec = db.exec (query, (n_columns, values, column_names) => {
				var subitem = new Benchwell.HttpItem ();

				subitem.id = int64.parse (values[0]);
				subitem.name = values[1];
				subitem.parent_id = int64.parse (values[2]);
				subitem.is_folder = values[3] == "1";
				subitem.sort = int.parse (values[4]);
				subitem.http_collection_id = int64.parse (values[5]);
				subitem.method = values[6];

				items += subitem;
				return 0;
			}, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.GET_CONNECTIONS(errmsg);
			}

			item.items = items;
			return;
		}

		// request
		var query = """SELECT ifnull(method,""), ifnull(url,""), ifnull(body, ""), ifnull(mime,"")
				FROM http_items
				WHERE id = %lld""".printf (item.id);
		var ec = db.exec (query, (n_columns, values, column_names) => {
			item.method = values[0];
			item.url = values[1];
			item.body = values[2];
			item.mime = values[3];
			return 0;
		}, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		Benchwell.HttpKv[] kvs = {};
		query = """SELECT id, ifnull(key, ""), ifnull(value, ""), type, sort, enabled
			FROM http_kvs
			WHERE http_items_id = %lld
			ORDER BY sort ASC""".printf (item.id);
		ec = db.exec (query, (n_columns, values, column_names) => {
			var kv = new Benchwell.HttpKv ();
			kv.id = int64.parse (values[0]);
			kv._key = values[1];
			kv._val = values[2];
			kv.type = values[3];
			kv.sort = int.parse (values[4]);
			kv._enabled = values[5] == "1";
			kvs += kv;
			return 0;
		}, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		Benchwell.HttpKv[] headers = {};
		Benchwell.HttpKv[] query_params = {};
		foreach (var kv in kvs) {
			if (kv.type == "header") {
				headers += kv;
				continue;
			}
			query_params += kv;
		}
		item.headers = headers;
		item.query_params = query_params;
		item.loaded = true;
	}

	public static void load_environments () {
		string errmsg;
		var query = """SELECT *
						FROM environments
						""";

		var ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.Environment ();

			item.id = int64.parse (values[0]);
			item.name = values[1];

			environments.append (item);
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_ENVIRONMENTS(errmsg);
		}

		query = """SELECT *
						FROM environment_variables
						""";

		Benchwell.EnvVar[] variables = {};
		ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.EnvVar ();

			item.id = int64.parse (values[0]);
			item._key = values[1];
			item._val = values[2];
			item._enabled = values[3] == "1";
			item.environment_id = int64.parse (values[4]);

			variables  += item;
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_ENVIRONMENTS(errmsg);
		}

		environments.foreach ((env) => {
			Benchwell.EnvVar[] envvars = {};
			foreach (var v in variables) {
				if (env.id == v.environment_id) {
					envvars += v;
				}
			}
			env.variables = envvars;
		});
	}
}
