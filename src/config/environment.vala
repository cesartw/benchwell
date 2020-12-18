public class Benchwell.Environment : Object {
	public int64  id;
	public string name;
	public Benchwell.EnvVar[] variables;

	public signal Benchwell.EnvVar variable_added (Benchwell.EnvVar envvar);

	public Regex regex = /({{\s*([a-zA-Z0-9]+)\s*}})/;

	public Environment () {}

	public string interpolate (string s) {
		MatchInfo info;
		string result = s;

		for (regex.match (s, 0, out info); info.matches () ; info.next ()) {
			for (var i = info.get_match_count () - 1; i > 0; i-=2) {
				var var_name = info.fetch (i);
				var to_replace = info.fetch (i-1);

				foreach (var envvar in variables) {
					if (envvar.key == var_name) {
						result = result.replace (to_replace, envvar.val);
					}
				}
			}
		}

		return result;
	}

	public Benchwell.EnvVar add_variable () throws ConfigError {
		var envvar = new Benchwell.EnvVar ();
		envvar.save ();

		var tmp = variables;
		tmp += envvar;
		variables = tmp;

		variable_added (envvar);
		return envvar;
	}

	public void save () throws ConfigError {
		save_environment ();
		foreach (var envvar in variables) {
			envvar.environment_id = this.id;
			envvar.save ();
		}
	}

	private void save_environment () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (this.id > 0) {
			 prepared_query_str = """
				UPDATE environments
					SET name = $NAME
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO environments(name)
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
			stmt.bind_int64 (param_position, this.id);
		}

		param_position = stmt.bind_parameter_index ("$NAME");
		stmt.bind_text (param_position, this.name);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
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
		string prepared_query_str = """DELETE FROM environments WHERE ID = $ID""";

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

		Config.remove_environment (this);
	}
}

public class Benchwell.EnvVar : Object, Benchwell.KeyValueI {
	public int64  id      { get; set; }
	public string key     { get; set; }
	public string val     { get; set; }
	public bool   enabled { get; set; }
	public string type; // header | param
	public int    sort;
	public int64  environment_id;

	public void save () throws Benchwell.ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";
		if (this.key == null || this.key == "") {
			return;
		}

		if (this.id > 0) {
			 prepared_query_str = """
				UPDATE environment_variables
					SET key = $KEY, value = $VALUE, enabled = $ENABLED
				WHERE ID = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO environment_variables(key, value, enabled, environment_id)
				VALUES($KEY, $VALUE, $ENABLED, $ENVIRONMENT_ID)
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
			stmt.bind_int64 (param_position, this.id);
		} else {
			param_position = stmt.bind_parameter_index ("$ENVIRONMENT_ID");
			stmt.bind_int64 (param_position, this.environment_id);
		}

		param_position = stmt.bind_parameter_index ("$KEY");
		stmt.bind_text (param_position, this.key);

		param_position = stmt.bind_parameter_index ("$VALUE");
		stmt.bind_text (param_position, this.val);

		param_position = stmt.bind_parameter_index ("$ENABLED");
		stmt.bind_int (param_position, this.enabled ? 1 : 0);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			stderr.printf ("SQL: %s\n", stmt.expanded_sql());
			stderr.printf ("ERR: %s\n", errmsg);
			throw new ConfigError.SAVE_ENVVAR(errmsg);
		}

		if (this.id == 0) {
			this.id = Config.db.last_insert_rowid ();
		}
	}
}

