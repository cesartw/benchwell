public class Benchwell.Environment : Object {
	public int64  id;
	public string name { get; set; }
	public Benchwell.EnvVar[] variables;

	public signal void variable_added (Benchwell.EnvVar envvar);

	public Regex var_escape_regex  = /({{\s*([a-zA-Z0-9_]+)\s*}})/;
	public Regex func_escape_regex = /({%\s*([a-zA-Z0-9_]+)\s*(.*)%})/;

	private bool no_auto_save;

	public Environment () {
		notify["name"].connect (on_save);
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
		if (no_auto_save) {
			return;
		}

		try {
			save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	private void save () throws ConfigError {
		Sqlite.Statement stmt;
		string prepared_query_str = "";

		if (this.id > 0) {
			 prepared_query_str = """
				UPDATE environments
					SET name = $NAME
				WHERE rowid = $ID
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
			throw new ConfigError.STORE(errmsg);
		}

		if (this.id == 0) {
			this.id = Config.db.last_insert_rowid ();
		}
	}

	public void clone () throws ConfigError {
		var env = new Benchwell.Environment ();
		env.touch_without_save (() => {
			env.name = @"$name Copy";
		});
		env.save ();

		foreach (var kv in variables) {
			env.add_variable (kv.key, kv.val);
		}

		Config.add_environment (env);
	}

	public void remove () throws ConfigError {
		if (this.id == 0) {
			return;
		}

		Sqlite.Statement stmt;
		string prepared_query_str = """DELETE FROM environments WHERE rowid = $ID""";

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
			throw new ConfigError.STORE(errmsg);
		}

		Config.remove_environment (this);
	}

	public string interpolate (string s) {
		var result = interpolate_variables (s);
		result = interpolate_functions (result);

		return result;
	}

	public string dry_interpolate (string s) {
		var result = interpolate_variables (s);
		result = interpolate_functions (result);

		return result;
	}

	public string interpolate_variables (string s) {
		MatchInfo info;
		string result = s;

		try {
			for (var_escape_regex.match (s, 0, out info); info.matches () ; info.next ()) {
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
		} catch (GLib.RegexError err) {
			stderr.printf (err.message);
		}

		return result;
	}

	public string interpolate_functions (string s) {
		MatchInfo info;
		string result = s;

		try {
			for (func_escape_regex.match (s, 0, out info); info.matches () ; info.next ()) {
				var raw_params = info.fetch (3);
				var func_name = info.fetch (2);
				var to_replace = info.fetch (1);

				try {
					foreach (var plugin in Config.plugins) {
						if ( plugin.name != func_name ) {
							continue;
						}

						var parameters = plugin.parse_params (raw_params, this);
						var plugin_result = plugin.callv (parameters);
						result = result.replace (to_replace, plugin_result);
						break;
					}
				} catch (Benchwell.PluginError err) {
					stderr.printf (err.message);
				}
			}
		} catch (GLib.RegexError err) {
			stderr.printf (err.message);
		}

		return result;
	}

	public string dry_interpolate_functions (string s) {
		MatchInfo info;
		string result = s;

		try {
			for (func_escape_regex.match (s, 0, out info); info.matches () ; info.next ()) {
				var func_name = info.fetch (2);
				var to_replace = info.fetch (1);

				foreach (var plugin in Config.plugins) {
					if ( plugin.name != func_name ) {
						continue;
					}

					result = result.replace (to_replace, @"{% $(func_name) %}");
					break;
				};
			}
		} catch (GLib.RegexError err) {
			stderr.printf (err.message);
		}

		return result;
	}

	public Benchwell.EnvVar add_variable (string key = "", string val = "") throws ConfigError {
		var envvar = new Benchwell.EnvVar ();
		envvar.touch_without_save (() => {
			envvar.key = key;
			envvar.val = val;
			envvar.environment_id = id;
		});
		envvar.save ();

		var tmp = variables;
		tmp += envvar;
		variables = tmp;

		variable_added (envvar);
		return envvar;
	}
}

public class Benchwell.EnvVar : Object, Benchwell.KeyValueI {
	public int64  id      { get; set; }
	public string key     { get; set; }
	public string val     { get; set; }
	public bool   enabled { get; set; }
	public Benchwell.KeyValueTypes kvtype { get; set; }
	public int    sort    { get; set; }
	public int64  environment_id;

	private bool no_auto_save;

	public EnvVar () {
		kvtype = Benchwell.KeyValueTypes.STRING;
		notify["key"].connect (on_save);
		notify["val"].connect (on_save);
		notify["enabled"].connect (on_save);
		notify["sort"].connect (on_save);
		notify["kvtype"].connect (on_save);
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
		Sqlite.Statement stmt;
		string prepared_query_str = "";
		if (this.key == null || this.key == "") {
			return;
		}

		if (this.id > 0) {
			 prepared_query_str = """
				UPDATE environment_variables
					SET key = $KEY, value = $VALUE, enabled = $ENABLED, kvtype = $KVTYPE
				WHERE rowid = $ID
			""";
		} else {
			 prepared_query_str = """
				INSERT INTO environment_variables(key, value, enabled, environment_id, kvtype)
				VALUES($KEY, $VALUE, $ENABLED, $ENVIRONMENT_ID, $KVTYPE)
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

		param_position = stmt.bind_parameter_index ("$KVTYPE");
		stmt.bind_int (param_position, this.kvtype);

		string errmsg = "";
		ec = Config.db.exec (stmt.expanded_sql(), null, out errmsg);
		if ( ec != Sqlite.OK ){
			throw new ConfigError.STORE(errmsg);
		}

		if (this.id == 0) {
			this.id = Config.db.last_insert_rowid ();
		}
	}
}



public class Benchwell.EnvvarCompletion : Object, Gtk.SourceCompletionProvider {
	public string get_name () {
		return _("Envaronment vars");
	}

	public bool match (Gtk.SourceCompletionContext context) {
		var end = context.iter;
		var start = context.iter;
		if (!start.backward_chars (2))
			return false;

		var s = context.completion.view.get_buffer ().get_text (start, end, false);

		return s == "{{";
	}

	public void populate (Gtk.SourceCompletionContext context) {
		GLib.List<Gtk.SourceCompletionProposal> proposals = null;
		foreach (var envvar in Config.environment.variables) {
			if (envvar.key == "")
				continue;
			proposals.append (new Benchwell.EnvvarCompletionProposal (envvar));
		}
		context.add_proposals (this, proposals, true);
	}
}

public class Benchwell.EnvvarCompletionProposal : Object, Gtk.SourceCompletionProposal {
	public weak Benchwell.EnvVar envvar { get; construct; }

	public EnvvarCompletionProposal (Benchwell.EnvVar envvar) {
		Object(
			envvar: envvar
		);
	}

	public bool equal (Gtk.SourceCompletionProposal other) {
		var other_of_this_type = other as Benchwell.EnvvarCompletionProposal;
		if (other_of_this_type == null)
			return false;

		return this.envvar.id == other_of_this_type.envvar.id;
	}

	public string get_text () {
		return @" $(envvar.key) }}";
	}

	// get_markup has priority
	public string get_label () {
		return "";
	}

	public unowned GLib.Icon? get_gicon () {
		return null;
	}

	public unowned Gdk.Pixbuf? get_icon () {
		return null;
	}

	public unowned string? get_icon_name () {
		return null;
	}

	public string? get_info () {
		return envvar.val;
	}

	public string get_markup () {
		return envvar.key;
	}
}
