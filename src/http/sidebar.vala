public class Benchwell.Http.HttpSideBar : Gtk.Box {
	public Gtk.TreeView treeview;
	public Gtk.TreeStore store;
	public Gtk.ComboBoxText collections_combo;

	public Gtk.Menu menu;
	public Gtk.MenuItem add_request_menu;
	public Gtk.MenuItem add_folder_menu;
	public Gtk.MenuItem delete_menu;
	public Gtk.MenuItem edit_menu;


	public weak Benchwell.HttpCollection? selected_collection;
	public weak Benchwell.HttpItem? selected_item;

	public signal void load_request(Benchwell.HttpItem item);

	private Gtk.CellRendererText name_renderer;
	private Gtk.TreeViewColumn name_column;

	public HttpSideBar () {
		Object (
			orientation: Gtk.Orientation.VERTICAL
		);

		// treeview
		treeview = new Gtk.TreeView ();
		treeview.headers_visible = false;
		treeview.show_expanders = false;
		treeview.enable_tree_lines = false;
		treeview.reorderable = true;
		treeview.button_release_event.connect (on_button_release_event);
		//var selection = treeview.get_selection ();
		//selection.changed.connect (on_selection_changed);

		var image_column = new Gtk.TreeViewColumn.with_attributes("image", new Gtk.CellRendererPixbuf (), "pixbuf", Benchwell.Http.Columns.ICON);

		name_renderer = new Gtk.CellRendererText ();
		name_column = new Gtk.TreeViewColumn.with_attributes("name", name_renderer, "text", Benchwell.Http.Columns.TEXT);

		var method_renderer = new Gtk.CellRendererText ();
		var method_column = new Gtk.TreeViewColumn.with_attributes("method", method_renderer, "text", Benchwell.Http.Columns.METHOD);
		method_column.set_cell_data_func (method_renderer, (cell_layout, cell, tree_model, iter) => {
			GLib.Value val;
			tree_model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val.get_object () as Benchwell.HttpItem;
			if (item == null) {
				return;
			}

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

		add_folder_menu = new Benchwell.MenuItem(_("New folder"), "add");
		add_folder_menu.show ();

		add_request_menu = new Benchwell.MenuItem(_("New request"), "add");
		add_request_menu.show ();

		delete_menu = new Benchwell.MenuItem(_("Delete"), "close");
		delete_menu.show ();

		edit_menu = new Benchwell.MenuItem(_("Edit"), "config");
		edit_menu.show ();

		menu.add (add_request_menu);
		menu.add (add_folder_menu);
		menu.add (edit_menu);
		menu.add (delete_menu);

		pack_start (collections_combo, false, false, 0);
		pack_start (treeview_sw, true, true, 0);

		collections_combo.changed.connect (on_collection_selected);
		treeview.row_activated.connect (on_load_item);

		var selected_collection_id = Config.settings.get_int64 (Benchwell.Settings.HTTP_COLLECTION_ID.to_string ());
		if (selected_collection_id > 0) {
			collections_combo.set_active_id (selected_collection_id.to_string ());
		}

		// SIGNALS
		// edit hack
		edit_menu.activate.connect ( () => {
			Gtk.TreeIter iter;
			treeview.get_selection ().get_selected (null, out iter);

			var path = store.get_path (iter);
			name_renderer.editable = true;
			treeview.set_cursor (path, name_column, true);
		});
		name_renderer.edited.connect ( () => {
			name_renderer.editable = false;
		});
		name_renderer.edited.connect (on_save_item_name);
		name_renderer.editing_canceled.connect ( () => {
			name_renderer.editable = false;
		});
		add_folder_menu.activate.connect (on_add_folder);
		add_request_menu.activate.connect (on_add_item);
		delete_menu.activate.connect (on_delete_item);

		var dnd = false;
		store.row_inserted.connect (() => {
			dnd = true;
		});

		store.row_changed.connect ((path, iter) => {
			if (!dnd) {
				return;
			}

			dnd = false;
			print (@"=======$(path)\n");
		});

		treeview.drag_data_received.connect ((context, x, y, selection_data, info, time) => {
			if (context.get_selected_action () != Gdk.DragAction.MOVE) {
				return;
			}
			Gtk.drag_finish( context, false, false, time);
		});
	}

	public unowned Benchwell.HttpItem? get_selected_item (out Gtk.TreeIter iter) {
		treeview.get_selection ().get_selected (null, out iter);

		GLib.Value val;
		store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
		return val.get_object () as Benchwell.HttpItem;
	}

	private void on_add_folder () {
		if (selected_collection == null) {
			return;
		}

		var item = new Benchwell.HttpItem ();
		item.is_folder = true;
		try {
			selected_collection.add_item (item);
		} catch (ConfigError err) {
			stderr.printf (err.message);
			return;
		}

		Gtk.TreeIter sibling;
		get_selected_item (out sibling);

		var iter = add_row (item, null, sibling);
		var path = store.get_path (iter);
		name_renderer.editable = true;
		treeview.set_cursor (path, name_column, true);
	}

	private void on_add_item () {
		if (selected_collection == null) {
			return;
		}

		Gtk.TreeIter parent;
		var selected_item = get_selected_item (out parent);
		var http_item_id = selected_item.id;
		Gtk.TreeIter? sibling = null;

		if (!selected_item.is_folder) {
			sibling = parent;
			http_item_id = selected_item.parent_id;
			store.iter_parent (out parent, parent);
		}

		var item = new Benchwell.HttpItem ();
		item.is_folder = false;
		item.parent_id = http_item_id;
		try {
			selected_collection.add_item (item);
		} catch (ConfigError err) {
			stderr.printf (err.message);
			return;
		}

		var iter = add_row (item, parent, sibling);
		var path = store.get_path (iter);
		name_renderer.editable = true;
		treeview.set_cursor (path, name_column, true);
		load_request (item);
	}

	private void on_delete_item () {
		Gtk.TreeIter iter;
		var item = get_selected_item (out iter);
		try {
			selected_collection.delete_item (item);
		} catch (ConfigError err) {
			stderr.printf (err.message);
			return;
		}
		store.remove (ref iter);
	}

	private void on_save_item_name (Gtk.CellRendererText renderer, string path, string new_text) {
		Gtk.TreeIter iter;
		var selected_item = get_selected_item (out iter);
		if (selected_item == null) {
			return;
		}

		selected_item.name = new_text;

		try {
			selected_item.simple_save ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
			return;
		}

		store.set_value (iter, Benchwell.Http.Columns.TEXT, new_text);
	}

	private void on_load_item () {
		Gtk.TreeIter iter;
		var selected_item = get_selected_item (out iter);
		if (selected_item == null) {
			return;
		}

		try {
			Config.load_full_item (selected_item);
		} catch (ConfigError err) {
			stderr.printf (err.message);
			return;
		}

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

			Config.settings.set_int64 (Benchwell.Settings.HTTP_COLLECTION_ID.to_string (), collection.id);

			try {
				Config.load_root_items (collection);
			} catch (Benchwell.ConfigError err) {
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

	private void build_tree (Gtk.TreeIter? parent, Benchwell.HttpItem[] items) {
		foreach (var item in items) {
			var folder_parent = add_row (item, parent, null);

			if (item.is_folder) {
				build_tree (folder_parent, item.items);
			}
		}
	}

	private Gtk.TreeIter add_row (Benchwell.HttpItem item, Gtk.TreeIter? parent, Gtk.TreeIter? sibling) {
		Gtk.TreeIter iter;
		store.insert_before (out iter, parent, sibling);

		if (item.is_folder) {
			try {
				var px = Gtk.IconTheme.get_default ().load_icon ("bw-directory", Gtk.IconSize.BUTTON, Gtk.IconLookupFlags.USE_BUILTIN);
				store.set_value (iter, Benchwell.Http.Columns.ICON, px);
			} catch (GLib.Error err) {
				stderr.printf (err.message);
			}
		} else {
			store.set_value (iter, Benchwell.Http.Columns.METHOD, item.method);
		}

		store.set_value (iter, Benchwell.Http.Columns.TEXT, item.name);
		store.set_value (iter, Benchwell.Http.Columns.ITEM, item);

		item.notify["method"].connect (() => {
			store.set_value (iter, Benchwell.Http.Columns.METHOD, item.method);
		});
		return iter;
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