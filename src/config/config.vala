namespace Benchwell {
	public static Benchwell._Config Config;
}

public errordomain Benchwell.ConfigError {
	GET_CONNECTIONS,
	GET_ENVIRONMENTS,
	SAVE_CONNECTION,
	SAVE_ENVVAR,
	DELETE_CONNECTION,
	ENVIRONMENTS
}

public class Benchwell._Config : Object {
	public Sqlite.Database db;
	public GLib.Settings settings;
	public Benchwell.Environment[] environments { set; get; }
	public List<Benchwell.Backend.Sql.ConnectionInfo> connections;
	public List<Benchwell.Backend.Sql.Query> queries;
	public Benchwell.HttpCollection[] http_collections;
	public Secret.Schema schema;
	public Benchwell.Http.Plugin plugins;

	private Benchwell.Environment? _environment;
	public Benchwell.Environment? environment {
		get { return _environment; }
		set { _environment = value; environment_changed (); }
	}

	public signal void environment_added (Benchwell.Environment env);
	public signal void environment_changed ();
	public signal void http_collection_added(HttpCollection collection);

	public _Config () throws Benchwell.ConfigError {
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
		plugins = new Benchwell.Http.Plugin ();

		load_environments ();
		load_connections ();
		load_http_collections ();
	}

	public Gtk.PositionType tab_position () {
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

	public Benchwell.Environment add_environment () throws ConfigError {
		var env = new Benchwell.Environment ();
		env.name = @"New environment #$(environments.length)";
		var tmp = environments;
		tmp += env;
		environments = tmp;
		environment_added (env);

		return env;
	}

	public void remove_environment (Benchwell.Environment env) {
		Benchwell.Environment[] tmp = {};

		for (var i = 0; i < environments.length; i++) {
			if (environments[i].id == env.id)	 {
				continue;
			}
			tmp += env;
		}

		environments = tmp;
	}

	public Benchwell.HttpCollection add_http_collection () throws ConfigError {
		var collection = new Benchwell.HttpCollection ();
		collection.name = @"New collection #$(http_collections.length)";
		collection.save ();

		var tmp = http_collections;
		tmp += collection;
		http_collections = tmp;
		http_collection_added (collection);

		return collection;
	}

	public void remove_http_collection (Benchwell.HttpCollection collection) {
		Benchwell.HttpCollection[] tmp = {};

		for (var i = 0; i < http_collections.length; i++) {
			if (http_collections[i].id == collection.id)	 {
				continue;
			}
			tmp += collection;
		}

		http_collections = tmp;
	}

	public void save_query (ref Benchwell.Backend.Sql.Query query) throws ConfigError {
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

	public void delete_query (Benchwell.Backend.Sql.Query query) throws ConfigError {
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

	public void save_connection (ref Benchwell.Backend.Sql.ConnectionInfo conn) throws ConfigError {
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

	public void delete_connection (Benchwell.Backend.Sql.ConnectionInfo c) throws ConfigError {
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

	private void load_connections () throws ConfigError {
		string errmsg;
		var ec = db.exec ("SELECT * FROM db_connections",
						  connections_cb,
						  out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}
		load_queries ();
	}

	private void load_queries () throws ConfigError {
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

	private int connections_cb(int n_columns, string[] values, string[] column_names){
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

	private int queries_cb(int n_columns, string[] values, string[] column_names){
		var query = new Benchwell.Backend.Sql.Query ();
		query.id = int.parse (values[0]);
		query.name = values[1];
		query.query = values[2];
		query.connection_id = int64.parse (values[3]);

		queries.append (query);
		return 0;
	}

	public async void encrypt (Benchwell.Backend.Sql.ConnectionInfo info) {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
		attributes["id"] = info.id.to_string ();
		attributes["schema"] = Constants.PROJECT_NAME;

		var key_name = Constants.PROJECT_NAME + "." + info.id.to_string ();

		bool result = yield Secret.password_storev (schema, attributes, Secret.COLLECTION_DEFAULT, key_name, info.password, null);

		if (! result) {
			debug ("Unable to store password for \"%s\" in libsecret keyring", key_name);
		}
	}

	public async string? decrypt (Benchwell.Backend.Sql.ConnectionInfo info) throws Error {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
		attributes["id"] = info.id.to_string ();
		attributes["schema"] = Constants.PROJECT_NAME;

		var key_name = Constants.PROJECT_NAME + "." + info.id.to_string ();

		Secret.Service srv = yield Secret.Service.get (Secret.ServiceFlags.OPEN_SESSION, null);
		bool ok = yield srv.ensure_session (null);
		if (!ok) {
			return null;
		}

		string? password = yield Secret.password_lookupv (schema, attributes, null);

		if (password == null) {
			print (@"Unable to fetch password in libsecret keyring for $(key_name)\n");
		}
		if (password == "") {
			print (@"no password $(key_name)\n");
		}

		return password;
    }

	// hack because password_lookpv doesn't trigger the unlock keychain popup
	public async void ping_dbus () throws Error {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
		attributes["id"] = "0";
		attributes["schema"] = Constants.PROJECT_NAME;

		var key_name = Constants.PROJECT_NAME + ".0";

		bool result = yield Secret.password_storev (schema, attributes, Secret.COLLECTION_DEFAULT, key_name, "none", null);

		if (! result) {
			debug ("Unable to store password for \"%s\" in libsecret keyring", key_name);
		}
    }

	private void load_http_collections () throws ConfigError {
		string errmsg;
		var query = """SELECT id, name, count
					 FROM http_collections
					 ORDER BY name
					""";
		var ec = db.exec (query, (n_columns, values, column_names) => {
			var collection = new Benchwell.HttpCollection ();
			collection.touch_without_save ( () => {
				collection.id = int.parse (values[0]);
				collection.name = values[1];
				collection.count = int.parse (values[2]);
			});

			HttpCollection[] tmp = http_collections;
			tmp += collection;
			http_collections = tmp;

			return 0;
		}, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}
	}

	public void load_root_items (Benchwell.HttpCollection collection) throws ConfigError {
		string errmsg;
		Benchwell.HttpItem[] items = {};
		var query = """SELECT id, name, is_folder, sort, http_collections_id, method
						FROM http_items
						WHERE http_collections_id = %lld AND (parent_id IS NULL OR parent_id = 0)
						ORDER BY sort ASC
						""".printf (collection.id);

		var ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.HttpItem ();

			item.touch_without_save ( () => {
				item.id = int64.parse (values[0]);
				item.name = values[1];
				item.is_folder = values[2] == "1";
				item.sort = int.parse (values[3]);
				item.http_collection_id = int64.parse (values[4]);
				item.method = values[5];
			});

			items += item;

			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		collection.items = items;
	}

	public void load_environments () throws Benchwell.ConfigError {
		string errmsg;
		var query = """SELECT *
						FROM environments
						""";

		var ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.Environment ();

			item.touch_without_save (() => {
				item.id = int64.parse (values[0]);
				item.name = values[1];
			});

			var tmp = environments;
			tmp += item;
			environments = tmp;
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_ENVIRONMENTS(errmsg);
		}

		query = """SELECT * FROM environment_variables""";

		Benchwell.EnvVar[] variables = {};
		ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.EnvVar ();

			item.touch_without_save (() => {
				item.id = int64.parse (values[0]);
				item.key = values[1];
				item.val = values[2];
				item.enabled = values[3] == "1";
				item.environment_id = int64.parse (values[4]);
			});

			variables  += item;
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_ENVIRONMENTS(errmsg);
		}

		for (var i = 0; i < environments.length; i++) {
			var env = environments[i];
			Benchwell.EnvVar[] envvars = {};
			foreach (var v in variables) {
				if (env.id == v.environment_id) {
					envvars += v;
				}
			}
			env.variables = envvars;
		}
	}

	public void save_httpitem (Benchwell.HttpItem item) throws Error {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (item.id > 0) {
			 prepared_query_str = """
				UPDATE http_items
					SET name = $NAME, description = $DESCRIPTION,
						parent_id = $PARENT_ID, is_folder = $IS_FOLDER,
						sort = $SORT, http_collections_id = $HTTP_COLLECTION_ID,
						method = $METHOD, url = $URL,
						body = $BODY, mime = $MIME
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO http_items(name, description, parent_id, is_folder,
					sort, http_collections_id, method, url, body, mime)
				VALUES($NAME, $DESCRIPTION, $IS_PARENT, $IS_FOLDER, $SORT,
					$HTTP_COLLECTION_ID, $METHOD, $URL, $BODY, $MIME)
			""";
		}

		var ec = db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", db.errcode (), db.errmsg ());
			return;
		}

		int param_position;
		if (item.id > 0) {
			param_position = stmt.bind_parameter_index ("$ID");
			stmt.bind_int64 (param_position, item.id);
		}

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, item.name);

		param_position = stmt.bind_parameter_index ("$DESCRIPTION");
		stmt.bind_text (param_position, item.description);

		param_position = stmt.bind_parameter_index ("$PARENT_ID");
		stmt.bind_int64 (param_position, item.parent_id);

		param_position = stmt.bind_parameter_index ("$IS_FOLDER");
		stmt.bind_int (param_position, item.is_folder ? 1 : 0);

		param_position = stmt.bind_parameter_index ("$SORT");
		stmt.bind_int (param_position, item.sort);

		param_position = stmt.bind_parameter_index ("$HTTP_COLLECTION_ID");
		stmt.bind_int64 (param_position, item.http_collection_id);

		param_position = stmt.bind_parameter_index ("$METHOD");
		stmt.bind_text (param_position, item.method);

		param_position = stmt.bind_parameter_index ("$URL");
		stmt.bind_text (param_position, item.url);

		param_position = stmt.bind_parameter_index ("$BODY");
		stmt.bind_text (param_position, item.body);

		param_position = stmt.bind_parameter_index ("$MIME");
		stmt.bind_text (param_position, item.mime);

		string errmsg = "";
		ec = db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}

		if (item.id == 0) {
			item.id = db.last_insert_rowid ();
			//Config.environments.append (item);
		}
	}
}
