public class Benchwell.HttpCollection : Object {
	public int64      id;
	public string     name;
	public int        count;
	public HttpItem[] items;

	public signal Benchwell.HttpItem item_added (Benchwell.HttpItem item);

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

	public signal Benchwell.HttpKv header_added (Benchwell.HttpKv kv);
	public signal Benchwell.HttpKv query_param_added (Benchwell.HttpKv kv);

	public void save () throws Benchwell.ConfigError {
		simple_save ();
		for (var i = 0; i < headers.length; i++) {
			headers[i].save ();
		}
		for (var i = 0; i < query_params.length; i++) {
			query_params[i].save ();
		}
	}

	public void simple_save () throws Benchwell.ConfigError {
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
		if (method == null) {
			method = "GET";
		}
		if (mime == null) {
			mime = "";
		}

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
}

public class Benchwell.HttpKv : Object, Benchwell.KeyValueI {
	public int64  id;
	public string key     { get; set; }
	public string val     { get; set; }
	public bool   enabled { get; set; }
	public string type; // header | param
	public int    sort;
	public int64  http_item_id;

	public void save () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

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
