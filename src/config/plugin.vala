public errordomain Benchwell.PluginError {
	PARSE_ERROR,
	UNKNOWN_FUNC,
	UNKNOWN_VAR
}

public interface Benchwell.Plugin : Object {
	public abstract string name { get; construct; }

	// parse_params from text inputs
	// {% myfunc param1 "param2" 3 'param3' true false @some_var %}
	//    paramaters at 0(bare string), 1, 3 are strings
	//    parameter at 2 is a double
	//    parameter at 4 and 5 are booleans
	//    paramter at 5 is another variable in the current environment.
	//      plugins will get the value of the var which is an string. if the doesn't exist, "" is provided
	//
	// TODO: support environment variables and functions
	// Given that the environment has the following variables
	//     token = "123"
	//     user  = 1
	// When we have
	// // variables
	// {% myfunc param1 "param2" 3 'param3' true false @token \@token "ehlo @token" "ehlo \@token" %}
	//    parameter at 5 must converted to '123'
	//    parameter at 6 must converted to '\@token' (considered a bare string)
	//    parameter at 7 must converted to 'ehlo 123'
	//    parameter at 7 must converted to 'ehlo 123'
	//    parameter at 8 must converted to 'ehlo @token'
	//
	// // functions
	// {% myfunc $(myotherfunc ...see param spec...) %}
	public virtual GLib.Value[]? parse_params (string raw_params, Benchwell.Environment? env = null) throws Benchwell.PluginError {
		GLib.Value[]? parameters = {};

		bool in_single_quote = false;
		bool in_double_quote = false;
		bool in_number       = false;
		bool in_bare_string  = false;
		bool in_escape       = false;
		bool in_var          = false;

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
				if (in_number || in_var) {
					throw new Benchwell.PluginError.PARSE_ERROR (@"unexpected \" after `$(current.str)`");
				}

				in_double_quote = !in_double_quote;
				if ( !in_double_quote ) {
					var val = GLib.Value (typeof (string));
					val.set_string (current.str);
					parameters += val;
					current = new StringBuilder ();
				}
				continue;
			case '\'':
				if (in_number) {
					throw new Benchwell.PluginError.PARSE_ERROR (@"unexpected ' after `$(current.str)`");
				}

				in_single_quote = !in_single_quote;
				if ( !in_single_quote ) {
					var val = GLib.Value (typeof (string));
					val.set_string (current.str);
					parameters += val;
					current = new StringBuilder ();
				}
				continue;
			case '\\':
				if (in_number) {
					throw new Benchwell.PluginError.PARSE_ERROR (@"unexpected \\ after `$(current.str)`");
				}

				in_escape = true;
				continue;
			case '0', '1', '2', '3' ,'4', '5', '6', '7', '8', '9', '.':
				in_number = !in_bare_string && !in_double_quote && !in_single_quote && !in_var;
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
						var val = GLib.Value (typeof (bool));
						val.set_boolean (current.str == "true");
						parameters += val;
					} else {
						var val = GLib.Value (typeof (string));
						val.set_string (current.str);
						parameters += val;
					}

					current = new StringBuilder ();
					continue;
				}

				if ( in_double_quote || in_single_quote ) {
					current.append_unichar (c);
					continue;
				}

				if ( in_number ) {
					var val = GLib.Value (typeof (double));
					val.set_double (double.parse (current.str));
					parameters += val;

					in_number = false;
					current = new StringBuilder ();
					continue;
				}

				if ( in_var ) {
					var val = GLib.Value (typeof (string));

					string var_name = current.str;
					string var_val = "";
					if (env != null) {
						for (var envi = 0; envi < env.variables.length; envi++) {
							if (env.variables[envi].key == var_name) {
								var_val = env.variables[envi].val;
								break;
								//var_val = env.interpolate_functions(var_val);
								//var_val = env.interpolate_variables(var_val);
							}
						}
					}

					val.set_string (var_val);
					parameters += val;

					current = new StringBuilder ();
					continue;
				}

				break;
			case '@':
				// TODO: allow inline interpolation
				in_var = !in_bare_string && !in_double_quote && !in_single_quote && !in_var;
				break;
			default:
				if ( !in_single_quote && !in_double_quote && !in_number && !in_var )
					in_bare_string = true;
				current.append_unichar (c);
				continue;
			}
		}

		if ( in_bare_string ) {
			if ( current.str == "true" || current.str == "false" ) {
				var val = GLib.Value (typeof (bool));
				val.set_boolean (current.str == "true");
				parameters += val;
			} else {
				var val = GLib.Value (typeof (string));
				val.set_string (current.str);
				parameters += val;
			}
		}

		if ( in_number ) {
			var val = GLib.Value (typeof (double));
			val.set_double (double.parse (current.str));
			parameters += val;
		}

		return parameters;
	}

	public abstract string callv (GLib.Value[] parameters);
}

public class Benchwell.JSPlugin : Object, Benchwell.Plugin {
	public string name      { get; construct; }
	public JSC.Value call   { get; construct; }
	public JSC.Context ctx  { get; construct; }

	// example of JSC https://github.com/fread-ink/fread.ui/blob/master/web_extensions/fread.c
	protected JSPlugin (owned JSC.Context ctx, string name, owned JSC.Value call) {
		Object(
			name: name,
			ctx: ctx,
			call: call
		);
	}

	public static Plugin[] load () {
		Plugin[] plugins = {};

		try {
			string folder = GLib.Environment.get_user_config_dir () + "/benchwell/plugins";
			var directory = File.new_for_path (folder);
			if ( !directory.query_exists () ) {
				directory.make_directory ();
			}
			var enumerator = directory.enumerate_children (FileAttribute.STANDARD_NAME, 0);

			FileInfo file_info;
			while ((file_info = enumerator.next_file ()) != null) {
				var jsctx = new JSC.Context ();
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

				plugins += new JSPlugin (jsctx, file_name, call);
			}
		} catch (GLib.IOError err) {
			stderr.printf ("error %s\n", err.message);
		} catch (GLib.Error err) {
			stderr.printf ("error %s\n", err.message);
		}


		return plugins;
	}

	public string callv (GLib.Value[] parameters) {
		JSC.Value[] jsparams = {};

		foreach (var val in parameters) {
			switch (val.type ()) {
				case GLib.Type.STRING:
					jsparams += new JSC.Value.string (ctx, val.get_string ());
					break;
				case GLib.Type.DOUBLE:
					jsparams += new JSC.Value.number (ctx, val.get_double ());
					break;
				case GLib.Type.BOOLEAN:
					jsparams += new JSC.Value.boolean (ctx, val.get_boolean ());
					break;
			}
		}

		return call.function_callv (jsparams).to_string ();
	}
}


public class Benchwell.BuiltinPlugin : Object, Benchwell.Plugin {
	public delegate string call (GLib.Value[] parameters);
	public string name { get; construct; }
	private call* f;

	protected BuiltinPlugin (string name, call f) {
		Object(
			name: name
		);
		this.f = &f;
	}

	public string callv (GLib.Value[] parameters) {
		return ((call)*f)(parameters);
	}

	public static Benchwell.Plugin[] load () {
		Benchwell.Plugin[] plugins = {};

		// BASE64
		plugins += new Benchwell.BuiltinPlugin("base64", (parameters) => {
			if (parameters.length == 0) {
				return "";
			}
			if (parameters[0].type () != GLib.Type.STRING) {
				return "";
			}

			return GLib.Base64.encode (parameters[0].get_string ().data);
		});
		/////////

		// ENCODE URL
		plugins += new Benchwell.BuiltinPlugin("url_encode", (parameters) => {
			if (parameters.length == 0) {
				return "";
			}
			if (parameters[0].type () != GLib.Type.STRING) {
				return "";
			}

			var handle = new Curl.EasyHandle ();
			var s = parameters[0].get_string ();
			return handle.escape (s, s.length);
		});
		/////////

		return plugins;
	}
}

