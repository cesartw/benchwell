public class Benchwell.ConnectionInfo {
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
		adapter     = "";
		ttype       = "";
		name        = "";
		socket      = "";
		file        = "";
		host        = "";
		port        = 0;
		user        = "";
		password    = "";
		database    = "";
		sshhost     = "";
		sshagent    = "";
		options     = "";
		encrypted   = false;
	}

	public string to_string() {
		return name;
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

}

public class Benchwell.Query : Object {
	public int64 id           { get; set; }
	public string name        { get; set; }
	public string query       { get; set; }
	public int64 connection_id { get; set; }
}
