public class Benchwell.ConnectionInfo : Object {
	public int64  id;
	public string adapter  { get; set; }
	public string ttype    { get; set; }
	public string name     { get; set; }
	public string socket   { get; set; }
	public string file     { get; set; }
	public string host     { get; set; }
	public int    port     { get; set; }
	public string user     { get; set; }
	public string password { get; set; }
	public string database { get; set; }
	public string sshhost  { get; set; }
	public string sshagent { get; set; }
	public string options  { get; set; }
	public bool encrypted  { get; set; }
	public Query[] queries;

	private bool no_auto_save;

	public ConnectionInfo () {
		Object(
			adapter: "",
			ttype: "",
			name: "",
			socket: "",
			file: "",
			host: "",
			port: 0,
			user: "",
			password: "",
			database: "",
			sshhost : "",
			sshagent: "",
			options : "",
			encrypted: false
		);

		notify["adapter"].connect (on_save);
		notify["ttype"].connect (on_save);
		notify["name"].connect (on_save);
		notify["socket"].connect (on_save);
		notify["file"].connect (on_save);
		notify["host"].connect (on_save);
		notify["port"].connect (on_save);
		notify["user"].connect (on_save);
		notify["password"].connect (on_save);
		notify["database"].connect (on_save);
		notify["sshhost"].connect (on_save);
		notify["sshagent"].connect (on_save);
		notify["options"].connect (on_save);
		notify["encrypted"].connect (on_save);
	}

	public string to_string() {
		return name;
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

	public void save () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (id > 0) {
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

		var ec = Config.db.prepare_v2 (prepared_query_str, prepared_query_str.length, out stmt);
		if (ec != Sqlite.OK) {
			stderr.printf ("Error: %d: %s\n", Config.db.errcode (), Config.db.errmsg ());
			return;
		}

		if (id > 0) {
			int param_position = stmt.bind_parameter_index ("$ID");
			assert (param_position > 0);
			stmt.bind_int64 (param_position, id);
		}

		int param_position = stmt.bind_parameter_index ("$ADAPTER");
		stmt.bind_text (param_position, adapter);

		param_position = stmt.bind_parameter_index ("$TYPE");
		stmt.bind_text (param_position, ttype);

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, name);

		param_position = stmt.bind_parameter_index ("$SOCKET");
		stmt.bind_text (param_position, socket);

		param_position = stmt.bind_parameter_index ("$FILE");
		stmt.bind_text (param_position, file);

		param_position = stmt.bind_parameter_index ("$HOST");
		stmt.bind_text (param_position, host);

		param_position = stmt.bind_parameter_index ("$PORT");
		stmt.bind_int (param_position, port);

		param_position = stmt.bind_parameter_index ("$USER");
		stmt.bind_text (param_position, user);

		param_position = stmt.bind_parameter_index ("$DATABASE");
		stmt.bind_text (param_position, database);

		param_position = stmt.bind_parameter_index ("$ENC");
		stmt.bind_int (param_position, encrypted ? 1 : 0);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new Benchwell.ConfigError.SAVE_CONNECTION(errmsg);
		}

		if (id == 0) {
			id = Config.db.last_insert_rowid ();
		}

		Config.encrypt (this);
	}

	public void remove () throws Benchwell.ConfigError {
		if ( id == 0 ) {
			return;
		}

		string errmsg = "";
		var ec = Config.db.exec (@"DELETE FROM db_connections WHERE ID = $(id)", null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new Benchwell.ConfigError.DELETE_CONNECTION (errmsg);
		}
	}

	public Benchwell.Query add_query (owned Benchwell.Query? query = null) throws ConfigError {
		if (query == null) {
			query = new Benchwell.Query ();
		}

		query.connection_id = id;

		var tmp = queries;
		tmp += query;
		queries = tmp;

		return query;
	}

	public void remove_query (Benchwell.Query query) throws ConfigError {
		query.remove ();

		Query[] list = {};
		for (var i = 0; i < queries.length; i++) {
			if (query.id == queries[i].id) {
				continue;
			}
			list += query;
		}

		queries = list;
	}
}

public class Benchwell.Query : Object {
	public int64 id            { get; set; }
	public string name         { get; set; }
	public string query        { get; set; }
	public string query_type   { get; set; }
	public int64 connection_id { get; set; }

	private bool no_auto_save;

	public Query () {
		name = "";
		query = "";
		query_type = "fav";
		notify["name"].connect (on_save);
		notify["query"].connect (on_save);
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

	public void save () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (query == null || name == null) {
			return;
		}

		if (id > 0) {
			 prepared_query_str = """
				UPDATE db_queries
					SET name = $NAME, query = $query
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO db_queries(name, query, query_type, connections_id)
				VALUES($NAME, $QUERY, $QUERY_TYPE, $CONNECTION_ID)
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
			assert (param_position > 0);
			stmt.bind_int64 (param_position, this.id);
		} else {
			param_position = stmt.bind_parameter_index ("$CONNECTION_ID");
			stmt.bind_int64 (param_position, this.connection_id);
		}

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, this.name);

		param_position = stmt.bind_parameter_index ("$QUERY");
		stmt.bind_text (param_position, this.query);

		param_position = stmt.bind_parameter_index ("$QUERY_TYPE");
		stmt.bind_text (param_position, this.query_type);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_CONNECTION(errmsg);
		}

		if (id == 0) {
			id = Config.db.last_insert_rowid ();
		}
	}

	public void remove () throws ConfigError {
		if (id == 0) {
			return;
		}

		string errmsg = "";
		var ec = Config.db.exec (@"DELETE FROM db_queries WHERE ID = $(id)", null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.DELETE_CONNECTION (errmsg);
		}
	}
}
