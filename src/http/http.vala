enum Benchwell.HttpColumns {
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
	public Benchwell.Http.KeyValues  headers;
	public Benchwell.Http.KeyValues  query_params;
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

		headers = new Benchwell.Http.KeyValues ();
		headers.show ();

		query_params = new Benchwell.Http.KeyValues ();
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

		address.send_btn.clicked.connect (() => {
			var url = address.address.get_text ();
			var method = address.method_combo.get_active_id ();

			var session = new Soup.Session ();
			var message = new Soup.Message (method, url);
			if (message == null) {
				return;
			}

			session.send_message (message);

			response.get_buffer ().set_text ((string) message.response_body.flatten ().data);
		});
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

public class Benchwell.Http.HttpSideBar : Gtk.Box {
	public Gtk.TreeView treeview;
	public Gtk.TreeStore store;
	public Gtk.ComboBoxText collections_combo;

	public Gtk.Menu menu;
	public Gtk.MenuItem new_request_menu;
	public Gtk.MenuItem new_folder_menu;
	public Gtk.MenuItem delete_menu;
	public Gtk.MenuItem edit_menu;

	public weak Benchwell.HttpCollection? selected_collection;
	public weak Benchwell.HttpItem? selected_item;

	public signal void load_request(unowned Benchwell.HttpItem item);

	public HttpSideBar () {
		Object (
			orientation: Gtk.Orientation.VERTICAL
		);

		// treeview
		treeview = new Gtk.TreeView ();
		treeview.headers_visible = false;
		treeview.show_expanders = false;
		treeview.enable_tree_lines = false;
		treeview.button_release_event.connect (on_button_release_event);
		//var selection = treeview.get_selection ();
		//selection.changed.connect (on_selection_changed);

		var image_column = new Gtk.TreeViewColumn.with_attributes("image", new Gtk.CellRendererPixbuf (), "pixbuf", Benchwell.HttpColumns.ICON);

		var name_renderer = new Gtk.CellRendererText ();
		var name_column = new Gtk.TreeViewColumn.with_attributes("name", name_renderer, "text", Benchwell.HttpColumns.TEXT);

		var method_renderer = new Gtk.CellRendererText ();
		var method_column = new Gtk.TreeViewColumn.with_attributes("method", method_renderer, "text", Benchwell.HttpColumns.METHOD);
		method_column.set_cell_data_func (method_renderer, (cell_layout, cell, tree_model, iter) => {
			GLib.Value val;
			tree_model.get_value (iter, Benchwell.HttpColumns.ITEM, out val);
			var item = val.get_object () as Benchwell.HttpItem;

			if (!item.is_folder) {
				var color = Benchwell.Colors.parse (item.method);
				if (color == null) {
					return;
				}
				cell.set_property ("markup", @"<span foreground=\"$color\">$(item.method)</span>");
			}
		});

		treeview.append_column (image_column);
		treeview.append_column (name_column);
		treeview.append_column (method_column);
		treeview.expander_column = image_column;

		store = new Gtk.TreeStore (4, GLib.Type.OBJECT, GLib.Type.OBJECT, GLib.Type.STRING, GLib.Type.STRING);
		treeview.set_model (store);
		treeview.show ();
		var treeview_sw = new Gtk.ScrolledWindow (null, null);
		treeview_sw.add (treeview);
		treeview_sw.show ();
		///////////

		collections_combo = new Gtk.ComboBoxText ();
		collections_combo.append ("", "");
		foreach (var collection in Config.http_collections) {
			collections_combo.append (collection.id.to_string (), collection.name);
		}
		collections_combo.show ();

		menu = new Gtk.Menu ();

		new_folder_menu = new Benchwell.MenuItem(_("New folder"), "add");
		new_folder_menu.show ();

		new_request_menu = new Benchwell.MenuItem(_("New request"), "add");
		new_request_menu.show ();

		delete_menu = new Benchwell.MenuItem(_("Delete"), "close");
		delete_menu.show ();

		edit_menu = new Benchwell.MenuItem(_("Edit"), "config");
		edit_menu.show ();

		menu.add (new_request_menu);
		menu.add (new_folder_menu);
		menu.add (edit_menu);
		menu.add (delete_menu);

		pack_start( collections_combo, false, false, 0);
		pack_start( treeview_sw, true, true, 0);

		collections_combo.changed.connect (on_collection_selected);
		treeview.row_activated.connect (on_row_activated);
	}

	private void on_row_activated () {
		Gtk.TreeIter iter;
		treeview.get_selection ().get_selected (null, out iter);

		GLib.Value val;
		store.get_value (iter, Benchwell.HttpColumns.ITEM, out val);
		var selected_item = val.get_object () as Benchwell.HttpItem;

		if (selected_item == null) {
			return;
		}

		Config.load_full_item (selected_item);

		if (selected_item.is_folder) {
			// unload items
			// // collapse row
			var selected_path = store.get_path (iter);
			if (treeview.is_row_expanded (selected_path)) {
				treeview.collapse_row (selected_path);
				return;
			}

			// // remove items
			Gtk.TreeIter child;
			store.iter_children (out child, iter);
			while (store.iter_is_valid (child)) {
				store.remove(ref child);
			}

			build_tree (iter, selected_item.items);
			treeview.expand_row (selected_path, false);
			return;
		}

		load_request (selected_item);
	}

	private void on_collection_selected () {
		var collection_id = int64.parse (collections_combo.get_active_id ());

		foreach(var collection in Config.http_collections) {
			if (collection.id != collection_id){
				continue;
			}

			try {
				Config.load_root_items (collection);
			} catch (Benchwell.Backend.Sql.Error err) {
				//result_view.show_alert (err.message);
				stderr.printf (err.message);
				return;
			}

			load_collection (collection);

			selected_collection = collection;

			break;
		}
	}

	private void load_collection (Benchwell.HttpCollection collection) {
		store.clear ();
		build_tree (null, collection.items);
	}

	private void build_tree (Gtk.TreeIter? iter, unowned Benchwell.HttpItem[] items) {
		foreach (var item in items) {
			var this_iter = add_row (iter, item);

			if (item.is_folder) {
				build_tree (this_iter, item.items);
			}
		}
	}

	private Gtk.TreeIter add_row (Gtk.TreeIter? parent_iter, Benchwell.HttpItem item) {
		Gtk.TreeIter iter;
		store.append (out iter, parent_iter);

		if (item.is_folder) {
			var px = Gtk.IconTheme.get_default ().load_icon ("bw-directory", Gtk.IconSize.BUTTON, Gtk.IconLookupFlags.USE_BUILTIN);
			store.set_value (iter, Benchwell.HttpColumns.ICON, px);
		} else {
			store.set_value (iter, Benchwell.HttpColumns.METHOD, item.method);
		}

		store.set_value (iter, Benchwell.HttpColumns.TEXT, item.name);
		store.set_value (iter, Benchwell.HttpColumns.ITEM, item);

		return iter;
	}

	private void on_selection_changed () {
	}

	private bool on_button_release_event (Gtk.Widget w, Gdk.EventButton event) {
		if (event.button != Gdk.BUTTON_SECONDARY) {
			return false;
		}

		Gtk.TreePath path;
		treeview.get_path_at_pos ((int) event.x, (int) event.y , out path, null, null, null);

		treeview.get_selection ().select_path (path);

		menu.popup_at_pointer (event);
		return true;
	}
}

public class Benchwell.Http.HttpAddressBar : Gtk.Box {
	public Gtk.ComboBoxText method_combo;
	public Gtk.Entry address;
	public Gtk.Button send_btn;
	public Benchwell.OptButton save_btn;

	public HttpAddressBar () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		method_combo = new Gtk.ComboBoxText ();
		for (var i = 0; i < Benchwell.Methods.length; i++) {
			method_combo.append (Benchwell.Methods[i], Benchwell.Methods[i]);
		}
		method_combo.set_active (0);
		method_combo.show ();

		address = new Gtk.Entry ();
		address.placeholder_text = "http://localhost/path.json";
		address.show ();

		send_btn = new Gtk.Button.with_label (_("SEND"));
		send_btn.get_style_context ().add_class ("suggested-action");
		send_btn.show ();

		// TODO: add to window
		var save_as_action = new GLib.SimpleAction ("win.saveas", null);
		save_btn = new Benchwell.OptButton(_("SAVE"), _("Save as"), "win.saveas");
		save_btn.show ();

		pack_start(method_combo, false, false, 0);
		pack_start(address, true, true, 0);
		pack_end(save_btn, false, false, 0);
		pack_end(send_btn, false, false, 0);
	}

	public void set_request (unowned Benchwell.HttpItem item) {
		address.set_text (item.url);
		method_combo.set_active_id (item.method);
	}
}

public class Benchwell.Http.KeyValues : Gtk.Box {
	public KeyValues () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);

		add (null);
	}

	public void add (KeyValueI? kvi) {
		var kv = new Benchwell.Http.KeyValue ();
		kv.show ();

		if (kvi != null) {
			kv.key.set_text (kvi.key ());
			kv.val.set_text (kvi.val ());
			kv.enabled.set_active (kvi.enabled ());
		}

		kv.key.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (null);
		});
		kv.val.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (null);
		});
		kv.enabled.state_set.connect ((b) => {
			if (get_children ().index (kv) == get_children ().length () - 1) {
				add (null);
			}

			return false;
		});

		kv.remove_btn.clicked.connect( () => {
			remove(kv);
			if (get_children ().length () == 0) {
				add (null);
			}
		});

		pack_start (kv, false, false, 0);
	}

	public void clear () {
		get_children ().foreach ( (c) => {
			remove (c);
		});
	}
}

public class Benchwell.Http.KeyValue : Gtk.Box {
	public Gtk.Switch enabled;
	public Gtk.Entry        key;
	public Gtk.Entry        val;
	public Benchwell.Button remove_btn;

	public KeyValue () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		key = new Gtk.Entry ();
		key.placeholder_text = _("Name");
		key.show ();

		val = new Gtk.Entry ();
		val.placeholder_text = _("Value");
		val.show ();

		remove_btn = new Benchwell.Button ("close", Gtk.IconSize.BUTTON);
		remove_btn.show ();

		enabled = new Gtk.Switch ();
		enabled.valign = Gtk.Align.CENTER;
		enabled.vexpand = false;
		enabled.set_active (true);
		enabled.show ();

		pack_start (key, true, true, 0);
		pack_start (val, true, true, 0);
		pack_end (remove_btn, false, false, 0);
		pack_end (enabled, false, false, 5);
	}
}
