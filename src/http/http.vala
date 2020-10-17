enum Benchwell.Http.Columns {
	ITEM,
	ICON,
	TEXT,
	METHOD
}

public class Benchwell.Http.Http : Gtk.Paned {
	public Benchwell.ApplicationWindow    window { get; construct; }
	public string                         title;
	public Benchwell.Http.HttpSideBar    sidebar;
	public Benchwell.Http.HttpAddressBar address;

	// request
	public Benchwell.SourceView       body;
	public Gtk.ComboBoxText           mime;
	public Gtk.Label                  body_size;
	public Benchwell.KeyValues  headers;
	public Benchwell.KeyValues  query_params;
	//////////

	// response
	public Benchwell.SourceView response;
	public Gtk.Label            status;
	public Gtk.Label            duration;
	public Gtk.Label            response_size;
	///////////

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
			Gtk.TextIter start, end;

			buff.get_start_iter (out start);
			buff.get_end_iter (out end);
			var txt = buff.get_text (start, end, false);
			body_size.set_text (@"$(txt.length/2014)KB");
		});
		var body_sw = new Gtk.ScrolledWindow (null, null);
		body_sw.add (body);
		body_sw.show ();

		mime = new Gtk.ComboBoxText ();
		mime.show ();

		body_size = new Gtk.Label ("0KB");
		body_size.show ();

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

		status    = new Gtk.Label ("200 OK");
		status.show ();

		duration  = new Gtk.Label ("0ms");
		duration.show ();

		response_size  = new Gtk.Label ("0KB");
		response_size.show ();

		details_box.pack_start (status, false, false, 0);
		details_box.pack_start (duration, false, false, 0);
		details_box.pack_end (response_size, false, false, 0);

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

		address.send_btn.clicked.connect (on_send);
	}

	private void on_send () {
		var url = Config.environments.nth_data (0).interpolate (address.address.get_text ());
		var method = address.method_combo.get_active_id ();

		var session = new Soup.Session ();
		var message = new Soup.Message (method, url);
		if (message == null) {
			return;
		}

		session.send_message (message);

		response.get_buffer ().set_text ((string) message.response_body.flatten ().data);
	}

	private void on_load_request (unowned Benchwell.HttpItem item) {
		address.set_request (item);
		body.get_buffer ().set_text (item.body);
		mime.set_active_id (item.mime);

		headers.clear ();
		query_params.clear ();
		foreach (var h in item.headers) {
			headers.add (h);
		}
		foreach (var q in item.query_params) {
			query_params.add (q);
		}

		if (item.headers.length == 0) {
			headers.add (null);
		}

		if (item.query_params.length == 0) {
			query_params.add (null);
		}
	}
}
