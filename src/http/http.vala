public class Benchwell.Http.Result : Object {
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

public class Benchwell.Http.Http : Gtk.Paned {
	public Benchwell.ApplicationWindow   window { get; construct; }
	public string                        title  { get; set; }
	public Benchwell.Http.SideBar    sidebar;
	public Benchwell.Http.AddressBar address;
	public Benchwell.Http.Overlay overlay;

	// request
	public Gtk.Stack body_stack;
	public Benchwell.SourceView body;
	public Benchwell.KeyValues  body_fields;
	public Benchwell.ComboTab mime_switch;
	//public Gtk.ComboBoxText     mime;
	public Benchwell.KeyValues  headers;
	public Benchwell.KeyValues  query_params;
	public Gtk.TextView  note;
	//////////

	// response
	public Benchwell.SourceView response;
	public Benchwell.SourceView response_headers;
	public Gtk.Label            status_label;
	public Gtk.Label            duration_label;
	public Gtk.Label            response_size_label;
	public Benchwell.Http.HistoryPopover history_popover;
	///////////

	public Benchwell.HttpItem? item;

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

		sidebar = new Benchwell.Http.SideBar (window);
		sidebar.show ();

		address = new Benchwell.Http.AddressBar ();
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

		headers = new Benchwell.KeyValues (Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE);
		headers.show ();

		var headersw = new Gtk.ScrolledWindow (null, null);
		headersw.add (headers);
		headersw.show ();

		query_params = new Benchwell.KeyValues (Benchwell.KeyValueTypes.STRING|Benchwell.KeyValueTypes.MULTILINE);
		query_params.show ();
		var query_paramsw = new Gtk.ScrolledWindow (null, null);
		query_paramsw.add (query_params);
		query_paramsw.show ();

		var params_label = new Gtk.Label (_("Params"));
		params_label.show ();

		var headers_label = new Gtk.Label (_("Headers"));
		headers_label.show ();

		note = new Gtk.TextView ();
		note.show ();

		var notesw = new Gtk.ScrolledWindow (null, null);
		notesw.add (note);
		notesw.show ();

		var note_label = new Gtk.Label (_("Note"));
		note_label.show ();


		mime_switch = new Benchwell.ComboTab (_("Body"), true);
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
		body_notebook.append_page (query_paramsw, params_label);
		body_notebook.append_page (headersw, headers_label);
		body_notebook.append_page (notesw, note_label);
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

		history_popover = new Benchwell.Http.HistoryPopover (history_label, null);
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
		ws_paned.pack1 (body_notebook, true, false);
		ws_paned.pack2 (response_box, true, true);

		var ws_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		ws_box.vexpand = true;
		ws_box.hexpand = true;
		ws_box.show ();

		ws_box.pack_start (address, false, false, 0);
		ws_box.pack_start (ws_paned, true, true, 0);

		overlay = new Benchwell.Http.Overlay ();
		overlay.add (ws_box);
		overlay.show ();

		pack1 (sidebar, false, true);
		pack2 (overlay, false, true);

		// SIGNALS
		mime_switch.combo.changed.connect (on_mime_switch_change);

		sidebar.item_activated.connect (on_item_activated);
		sidebar.item_removed.connect (on_item_removed);

		address.send_btn.btn.clicked.connect (on_send);
		//address.send_btn.menu_btn.activate.connect (on_save_as);
		window.copy_curl_action.activate.connect (on_copy_curl);
		window.saveas.activate.connect (on_save_as);

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

			var kv = kvi as Benchwell.KeyValue;
			if (kv == null) {
				return;
			}
			try {
				(kv.keyvalue as HttpKv).delete ();
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

			var kv = kvi as Benchwell.KeyValue;
			if (kv == null) {
				return;
			}
			try {
				(kv.keyvalue as HttpKv).delete ();
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

			var kv = kvi as Benchwell.KeyValue;
			if (kv == null) {
				return;
			}
			try {
				(kv.keyvalue as HttpKv).delete ();
			} catch (ConfigError err) {
				Config.show_alert (this, err.message);
			}
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

		item = new Benchwell.HttpItem ();
		headers.add (item.add_header ());
		query_params.add (item.add_param ());
		note.buffer.changed.connect (() => {
			item.description = note.buffer.text;
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
		sidebar.touch (item);
	}

	private void on_save_as () {
		var selector = new SaveAsSelector (window);
		selector.show ();
		var dialog = new Gtk.Dialog.with_buttons (_("Save As"), window,
									Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL);
		var ok_button = dialog.add_button (_("Ok"), Gtk.ResponseType.OK);
		dialog.add_button (_("Cancel"), Gtk.ResponseType.CANCEL);
		dialog.set_default_size (250, 130);
		dialog.get_content_area ().spacing = 5;
		dialog.get_content_area ().add (selector);
		ok_button.sensitive = false;
		selector.changed.connect (() => {
			ok_button.sensitive = selector.collection != null && selector.get_name () != "";
		});

		var resp = (Gtk.ResponseType) dialog.run ();
		if (resp != Gtk.ResponseType.OK) {
			dialog.destroy ();
			return;
		}

		int64? collection_id = null;
		if (selector.collection != null) {
			collection_id = selector.collection.id;
		}
		int64? folder_id = null;
		if (selector.folder != null) {
			folder_id = selector.folder.id;
		}
		var name = selector.get_name ();
		dialog.destroy ();

		var new_item = new Benchwell.HttpItem ();
		new_item.touch_without_save (() => {
			if (folder_id != null)
				new_item.parent_id = folder_id;
			new_item.name = name;
			new_item.http_collection_id = collection_id;
			new_item.body = item.body;

			new_item.description = item.description;
			new_item.method = item.method;
			new_item.url = item.url;
			new_item.body = item.body;
			new_item.mime = item.mime;
		});

		var collection = Config.get_http_collection_by_id (collection_id);
		try {
			collection.add_item (new_item);
			foreach (var h in item.headers)
				new_item.add_header (h.key, h.val);
			foreach (var p in item.query_params)
				new_item.add_param (p.key, p.val);
			foreach (var p in item.form_params)
				new_item.add_form_param (p.key, p.val, p.kvtype);
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (this, err.message);
		}

		dialog.destroy ();
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
		builder.append (@"curl --request $method \\\n");
		builder.append (@"     --url $url \\\n");

		for (var i = 0; i < headers.length; i++) {
			var h = headers[i].replace ("'", "\\'");
			builder.append (@"     --header '$(h)'");
			if (i < headers.length-1)
				builder.append (@" \\\n");
		}

		if (body.length > 0) {
			builder.append ("\n");
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
		}

		var cmd = builder.str;
		cmd = cmd.substring (0, cmd.length - 2); // following \\\n

		var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
		cb.set_text (builder.str, builder.str.length);
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

		url = Config.environments.selected.interpolate (address.address.get_text ());
		method = address.method_combo.get_active_id ();

		var kv_params = query_params.get_kvs ();

		// TODO: Improve URL parsing
		var builder = new StringBuilder ();
		if (url.index_of ("?") == -1 && kv_params.length > 0)
			builder.append ("?");

		for (var i = 0; i < kv_params.length; i++) {
			var key = Config.environments.selected.interpolate (kv_params[i].key);
			var val = Config.environments.selected.interpolate (kv_params[i].val);
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
		l_headers_array += "User-Agent: Benchwell";

		var content_type = "";
		for (var i = 0; i < kv_headers.length; i++) {
			var key = Config.environments.selected.interpolate (kv_headers[i].key);
			var val = Config.environments.selected.interpolate (kv_headers[i].val);
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
					var key = Config.environments.selected.interpolate (form_fields[i].key);
					var val = Config.environments.selected.interpolate (form_fields[i].val);
					key = handle.escape (key, key.length);
					val = handle.escape (val, val.length);
					l_raw_body += @"$key=$val";
				}

				break;
			case "multipart/form-data":
				var form_fields = body_fields.get_kvs ();

				for (var i = 0; i < form_fields.length; i++) {
					var key = Config.environments.selected.interpolate (form_fields[i].key);
					var val = Config.environments.selected.interpolate (form_fields[i].val);

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
				l_raw_body += Config.environments.selected.interpolate (body.get_text ());
				break;
		}

		if (content_type != "")
			l_headers_array += @"Content-Type: $(content_type)";
		headers_array = l_headers_array;
		raw_body = l_raw_body;
		return true;
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
				curl_body = Config.environments.selected.interpolate (curl_body);
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

	private async Result? perform (string method, string url, string body, owned Curl.SList headers) {
		overlay.start ();
		bool canceled = false;
		var cancel_handler_id = overlay.cancel.connect (() => {
			canceled = true;
		});

		SourceFunc callback = perform.callback;
		Result? result = new Result ();
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
				case "PUT":
					handle.setopt (Curl.Option.PUT, true);
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
			handle.setopt (Curl.Option.FOLLOWLOCATION, Config.settings.http_follow_redirect);
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
		note.get_buffer ().set_text (item.description == null ? "" : item.description);

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
			title = _("HTTP");

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
		if (Config.environments.selected != null) {
			interpolator = Config.environments.selected.dry_interpolate;
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

	public void set_response (Result? result) {
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

[Compact]
private struct buffer_s {
	uchar[] buffer;
}

[Compact]
private struct buffer_s2 {
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
