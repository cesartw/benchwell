enum Benchwell.Http.Columns {
	ICON,
	TEXT,
	METHOD,
	ITEM
}

public class Benchwell.Http.HttpSideBar : Gtk.Box {
	public Benchwell.ApplicationWindow    window { get; construct; }
	public Gtk.TreeView treeview;
	public Benchwell.HttpStore store;
	public Gtk.ComboBoxText collections_combo;

	public Gtk.Menu menu;
	public Gtk.MenuItem add_request_menu;
	public Gtk.MenuItem add_folder_menu;
	public Gtk.MenuItem delete_menu;
	public Gtk.MenuItem edit_menu;
	public Gtk.MenuItem clone_request_menu;


	public weak Benchwell.HttpCollection? selected_collection;
	public weak Benchwell.HttpItem? selected_item;

	public signal void item_activated (Benchwell.HttpItem item);
	public signal void item_removed (Benchwell.HttpItem item);

	private Gtk.CellRendererText name_renderer;
	private Gtk.TreeViewColumn name_column;
	private int cursor_x;
	private int cursor_y;

	public HttpSideBar (Benchwell.ApplicationWindow window) {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			window: window
		);

		// treeview
		treeview = new Gtk.TreeView ();
		treeview.margin_top = 10;
		treeview.margin_bottom = 10;
		treeview.headers_visible = false;
		treeview.show_expanders = true;
		treeview.enable_tree_lines = true;
		treeview.search_column = Benchwell.Http.Columns.TEXT;
		treeview.enable_search = true;
		treeview.reorderable = true; // would be nice
		treeview.button_release_event.connect (on_button_release_event);
		treeview.activate_on_single_click = Config.settings.get_boolean ("http-sigle-click-activate");

		store = new Benchwell.HttpStore.newv ({GLib.Type.OBJECT, GLib.Type.STRING, GLib.Type.STRING, GLib.Type.OBJECT});

		name_renderer = new Gtk.CellRendererText ();
		var icon_renderer = new Gtk.CellRendererPixbuf ();
		name_column = new Gtk.TreeViewColumn ();
		name_column.title = _("Name");
		// NOTE: must pack_* before add_attribute
		name_column.pack_start (icon_renderer, false);
		name_column.pack_start (name_renderer, true);
		name_column.add_attribute (icon_renderer, "pixbuf", Benchwell.Http.Columns.ICON);
		name_column.add_attribute (name_renderer, "text", Benchwell.Http.Columns.TEXT);
		name_column.resizable = true;

		var method_renderer = new Gtk.CellRendererText ();
		var method_column = new Gtk.TreeViewColumn ();
		method_column.title = _("Method");
		method_column.pack_start (method_renderer, false);
		method_column.add_attribute(method_renderer, "text", Benchwell.Http.Columns.METHOD);
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

		treeview.append_column (name_column);
		treeview.append_column (method_column);

		treeview.set_model (store);
		treeview.show ();
		var treeview_sw = new Gtk.ScrolledWindow (null, null);
		treeview_sw.add (treeview);
		treeview_sw.show ();
		///////////

		collections_combo = new Gtk.ComboBoxText ();
		collections_combo.append ("new", _("Add collection"));
		foreach (var collection in Config.http_collections) {
			collections_combo.append (collection.id.to_string (), collection.name);
		}
		collections_combo.show ();
		collections_combo.name = "HttpCollectionSelect";

		menu = new Gtk.Menu ();

		add_folder_menu = new Benchwell.MenuItem(_("New folder"), "directory");
		add_folder_menu.show ();

		add_request_menu = new Benchwell.MenuItem(_("New request"), "add");
		add_request_menu.show ();

		clone_request_menu = new Benchwell.MenuItem(_("Clone request"), "copy");
		clone_request_menu.show ();

		delete_menu = new Benchwell.MenuItem(_("Delete"), "close");
		delete_menu.show ();

		edit_menu = new Benchwell.MenuItem(_("Edit"), "config");
		edit_menu.show ();

		menu.add (add_request_menu);
		menu.add (add_folder_menu);
		menu.add (clone_request_menu);
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
		clone_request_menu.activate.connect (on_clone_request);

		treeview.drag_end.connect (on_resort);
		// see comment on row_drop_possible
		treeview.drag_motion.connect ( (ctx, x, y, time) => {
			int cellx, celly;
			Gtk.TreePath path;
			treeview.get_path_at_pos (x,y, out path, null, out cellx, out celly);
			store.drop_path = path;

			return false;
		});
	}

	public void on_resort () {
		GLib.MainLoop loop = new GLib.MainLoop ();
		resort.begin ((obj, res) => {
			loop.quit ();
		});
		loop.run ();
	}

	public async void resort () {
		var sort = 0;
		store.foreach ((model, path, iter) => {
			GLib.Value val;
			store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val.get_object () as Benchwell.HttpItem;

			// PARENT CHANGE
			Gtk.TreeIter? parent_iter = null;
			store.iter_parent (out parent_iter, iter);
			if (parent_iter != null) {
				store.get_value (parent_iter, Benchwell.Http.Columns.ITEM, out val);
				var parent_item = val.get_object () as Benchwell.HttpItem;
				if (parent_item != null && parent_item.is_folder) {
					item.parent_id = parent_item.id; // notify isn't working here... why?
				}
			} else {
				item.parent_id = 0;
			}
			////////////////

			item.sort = sort;
			try {
				item.save ();
			} catch (Benchwell.ConfigError err) {
				stderr.printf ("saving item %s", err.message);
			}
			sort++;

			return false;
		});
	}

	public unowned Benchwell.HttpItem? get_selected_item (out Gtk.TreeIter iter) {
		treeview.get_selection ().get_selected (null, out iter);

		GLib.Value val;
		store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
		return val.get_object () as Benchwell.HttpItem;
	}

	private void on_clone_request () {
		Gtk.TreeIter iter;
		var selected_item = get_selected_item (out iter);
		if (selected_item == null || selected_item.is_folder) {
			return;
		}

		//Gtk.TreeIter parent;

		//var selected_item = get_selected_item (out parent);
		//int64? http_item_id = null;
		//if ( selected_item  != null )
			//http_item_id = selected_item.id;

		//Gtk.TreeIter? sibling = null;

		//if (selected_item != null && !selected_item.is_folder) {
			//sibling = parent;
			//http_item_id = selected_item.parent_id;
			//store.iter_parent (out parent, parent);
		//}

		var new_item = selected_collection.clone_item (selected_item);

		iter = add_row (new_item, null, iter);
		var path = store.get_path (iter);
		name_renderer.editable = true;
		treeview.expand_to_path (path);
		treeview.set_cursor (path, name_column, true);
		item_activated (new_item);
	}

	private void on_add_folder () {
		int cell_x, cell_y;
		Gtk.TreePath? path;
		Gtk.TreeIter? parent = null;
		Gtk.TreeIter? sibling = null;
		GLib.Value val;
		Benchwell.HttpItem? parent_item = null;

		treeview.get_path_at_pos (cursor_x, cursor_y, out path, null, out cell_x, out cell_y);
		if (path != null) {
			var ok = store.get_iter (out parent, path);
			if (ok) {
				store.get_value (parent, Benchwell.Http.Columns.ITEM, out val);
				parent_item = val.get_object () as Benchwell.HttpItem;

				if (!parent_item.is_folder) {
					sibling = parent;
					store.iter_parent (out parent, parent);
					store.get_value (parent, Benchwell.Http.Columns.ITEM, out val);
					parent_item = val.get_object () as Benchwell.HttpItem;
				}
			}
		}

		if (parent_item != null) {
			var item = new Benchwell.HttpItem ();
			item.touch_without_save (() => {
				item.is_folder = true;
				item.parent_id = parent_item.id;
				try {
					selected_collection.add_item (item);
				} catch (ConfigError err) {
					stderr.printf (err.message);
					return;
				}
			});

			var iter = add_row (item, parent, sibling);
			path = store.get_path (iter);
			name_renderer.editable = true;
			treeview.expand_to_path (path);
			treeview.set_cursor (path, name_column, true);
		}
	}

	private void on_add_item () {
		if (selected_collection == null) {
			return;
		}

		Gtk.TreeIter parent;

		var selected_item = get_selected_item (out parent);
		int64? http_item_id = null;
		if ( selected_item  != null )
			http_item_id = selected_item.id;

		Gtk.TreeIter? sibling = null;

		if (selected_item != null && !selected_item.is_folder) {
			sibling = parent;
			http_item_id = selected_item.parent_id;
			store.iter_parent (out parent, parent);
		}

		var item = new Benchwell.HttpItem ();
		item.touch_without_save (() => {
			item.is_folder = false;
			if (http_item_id != null)
				item.parent_id = http_item_id;

			try {
				selected_collection.add_item (item);
			} catch (ConfigError err) {
				stderr.printf (err.message);
				return;
			}
		});

		var iter = add_row (item, parent, sibling);
		var path = store.get_path (iter);
		name_renderer.editable = true;
		treeview.expand_to_path (path);
		treeview.set_cursor (path, name_column, true);
		item_activated (item);
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
		item_removed (item);
	}

	private void on_save_item_name (Gtk.CellRendererText renderer, string path, string new_text) {
		Gtk.TreeIter iter;
		var selected_item = get_selected_item (out iter);
		if (selected_item == null) {
			return;
		}

		selected_item.name = new_text;

		try {
			selected_item.save ();
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

		if (selected_item.is_folder) {
			Config.settings.set_int64 (Benchwell.Settings.HTTP_ITEM_ID.to_string (), selected_item.id);

			var selected_path = store.get_path (iter);
			if (treeview.is_row_expanded (selected_path)) {
				treeview.collapse_row (selected_path);
				return;
			}

			// // remove items
			//Gtk.TreeIter child;
			//store.iter_children (out child, iter);
			//store.remove(ref child);

			//build_tree (iter, selected_item.items);
			treeview.expand_row (selected_path, false);
			return;
		}

		item_activated (selected_item);
	}

	private void on_collection_selected () {
		if (collections_combo.get_active_id () == "new") {
			var dialog = new Gtk.Dialog.with_buttons (_("Add collection"), window,
										Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
										_("Save"), Gtk.ResponseType.OK,
										_("Cancel"), Gtk.ResponseType.CANCEL);
			dialog.set_default_size (250, 130);

			var label = new Gtk.Label (_("Enter collection name"));
			label.show ();

			var entry = new Gtk.Entry ();
			entry.show ();

			var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 10);
			box.show ();

			box.pack_start (label, true, true, 0);
			box.pack_start (entry, true, true, 0);

			dialog.get_content_area ().add (box);

			var resp = (Gtk.ResponseType) dialog.run ();
			var name = entry.get_text ();
			dialog.destroy ();

			if (resp != Gtk.ResponseType.OK) {
				return;
			}

			try {
				var collection = Config.add_http_collection ();
				collection.name = name;

				collections_combo.append (collection.id.to_string (), collection.name);
				collections_combo.set_active_id (collection.id.to_string ());
			} catch (Benchwell.ConfigError err) {
				stderr.printf (err.message);
			}
		}

		var collection_id = int64.parse (collections_combo.get_active_id ());

		foreach(var collection in Config.http_collections) {
			if (collection.id != collection_id){
				continue;
			}

			store.clear ();

			Config.settings.set_int64 (Benchwell.Settings.HTTP_COLLECTION_ID.to_string (), collection.id);

			try {
				Config.load_http_items (collection);
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
		var saved_selected_item_id = Config.settings.get_int64 (Benchwell.Settings.HTTP_ITEM_ID.to_string ());
		build_tree (null, collection.items, saved_selected_item_id);
	}

	private void build_tree (Gtk.TreeIter? parent, Benchwell.HttpItem[] items, int64 auto_expand = 0) {
		foreach (var item in items) {
			var folder_parent = add_row (item, parent, null);

			if (auto_expand > 0 && item.id == auto_expand) {
				build_tree (folder_parent, item.items, 0);
				treeview.get_selection ().select_path (store.get_path (folder_parent)); // not really a folder
				on_load_item ();
				if (!item.is_folder)
					item_activated (item);
				continue; // on_load_item will build the tree for us
			}

			if (item.is_folder) {
				build_tree (folder_parent, item.items, 0);

				if (auto_expand > 0 && item.id == auto_expand) {
					treeview.expand_row (store.get_path(folder_parent), false);
				}
			}
		}
	}

	private Gtk.TreeIter add_row (Benchwell.HttpItem item, Gtk.TreeIter? parent = null, Gtk.TreeIter? sibling = null) {
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
			if (item.is_folder)
				store.set_value (iter, Benchwell.Http.Columns.METHOD, item.method);
		});
		return iter;
	}

	private bool on_button_release_event (Gtk.Widget w, Gdk.EventButton event) {
		if (event.button != Gdk.BUTTON_SECONDARY) {
			return false;
		}

		cursor_x = (int) event.x;
		cursor_y = (int) event.y;

		Gtk.TreePath path;
		treeview.get_path_at_pos (cursor_x, cursor_y , out path, null, null, null);
		if (path == null) {
			delete_menu.sensitive = false;
			edit_menu.sensitive = false;
			clone_request_menu.sensitive = false;
			menu.popup_at_pointer (event);
			return true;
		}

		treeview.get_selection ().select_path (path);

		Gtk.TreeIter iter;
		var item = get_selected_item (out iter);

		var enabled = item != null;

		delete_menu.sensitive = enabled;
		edit_menu.sensitive = enabled;
		clone_request_menu.sensitive = enabled && !item.is_folder;

		menu.popup_at_pointer (event);
		return true;
	}
}

public class Benchwell.HttpStore : Gtk.TreeStore, Gtk.TreeDragDest {
	// see comment on row_drop_possible
	public Gtk.TreePath? drop_path;

	public HttpStore.newv (GLib.Type[] types) {
		Object();
		set_column_types (types);
	}

	// NOTE: the event that calls this method seems to bubble throught the tree.
	//       so even when it says "no drop allowed", the propagation may cause that a
	//       parent accepts the drop thus moving the item to an unexpected place
	//
	//       We need a why to know what's the actual drop zone
	public bool row_drop_possible (Gtk.TreePath dest, Gtk.SelectionData selection_data) {
		if (drop_path == null) {
			return false;
		}

		Gtk.TreeIter iter;
		GLib.Value val;
		Benchwell.HttpItem? item = null;

		Gtk.TreePath path; // item being dragged
		Gtk.tree_get_row_drag_data (selection_data, null, out path);

		var ok = get_iter (out iter, dest);
		if (!ok) {
			dest.up ();
			ok = get_iter (out iter, dest);
			if (!ok) {
				return false;
			}

			get_value (iter, Benchwell.Http.Columns.ITEM, out val);

			item = val.get_object () as Benchwell.HttpItem;
			if (item == null) {
				return false;
			}

			return item.is_folder;
		}

		return drop_path.get_depth () == dest.get_depth ();
	}
}
