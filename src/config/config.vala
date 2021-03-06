namespace Benchwell {
	public static _Config Config;
	public errordomain ConfigError {
		STORE,
		INIT
	}
}


public interface Benchwell.CollectionWithSelectedItem : Object {
	public abstract int64 id { get; set; }
	public abstract string name { get; set; }
}

public class Benchwell.CollectionWithSelected<T> {
	private T[] items;

	public weak T? _selected;
	public T? selected {
		owned get { return _selected; }
		set { _selected = value; selected_changed (); }
	}

	public signal void added (owned T item);
	public signal void removed (owned T item);
	public signal void selected_changed ();

	public CollectionWithSelected (owned T[] items) {
		this.items = items;
	}

	public int length () {
		return items.length;
	}

	public void select (T item) {
		for (var i = 0; i < items.length; i++) {
			if (item != items[i]) {
				continue;
			}

			selected = items[i];
			return;
		}
	}

	public unowned T at (int i)
		requires (i < items.length)
	{
		return items[i];
	}

	public T add (owned T item = null)
		requires (item != null)
	{
		var tmp = items;
		tmp += (owned) item;
		items = tmp;
		added (items[items.length -1]);

		return items[items.length -1];
	}

	public void remove (T item) {
		T[] tmp = {};

		for (var i = 0; i < items.length; i++) {
			if ((items[i] as CollectionWithSelectedItem).id == (item as CollectionWithSelectedItem).id) {
				continue;
			}
			tmp += items[i];
		}

		removed (item);
		item = tmp;
	}

	public delegate bool ForEachItem<T> (owned T item);
	public void for_each (ForEachItem f) {
		for (var i = 0; i < items.length; i++) {
			if (f (items[i])) {
				return;
			}
		}
	}
}

public class Benchwell._Config : Object {
	public Sqlite.Database db;
	public Benchwell.Settings settings;
	public Benchwell.CollectionWithSelected<Environment> environments;
	public Benchwell.CollectionWithSelected<ConnectionInfo> connections;
	public Benchwell.HttpCollection[] http_collections;
	public Secret.Schema schema;
	public Benchwell.Plugin[] plugins;
	public Json.Node filters;
	public HashTable<string?, bool?> http_tree_state;

	public signal void http_collection_added (HttpCollection collection);

	private string APP_DIR = GLib.Environment.get_user_config_dir () + "/benchwell";

	public _Config () throws Benchwell.ConfigError {
		settings = new Benchwell.Settings ();

		create_app_dir ();
		init_db ();

		schema = new Secret.Schema (Constants.PROJECT_NAME, Secret.SchemaFlags.NONE,
								 "id", Secret.SchemaAttributeType.INTEGER,
								 "schema", Secret.SchemaAttributeType.STRING);

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

	private void create_app_dir () throws GLib.Error {
		var dir = File.new_for_path (APP_DIR);

		if (dir.query_exists ()) {
			return;
		}

		dir.make_directory ();
	}

	private void init_db () throws Benchwell.ConfigError {
		string dbpath = @"$APP_DIR/config.db";
		stdout.printf ("Using config db: %s\n", dbpath);

		var dbfile = File.new_for_path (dbpath);
		var should_create_schema = !dbfile.query_exists ();

		int ec = Sqlite.Database.open_v2 (dbpath, out db, Sqlite.OPEN_READWRITE|Sqlite.OPEN_CREATE);
		if (ec != Sqlite.OK) {
			stderr.printf ("could not open config database: %d: %s\n", db.errcode (), db.errmsg ());
		}

		if (should_create_schema) {
			create_schema ();
		}
	}

	public void create_schema () throws Benchwell.ConfigError {
		stderr.printf ("Creating config file");
		var sch_environments = """
			CREATE TABLE "environments" (
				name TEXT NOT NULL
			);
		""";

		var sch_environment_variables = """
			CREATE TABLE "environment_variables" (
				key            TEXT NOT NULL,
				value          TEXT NOT NULL,
				enabled        BOOLEAN NOT NULL DEFAULT 1 CHECK (enabled IN (0,1)),
				environment_id INTEGER NOT NULL,
				kvtype      INTEGER NOT NULL DEFAULT 1
			);
		""";

		var sch_db_connections = """
			CREATE TABLE "db_connections" (
				name      TEXT NOT NULL,
				adapter   TEXT NOT NULL,
				type      TEXT NOT NULL,
				database  TEXT NULL,
				host      TEXT NULL,
				options   TEXT NULL,
				user      TEXT NULL,
				port      INTEGER NULL,
				encrypted BOOLEAN NOT NULL DEFAULT 0 CHECK (encrypted IN (0,1)),

				Socket    TEXT NULL,
				File      TEXT NULL,
				SshHost   TEXT NULL,
				SshAgent  TEXT NULL
			);
		""";

		var sch_db_queries = """
			CREATE TABLE "db_queries" (
				name           TEXT NOT NULL,
				query          TEXT NOT NULL,
				query_type     TEXT NOT NULL DEFAULT "fav", -- history
				connections_id INTEGER NOT NULL,
				created_at     INTEGER NOT NULL
			);
		""";

		var sch_http_collections = """
			CREATE TABLE "http_collections" (
				count integer default 0,
				name  TEXT NOT NULL
			);
		""";

		var sch_http_items = """
			CREATE TABLE "http_items" (
				name                TEXT NOT NULL,
				description         TEXT NOT NULL DEFAULT "",
				parent_id           INTEGER,
				is_folder           INTEGER,
				count               INTEGER default 0,
				sort                INTEGER NOT NULL,
				http_collections_id INTEGER NOT NULL,
				external_data       TEXT NOT NULL DEFAULT "",

				method TEXT DEFAULT "",
				url    TEXT DEFAULT "",
				body   TEXT DEFAULT "",
				mime   TEXT DEFAULT "json"
			);
		""";

		var sch_http_responses = """
			CREATE TABLE "http_responses" (
				http_items_id INTEGER NOT NULL,
				url           TEXT DEFAULT "",
				method        TEXT DEFAULT "",
				headers	      TEXT DEFAULT "",
				body          TEXT DEFAULT "",
				content_type  TEXT DEFAULT "",
				duration      INTEGER DEFAULT 0,
				code          INTEGER DEFAULT 0,
				created_at    INTEGER NOT NULL
			);
		""";

		var trigger_increment_http_collections_count = """
			CREATE TRIGGER increment_http_collections_count AFTER INSERT ON http_items
				BEGIN
					UPDATE http_collections SET count = count + 1 WHERE http_collections.rowid = NEW.http_collections_id;
				END;
		""";

		var trigger_decrement_http_collections_count = """
			CREATE TRIGGER decrement_http_collections_count AFTER DELETE ON http_items
				BEGIN
					UPDATE http_collections SET count = count - 1 WHERE http_collections.rowid = OLD.http_collections_id;
				END;
		""";

		var sch_http_kvs = """
			CREATE TABLE "http_kvs" (
				key           TEXT NOT NULL,
				value         TEXT NOT NULL,
				type          TEXT NOT NULL,
				kvtype     INTEGER NOT NULL DEFAULT 1,
				http_items_id INTEGER NOT NULL,
				sort          INTEGER NOT NULL,
				enabled       BOOLEAN NOT NULL DEFAULT 1 CHECK (enabled IN (0,1))
			);
		""";

		var sample_db_connections = """
			INSERT INTO db_connections(name, adapter, type, database, host, options, user, port, encrypted)
				  VALUES("localhost", "mysql", "tcp", "", "localhost", "", "", 3306, 0);
		""";

		string[] queries = {
			sch_environments,
			sch_environment_variables,
			sch_db_connections,
			sch_db_queries,
			sch_http_collections,
			sch_http_items,
			sch_http_responses,
			trigger_increment_http_collections_count,
			trigger_decrement_http_collections_count,
			sch_http_kvs,
			sample_db_connections
		};


		foreach (string q in queries) {
			var errmsg = "";
			var ec = db.exec (q, (n_columns, values, column_names) => {
				return 0;
			}, out errmsg);

			if ( ec != Sqlite.OK ){
				throw new Benchwell.ConfigError.INIT (@"unable to initialize config file: $errmsg $q");
			}
		}
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

	public void set_font (Gtk.Widget w, Pango.FontDescription font) {
		var size = font.get_size();
		if (!font.get_size_is_absolute ()) {
			size = size / Pango.SCALE;
		}

		try {
			var style_provider = new Gtk.CssProvider ();
			style_provider.load_from_data (
				"* { font-family: %s; font-size: %dpx}".printf (font.get_family (), size), -1
			);
			w.get_style_context ().add_provider (style_provider, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION);
		} catch (GLib.Error err) {
			debug ("updating font %s", err.message);
		}
	}

	public Benchwell.HttpCollection add_http_collection (Benchwell.HttpCollection? collection = null) throws ConfigError {
		var c = collection;
		if (c == null) {
			c = new Benchwell.HttpCollection ();
			c.name = @"New collection #$(http_collections.length)";
		}
		c.save ();

		var tmp = http_collections;
		tmp += c;
		http_collections = tmp;
		http_collection_added (c);

		return c;
	}

	public unowned Benchwell.HttpCollection? get_selected_http_collection () {
		foreach(weak Benchwell.HttpCollection collection in http_collections) {
			if (collection.id != settings.http_collection_id){
				continue;
			}

			return collection;
		}

		return null;
	}

	public unowned Benchwell.HttpCollection? get_http_collection_by_id (int64 id) {
		foreach(weak Benchwell.HttpCollection collection in http_collections) {
			if (collection.id != id){
				continue;
			}

			return collection;
		}

		return null;
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

	// ======== LOADERS ========
	private void load_connections () throws ConfigError {
		string errmsg;
		ConnectionInfo[] conns = {};
		var ec = db.exec ("SELECT rowid,* FROM db_connections", (n_columns, values, column_names) => {
			var conn = new Benchwell.ConnectionInfo ();
			conn.touch_without_save (() => {
				conn.id = int.parse (values[0]);
				conn.name = values[1];
				conn.adapter = values[2];
				conn.ttype = values[3];
				conn.database = values[4];
				conn.host = values[5];
				conn.user = values[7];
				conn.port = int.parse(values[8]);
				conn.encrypted = values[9] == "1";
				conn.socket = values[10];
				conn.file = values[11];
			});

			ConnectionInfo[] tmp = conns;
			tmp += conn;
			conns = tmp;
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		connections = new CollectionWithSelected<ConnectionInfo> (conns);

		Benchwell.Query[] queries = {};
		ec = db.exec ("SELECT rowid,* FROM db_queries WHERE query_type = 'fav'", (n_columns, values, column_names) => {
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

		connections.for_each ((item) => {
			var conn = item as ConnectionInfo;
			Benchwell.Query[] qq = {};
			foreach (Benchwell.Query query in queries) {
				if (query.connection_id == conn.id) {
					qq += query;
				}
			}
			conn.queries = qq;
			return false;
		});
	}

	private void load_http_collections () throws ConfigError {
		string errmsg;
		var query = """SELECT rowid, name, count
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
		var query = """SELECT rowid, name, is_folder, sort, http_collections_id, method, parent_id
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
		var query = """SELECT rowid,*
						FROM environments
						""";

		Environment[] envs = {};
		var ec = db.exec (query, (n_columns, values, column_names) => {
			var item = new Benchwell.Environment ();

			item.touch_without_save (() => {
				item.id = int64.parse (values[0]);
				item.name = values[1];
			});

			envs += item;
			return 0;
		}, out errmsg);

		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		//environments = new Environments(envs);
		environments = new CollectionWithSelected<Environment> (envs);
		query = """SELECT rowid,* FROM environment_variables""";

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

		environments.for_each ((item) => {
			var env = item as Environment;
			Benchwell.EnvVar[] envvars = {};
			foreach (var v in variables) {
				if (env.id == v.environment_id) {
					envvars += v;
				}
			}
			env.variables = envvars;

			return false;
		});
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

			var enable_node = new Json.Node (Json.NodeType.VALUE);
			enable_node.set_string (stmt.enabled ? "true" : "false");
			table_node.add_element (enable_node);
		}

		Config.settings.set_string ("db-filters", Json.to_string (filters, false));
	}

	public void load_http_tree_state () {
		try {
			var tree_state = Json.from_string (settings.get_string ("http-tree-state"));
			http_tree_state = new HashTable<string?, bool?> (int64_hash, int64_equal);

			tree_state.get_object ().get_members ().foreach ((key) => {
				http_tree_state.insert (key, tree_state.get_object ().get_boolean_member (key));
			});
		} catch (GLib.Error err) {
			Config.show_alert (null, err.message);
		}
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

public class Benchwell.Settings : GLib.Settings {
	public Settings () {
		Object (
			schema_id: "io.benchwell"
		);
	}

	public int window_pos_x {
		get {
			return get_int ("window-pos-x");
		}
		set {
			set_int ("window-pos-x", value);
		}
	}

	public int window_pos_y {
		get {
			return get_int ("window-pos-y");
		}
		set {
			set_int ("window-pos-y", value);
		}
	}

	public int window_size_w {
		get {
			return get_int ("window-size-w");
		}
		set {
			set_int ("window-size-w", value);
		}
	}

	public int window_size_h {
		get {
			return get_int ("window-size-h");
		}
		set {
			set_int ("window-size-h", value);
		}
	}

	public Gtk.PositionType tab_position {
		get {
			switch (get_string ("tab-position")) {
				case "TOP":
					return Gtk.PositionType.TOP;
				case "BOTTOM":
					return Gtk.PositionType.BOTTOM;
				default:
					return Gtk.PositionType.TOP;
			}
		}
		set {
			switch (value) {
				case Gtk.PositionType.TOP:
					set_string ("tab_position", "TOP");
					break;
				case Gtk.PositionType.BOTTOM:
					set_string ("tab_position", "BOTTOM");
					break;
				default:
					set_string ("tab_position", "TOP");
					break;
			}
		}
	}

	public int64 environment_id {
		get {
			return get_int64 ("environment-id");
		}
		set {
			set_int64 ("environment-id", value);
		}
	}

	public int64 http_collection_id {
		get {
			return get_int64 ("http-collection-id");
		}
		set {
			set_int64 ("http-collection-id", value);
		}
	}

	public int64 http_item_id {
		get {
			return get_int64 ("http-item-id");
		}
		set {
			set_int64 ("http-item-id", value);
		}
	}

	public int http_history_limit {
		get {
			return get_int ("http-history-limit");
		}
		set {
			set_int ("http-history-limit", value);
		}
	}

	public string http_tree_state {
		owned get {
			return get_string ("http-tree-state");
		}
		set {
			set_string ("http-tree-state", value);
		}
	}

	public string http_font {
		owned get {
			return get_string ("http-font");
		}
		set {
			set_string ("http-font", value);
		}
	}

	public bool http_single_click_activate {
		get {
			return get_boolean ("http-single-click-activate");
		}
		set {
			set_boolean ("http-single-click-activate", value);
		}
	}

	public bool http_follow_redirect {
		get {
			return get_boolean ("http-follow-redirect");
		}
		set {
			set_boolean ("http-follow-redirect", value);
		}
	}

	public int db_query_history_limit {
		get {
			return get_int ("db-query-history-limit");
		}
		set {
			set_int ("db-query-history-limit", value);
		}
	}

	public string db_filter {
		owned get {
			return get_string ("db-filters");
		}
		set {
			set_string ("db-filters", value);
		}
	}

	public bool db_edit_panel {
		get {
			return get_boolean ("db-edit-panel");
		}
		set {
			set_boolean ("db-edit-panel", value);
		}
	}

	public bool dark_mode {
		get {
			return get_boolean ("dark-mode");
		}
		set {
			set_boolean ("dark-mode", value);
		}
	}

	public string editor_theme {
		owned get {
			return get_string ("editor-theme");
		}
		set {
			set_string ("editor-theme", value);
		}
	}

	public string editor_font {
		owned get {
			return get_string ("editor-font");
		}
		set {
			set_string ("editor-font", value);
		}
	}

	public int editor_tab_width {
		get {
			return get_int ("editor-tab-width");
		}
		set {
			set_int ("editor-tab-width", value);
		}
	}

	public bool editor_line_number {
		get {
			return get_boolean ("editor-line-number");
		}
		set {
			set_boolean ("editor-line-number", value);
		}
	}

	public bool editor_highlight_line {
		get {
			return get_boolean ("editor-highlight-line");
		}
		set {
			set_boolean ("editor-highlight-line", value);
		}
	}

	public bool editor_no_tabs {
		get {
			return get_boolean ("editor-no-tabs");
		}
		set {
			set_boolean ("editor-no-tabs", value);
		}
	}

	public int pomodoro_duration {
		get {
			return get_int ("pomodoro-duration");
		}
		set {
			set_int ("pomodoro-duration", value);
		}
	}

	public int pomodoro_break_duration {
		get {
			return get_int ("pomodoro-break-duration");
		}
		set {
			set_int ("pomodoro-break-duration", value);
		}
	}

	public int pomodoro_long_break_duration {
		get {
			return get_int ("pomodoro-long-break-duration");
		}
		set {
			set_int ("pomodoro-long-break-duration", value);
		}
	}

	public int pomodoro_set_count {
		get {
			return get_int ("pomodoro-set-count");
		}
		set {
			set_int ("pomodoro-set-count", value);
		}
	}
}

