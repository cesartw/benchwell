enum Benchwell.Http.Columns {
	ITEM,
	ICON,
	TEXT,
	METHOD
}

enum Benchwell.Http.CODES {
	OK,
	Created,
	Accepted,
	NoContent,
	BadRequest,
	Unauthorized,
	PaymentRequired,
	Forbidden,
	NotFound,
	MethodNotAllowed,
	RequestTimeout,
	Conflict,
	TooManyRequests,
	InternalServerError,
	NotImplemented,
	BadGateway,
	ServiceUnavailable,
	GatewayTimeout;

	public static Benchwell.Http.CODES? parse (uint code) {
		switch (code) {
			case 200:
				return OK;
			case 201:
				return Created;
			case 202:
				return Accepted;
			case 204:
				return NoContent;
			case 400:
				return BadRequest;
			case 401:
				return Unauthorized;
			case 402:
				return PaymentRequired;
			case 403:
				return Forbidden;
			case 404:
				return NotFound;
			case 405:
				return MethodNotAllowed;
			case 408:
				return RequestTimeout;
			case 409:
				return Conflict;
			case 429:
				return TooManyRequests;
			case 500:
				 return InternalServerError;
			case 501:
				 return NotImplemented;
			case 502:
				 return BadGateway;
			case 503:
				 return ServiceUnavailable;
			case 504:
				 return GatewayTimeout;
			default:
				 return null;
		}
	}

	public string to_string () {
		switch (this) {
			case OK:
				return "OK";
			case Created:
				return "Created";
			case Accepted:
				return "Accepted";
			case NoContent:
				return "No Content";
			case BadRequest:
				return "Bad Request";
			case Unauthorized:
				return "Unauthorized";
			case PaymentRequired:
				return "Payment Required";
			case Forbidden:
				return "Forbidden";
			case NotFound:
				return "NotFound";
			case MethodNotAllowed:
				return "Method Not Allowed";
			case RequestTimeout:
				return "Request Timeout";
			case Conflict:
				return "Conflict";
			case TooManyRequests:
				return "Too Many Requests";
			case InternalServerError:
				 return "Internal Server Error";
			case NotImplemented:
				 return "Not Implemented";
			case BadGateway:
				 return "Bad Gateway";
			case ServiceUnavailable:
				 return "Service Unavailable";
			case GatewayTimeout:
				 return "Gateway Timeout";
			default:
				 return "";
		}
	}
}

[Compact]
private struct buffer_s
{
	uchar[] buffer;
	size_t size_left;
}

[Compact]
private struct buffer_s2
{
	uchar[] buffer;
	size_t size_left;
}

private size_t WriteMemoryCallback(char* ptr, size_t size, size_t nmemb, void* data) {
	size_t total_size = size*nmemb;
	for(int i = 0; i<total_size; i++)
	{
		(( buffer_s* ) data).buffer+= ptr[i];
	}
	(( buffer_s* ) data).buffer+= 0;
	return total_size;
}

private size_t ReadMemoryCallback(char* dest, size_t size, size_t nmemb, void* data) {
	var wt = (( buffer_s2* ) data);
	size_t buffer_size = size*nmemb;
	if ( wt.size_left > 0) {
		size_t copy_this_much = wt.size_left;
		if (copy_this_much > buffer_size) {
			copy_this_much = buffer_size;
		}

		Posix.memcpy(dest, wt.buffer, copy_this_much);

		wt.buffer = wt.buffer[0:copy_this_much];
		wt.size_left -= copy_this_much;
		return copy_this_much;
	}

	return 0;
}

public class Benchwell.Http.Http : Gtk.Paned {
	public Benchwell.ApplicationWindow    window { get; construct; }
	public string                         title;
	public Benchwell.Http.HttpSideBar    sidebar;
	public Benchwell.Http.HttpAddressBar address;

	// request
	public Benchwell.SourceView body;
	public Gtk.ComboBoxText     mime;
	public Gtk.Label            body_size;
	public Benchwell.KeyValues  headers;
	public Benchwell.KeyValues  query_params;
	//////////

	// response
	public Benchwell.SourceView response;
	public Gtk.Label            status_label;
	public Gtk.Label            duration_label;
	public Gtk.Label            response_size_label;
	///////////

	public Benchwell.Environment env;
	public Benchwell.HttpItem? item;

	public Http(Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			wide_handle: true
		);

		title = _("HTTP");

		sidebar = new Benchwell.Http.HttpSideBar ();
		sidebar.show ();

		address = new Benchwell.Http.HttpAddressBar ();
		address.show ();

		// request
		body = new Benchwell.SourceView ();
		body.vexpand = true;
		body.hexpand = true;
		body.show_line_numbers = true;
		body.show_right_margin = true;
		body.auto_indent = true;
		body.show_line_marks = true;
		body.show_line_marks = true;
		body.highlight_current_line = true;
		body.show ();

		var buff = body.get_buffer();
		buff.changed.connect (() => {
			var txt = body.get_text ();
			body_size.set_text (@"$(txt.length/2014)KB");
		});
		var body_sw = new Gtk.ScrolledWindow (null, null);
		body_sw.add (body);
		body_sw.show ();

		mime = new Gtk.ComboBoxText ();
		mime.show ();

		body_size = new Gtk.Label ("0KB");
		body_size.show ();

		item = new Benchwell.HttpItem ();
		headers = new Benchwell.KeyValues ();
		headers.show ();

		query_params = new Benchwell.KeyValues ();
		query_params.show ();

		mime = new Gtk.ComboBoxText ();
		mime.append ("auto", "auto");
		mime.append ("none", "none");
		mime.append ("application/json", "application/json");
		mime.append ("application/html", "application/html");
		mime.append ("application/xml", "application/xml");
		mime.append ("application/yaml", "application/yaml");
		mime.set_active (0);
		mime.show ();

		var mime_options_box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		mime_options_box.pack_start (mime, false, false, 0);
		mime_options_box.pack_end (body_size, false, false, 0);
		mime_options_box.show ();

		var body_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		body_box.vexpand = true;
		body_box.hexpand = true;
		body_box.show ();

		body_box.pack_start(mime_options_box, false, false, 0);
		body_box.pack_start(body_sw, true, true, 0);

		var body_label = new Gtk.Label (_("Body"));
		body_label.show ();

		var params_label = new Gtk.Label (_("Params"));
		params_label.show ();

		var headers_label = new Gtk.Label (_("Headers"));
		headers_label.show ();

		var body_notebook = new Gtk.Notebook ();
		body_notebook.append_page (body_box, body_label);
		body_notebook.append_page (query_params, params_label);
		body_notebook.append_page (headers, headers_label);
		body_notebook.show ();
		//////////

		// response
		var response_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		response_box.show ();

		response  = new Benchwell.SourceView ();
		response.hexpand = true;
		response.vexpand = true;
		response.highlight_current_line = true;
		response.show_line_numbers = true;
		response.show_right_margin = true;
		response.auto_indent = true;
		response.show_line_marks = true;
		//response.get_buffer ().insert_text.connect (() => {
			//return false;
		//});
		response.show ();
		var response_sw = new Gtk.ScrolledWindow (null, null);
		response_sw.add (response);
		response_sw.show ();

		var details_box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		details_box.show ();

		status_label    = new Gtk.Label ("200 OK");
		status_label.show ();

		duration_label  = new Gtk.Label ("0ms");
		duration_label.show ();

		response_size_label  = new Gtk.Label ("0KB");
		response_size_label.show ();

		details_box.pack_start (status_label, false, false, 0);
		details_box.pack_start (duration_label, false, false, 0);
		details_box.pack_end (response_size_label, false, false, 0);

		response_box.pack_start (details_box, false, false, 0);
		response_box.pack_start (response_sw, true, true, 0);
		//////////

		var ws_paned = new Gtk.Paned (Gtk.Orientation.VERTICAL);
		ws_paned.vexpand = true;
		ws_paned.hexpand = true;
		ws_paned.wide_handle = true;
		ws_paned.show ();
		ws_paned.add1 (body_notebook);
		ws_paned.add2 (response_box);

		var ws_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		ws_box.vexpand = true;
		ws_box.hexpand = true;
		ws_box.show ();

		ws_box.pack_start (address, false, false, 0);
		ws_box.pack_start (ws_paned, true, true, 0);

		pack1 (sidebar, false, true);
		pack2 (ws_box, false, true);

		sidebar.load_request.connect (on_load_request);
		mime.changed.connect (() => {
			body.set_language_by_mime_type (mime.get_active_id ());
		});

		address.send_btn.btn.clicked.connect (on_send);
		//address.save_btn.clicked.connect (on_save);

		address.changed.connect (on_request_changed);
		address.changed.connect (build_interpolated_label);

		headers.changed.connect (on_request_changed);

		query_params.changed.connect (on_request_changed);
		query_params.changed.connect (build_interpolated_label);

		headers.row_added.connect (() => {
			if (item != null) {
				try {
					var kv = item.add_header ();
					return kv;
				} catch (ConfigError err) {
					stderr.printf (err.message);
				}
			}
			return new Benchwell.HttpKv ();
		});
		headers.row_removed.connect ((kvi) => {
			if (kvi == null) {
				return;
			}

			var kv = kvi as HttpKv;
			if (kv == null) {
				return;
			}
			try {
				kv.delete ();
			} catch (ConfigError err) {
				stderr.printf (err.message);
			}
		});

		query_params.row_added.connect (() => {
			if (item != null) {
				try {
					var kv = item.add_param ();
					return kv;
				} catch (ConfigError err) {
					stderr.printf (err.message);
				}
			}
			return new Benchwell.HttpKv ();
		});

		query_params.changed.connect (() => {
			address.update_url ();
		});

		query_params.row_removed.connect ((kvi) => {
			if (kvi == null) {
				return;
			}

			var kv = kvi as HttpKv;
			if (kv == null) {
				return;
			}
			try {
				kv.delete ();
			} catch (ConfigError err) {
				stderr.printf (err.message);
			}
		});
	}

	private void on_request_changed () {
		if (item == null) {
			return;
		}

		try {
			item.save ();
		} catch (ConfigError err){
			stderr.printf (err.message);
		}
	}

	// https://github.com/giuliopaci/ValaBindingsDevelopment/blob/master/libcurl-example.vala
	private void on_send () {
		var handle = new Curl.EasyHandle ();

		var url = env.interpolate (address.address.get_text ());
		var method = address.method_combo.get_active_id ();

		switch (method) {
			case "HEAD":
				handle.setopt (Curl.Option.HTTPGET, true);
				handle.setopt (Curl.Option.CUSTOMREQUEST, "HEAD");
				break;
			case "GET":
				handle.setopt (Curl.Option.HTTPGET, true);
				break;
			case "POST":
				handle.setopt (Curl.Option.POST, true);
				break;
			case "PATCH":
				handle.setopt (Curl.Option.POST, true);
				handle.setopt (Curl.Option.CUSTOMREQUEST, "PATCH");
				break;
			case "DELETE":
				handle.setopt (Curl.Option.POST, true);
				handle.setopt (Curl.Option.CUSTOMREQUEST, "DELETE");
				break;
		}

		handle.setopt (Curl.Option.URL, url);
		handle.setopt (Curl.Option.FOLLOWLOCATION, true);

		buffer_s tmp = buffer_s(){ buffer = new uchar[0] };
		handle.setopt(Curl.Option.WRITEFUNCTION, WriteMemoryCallback);
		handle.setopt(Curl.Option.WRITEDATA, ref tmp);

		string[] keys = {};
		string[] values = {};
		headers.get_kvs (out keys, out values);

		Curl.SList headers = null;
		for (var i = 0; i < keys.length; i++) {
			var val = env.interpolate (values[i]);
			headers = Curl.SList.append ((owned) headers, @"$(keys[i]): $(val)");
		}
		handle.setopt (Curl.Option.HTTPHEADER, headers);

		// BODY
		string raw_body = body.get_text ();
		buffer_s2 tmp_body = buffer_s2 () {
			buffer = new uchar[0],
			size_left = Posix.strlen (raw_body)
		};

		if (raw_body != "") {
			for (var i = 0; i<raw_body.length;i++){
				tmp_body.buffer += raw_body[i];
			}

			handle.setopt(Curl.Option.READFUNCTION, ReadMemoryCallback);
			handle.setopt(Curl.Option.READDATA, ref tmp_body);
		}

		// only connects to the host
		var now = get_real_time ();
		var code = handle.perform ();
		var then = get_real_time ();

		int http_code;
		weak string content_type; // otherwise causes double free
		handle.getinfo(Curl.Info.RESPONSE_CODE, out http_code);
		handle.getinfo(Curl.Info.CONTENT_TYPE, out content_type);
		var content = (string)tmp.buffer;
		var s = content_type.split(";")[0];

		var duration = then - now;
		set_response (http_code, (string) content, s, duration);
	}

	private void on_load_request (Benchwell.HttpItem item) {
		this.item = item;

		normalize_item (item);

		address.set_request (item);
		build_interpolated_label ();

		body.get_buffer ().set_text (item.body);
		mime.set_active_id (item.mime);

		response.get_buffer ().set_text ("");
		headers.clear ();
		query_params.clear ();

		foreach (var h in item.headers) {
			headers.add ((Benchwell.KeyValueI) h);
		}

		load_query_params ();

		try {
			if (item.headers.length == 0) {
				headers.add (item.add_header ());
			}

			if (item.query_params.length == 0) {
				query_params.add (item.add_param ());
			}
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	delegate string Interpolator (string s);

	private void build_interpolated_label () {
		Interpolator interpolator = (s) => { return s; };
		if (Config.environment != null) {
			interpolator = Config.environment.interpolate;
		}

		var interpolated_url = interpolator (item.url);

		var builder  = new StringBuilder ();
		builder.append (interpolated_url);
		if ( item.query_params.length > 0 && interpolated_url.index_of ("?") == -1) {
			builder.append ("?");
		}

		var handle = new Curl.EasyHandle ();
		for (var i = 0; i < item.query_params.length; i++) {
			if (item.query_params[i].key == "") {
				continue;
			}
			if (i != 0) {
				builder.append("&");
			}
			var key = handle.escape(item.query_params[i].key, item.query_params[i].key.length);
			var val = interpolator (item.query_params[i].val);
			val = handle.escape(val, val.length);
			builder.append ("%s=%s".printf (key, val));
		}

		interpolated_url = builder.str;
		address.address.tooltip_text = interpolated_url;
		address.address_label.set_text (interpolated_url);
	}

	private void load_query_params () {
		foreach (var q in item.query_params) {
			var found = false;
			query_params.get_children ().foreach ((i) => {
				var kv = i as KeyValue;
				if (kv.keyvalue.key == q.key) {
					kv.entry_val.text = q.val;
					found = true;
				}
			});
			if (!found) {
				query_params.add ((Benchwell.KeyValueI) q);
			}
		}
	}

	public void set_response(uint status, string raw_data, string content_type, int64 duration) {
		switch (content_type) {
			case "application/json":
				var json = new Json.Parser ();
				var generator = new Json.Generator ();

				generator.set_pretty (true);
				try {
					json.load_from_data (raw_data);
				} catch (GLib.Error err) {
					stderr.printf (err.message);
					return;
				}
				generator.set_root (json.get_root ());

				response.get_buffer ().set_text (generator.to_data (null));
				response.set_language ("json");
				break;
			default:
				response.get_buffer ().set_text (raw_data);
				response.set_language (null);
				break;
		}

		if (raw_data.length < 1024) {
			response_size_label.set_text (@"$(raw_data.length)B");
		} else {
			response_size_label.set_text (@"$(raw_data.length/1024)kB");
		}

		var code = Benchwell.Http.CODES.parse (status);
		if (code == null) {
			status_label.set_text (@"$(status)");
		} else {
			status_label.set_text (@"$(status) $(code)");
		}

		duration_label.set_text (@"$(duration/1000)ms");
	}

	private void normalize_item (Benchwell.HttpItem item){
		string[] keys, values;
		item.url = parse_url (item.url, out keys, out values);

		for (var i = 0; i < keys.length; i++) {
			var at = -1;
			for (var ii = 0; ii < item.query_params.length; ii++) {
				if (item.query_params[ii].key == keys[i]) {
					at = ii;
					break;
				}
			}
			if (at != -1){
				item.query_params[at].val = values[at];
				try {
					item.query_params[at].save ();
				} catch (ConfigError err) {
					stderr.printf (err.message);
				}
				continue;
			}
			if (keys[i] == null || keys[i] == "") {
				continue;
			}

			try {
				var kv = item.add_param ();
				kv.key = keys[i];
				kv.val = values[i];
				kv.save ();
			} catch (ConfigError err) {
				stderr.printf (err.message);
			}
		}

		try {
			item.simple_save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	private string parse_url (string url, out string[] keys, out string[] values) {
		keys = {};
		values = {};
		var params_at = url.index_of ("?");
		if (params_at == -1) {
			return url;
		}

		var base_url = url.substring (0, params_at);
		var query_string = url.substring (params_at + 1, -1);

		string[] _keys = {}, _values = {};
		var kv_strings = query_string.split ("&");
		for (var i = 0; i < kv_strings.length; i++) {
			var a = kv_strings[i].split("=");
			_keys += a[0];
			if (a[1] == null)  {
				_values += "";
				continue;
			}
			_values += a[1];
		}
		keys = _keys;
		values = _values;

		return base_url;
	}
}
