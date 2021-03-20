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

// https://github.com/giuliopaci/ValaBindingsDevelopment/blob/master/libcurl-example.vala
private size_t ReadResponseCallback (char* ptr, size_t size, size_t nmemb, void* data) {
	size_t total_size = size*nmemb;
	var buffer = (( buffer_s* ) data);
	// remove the termination char(0)
	if (buffer.buffer.length > 0 && buffer.buffer[buffer.buffer.length - 1] == 0) {
		buffer.buffer = buffer.buffer[:buffer.buffer.length - 1];
	}

	for(int i = 0; i<total_size; i++)
	{
		buffer.buffer+= ptr[i];
	}
	buffer.buffer+= 0;
	return total_size;
}

private size_t WriteRequestCallback (char* dest, size_t size, size_t nmemb, void* data) {
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

private size_t ReadHeaderCallback (char *dest, size_t size, size_t nmemb, void *data) {
	var wt = ((HashTable<string,string>) data);
	 //size_t numbytes = size * nmemb;
	//printf("%.*s\n", numbytes, dest);
	var header = (string)dest;
	var at = header.index_of (":", 0);
	if ( at != -1 ){
		var key = header[0:at];
		var val = header[at+2:header.length];
		wt.insert (key, val);
	}
	return size * nmemb;
}

private int RequestProgressCallback (void *clientp, double dltotal, double dlnow, double ultotal, double ulnow) {
	var canceled = ((bool*) clientp);
	if (*canceled) {
		return 1;
	}
	return 0;
}

public class Benchwell.HttpResult : Object {
	public string Body;
	public HashTable<string, string> Headers;
	public int Status;
	public int64 Duration;

	public string headers_string () {
		var builder = new StringBuilder ();
		Headers.foreach ((key, val) => {
			builder.append (@"$key: $val");
		});
		return builder.str;
	}
}

public class Benchwell.Http.Http : Gtk.Paned {
	public Benchwell.ApplicationWindow   window { get; construct; }
	public string                        title  { get; set; }
	public Benchwell.Http.HttpSideBar    sidebar;
	public Benchwell.Http.HttpAddressBar address;
	public Benchwell.HttpOverlay overlay;

	// request
	public Benchwell.SourceView body;
	//public Gtk.ComboBoxText     mime;
	public Benchwell.KeyValues  headers;
	public Benchwell.KeyValues  query_params;
	//////////

	// response
	public Benchwell.SourceView response;
	public Benchwell.SourceView response_headers;
	public Gtk.Label            status_label;
	public Gtk.Label            duration_label;
	public Gtk.Label            response_size_label;
	///////////

	public Benchwell.HttpItem? item;
	public Gtk.TreeIter? item_iter;

	public Http(Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			wide_handle: true
		);

		title = _("HTTP");

		sidebar = new Benchwell.Http.HttpSideBar (window);
		sidebar.show ();

		address = new Benchwell.Http.HttpAddressBar ();
		address.show ();

		// request
		body = new Benchwell.SourceView ();
		body.show ();

		var body_sw = new Gtk.ScrolledWindow (null, null);
		body_sw.add (body);
		body_sw.show ();

		item = new Benchwell.HttpItem ();
		headers = new Benchwell.KeyValues ();
		headers.show ();

		query_params = new Benchwell.KeyValues ();
		query_params.show ();

		var body_label = new Gtk.Label (_("Body"));
		body_label.show ();

		var params_label = new Gtk.Label (_("Params"));
		params_label.show ();

		var headers_label = new Gtk.Label (_("Headers"));
		headers_label.show ();

		var body_notebook = new Gtk.Notebook ();
		body_notebook.append_page (body_sw, body_label);
		body_notebook.append_page (query_params, params_label);
		body_notebook.append_page (headers, headers_label);
		body_notebook.show ();
		//////////

		// response
		var response_paned = new Gtk.Paned (Gtk.Orientation.VERTICAL);
		response_paned.show ();

		var response_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		response_box.show ();

		response_headers = new Benchwell.SourceView ();
		response_headers.editable = false;
		response_headers.margin_bottom = 10;
		response_headers.margin_top = 10;
		response_headers.show ();
		var response_headers_sw = new Gtk.ScrolledWindow (null, null);
		response_headers_sw.add (response_headers);
		response_headers_sw.show ();

		response  = new Benchwell.SourceView ();
		response.editable = false;
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

		details_box.get_style_context ().add_class ("responsemeta");
		details_box.pack_start (status_label, false, false, 0);
		details_box.pack_start (duration_label, false, false, 0);
		details_box.pack_end (response_size_label, false, false, 0);

		response_paned.pack1 (response_headers_sw, false, false);
		response_paned.pack2 (response_sw, true, true);
		response_box.pack_start (details_box, false, false, 0);
		response_box.pack_start (response_paned, true, true, 0);
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


		overlay = new Benchwell.HttpOverlay ();
		overlay.add (ws_box);
		overlay.show ();

		pack1 (sidebar, false, true);
		pack2 (overlay, false, true);

		sidebar.item_activated.connect (on_item_activated);
		sidebar.item_removed.connect (on_item_removed);

		address.send_btn.btn.clicked.connect (on_send);
		address.send_btn.menu_btn.activate.connect (on_save_as);
		//address.save_btn.clicked.connect (on_save);

		address.changed.connect (on_request_changed);
		address.changed.connect (build_interpolated_label);

		headers.changed.connect (on_request_changed);

		query_params.changed.connect (on_request_changed);
		query_params.changed.connect (build_interpolated_label);

		body.get_buffer ().changed.connect (() => {
			if (item == null) {
				return;
			}
			Gtk.TextIter start, end;
			body.get_buffer ().get_start_iter (out start);
			body.get_buffer ().get_end_iter (out end);
			item.body = body.get_buffer ().get_text (start, end, false);
		});

		headers.row_wanted.connect (() => {
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

		query_params.row_wanted.connect (() => {
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

		headers.row_added.connect ((kv) => {
			if (kv.key.strip () .casefold () == "Content-Type".casefold ()) {
				body.set_language_by_mime_type (kv.val);
			}
		});
		headers.changed.connect (() => {
			headers.get_children ().foreach ((w) => {
				var kv = w as Benchwell.KeyValue;
				if (kv.keyvalue.key.strip ().casefold () == "Content-Type".casefold ()) {
					body.set_language_by_mime_type (kv.keyvalue.val);
				}
			});
		});
	}

	private void on_request_changed () {
		if (item_iter == null)
			return;

		// NOTE: hacky way to force external item changes to be updated in the sidebar
		sidebar.store.set_value (item_iter, Benchwell.Http.Columns.ITEM, item);
	}

	private void on_save_as () {
		print ("=======saveas\n");
	}

	// https://github.com/giuliopaci/ValaBindingsDevelopment/blob/master/libcurl-example.vala
	public void on_send () {
		response_headers.get_buffer ().set_text ("", 0);
		response.get_buffer ().set_text ("", 0);

		var handle = new Curl.EasyHandle ();

		var url = Config.environment.interpolate (address.address.get_text ());
		var method = address.method_combo.get_active_id ();

		string[] keys = {};
		string[] values = {};
		query_params.get_kvs (out keys, out values);

		// TODO: Improve URL parsing
		var builder = new StringBuilder ();
		if (url.index_of ("?") == -1)
			builder.append ("?");

		for (var i = 0; i < keys.length; i++) {
			var key = Config.environment.interpolate (keys[i]);
			var val = Config.environment.interpolate (values[i]);
			key = handle.escape (key, key.length);
			val = handle.escape (val, val.length);
			builder.append (@"$key=$val");
			if (i < keys.length - 1)
				builder.append ("&");
		}

		if (url.index_of ("?") != -1 && !url.has_suffix("&"))
			builder.prepend ("&");
		url += builder.str;

		keys = {};
		values = {};
		headers.get_kvs (out keys, out values);

		Curl.SList headers = null;
		for (var i = 0; i < keys.length; i++) {
			var val = Config.environment.interpolate (values[i]);
			headers = Curl.SList.append ((owned) headers, @"$(keys[i]): $(val)");
		}

		// BODY
		string raw_body = body.get_text ();
		raw_body = Config.environment.interpolate (raw_body);
		buffer_s2 tmp_body = buffer_s2 () {
			buffer = new uchar[0],
			size_left = Posix.strlen (raw_body)+1
		};

		perform.begin (method, url, raw_body, (owned) headers, (obj, res) => {
			try {
				var result = perform.end (res);
				if (result != null) {
					item.response_body = set_response (result);
					item.response_headers = result.headers_string();
					item.save_response ();
				}
			} catch (ThreadError e) {
				stderr.printf (@"performing request $(e.message)");
			}
			overlay.stop ();
		});
	}

	public async HttpResult? perform (string method, string url, string body, owned Curl.SList headers) {
		overlay.start ();
		bool canceled = false;
		var cancel_handler_id = overlay.cancel.connect (() => {
			canceled = true;
		});

		SourceFunc callback = perform.callback;
		HttpResult? result = new HttpResult ();

		ThreadFunc<bool> run = () => {
			var handle = new Curl.EasyHandle ();

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
			handle.setopt (Curl.Option.HTTPHEADER, headers);

			buffer_s tmp = buffer_s(){ buffer = new uchar[0] };
			handle.setopt(Curl.Option.WRITEFUNCTION, ReadResponseCallback);
			handle.setopt(Curl.Option.WRITEDATA, ref tmp);

			var resp_headers = new HashTable<string, string> (str_hash, str_equal);
			handle.setopt (Curl.Option.HEADERFUNCTION, ReadHeaderCallback);
			handle.setopt (Curl.Option.HEADERDATA, resp_headers);
			handle.setopt (Curl.Option.PROGRESSFUNCTION, RequestProgressCallback);
			handle.setopt (Curl.Option.PROGRESSDATA, ref canceled);
			handle.setopt (Curl.Option.NOPROGRESS, 0);

			buffer_s2 tmp_body = buffer_s2 () {
				buffer = new uchar[0],
				size_left = Posix.strlen (body)+1
			};
			if (body != "") {
				for (var i = 0; i < body.length; i++){
					tmp_body.buffer += body[i];
				}
				tmp_body.buffer += 0;
				handle.setopt (Curl.Option.READFUNCTION, WriteRequestCallback);
				handle.setopt (Curl.Option.READDATA, ref tmp_body);
			}

			var then = get_real_time ();
			var code = handle.perform ();
			var now = get_real_time ();

			switch (code) {
				case Curl.Code.OK:
					handle.getinfo(Curl.Info.RESPONSE_CODE, out result.Status);
					result.Body = (string)tmp.buffer;
					result.Duration = now - then;
					result.Headers = resp_headers;
					break;
				case Curl.Code.URL_MALFORMAT:
					stderr.printf (@"========$(url)\n");
					break;
				case Curl.Code.ABORTED_BY_CALLBACK:
					result = null;
					break;
				default:
					stderr.printf (@"========$((int)code)\n");
					break;
			}

			Idle.add((owned) callback);
			return true;
		};
		new Thread<bool>("benchwell-http", run);
		yield;
		overlay.disconnect (cancel_handler_id);
		return result;
	}

	private void on_item_activated (Benchwell.HttpItem item, Gtk.TreeIter iter) {
		this.item = item;
		this.item_iter = iter;
		try {
			this.item.load_full_item ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
			return;
		}

		title = item.name;

		normalize_item (this.item);

		address.set_request (this.item);
		build_interpolated_label ();

		this.item.touch_without_save (() => {
			if (this.item.body != null)
				body.get_buffer ().text = this.item.body;
			else
				body.get_buffer ().text = "";
		});

		response.get_buffer ().text = item.response_body;
		response_headers.get_buffer ().text = item.response_headers;
		headers.clear ();
		query_params.clear ();

		foreach (var h in this.item.headers) {
			headers.add ((Benchwell.KeyValueI) h);
		}

		load_query_params ();

		try {
			if (this.item.headers.length == 0) {
				headers.add (item.add_header ());
			}

			if (this.item.query_params.length == 0) {
				query_params.add (this.item.add_param ());
			}
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	private void on_item_removed (Benchwell.HttpItem removed_item) {
		if (item == null) {
			return;
		}

		if (item.id == removed_item.id) {
			item = null;

			address.address.text = "";
			address.method_combo.set_active_id ("GET");
			response_headers.get_buffer ().text = "";
			response.get_buffer ().text = "";
			body.get_buffer ().text = "";
			headers.clear ();
			query_params.clear ();
		}
	}

	delegate string Interpolator (string s);

	private void build_interpolated_label () {
		if (item == null) {
			address.address.tooltip_text = "";
			address.address_label.set_text ("");
			return;
		}
		Interpolator interpolator = (s) => { return s; };
		if (Config.environment != null) {
			interpolator = Config.environment.dry_interpolate;
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

	public string set_response (HttpResult result) {
		var content_type_header = result.Headers.get("Content-Type");
		string content_type = "";
		if (content_type_header != null)
			content_type = content_type_header.split(";")[0];

		if (result.Body == null)
			result.Body = "";

		Gtk.TextIter iter;
		response_headers.get_buffer ().get_end_iter (out iter);
		result.Headers.foreach ((key, val) => {
			var line = @"$key: $val";
			response_headers.get_buffer ().insert (ref iter, line, line.length);
		});
		response_headers.get_buffer ().insert (ref iter, "\n", 1);

		var formatted_body = "";
		Gtk.TextIter start_iter;
		response.get_buffer ().get_start_iter (out start_iter);
		if ( result.Body != "" ) {
			switch (content_type) {
				case "application/json":
					var json = new Json.Parser ();
					var generator = new Json.Generator ();
					generator.indent = 4;
					generator.pretty = true;
					try {
						json.load_from_data (result.Body);
						generator.set_root (json.get_root ());

						formatted_body = generator.to_data (null);
						response.get_buffer ().insert_text (ref start_iter, formatted_body, formatted_body.length);
						response.set_language ("json");
					} catch (GLib.Error err) {
						stderr.printf (err.message);
						formatted_body = result.Body;
						response.get_buffer ().insert_text (ref start_iter, formatted_body, formatted_body.length);
					}
					break;
				default:
					response.set_language (null);
					formatted_body = result.Body;
					response.get_buffer ().insert_text (ref start_iter, formatted_body, formatted_body.length);
					break;
			}
		} else {
			Gtk.TextIter end_iter;
			response.get_buffer ().get_end_iter (out end_iter);
			response.get_buffer ().delete (ref start_iter, ref end_iter);
		}

		if (result.Body.length < 1024) {
			response_size_label.set_text (@"$(result.Body.length)B");
		} else {
			response_size_label.set_text (@"$(result.Body.length/1024)kB");
		}

		var code = Benchwell.Http.CODES.parse (result.Status);
		if (code == null) {
			status_label.set_text (@"$(result.Status)");
		} else {
			status_label.set_text (@"$(result.Status) $(code)");
		}

		status_label.get_style_context ().remove_class("ok");
		status_label.get_style_context ().remove_class("warning");
		status_label.get_style_context ().remove_class("bad");

		if (result.Status >= 200 && result.Status < 300) {
			status_label.get_style_context ().add_class("ok");
		}
		if (result.Status >= 300 && result.Status < 500) {
			status_label.get_style_context ().add_class("warning");
		}
		if (result.Status >= 500) {
			status_label.get_style_context ().add_class("bad");
		}

		duration_label.set_text (@"$(result.Duration/1000)ms");
		return formatted_body;
	}

	private void normalize_item (Benchwell.HttpItem item){
		string[] keys = null, values = null;

		item.touch_without_save (() => {
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
					continue;
				}
				if (keys[i] == null || keys[i] == "") {
					continue;
				}

				try {
					var kv = item.add_param ();
					kv.key = keys[i];
					kv.val = values[i];
					//kv.save ();
				} catch (ConfigError err) {
					stderr.printf (err.message);
				}
			}
		});

		//try {
			//item.save ();
		//} catch (ConfigError err) {
			//stderr.printf (err.message);
		//}
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

public class Benchwell.HttpOverlay : Gtk.Overlay {
	public Gtk.Button btn_cancel;
	public Gtk.Spinner spinner;
	public Gtk.Box box;

	public signal void cancel ();

	public HttpOverlay () {
		Object(
			name: "HttpOverlay"
		);

		btn_cancel = new Gtk.Button.with_label (_("Cancel"));
		btn_cancel.set_size_request (100, 30);
		btn_cancel.show ();

		spinner = new Gtk.Spinner ();
		spinner.show ();

		box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		box.get_style_context ().add_class ("overlay-bg");

		var center_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		center_box.set_size_request (100, 150);
		center_box.valign = Gtk.Align.CENTER;
		center_box.halign = Gtk.Align.CENTER;
		center_box.vexpand = true;
		center_box.hexpand = true;
		center_box.show ();

		box.add (center_box);

		center_box.pack_start (spinner, true, true, 0);
		center_box.pack_start (btn_cancel, false, false, 0);
		add_overlay (box);

		btn_cancel.clicked.connect (on_cancel);
	}

	private void on_cancel () {
		spinner.stop ();
		box.hide ();
		cancel ();
	}

	public void start () {
		box.show ();
		spinner.start ();
	}

	public void stop () {
		spinner.stop ();
		box.hide ();
	}
}
