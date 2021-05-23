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
}

[Compact]
private struct buffer_s2
{
	uchar[] buffer;
	size_t size_left;
}

public string RandomString (int length, string char_set = "abcdefghijklmnopqrstuvwxyz0123456789") {
	var builder = new StringBuilder ();

	for (var i = 0; i < length; i++) {
		builder.append_c (char_set[Random.int_range(0, length - 1)]);
	}

	return builder.str;
}

public class Benchwell.HttpResult : Object {
	public string method;
	public string url;
	public string body;
	public string headers;
	public int status;
	public int64 duration;
	public string content_type;
	public Curl.Code code;
	public DateTime created_at;
}

public class Benchwell.CBNotebookTab : Gtk.Box {
	public Gtk.ComboBoxText combo;
	public Gtk.Label label;
	public bool enabled  { get; set; }

	public CBNotebookTab (string l, bool enabled = false) {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 0
		);

		combo = new Gtk.ComboBoxText ();
		label = new Gtk.Label (l);

		pack_start (combo, true, true, 0);
		pack_start (label, true, true, 0);
		combo.changed.connect (() => {
			label.set_text (combo.get_active_text ());
		});

		this.enabled = enabled;
		on_toggle ();
		notify["enabled"].connect (on_toggle);
	}

	private void on_toggle () {
		if (enabled) {
			combo.show ();
			label.hide ();
			return;
		}
		combo.hide ();
		label.show ();
	}
}

public class Benchwell.Http.Http : Gtk.Paned {
	public Benchwell.ApplicationWindow   window { get; construct; }
	public string                        title  { get; set; }
	public Benchwell.Http.HttpSideBar    sidebar;
	public Benchwell.Http.HttpAddressBar address;
	public Benchwell.HttpOverlay overlay;

	// request
	public Gtk.Stack body_stack;
	public Benchwell.SourceView body;
	public Benchwell.KeyValues  body_fields;
	public Benchwell.CBNotebookTab mime_switch;
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
	public Benchwell.HttpHistoryPopover history_popover;
	///////////

	public Benchwell.HttpItem? item;
	public Gtk.TreeIter? item_iter;

	private bool loading = false;
	private Regex kvrg;

	public Http(Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			wide_handle: true
		);
		kvrg = new Regex ("(.+)=(@?)(.+)");

		title = _("HTTP");

		sidebar = new Benchwell.Http.HttpSideBar (window);
		sidebar.show ();

		address = new Benchwell.Http.HttpAddressBar ();
		address.show ();

		// request
		body = new Benchwell.SourceView ();
		body.show ();

		body_fields = new Benchwell.KeyValues (Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE|Benchwell.KeyValueTypes.FILE);
		body_fields.show ();

		var body_sw = new Gtk.ScrolledWindow (null, null);
		body_sw.add (body);
		body_sw.show ();

		var body_fields_sw = new Gtk.ScrolledWindow (null, null);
		body_fields_sw.add (body_fields);
		body_fields_sw.show ();

		body_stack = new Gtk.Stack ();
		body_stack.add_named (body_sw, "editor");
		body_stack.add_named (body_fields_sw, "fields");
		body_stack.set_visible_child_name ("editor");
		body_stack.show ();

		item = new Benchwell.HttpItem ();
		headers = new Benchwell.KeyValues (Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE);
		headers.show ();

		query_params = new Benchwell.KeyValues (Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE);
		query_params.show ();

		var params_label = new Gtk.Label (_("Params"));
		params_label.show ();

		var headers_label = new Gtk.Label (_("Headers"));
		headers_label.show ();

		mime_switch = new Benchwell.CBNotebookTab (_("Body"), true);
		mime_switch.combo.append("plain/text", "Other");
		mime_switch.combo.append("application/json", "JSON");
		mime_switch.combo.append("application/xml", "XML");
		mime_switch.combo.append("multipart/form-data", "Multipart");
		mime_switch.combo.append("application/x-www-form-urlencoded", "Form URL encoded");
		mime_switch.combo.append("application/yaml", "YAML");
		mime_switch.combo.set_active_id ("plain/text");
		mime_switch.show ();

		var body_notebook = new Gtk.Notebook ();
		body_notebook.append_page (body_stack, mime_switch);
		body_notebook.append_page (query_params, params_label);
		body_notebook.append_page (headers, headers_label);
		body_notebook.show ();

		body_notebook.switch_page.connect ((page, page_num) => {
			mime_switch.enabled = page_num == 0;
		});
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
		status_label.get_style_context ().add_class ("tag");
		status_label.show ();

		duration_label  = new Gtk.Label ("0ms");
		duration_label.get_style_context ().add_class ("tag");
		duration_label.show ();

		response_size_label  = new Gtk.Label ("0KB");
		response_size_label.get_style_context ().add_class ("tag");
		response_size_label.show ();

		var history_label = new Benchwell.Label (_("History"));
		history_label.show ();

		history_popover = new Benchwell.HttpHistoryPopover (history_label, null);
		history_popover.position = Gtk.PositionType.BOTTOM;
		history_label.clicked.connect ((e) => {
			history_popover.show ();
		});
		history_popover.result_activated.connect ((result) => {
			set_response (result);
		});

		details_box.get_style_context ().add_class ("responsemeta");
		details_box.pack_start (status_label, false, false, 0);
		details_box.pack_start (duration_label, false, false, 0);
		details_box.pack_start (response_size_label, false, false, 0);
		details_box.pack_end (history_label, false, false, 0);

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

		// SIGNALS
		mime_switch.combo.changed.connect (on_mime_switch_change);

		sidebar.item_activated.connect (on_item_activated);
		sidebar.item_removed.connect (on_item_removed);

		address.send_btn.btn.clicked.connect (on_send);
		address.send_btn.menu_btn.activate.connect (on_save_as);
		//address.save_btn.clicked.connect (on_save);
		window.copy_curl_action.activate.connect (on_copy_curl);

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

		headers.no_row_left.connect (() => {
			if (item == null) {
				//return new Benchwell.HttpKv ();
				return ;
			}

			try {
				var kv = item.add_header ();
				headers.add (kv);
			} catch (ConfigError err) {
				Config.show_alert (this, err.message);
			}
			//return new Benchwell.HttpKv ();
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
				Config.show_alert (this, err.message);
			}
		});

		query_params.no_row_left.connect (() => {
			if (item == null) {
				//return new Benchwell.HttpKv ();
				return;
			}

			try {
				var kv = item.add_param ();
				query_params.add (kv);
			} catch (ConfigError err) {
				Config.show_alert (this, err.message);
			}

			//return new Benchwell.HttpKv ();
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
				Config.show_alert (this, err.message);
			}
		});


		body_fields.no_row_left.connect (() => {
			if (item == null) {
				return;
			}

			try {
				var kv = item.add_form_param ();
				body_fields.add (kv);
			} catch (ConfigError err) {
				Config.show_alert (this, err.message);
			}
		});

		body_fields.row_removed.connect ((kvi) => {
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
				Config.show_alert (this, err.message);
			}
		});

		query_params.changed.connect (() => {
			address.update_url ();
		});

		headers.row_added.connect ((kv) => {
			if (kv.keyvalue.key.strip ().casefold () == "Content-Type".casefold ()) {
				body.set_language_by_mime_type (kv.keyvalue.val);
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

	private void on_mime_switch_change () {
		if (loading) {
			return;
		}

		var found = false;
		headers.get_children ().foreach ((w) => {
			var kv = w as Benchwell.KeyValue;
			if (kv.keyvalue.key.strip ().casefold () == "Content-Type".casefold ()) {
				kv.entry_val.text = mime_switch.combo.get_active_id ();
				found = true;
			}
		});

		if (!found) {
			Benchwell.HttpKv? kv = null;
			if (item != null) {
				try {
					kv = item.add_header ();
				} catch (ConfigError err) {
					Config.show_alert (this, err.message);
				}
			} else {
				kv = new Benchwell.HttpKv ();
			}

			kv.key = "Content-Type";
			kv.val = mime_switch.combo.get_active_id ();

			headers.add (kv);
		}

		if (item != null) {
			item.mime = mime_switch.combo.get_active_id ();
		}

		switch (mime_switch.combo.get_active_id ()) {
			case "application/x-www-form-urlencoded":
				body_stack.set_visible_child_name ("fields");
				body_fields.supported_types = Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE;
				break;
			case "multipart/form-data":
				body_stack.set_visible_child_name ("fields");
				body_fields.supported_types = Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE|Benchwell.KeyValueTypes.FILE;
				break;
			default:
				body_stack.set_visible_child_name ("editor");
				break;
		}

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

	private void on_copy_curl () {
		string mime = "";
		string method = "";
		string url = "";
		string[] headers = null;
		string[] body = null;

		var ok = build_request (out mime, out method, out url, out body, out headers);
		if (!ok) {
			Config.show_alert (this, _("couldn't build request"));
			return;
		}

		var builder = new StringBuilder ();
		builder.append (@"curl --request $method\\\n");
		builder.append (@"     --url $url\\\n");

		for (var i = 0; i < headers.length; i++) {
			var h = headers[i].replace ("'", "\\'");
			builder.append (@"     --header '$(h)'");
			if (i < headers.length - 1)
				builder.append (@"\\\n");
		}

		switch (mime) {
			case "application/x-www-form-urlencoded":
				for (var i = 0; i < body.length; i++) {
					var b = body[0].replace ("'", "\\'");
					builder.append (@"     --data '$(b)'\\\n");
				}
				break;
			case "multipart/form-data":
				for (var i = 0; i < body.length; i++) {
					var b = body[0].replace ("'", "\\'");
					builder.append (@"     --form '$(b)'\\\n");
				}
				break;
			default:
				var b = body[0].replace ("'", "\\'");
				builder.append (@"     --data '$(b)'\\\n");
				break;
		}

		var cmd = builder.str;
		cmd = cmd.substring (0, cmd.length - 2); // following \\\n

		var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
		cb.set_text (builder.str, builder.str.length);
	}

	// https://github.com/giuliopaci/ValaBindingsDevelopment/blob/master/libcurl-example.vala
	private void on_send () {
		string mime;
		string method;
		string url;
		string[] raw_body;
		string[] headers_array;

		response_headers.get_buffer ().set_text ("", 0);
		response.get_buffer ().set_text ("", 0);

		var ok = build_request (out mime, out method, out url, out raw_body, out headers_array);
		if (!ok) {
			Config.show_alert (this, _("couldn't build request"));
			return;
		}

		// curl headers
		Curl.SList headers = null;
		foreach(var h in headers_array) {
			headers = Curl.SList.append ((owned) headers, h);
		}

		// curl BODY
		string boundary = "";
		string curl_body = "";
		string content_type = "";
		switch (mime) {
			case "application/x-www-form-urlencoded":
				var body_builder = new StringBuilder ();

				for (var i = 0; i < raw_body.length; i++) {
					body_builder.append (raw_body[i]);
					if (i < raw_body.length - 1)
						body_builder.append ("&");
				}

				curl_body = body_builder.str;
				break;
			case "multipart/form-data":
				boundary = RandomString(60);
				content_type = @"multipart/form-data; boundary=$(boundary)";
				var body_builder = new StringBuilder ();

				body_builder.append ("--");
				body_builder.append (boundary);
				for (var i = 0; i < raw_body.length; i++) {
					body_builder.append ("\n");

					var parts = kvrg.split (raw_body[i]);
					var key = parts[1];
					var val = parts[3];

					if (parts[2] == "@") {
						var file = File.new_for_path (val);
						if (!file.query_exists ()) {
							Config.show_alert (this, @"$(val) doesn't exist");
							return;
						}
						if (file.query_file_type (0) == FileType.DIRECTORY) {
							Config.show_alert (this, @"$(val) is a directory");
							return;
						}

						GLib.FileInfo file_info = null;
						uint8[] content;
						try {
							file_info = file.query_info ("*", FileQueryInfoFlags.NONE);

							ok = GLib.FileUtils.get_data (val, out content);
							if (!ok) {
								Config.show_alert (this, @"could not read $(val)");
								return;
							}
						} catch (GLib.Error err) {
							Config.show_alert (this, err.message);
							return;
						}

						body_builder.append ("Content-Disposition: form-data; name=\"").
							append (key).
							append ("; filename=\"").
							append (file.get_basename ()).
							append ("\"\n").
							append ("Content-Type: ").
							append (file_info.get_content_type ()).
							append ("\n").
							append ("Content-Transfer-Encoding: base64\n\n").
							append (GLib.Base64.encode (content)).
							append ("\n");

						break;
					} else {
						body_builder.append ("Content-Disposition: form-data; name=\"").
							append (key).
							append ("\"\n").
							append ("Content-Type: text/plain\n\n").
							//append (handle.escape (val, val.length)).
							append (val).
							append ("\n");
						break;
					}

					body_builder.append ("--");
					body_builder.append (boundary);
				}
				body_builder.append ("--\n");

				curl_body = body_builder.str;
				break;
			default:
				curl_body = body.get_text ();
				curl_body = Config.environment.interpolate (curl_body);
				break;
		}

		perform.begin (method, url, curl_body, (owned) headers, (obj, res) => {
			try {
				var result = perform.end (res);
				switch (result.code) {
					case Curl.Code.OK:
						item.save_response (result);
						set_response (result);
						break;
					default:
						Config.show_alert (this, Curl.Global.strerror (result.code));
						break;
				}
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
			}
			overlay.stop ();
		});
	}

	private bool build_request (
		out string mime,
		out string method,
		out string url,
		out string[] raw_body,
		out string[] headers_array
	) {
		headers_array = null;
		raw_body = null;
		mime = mime_switch.combo.get_active_id ();

		var handle = new Curl.EasyHandle ();

		url = Config.environment.interpolate (address.address.get_text ());
		method = address.method_combo.get_active_id ();

		var kv_params = query_params.get_kvs ();

		// TODO: Improve URL parsing
		var builder = new StringBuilder ();
		if (url.index_of ("?") == -1 && kv_params.length > 0)
			builder.append ("?");

		for (var i = 0; i < kv_params.length; i++) {
			var key = Config.environment.interpolate (kv_params[i].key);
			var val = Config.environment.interpolate (kv_params[i].val);
			key = handle.escape (key, key.length);
			val = handle.escape (val, val.length);
			builder.append (@"$key=$val");
			if (i < kv_params.length - 1)
				builder.append ("&");
		}

		if (url.index_of ("?") != -1 && !url.has_suffix("&"))
			builder.prepend ("&");
		url += builder.str;

		var kv_headers = headers.get_kvs ();

		// NOTE: can't tranfer ownership of null
		string[] l_headers_array = null;
		l_headers_array += "X-Powered-by: Benchwell";

		var content_type = "";
		for (var i = 0; i < kv_headers.length; i++) {
			var key = Config.environment.interpolate (kv_headers[i].key);
			var val = Config.environment.interpolate (kv_headers[i].val);
			// NOTE: delayed appending content-type because multipart/form-data needs to include the bounday
			if (key.strip ().casefold () == "Content-Type".casefold ()) {
				content_type = val;
				continue;
			}
			l_headers_array += @"$(kv_headers[i].key): $(val)";
		}

		// BODY
		string boundary = "";
		string[] l_raw_body = null;
		switch (mime_switch.combo.get_active_id ()) {
			case "application/x-www-form-urlencoded":
				var form_fields = body_fields.get_kvs ();

				for (var i = 0; i < form_fields.length; i++) {
					var key = Config.environment.interpolate (form_fields[i].key);
					var val = Config.environment.interpolate (form_fields[i].val);
					key = handle.escape (key, key.length);
					val = handle.escape (val, val.length);
					l_raw_body += @"$key=$val";
				}

				break;
			case "multipart/form-data":
				var form_fields = body_fields.get_kvs ();

				for (var i = 0; i < form_fields.length; i++) {
					var key = Config.environment.interpolate (form_fields[i].key);
					var val = Config.environment.interpolate (form_fields[i].val);

					switch (form_fields[i].kvtype) {
						case Benchwell.KeyValueTypes.FILE:
							var file = File.new_for_path (val);
							if (!file.query_exists ()) {
								Config.show_alert (this, @"$(val) doesn't exist");
								return false;
							}
							if (file.query_file_type (0) == FileType.DIRECTORY) {
								Config.show_alert (this, @"$(val) is a directory");
								return false;
							}

							l_raw_body += @"$key=@$val";
							break;
						default:
							l_raw_body += @"$key=$val";
							break;
					}
				}

				break;
			default:
				l_raw_body += Config.environment.interpolate (body.get_text ());
				break;
		}

		if (content_type != "")
			l_headers_array += @"Content-Type: $(content_type)";
		headers_array = l_headers_array;
		raw_body = l_raw_body;
		return true;
	}

	private async HttpResult? perform (string method, string url, string body, owned Curl.SList headers) {
		overlay.start ();
		bool canceled = false;
		var cancel_handler_id = overlay.cancel.connect (() => {
			canceled = true;
		});

		SourceFunc callback = perform.callback;
		HttpResult? result = new HttpResult ();
		result.method = method;
		result.url = url;

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

			buffer_s2 tmp_body;
			tmp_body = buffer_s2 () {
				buffer = new uchar[0],
				size_left = Posix.strlen (body) //+1
			};
			for (var i = 0; i < body.length; i++){
				tmp_body.buffer += body[i];
			}
			//tmp_body.buffer += 0;
			handle.setopt (Curl.Option.READFUNCTION, WriteRequestCallback);
			handle.setopt (Curl.Option.READDATA, ref tmp_body);

			var then = get_real_time ();
			result.code = handle.perform ();
			var now = get_real_time ();

			if (result.code == Curl.Code.OK) {
				handle.getinfo(Curl.Info.RESPONSE_CODE, out result.status);
				result.body = (string)tmp.buffer;
				result.duration = now - then;
				var headers_string = "";
				resp_headers.foreach ((key, val) => {
					headers_string += @"$key: $val";
					if (key.casefold () == "Content-Type".casefold())
						result.content_type = val.split(";")[0];
				});
				result.headers = headers_string;
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
		loading = true;
		Config.settings.http_item_id = item.id;

		this.item = item;
		this.item_iter = iter;
		try {
			this.item.load_full_item ();
		} catch (ConfigError err) {
			Config.show_alert (this, err.message);
			loading = false;
			return;
		}

		history_popover.disconnect_from_item ();
		history_popover.item = this.item;
		title = this.item.name;

		mime_switch.combo.set_active_id (item.mime);

		normalize_item (this.item);

		address.set_request (this.item);
		build_interpolated_label ();

		this.item.touch_without_save (() => {
			if (this.item.body != null) {
				body.get_buffer ().text = this.item.body;
			} else {
				body.get_buffer ().text = "";
			}
		});

		set_response (item.last_response ());

		headers.clear ();
		query_params.clear ();
		body_fields.clear ();

		foreach (var h in this.item.headers) {
			headers.add ((Benchwell.KeyValueI) h);
		}

		foreach (var h in this.item.form_params) {
			body_fields.add ((Benchwell.KeyValueI) h);
		}

		load_query_params ();

		try {
			if (this.item.headers.length == 0) {
				headers.add (item.add_header ());
			}

			if (this.item.query_params.length == 0) {
				query_params.add (this.item.add_param ());
			}

			if (this.item.form_params.length == 0) {
				body_fields.add (this.item.add_form_param ());
			}
		} catch (ConfigError err) {
			Config.show_alert (this, err.message);
		}

		if (this.item.mime == "application/x-www-form-urlencoded" ||
			this.item.mime == "multipart/form-data") {
			body_stack.set_visible_child_name ("fields");
		} else {
			body_stack.set_visible_child_name ("editor");
		}
		loading = false;
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

	public void set_response (HttpResult? result) {
		if (result == null)
			return;

		if (result.body == null)
			result.body = "";

		response_headers.get_buffer ().text = result.headers;

		var formatted_body = "";
		Gtk.TextIter start_iter;
		response.get_buffer ().get_start_iter (out start_iter);
		if ( result.body != "" ) {
			switch (result.content_type) {
				case "application/json":
					var json = new Json.Parser ();
					var generator = new Json.Generator ();
					generator.indent = 4;
					generator.pretty = true;
					try {
						json.load_from_data (result.body);
						generator.set_root (json.get_root ());

						formatted_body = generator.to_data (null);
						response.get_buffer ().insert_text (ref start_iter, formatted_body, formatted_body.length);
						response.set_language ("json");
					} catch (GLib.Error err) {
						Config.show_alert (this, err.message);
						formatted_body = result.body;
						response.get_buffer ().insert_text (ref start_iter, formatted_body, formatted_body.length);
					}
					break;
				default:
					response.set_language (null);
					formatted_body = result.body;
					response.get_buffer ().insert_text (ref start_iter, formatted_body, formatted_body.length);
					break;
			}
		} else {
			Gtk.TextIter end_iter;
			response.get_buffer ().get_end_iter (out end_iter);
			response.get_buffer ().delete (ref start_iter, ref end_iter);
		}

		if (result.body.length < 1024) {
			response_size_label.set_text (@"$(result.body.length)B");
		} else {
			response_size_label.set_text (@"$(result.body.length/1024)kB");
		}

		var code = Benchwell.Http.CODES.parse (result.status);
		if (code == null) {
			status_label.set_text (@"$(result.status)");
		} else {
			status_label.set_text (@"$(result.status) $(code)");
		}

		status_label.get_style_context ().remove_class("ok");
		status_label.get_style_context ().remove_class("warning");
		status_label.get_style_context ().remove_class("bad");

		if (result.status >= 200 && result.status < 300) {
			status_label.get_style_context ().add_class("ok");
		}
		if (result.status >= 300 && result.status < 500) {
			status_label.get_style_context ().add_class("warning");
		}
		if (result.status >= 500) {
			status_label.get_style_context ().add_class("bad");
		}

		duration_label.set_text (@"$(result.duration/1000)ms");
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
					Config.show_alert (this, err.message);
				}
			}
		});

		//try {
			//item.save ();
		//} catch (ConfigError err) {
			//Config.show_alert (this, err.message);
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

public class Benchwell.HttpHistoryPopover : Gtk.Popover {
	private Gtk.Grid grid;
	public weak Benchwell.HttpItem? item { owned get; set; }

	public signal void result_activated (Benchwell.HttpResult result);

	public HttpHistoryPopover (Gtk.Widget relative_to, Benchwell.HttpItem? item) {
		Object (
			relative_to: relative_to,
			item: item
		);

		build_results ();
		notify["item"].connect ((sender, property) => {
			if (this.item == null)
				return;
			this.item.response_added.connect (build_results);
			build_results ();
		});
	}

	public void disconnect_from_item () {
		if (item == null)
			return;
		item.response_added.disconnect (build_results);
	}

	private void build_results () {
		if (item == null) {
			return;
		}

		if (grid != null) {
			remove (grid);
			grid.destroy ();
		}

		grid = new Gtk.Grid ();
		grid.row_spacing = 5;
		grid.column_spacing = 5;
		add (grid);

		var i = 0;
		foreach (var response in item.responses) {
			add_response (i, response);
			i++;
		}

		grid.show_all ();
	}

	private void add_response (int at, Benchwell.HttpResult response) {
		var status_label = new Benchwell.Label ("0");
		status_label.set_text (response.status.to_string ());
		status_label.get_style_context ().add_class ("tag");

		if (response.status >= 200 && response.status < 300) {
			status_label.get_style_context ().add_class("ok");
		}
		if (response.status >= 300 && response.status < 500) {
			status_label.get_style_context ().add_class("warning");
		}
		if (response.status >= 500) {
			status_label.get_style_context ().add_class("bad");
		}

		var duration_label  = new Benchwell.Label (@"$(response.duration/1000)ms");
		duration_label.get_style_context ().add_class ("tag");

		var size_label  = new Benchwell.Label ("0KB");
		size_label.get_style_context ().add_class ("tag");
		if (response.body.length < 1024) {
			size_label.set_text (@"$(response.body.length)B");
		} else {
			size_label.set_text (@"$(response.body.length/1024)kB");
		}

		var url_label = new Benchwell.Label (@"$(response.method) $(response.url)");
		url_label.max_width_chars = 30;
		url_label.ellipsize = Pango.EllipsizeMode.END;
		url_label.tooltip_text = response.url;

		var time_fmt = "%Y-%m-%d %H:%M:%S";
		var now = new DateTime.now_local ();
		if (now.get_day_of_year () == response.created_at.get_day_of_year ())
			time_fmt = "%H:%M:%S";

		var time_label = new Benchwell.Label (response.created_at.format (time_fmt));

		grid.attach (status_label, 0, at, 1, 1);
		grid.attach (duration_label, 1, at, 1, 1);
		grid.attach (size_label, 2, at, 1, 1);
		grid.attach (url_label, 3, at, 1, 1);
		grid.attach (time_label, 4, at, 1, 1);

		status_label.clicked.connect ( () => { result_activated (response); });
		duration_label.clicked.connect ( () => { result_activated (response); });
		size_label.clicked.connect ( () => { result_activated (response); });
		url_label.clicked.connect ( () => { result_activated (response); });
		time_label.clicked.connect ( () => { result_activated (response); });
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
		var key = header[0:at].strip ().casefold ();
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
