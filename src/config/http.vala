public delegate void NoUpdateFunc ();

public class Benchwell.HttpCollection : Object {
	public int64      id;
	public string     name { get; set; }
	public int        count;
	public HttpItem[] items;

	private bool no_auto_save;

	public signal Benchwell.HttpItem item_added (Benchwell.HttpItem item);

	public HttpCollection () {
		notify["name"].connect (on_save);
	}

	private void on_save (Object obj, ParamSpec spec) {
		if (no_auto_save) {
			return;
		}

		try {
			save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	public void save () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (this.id > 0) {
			 prepared_query_str = """
				UPDATE http_collections
				SET name = $NAME
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO http_collections(name)
				VALUES($NAME)
			""";
		}

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		int param_position;
		if (this.id > 0) {
			param_position = stmt.bind_parameter_index ("$ID");
			assert (param_position > 0);
			stmt.bind_int64 (param_position, this.id);
		}

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, this.name);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_CONNECTION(errmsg);
		}

		if (this.id == 0) {
			this.id = Config.db.last_insert_rowid ();
		}
	}

	public void touch_without_save (NoUpdateFunc f) {
		no_auto_save = true;
		f ();
		no_auto_save = false;
	}

	public Benchwell.HttpItem add_item (owned Benchwell.HttpItem? item) throws ConfigError {
		if (item == null) {
			item = new Benchwell.HttpItem ();
		}
		item.http_collection_id = id;

		item.save ();

		if (item.parent_id == 0) {
			var tmp = items;
			tmp += item;
			items = tmp;
		} else {
			for (var i = 0; i < items.length; i++) {
				if (items[i].id == item.parent_id) {
					var tmp = items[i].items;
					tmp += item;
					items[i].items = tmp;
					break;
				}
			}
		}

		item_added (item);
		return item;
	}

	public void remove () throws ConfigError {
		if (this.id == 0) {
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = """DELETE FROM http_collections WHERE ID = $ID""";

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		int param_position = stmt.bind_parameter_index ("$ID");
		stmt.bind_int64 (param_position, this.id);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}

		Config.remove_http_collection (this);
	}

	public void delete_item (Benchwell.HttpItem item) throws ConfigError {
		item.delete ();

		HttpItem[] list = {};
		for (var i = 0; i < items.length; i++) {
			if (item.id == items[i].id) {
				continue;
			}
			list += item;
		}

		items = list;
	}
}

public class Benchwell.HttpItem : Object {
	public int64                id;
	public int64                parent_id;
	public string               name { get; set; }
	public string               description { get; set; }
	public bool                 is_folder;
	public int64                http_collection_id;
	public int                  sort { get; set; }
	public int64                count;
	public Benchwell.HttpItem[] items;

	public string             method { get; set; }
	public string             url { get; set; }
	public string             body { get; set; }
	public string             mime { get; set; }
	public Benchwell.HttpKv[] headers;
	public Benchwell.HttpKv[] query_params;

	internal bool                  loaded;
	private bool no_auto_save;

	public signal Benchwell.HttpKv header_added (Benchwell.HttpKv kv);
	public signal Benchwell.HttpKv query_param_added (Benchwell.HttpKv kv);

	public HttpItem () {
		Object ();
		notify["name"].connect (on_save);
		notify["description"].connect (on_save);
		notify["sort"].connect (on_save);
		notify["method"].connect (on_save);
		notify["url"].connect (on_save);
		notify["body"].connect (on_save);
		notify["mime"].connect (on_save);
	}

	public void touch_without_save (NoUpdateFunc f) {
		no_auto_save = true;
		f ();
		no_auto_save = false;
	}

	private void on_save (Object obj, ParamSpec spec) {
		if (no_auto_save) {
			return;
		}

		try {
			save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	public void save () throws Benchwell.ConfigError {
		//if (parent_id == 0 && !is_folder) {
			//stderr.printf ("a request with no parent_id\n");
			//return;
		//}

		if (http_collection_id == 0) {
			stderr.printf ("an item with no http_collection_id\n");
			return;
		}

		no_auto_save = true;
		if (name == null) {
			if (is_folder) {
				name = _("New folder");
			} else {
				name = _("New request");
			}
		}
		if (url == null) {
			url = "";
		}
		if ((method == null || method == "") && !is_folder) {
			method = "GET";
		}
		if (mime == null && is_folder) {
			mime = "";
		}
		no_auto_save = false;

		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (id > 0) {
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
				VALUES($NAME, $DESCRIPTION, $PARENT_ID, $IS_FOLDER, $SORT,
					$HTTP_COLLECTION_ID, $METHOD, $URL, $BODY, $MIME)
			""";
		}

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		int param_position;
		if (id > 0) {
			param_position = stmt.bind_parameter_index ("$ID");
			stmt.bind_int64 (param_position, id);
		}

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, name);

		param_position = stmt.bind_parameter_index ("$DESCRIPTION");
		stmt.bind_text (param_position, description == null ? "" : description);

		param_position = stmt.bind_parameter_index ("$PARENT_ID");
		stmt.bind_int64 (param_position, parent_id);

		param_position = stmt.bind_parameter_index ("$IS_FOLDER");
		stmt.bind_int (param_position, is_folder ? 1 : 0);

		param_position = stmt.bind_parameter_index ("$SORT");
		stmt.bind_int (param_position, sort);

		param_position = stmt.bind_parameter_index ("$HTTP_COLLECTION_ID");
		stmt.bind_int64 (param_position, http_collection_id);

		param_position = stmt.bind_parameter_index ("$METHOD");
		stmt.bind_text (param_position, method);

		param_position = stmt.bind_parameter_index ("$URL");
		stmt.bind_text (param_position, url);

		param_position = stmt.bind_parameter_index ("$BODY");
		stmt.bind_text (param_position, body);

		param_position = stmt.bind_parameter_index ("$MIME");
		stmt.bind_text (param_position, mime);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}

		if (id == 0) {
			id = Config.db.last_insert_rowid ();
		}
	}

	public Benchwell.HttpKv add_header () throws ConfigError {
		var kv = new Benchwell.HttpKv ();
		kv.key = "";
		kv.val = "";
		kv.type = "header";
		kv.http_item_id = id;
		kv.save ();

		var tmp = headers;
		tmp += kv;
		headers = tmp;

		header_added (kv);
		return kv;
	}

	public Benchwell.HttpKv add_param () throws ConfigError {
		var kv = new Benchwell.HttpKv ();
		kv.key = "";
		kv.val = "";
		kv.type = "param";
		kv.http_item_id = id;
		kv.save ();

		var tmp = query_params;
		tmp += kv;
		query_params = tmp;

		query_param_added (kv);
		return kv;
	}

	public void delete () throws ConfigError {
		if (id == 0){
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = "DELETE FROM http_items WHERE id = $ID";

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		int param_position;
		param_position = stmt.bind_parameter_index ("$ID");
		stmt.bind_int64 (param_position, id);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}

		// key values
		prepared_query_str = "DELETE FROM http_kvs WHERE http_items_id = $ID";

		ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		param_position = stmt.bind_parameter_index ("$ID");
		stmt.bind_int64 (param_position, id);

		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}
	}

	public void delete_query_param (Benchwell.HttpKv kv) {
		HttpKv[] new_query_params = {};
		for (var i = 0; i < query_params.length; i++) {
			if (query_params[i].key ==	kv.key) {
				try {
					kv.delete ();
				} catch (ConfigError err) {
					stderr.printf (err.message);
				}
				continue;
			}
			new_query_params += kv;
		}

		query_params = new_query_params;
	}

	public void load_full_item () throws Benchwell.ConfigError {
		touch_without_save (() => {
			if (loaded) {
				return;
			}

			string errmsg = "";

			Benchwell.HttpItem[] new_items = {};
			// folder
			if (is_folder) {
				var query = """SELECT id, name, parent_id, is_folder, sort,
											http_collections_id, method
								FROM http_items
								WHERE http_collections_id = %lld AND parent_id = %lld
								ORDER BY sort ASC
								""".printf (http_collection_id, id);
				var ec = Config.db.exec (query, (n_columns, values, column_names) => {
					var subitem = new Benchwell.HttpItem ();

					subitem.touch_without_save (() => {
						subitem.id = int64.parse (values[0]);
						subitem.name = values[1];
						subitem.parent_id = int64.parse (values[2]);
						subitem.is_folder = values[3] == "1";
						subitem.sort = int.parse (values[4]);
						subitem.http_collection_id = int64.parse (values[5]);
						subitem.method = values[6];
					});

					new_items += subitem;
					return 0;
				}, out errmsg);
				if ( ec != Sqlite.OK ){
					throw new ConfigError.GET_CONNECTIONS(errmsg);
				}

				items = new_items;
				return;
			}

			// request
			var query = """SELECT ifnull(method,""), ifnull(url,""), ifnull(body, ""), ifnull(mime,"")
					FROM http_items
					WHERE id = %lld""".printf (id);
			var ec = Config.db.exec (query, (n_columns, values, column_names) => {
				touch_without_save (() => {
					method = values[0];
					url = values[1];
					body = values[2];
					mime = values[3];
				});
				return 0;
			}, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.GET_CONNECTIONS(errmsg);
			}

			Benchwell.HttpKv[] kvs = {};
			query = """SELECT id, ifnull(key, ""), ifnull(value, ""), type, sort, enabled, http_items_id
				FROM http_kvs
				WHERE http_items_id = %lld
				ORDER BY sort ASC""".printf (id);
			ec = Config.db.exec (query, (n_columns, values, column_names) => {
				var kv = new Benchwell.HttpKv ();

				kv.touch_without_save (() => {
					kv.id = int64.parse (values[0]);
					kv.key = values[1];
					kv.val = values[2];
					kv.type = values[3];
					kv.sort = int.parse (values[4]);
					kv.enabled = values[5] == "1";
					kv.http_item_id = int64.parse (values[6]);
				});

				kvs += kv;
				return 0;
			}, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.GET_CONNECTIONS(errmsg);
			}

			Benchwell.HttpKv[] new_headers = {};
			Benchwell.HttpKv[] new_query_params = {};
			foreach (var kv in kvs) {
				if (kv.type == "header") {
					new_headers += kv;
					continue;
				}
				new_query_params += kv;
			}
			headers = new_headers;
			query_params = new_query_params;
			loaded = true;
		});
	}
}

public class Benchwell.HttpKv : Object, Benchwell.KeyValueI {
	public int64  id      { get; set; }
	public string key     { get; set; }
	public string val     { get; set; }
	public bool   enabled { get; set; }
	public string type; // header | param
	public int    sort;
	public int64  http_item_id;
	public bool   no_auto_save;

	public HttpKv () {
		notify["key"].connect (on_save);
		notify["val"].connect (on_save);
		notify["enabled"].connect (on_save);
	}

	private void on_save (Object obj, ParamSpec spec) {
		if (no_auto_save) {
			return;
		}

		try {
			save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	public void touch_without_save (NoUpdateFunc f) {
		no_auto_save = true;
		f ();
		no_auto_save = false;
	}

	public void save () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (val == null) {
			val = "";
		}

		if (key == null) {
			return;
		}
		if (type == null) {
			return;
		}

		if (this.id > 0) {
			 prepared_query_str = """
				UPDATE http_kvs
				SET key = $KEY,
					value = $VALUE,
					type = $TYPE,
					http_items_id = $HTTP_ITEM_ID,
					sort = $SORT,
					enabled = $ENABLED
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO http_kvs(key, value, type, http_items_id, sort, enabled)
				VALUES($KEY, $VALUE, $TYPE, $HTTP_ITEM_ID, $SORT, $ENABLED)
			""";
		}

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		int param_position;
		if (this.id > 0) {
			param_position = stmt.bind_parameter_index ("$ID");
			assert (param_position > 0);
			stmt.bind_int64 (param_position, this.id);
		}

		param_position = stmt.bind_parameter_index ("$KEY");
		stmt.bind_text (param_position, this.key);

		param_position = stmt.bind_parameter_index ("$VALUE");
		stmt.bind_text (param_position, this.val);

		param_position = stmt.bind_parameter_index ("$TYPE");
		stmt.bind_text (param_position, this.type);

		param_position = stmt.bind_parameter_index ("$HTTP_ITEM_ID");
		stmt.bind_int64 (param_position, this.http_item_id);

		param_position = stmt.bind_parameter_index ("$SORT");
		stmt.bind_int (param_position, this.sort);

		param_position = stmt.bind_parameter_index ("$ENABLED");
		stmt.bind_int (param_position, this.enabled ? 1 : 0);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_CONNECTION(errmsg);
		}

		if (this.id == 0) {
			this.id = Config.db.last_insert_rowid ();
		}
	}

	public void delete () throws ConfigError {
		if (id == 0){
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = "DELETE FROM http_kvs WHERE id = $ID";

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		var param_position = stmt.bind_parameter_index ("$ID");
		stmt.bind_int64 (param_position, id);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}
	}
}
