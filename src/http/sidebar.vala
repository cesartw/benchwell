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

	public signal void load_request(Benchwell.HttpItem item);

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

		var image_column = new Gtk.TreeViewColumn.with_attributes("image", new Gtk.CellRendererPixbuf (), "pixbuf", Benchwell.Http.Columns.ICON);

		var name_renderer = new Gtk.CellRendererText ();
		var name_column = new Gtk.TreeViewColumn.with_attributes("name", name_renderer, "text", Benchwell.Http.Columns.TEXT);

		var method_renderer = new Gtk.CellRendererText ();
		var method_column = new Gtk.TreeViewColumn.with_attributes("method", method_renderer, "text", Benchwell.Http.Columns.METHOD);
		method_column.set_cell_data_func (method_renderer, (cell_layout, cell, tree_model, iter) => {
			GLib.Value val;
			tree_model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
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
		store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
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

	private void build_tree (Gtk.TreeIter? iter, Benchwell.HttpItem[] items) {
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
			store.set_value (iter, Benchwell.Http.Columns.ICON, px);
		} else {
			store.set_value (iter, Benchwell.Http.Columns.METHOD, item.method);
		}

		store.set_value (iter, Benchwell.Http.Columns.TEXT, item.name);
		store.set_value (iter, Benchwell.Http.Columns.ITEM, item);

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

