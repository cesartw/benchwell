public class Benchwell.Http.Plugin {
	public HashTable<string, JSC.Value> plugins;
	private JSC.Context jsctx;

	// example of JSC https://github.com/fread-ink/fread.ui/blob/master/web_extensions/fread.c
	public Plugin () {
		plugins = new HashTable<string, JSC.Value> (str_hash, str_equal);
		jsctx = new JSC.Context ();

		try {
			string folder = GLib.Environment.get_user_config_dir () + "/benchwell/plugins";
			var directory = File.new_for_path (folder);
			if ( !directory.query_exists () ) {
				directory.make_directory ();
			}
			var enumerator = directory.enumerate_children (FileAttribute.STANDARD_NAME, 0);

			FileInfo file_info;
			while ((file_info = enumerator.next_file ()) != null) {
				var file_name = file_info.get_name ();
				var file_path = folder + "/" + file_name;
				var file = File.new_for_path (file_path);
				var stream = new DataInputStream (file.read ());

				string line;
				string data = "";
				while ((line = stream.read_line ()) != null) {
					data += line + "\n";
				}

				jsctx.evaluate (data, data.length);
				var exception = jsctx.get_exception ();
				if (exception != null) {
					stderr.printf ("%s\n%s\n", file_path, exception.to_string ());
					stderr.printf ("at: %s#%zu:%zu \n", file_name, exception.get_line_number (), exception.get_column_number ());
					var backtrace = exception.get_backtrace_string ();
					if (backtrace != null) {
						stderr.printf ("==========EXCEPTION==========\n%s\n=============================\n", backtrace);
					}
					continue;
				}

				var call = jsctx.get_value ("call");

				if ( !call.is_function ()) {
					stderr.printf ("%s must define function `call`\n", file_path);
					continue;
				}

				if ( file_name.has_suffix (".js") ) {
					file_name = file_name.substring (0, file_name.length - 3);
				}

				plugins.insert (file_name, call);
			}
		} catch (Error e) {
			stderr.printf ("error %s\n", e.message);
		}
	}

	public JSC.Value new_string (string s) {
		return new JSC.Value.string (jsctx, s);
	}

	public JSC.Value new_number (string s) {
		return new JSC.Value.number (jsctx, double.parse (s));
	}

	public JSC.Value new_bool (string s) {
		return new JSC.Value.boolean (jsctx, bool.parse (s));
	}

	public JSC.Value[]? parse_params (string raw_params) {
		JSC.Value[]? parameters = {};

		bool in_single_quote = false;
		bool in_double_quote = false;
		bool in_number       = false;
		bool in_bare_string  = false;
		bool in_escape       = false;

		var current = new StringBuilder ();
		for (var i = 0; i < raw_params.length; i++) {
			unichar c = raw_params[i];

			if (in_escape) {
				in_escape = false;
				current.append_unichar (c);
				continue;
			}

			if (in_bare_string && c != ' ' && c != '\\') {
				current.append_unichar (c);
				continue;
			}

			switch (c) {
			case '"':
				in_double_quote = !in_double_quote;
				if ( !in_double_quote ) {
					parameters += Config.plugins.new_string(current.str);
					current = new StringBuilder ();
				}
				continue;
			case '\'':
				in_single_quote = !in_single_quote;
				if ( !in_single_quote ) {
					parameters += Config.plugins.new_string(current.str);
					current = new StringBuilder ();
				}
				continue;
			case '\\':
				in_escape = true;
				continue;
			case '0', '1', '2', '3' ,'4', '5', '6', '7', '8', '9', '.':
				in_number = true;
				current.append_unichar (c);
				continue;
			case ' ':
				if ( in_bare_string ) {
					if ( in_escape ) {
						current.append_unichar (c);
						in_escape = false;
						continue;
					}

					in_bare_string = false;

					if ( current.str == "true" || current.str == "false" ) {
						parameters += Config.plugins.new_bool(current.str);
					} else {
						parameters += Config.plugins.new_string(current.str);
					}

					current = new StringBuilder ();
					continue;
				}

				if ( in_double_quote || in_single_quote ) {
					current.append_unichar (c);
					continue;
				}

				if ( in_number ) {
					parameters += Config.plugins.new_number(current.str);
					in_number = false;
					current = new StringBuilder ();
					continue;
				}

				break;
			default:
				if ( !in_single_quote && !in_double_quote )
					in_bare_string = true;
				current.append_unichar (c);
				continue;
			}
		}

		if ( in_bare_string ) {
			if ( current.str == "true" || current.str == "false" ) {
				parameters += Config.plugins.new_bool(current.str);
			} else {
				parameters += Config.plugins.new_string(current.str);
			}
		}

		if ( in_number ) {
			parameters += Config.plugins.new_number(current.str);
		}

		return parameters;
	}
}
