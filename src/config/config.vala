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
	public Benchwell.ConnectionInfo[] connections;
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
	public signal void connection_added(ConnectionInfo connection);

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

	public Benchwell.ConnectionInfo add_connection () throws ConfigError {
		var connection = new Benchwell.ConnectionInfo ();
		connection.name = @"New connection #$(connections.length)";
		connection.save ();

		var tmp = connections;
		tmp += connection;
		connections = tmp;
		connection_added (connection);

		return connection;
	}

	public void remove_connection (Benchwell.ConnectionInfo connection) {
		Benchwell.ConnectionInfo[] tmp = {};

		foreach (Benchwell.ConnectionInfo c in Config.connections) {
			if (c.id == connection.id)	 {
				continue;
			}
			tmp += c;
		}

		connections = tmp;
	}

	// ======== LOADERS ========
	private void load_connections () throws ConfigError {
		string errmsg;
		var ec = db.exec ("SELECT * FROM db_connections", (n_columns, values, column_names) => {
			var info = new Benchwell.ConnectionInfo ();
			info.touch_without_save (() => {
				//info.to
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
			});

			ConnectionInfo[] tmp = connections;
			tmp += info;
			connections = tmp;
			return 0;
		}, out errmsg);


		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		Benchwell.Query[] queries = {};
		ec = db.exec ("SELECT * FROM db_queries WHERE query_type = 'fav'", (n_columns, values, column_names) => {
			var query = new Benchwell.Query ();
			query.touch_without_save (() => {
				query.id = int.parse (values[0]);
				query.name = values[1];
				query.query = values[2];
				query.query_type = values[3];
				query.connection_id = int64.parse (values[4]);
				query.created_at = new DateTime.from_unix_local (int64.parse (values[5]));
			});

			queries += query;
			return 0;
		}, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.GET_CONNECTIONS(errmsg);
		}

		foreach (Benchwell.ConnectionInfo conn in connections) {
			Benchwell.Query[] qq = {};
			foreach (Benchwell.Query query in queries) {
				if (query.connection_id == conn.id) {
					qq += query;
				}
			}
			conn.queries = qq;
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
	// =========================

	public async void encrypt (Benchwell.ConnectionInfo info) {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
		attributes["id"] = info.id.to_string ();
		attributes["schema"] = Constants.PROJECT_NAME;

		var key_name = Constants.PROJECT_NAME + "." + info.id.to_string ();

		bool result = yield Secret.password_storev (schema, attributes, Secret.COLLECTION_DEFAULT, key_name, info.password, null);

		if (! result) {
			debug ("Unable to store password for \"%s\" in libsecret keyring", key_name);
		}
	}

	public async string? decrypt (Benchwell.ConnectionInfo info) throws GLib.Error {
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
	public async void ping_dbus () throws GLib.Error {
		var attributes = new GLib.HashTable<string, string> (str_hash, str_equal);
		attributes["id"] = "0";
		attributes["schema"] = Constants.PROJECT_NAME;

		var key_name = Constants.PROJECT_NAME + ".0";

		bool result = yield Secret.password_storev (schema, attributes, Secret.COLLECTION_DEFAULT, key_name, "none", null);

		if (! result) {
			debug ("Unable to store password for \"%s\" in libsecret keyring", key_name);
		}
    }
}
