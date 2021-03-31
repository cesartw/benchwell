namespace Benchwell {
	public static Benchwell._Config Config;
}

public errordomain Benchwell.ConfigError {
	STORE
}

public class Benchwell._Config : Object {
	public Sqlite.Database db;
	public GLib.Settings settings;
	public Benchwell.Environment[] environments;
	public Benchwell.ConnectionInfo[] connections;
	public Benchwell.HttpCollection[] http_collections;
	public Secret.Schema schema;
	public Benchwell.Plugin[] plugins;
	public Json.Node filters;
	public HashTable<int64?, bool?> http_tree_state;

	private Benchwell.Environment? _environment;
	public Benchwell.Environment? environment {
		get { return _environment; }
		set { _environment = value; environment_changed (); }
	}

	public signal void environment_added (Benchwell.Environment env);
	public signal void environment_removed (Benchwell.Environment env);
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

		var allplugins = Benchwell.JSPlugin.load ();
		foreach (Benchwell.Plugin p in Benchwell.BuiltinPlugin.load ()) {
			allplugins += p;
		}
		plugins = allplugins;

		load_environments ();
		load_connections ();
		load_http_collections ();
		load_filters ();
		load_http_tree_state ();
	}

	public void show_alert (Gtk.Widget? w, string message, Gtk.MessageType type = Gtk.MessageType.ERROR, bool autohide = false, int timeout = 0) {
		if (w == null) {
		// TODO: show dialog
			stderr.printf (message);
			return;
		}
		var aw = w.get_toplevel () as Gtk.Window as Gtk.ApplicationWindow as Benchwell.ApplicationWindow;
		if (aw == null) {
			stderr.printf (message);
			return;
		}
		aw.show_alert (message, type, autohide, timeout);
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

	public Benchwell.Environment add_environment (Benchwell.Environment? env = null) throws ConfigError {
		var e = env;
		if (env == null) {
			e = new Benchwell.Environment ();
			e.name = @"New environment #$(environments.length)";
		}

		var tmp = environments;
		tmp += e;
		environments = tmp;
		environment_added (e);

		return e;
	}

	public void remove_environment (Benchwell.Environment env) {
		Benchwell.Environment[] tmp = {};

		for (var i = 0; i < environments.length; i++) {
			if (environments[i].id == env.id)	 {
				continue;
			}
			tmp += environments[i];
		}

		environment_removed (env);
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
			throw new ConfigError.STORE(errmsg);
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
			throw new ConfigError.STORE(errmsg);
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
			throw new ConfigError.STORE(errmsg);
		}
	}

	public void load_http_items (Benchwell.HttpCollection collection) throws ConfigError {
		string errmsg;
		Benchwell.HttpItem[] items = {};
		var query = """SELECT id, name, is_folder, sort, http_collections_id, method, parent_id
						FROM http_items
						WHERE http_collections_id = %lld
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
				item.parent_id = int64.parse (values[6]);
			});

			items += item;

			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		Benchwell.HttpItem[] mapped_items = {};

		foreach (var item in items) {
			if (item.parent_id == 0) {
				// ROOT ITEM
				mapped_items += item;
			}

			Benchwell.HttpItem[] children = {};
			foreach (var child in items) {
				if (child.parent_id == item.id)
					children += child;
			}

			item.items = children;
		}

		collection.items = mapped_items;
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
			throw new ConfigError.STORE(errmsg);
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
				item.kvtype = (Benchwell.KeyValueTypes)int64.parse (values[5]);
			});

			variables  += item;
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
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

	public string[]? get_table_filters (ConnectionInfo info, string table_name) {
		string[] result = null;

		if (filters == null)
			return null;

		if (!filters.get_object ().has_member (info.id.to_string ()))
			return null;

		var conn_node = filters.get_object ().get_object_member (info.id.to_string ());
		if (!conn_node.has_member (table_name))
			return null;

		var table_node = conn_node.get_array_member (table_name);
		table_node.foreach_element ((array, index, node) => {
			result += node.get_string ();
		});

		return result;
	}

	public void save_filters (ConnectionInfo info, string table_name, Benchwell.CondStmt[] stmts) {
		var conn_node = filters.get_object ().get_object_member (info.id.to_string ());
		if (conn_node == null) {
			filters.get_object ().set_object_member (info.id.to_string (), new Json.Object ());
			conn_node = filters.get_object ().get_object_member (info.id.to_string ());
		}

		// rebuild all conditions
		conn_node.remove_member (table_name);
		conn_node.set_array_member (table_name, new Json.Array ());
		var table_node = conn_node.get_array_member (table_name);

		foreach (Benchwell.CondStmt stmt in stmts) {
			var name_node = new Json.Node (Json.NodeType.VALUE);
			name_node.set_string (stmt.field.name);
			table_node.add_element (name_node);

			var op_node = new Json.Node (Json.NodeType.VALUE);
			op_node.set_string (stmt.op.to_string ());
			table_node.add_element (op_node);

			var val_node = new Json.Node (Json.NodeType.VALUE);
			val_node.set_string (stmt.val);
			table_node.add_element (val_node);
		}

		Config.settings.set_string ("db-filters", Json.to_string (filters, false));
	}

	public void load_http_tree_state () {
		var tree_state = Json.from_string (settings.get_string ("http-tree-state"));
		http_tree_state = new HashTable<int64?, bool?> (int64_hash, int64_equal);

		tree_state.get_object ().get_members ().foreach ((key) => {
			http_tree_state.insert (int64.parse (key), tree_state.get_object ().get_boolean_member (key));
		});
	}

	public void save_http_tree_state () {
		var main_node = new Json.Node (Json.NodeType.OBJECT);
		main_node.set_object (new Json.Object ());
		http_tree_state.foreach ((key, val) => {
			main_node.get_object ().set_boolean_member (key.to_string (), val);
		});

		settings.set_string ("http-tree-state", Json.to_string (main_node, false));
	}

	public void load_filters () {
		try {
			filters = Json.from_string (settings.get_string ("db-filters"));
		} catch (GLib.Error err) {
			stderr.printf ("Loading saved filters: %s", err.message);
		}
	}
}
