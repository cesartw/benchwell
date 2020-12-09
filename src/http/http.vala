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

	public Http(Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			wide_handle: true
		);

		env = Config.environments.nth_data (0);

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

		address.send_btn.clicked.connect (on_send);
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
		address.set_request (item);
		body.get_buffer ().set_text (item.body);
		mime.set_active_id (item.mime);

		response.get_buffer ().set_text ("");
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

	public void set_response(uint status, string raw_data, string content_type, int64 duration) {
		print(content_type);

		switch (content_type) {
			case "application/json":
				var json = new Json.Parser ();
				var generator = new Json.Generator ();

				generator.set_pretty (true);
				json.load_from_data (raw_data);
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
}
