public delegate void NoUpdateFunc () throws Benchwell.ConfigError;

public class Benchwell.HttpCollection : Object {
	public int64      id;
	public string     name { get; set; }
	public int        count;
	public HttpItem[] items;

	private bool no_auto_save;

	public signal void item_added (Benchwell.HttpItem item);

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
				WHERE rowid = $ID
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
			throw new ConfigError.STORE(errmsg);
		}

		if (this.id == 0) {
			this.id = Config.db.last_insert_rowid ();
		}
	}

	public void remove () throws ConfigError {
		if (this.id == 0) {
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = """DELETE FROM http_collections WHERE rowid = $ID""";

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
			throw new ConfigError.STORE(errmsg);
		}

		Config.remove_http_collection (this);
	}

	public void touch_without_save (NoUpdateFunc f) {
		no_auto_save = true;
		try {
			f ();
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (null, err.message);
		}
		no_auto_save = false;
	}

	public Benchwell.HttpItem add_item (owned Benchwell.HttpItem? item) throws ConfigError {
		if (item == null) {
			item = new Benchwell.HttpItem ();
		}
		item.http_collection_id = id;

		item.save_all ();

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

	public Benchwell.HttpItem clone_item (Benchwell.HttpItem? item) throws Benchwell.ConfigError {
		item.load_full_item ();

		var new_item = add_item (null);
		new_item.touch_without_save (() => {
			new_item.is_folder = false;
			new_item.parent_id = item.parent_id;
			new_item.name = item.name + _(" Copy");
			new_item.description = item.description;
			new_item.http_collection_id = item.http_collection_id;
			new_item.method = item.method;
			new_item.url = item.url;
			new_item.body = item.body;
			new_item.mime = item.mime;

			foreach (HttpKv kv in item.headers) {
				new_item.add_header (kv.key, kv.val);
			}
			foreach (HttpKv param in item.query_params) {
				new_item.add_param (param.key, param.val);
			}
		});

		new_item.save_all ();
		new_item.load_full_item ();

		return new_item;
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
	public int64                parent_id { get; set; }
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
	public Benchwell.HttpKv[] form_params;
	public Benchwell.Http.Result[] responses;

	internal bool              loaded;
	private bool               no_auto_save;

	public signal Benchwell.HttpKv header_added (Benchwell.HttpKv kv);
	public signal Benchwell.HttpKv query_param_added (Benchwell.HttpKv kv);
	public signal Benchwell.HttpKv form_param_added (Benchwell.HttpKv kv);
	public signal void response_added (owned Benchwell.Http.Result response);

	public HttpItem () {
		Object ();

		headers = {};
		query_params = {};
		form_params = {};
		notify["name"].connect (on_save);
		notify["description"].connect (on_save);
		notify["sort"].connect (on_save);
		notify["method"].connect (on_save);
		notify["parent_id"].connect (on_save);

		notify["url"].connect (on_save_body);
		notify["body"].connect (on_save_body);
		notify["mime"].connect (on_save_body);
	}

	public void touch_without_save (NoUpdateFunc f) {
		no_auto_save = true;
		try {
			f ();
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (null, err.message);
		}
		no_auto_save = false;
	}

	private void on_save (Object obj, ParamSpec spec) {
		if (no_auto_save || http_collection_id == 0) {
			return;
		}

		try {
			save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	private void on_save_body (Object obj, ParamSpec spec) {
		if (no_auto_save || !loaded) {
			return;
		}

		try {
			save_body ();
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
						method = $METHOD
				WHERE rowid = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO http_items(name, description, parent_id, is_folder,
					sort, http_collections_id, method)
				VALUES($NAME, $DESCRIPTION, $PARENT_ID, $IS_FOLDER, $SORT,
					$HTTP_COLLECTION_ID, $METHOD)
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

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		if (id == 0) {
			id = Config.db.last_insert_rowid ();
		}
	}

	public void save_body () throws Benchwell.ConfigError {
		no_auto_save = true;
		if (mime == null && is_folder) {
			mime = "";
		}
		no_auto_save = false;

		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (id > 0) {
			 prepared_query_str = """
				UPDATE http_items
				SET url = $URL, body = $BODY, mime = $MIME
				WHERE rowid = $ID
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

		param_position = stmt.bind_parameter_index ("$URL");
		stmt.bind_text (param_position, url);

		param_position = stmt.bind_parameter_index ("$BODY");
		stmt.bind_text (param_position, body);

		param_position = stmt.bind_parameter_index ("$MIME");
		stmt.bind_text (param_position, mime);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		if (id == 0) {
			id = Config.db.last_insert_rowid ();
		}
	}

	public void save_all () throws Benchwell.ConfigError {
		save ();
		save_body ();
	}

	public void save_response (Benchwell.Http.Result result) throws Benchwell.ConfigError {
		if (id == 0) {
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = "";
		prepared_query_str = """
			INSERT INTO http_responses (http_items_id, method, url, content_type, body, headers, duration, code, created_at)
			VALUES($HTTP_ITEM_ID, $METHOD, $URL, $CONTENT_TYPE, $BODY, $HEADERS, $DURATION, $CODE, $CREATED_AT)
		""";

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		int param_position;

		param_position = stmt.bind_parameter_index ("$HTTP_ITEM_ID");
		stmt.bind_int64 (param_position, id);

		param_position = stmt.bind_parameter_index ("$METHOD");
		stmt.bind_text (param_position, result.method);

		param_position = stmt.bind_parameter_index ("$URL");
		stmt.bind_text (param_position, result.url);

		param_position = stmt.bind_parameter_index ("$CONTENT_TYPE");
		stmt.bind_text (param_position, result.content_type);

		param_position = stmt.bind_parameter_index ("$BODY");
		stmt.bind_text (param_position, result.body);

		param_position = stmt.bind_parameter_index ("$HEADERS");
		stmt.bind_text (param_position, result.headers);

		param_position = stmt.bind_parameter_index ("$DURATION");
		stmt.bind_int64 (param_position, result.duration);

		param_position = stmt.bind_parameter_index ("$CODE");
		stmt.bind_int (param_position, result.status);

		result.created_at = new DateTime.now_local ();
		param_position = stmt.bind_parameter_index ("$CREATED_AT");
		stmt.bind_int64 (param_position, result.created_at.to_unix ());

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		Benchwell.Http.Result[] l_responses = {result};
		foreach (var response in responses) {
			l_responses += response;
			if (l_responses.length == Config.settings.http_history_limit)
				break;
		}
		responses = l_responses;


		prepared_query_str = """DELETE FROM http_responses
			WHERE
			http_items_id = $HTTP_ITEM_ID AND
			rowid NOT IN (
				SELECT rowid
				FROM http_responses
				WHERE http_items_id = $HTTP_ITEM_ID
				ORDER BY rowid DESC
				LIMIT $LIMIT
			)""";

		ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		param_position = stmt.bind_parameter_index ("$HTTP_ITEM_ID");
		stmt.bind_int64 (param_position, id);

		param_position = stmt.bind_parameter_index ("$LIMIT");
		stmt.bind_int (param_position, Config.settings.http_history_limit);

		errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		response_added (result);
	}

	public Benchwell.HttpKv add_header (string key = "", string val = "") throws ConfigError {
		var kv = new Benchwell.HttpKv ();
		kv.touch_without_save (() => {
			kv.key = key;
			kv.val = val;
			kv.type = "header";
			kv.http_item_id = id;
		});
		kv.save ();

		var tmp = headers;
		tmp += kv;
		headers = tmp;

		header_added (kv);
		return kv;
	}

	public Benchwell.HttpKv add_param (string key = "", string val = "") throws ConfigError {
		var kv = new Benchwell.HttpKv ();
		kv.touch_without_save (() => {
			kv.key = key;
			kv.val = val;
			kv.type = "param";
			kv.http_item_id = id;
		});
		if (http_collection_id != 0)
			kv.save ();

		var tmp = query_params;
		tmp += kv;
		query_params = tmp;

		query_param_added (kv);
		return kv;
	}

	public Benchwell.HttpKv add_form_param (string key = "", string val = "", KeyValueTypes kvtype = KeyValueTypes.STRING) throws ConfigError {
		var kv = new Benchwell.HttpKv ();
		kv.touch_without_save (() => {
			kv.key = key;
			kv.val = val;
			kv.type = "form_param";
			kv.kvtype = kvtype;
			kv.http_item_id = id;
		});
		if (http_collection_id != 0)
			kv.save ();

		var tmp = form_params;
		tmp += kv;
		form_params = tmp;

		form_param_added (kv);
		return kv;
	}

	public void delete () throws ConfigError {
		if (id == 0){
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = "DELETE FROM http_items WHERE rowid = $ID";

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
			throw new ConfigError.STORE(errmsg);
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
			throw new ConfigError.STORE(errmsg);
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

			// request
			var query = """SELECT ifnull(method,"GET"), ifnull(url,""), ifnull(body, ""), ifnull(mime,""), ifnull(description,"")
					FROM http_items
					WHERE rowid = %lld""".printf (id);
			var ec = Config.db.exec (query, (n_columns, values, column_names) => {
				method = values[0];
				url = values[1];
				body = values[2];
				mime = values[3];
				description = values[4];
				return 0;
			}, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.STORE(errmsg);
			}

			// responses
			Benchwell.Http.Result[] l_responses = {};
			query = """SELECT method, url, headers, body, content_type, duration, code, created_at
				FROM http_responses
				WHERE http_items_id = %lld
				ORDER BY created_at DESC""".printf (id);

			ec = Config.db.exec (query, (n_columns, values, column_names) => {
				var r = new Benchwell.Http.Result ();

				r.method = values[0] ?? "";
				r.url = values[1] ?? "";
				r.headers = values[2] ?? "";
				r.body = values[3] ?? "";
				r.content_type = values[4] ?? "";
				r.duration = int64.parse (values[5]);
				r.status = int.parse (values[6]);
				r.created_at = new DateTime.from_unix_local (int64.parse (values[7]));

				l_responses += r;
				return 0;
			}, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.STORE(errmsg);
			}
			responses = l_responses;

			Benchwell.HttpKv[] kvs = {};
			query = """SELECT rowid, ifnull(key, ""), ifnull(value, ""), type, sort, enabled, http_items_id, kvtype
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
					kv.kvtype = (Benchwell.KeyValueTypes)int.parse (values[7]);
				});

				kvs += kv;
				return 0;
			}, out errmsg);
			if ( ec != Sqlite.OK ){
				throw new ConfigError.STORE(errmsg);
			}

			Benchwell.HttpKv[] new_headers = {};
			Benchwell.HttpKv[] new_query_params = {};
			Benchwell.HttpKv[] new_form_params = {};
			foreach (var kv in kvs) {
				switch (kv.type) {
					case "header":
						new_headers += kv;
						break;
					case "form_param":
						new_form_params += kv;
						break;
					case "param":
						new_query_params += kv;
						break;
				}
			}
			headers = new_headers;
			query_params = new_query_params;
			form_params = new_form_params;
			loaded = true;
		});
	}

	public Benchwell.Http.Result? last_response () {
		if (responses.length == 0)
			return null;
		return responses[0];
	}
}

public class Benchwell.HttpKv : Object, Benchwell.KeyValueI {
	public int64  id      { get; set; }
	public string key     { get; set; }
	public string val     { get; set; }
	public bool   enabled { get; set; }
	public Benchwell.KeyValueTypes kvtype { get; set; }

	public string type; // header | param
	public int    sort;
	public int64  http_item_id;
	public bool   no_auto_save;

	public HttpKv () {
		Object ();

		enabled = true;
		kvtype = Benchwell.KeyValueTypes.STRING;
		notify["key"].connect (on_save);
		notify["val"].connect (on_save);
		notify["enabled"].connect (on_save);
		notify["kvtype"].connect (on_save);
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
		try {
			f ();
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (null, err.message);
		}
		no_auto_save = false;
	}

	public void save () throws Benchwell.ConfigError {
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

		if (http_item_id == 0) {
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
					enabled = $ENABLED,
					kvtype = $KVTYPE
				WHERE rowid = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO http_kvs(key, value, type, http_items_id, sort, enabled, kvtype)
				VALUES($KEY, $VALUE, $TYPE, $HTTP_ITEM_ID, $SORT, $ENABLED, $KVTYPE)
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

		param_position = stmt.bind_parameter_index ("$KVTYPE");
		stmt.bind_int (param_position, kvtype);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
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
		string prepared_query_str = "DELETE FROM http_kvs WHERE rowid = $ID";

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
			throw new ConfigError.STORE(errmsg);
		}
	}
}
