public class Benchwell.Views.Http : Gtk.Paned {
	public Benchwell.ApplicationWindow    window { get; construct; }
	public string                         title;
	public Benchwell.Views.HttpSideBar    sidebar;
	public Benchwell.Views.HttpAddressBar address;

	// request
	public Benchwell.SourceView body;
	public Gtk.ComboBoxText           mime;
	public Gtk.Label                  body_size;
	public Benchwell.Views.KeyValues  headers;
	public Benchwell.Views.KeyValues  query_params;
	//////////

	// response
	public Benchwell.SourceView response;
	public Gtk.Label            status;
	public Gtk.Label            duration;
	public Gtk.Label            respSize;
	///////////

	public Http(Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			orientation: Gtk.Orientation.HORIZONTAL,
			wide_handle: true
		);

		title = _("HTTP");

		sidebar = new Benchwell.Views.HttpSideBar ();
		sidebar.show ();

		address = new Benchwell.Views.HttpAddressBar ();
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

		headers = new Benchwell.Views.KeyValues ();
		headers.show ();

		query_params = new Benchwell.Views.KeyValues ();
		query_params.show ();

		mime = new Gtk.ComboBoxText ();
		mime.append ("json", "JSON");
		mime.append ("plain", "PLAIN");
		mime.append ("html", "HTML");
		mime.append ("xml", "XML");
		mime.append ("yaml", "YAML");
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

		// request
		response  = new Benchwell.SourceView ();
		status    = new Gtk.Label ("200 OK");
		duration  = new Gtk.Label ("0ms");
		respSize  = new Gtk.Label ("0KB");
		//////////

		var ws_paned = new Gtk.Paned (Gtk.Orientation.VERTICAL);
		ws_paned.vexpand = true;
		ws_paned.hexpand = true;
		ws_paned.wide_handle = true;
		ws_paned.show ();
		ws_paned.add1 (body_notebook);

		var ws_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		ws_box.vexpand = true;
		ws_box.hexpand = true;
		ws_box.show ();

		ws_box.pack_start (address, false, false, 0);
		ws_box.pack_start (ws_paned, true, true, 0);

		pack1 (sidebar, false, false);
		pack2 (ws_box, false, true);
	}
}

public class Benchwell.Views.HttpSideBar : Gtk.Box {
	public Gtk.TreeView treeview;
	public Gtk.TreeStore store;
	public Gtk.ComboBoxText collections_combo;

	public Gtk.Menu menu;
	public Gtk.MenuItem new_request_menu;
	public Gtk.MenuItem new_folder_menu;
	public Gtk.MenuItem delete_menu;
	public Gtk.MenuItem edit_menu;

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

		var image_column = new Gtk.TreeViewColumn.with_attributes("", new Gtk.CellRendererPixbuf (), "pixbuf", 0);
		treeview.append_column (image_column);
		treeview.expander_column = image_column;

		var text_renderer = new Gtk.CellRendererText ();
		var text_column = new Gtk.TreeViewColumn.with_attributes("", text_renderer, "text", 1);
		text_column.set_cell_data_func (text_renderer, (cell_layout, cell, tree_model, iter) => {
			GLib.Value val;
			tree_model.get_value (iter, 1, out val);
			var path = tree_model.get_path (iter);

			if ( val.holds (GLib.Type.STRING) && val.get_string () != "") {
				var color = Benchwell.Colors.parse (val.get_string ());
				cell.set_property ("markup", @"<span foreground=\"$color\">$(val.get_string ())</span>");
			}
		});
		treeview.append_column (text_column);

		store = new Gtk.TreeStore (4, GLib.Type.OBJECT, GLib.Type.STRING, GLib.Type.INT64, GLib.Type.STRING);
		treeview.set_model (store);
		treeview.show ();
		var treeview_sw = new Gtk.ScrolledWindow (null, null);
		treeview_sw.add (treeview);
		///////////

		collections_combo = new Gtk.ComboBoxText ();
		collections_combo.append ("", "");
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

public class Benchwell.Views.HttpAddressBar : Gtk.Box {
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
}

public interface KeyValueI {
	public abstract string name();
	public abstract string val();
	public abstract bool is_enabled();
	public abstract void set_name(string n);
	public abstract void set_val(string v);
	public abstract void set_enabled(bool e);
}

public class Benchwell.Views.KeyValues : Gtk.Box {
	public KeyValues () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);

		add (null);
	}

	public void add (KeyValueI? kvi) {
		var kv = new Benchwell.Views.KeyValue ();
		kv.show ();

		if (kvi != null) {
			kv.key.set_text (kvi.name ());
			kv.val.set_text (kvi.val ());
			kv.enabled.set_active (kvi.is_enabled ());
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

		pack_start (kv, false, false, 0);
	}
}

public class Benchwell.Views.KeyValue : Gtk.Box {
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
