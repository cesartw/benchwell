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

	public signal void item_activated (Benchwell.HttpItem item, Gtk.TreeIter iter);
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

		if (Config.settings.get_string("http-font") != "") {
			override_font (Pango.FontDescription.from_string (Config.settings.get_string("http-font")));
		}

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
		treeview.activate_on_single_click = Config.settings.get_boolean ("http-single-click-activate");

		store = new Benchwell.HttpStore.newv ({GLib.Type.OBJECT, GLib.Type.STRING, GLib.Type.STRING, GLib.Type.OBJECT});

		name_renderer = new Gtk.CellRendererText ();
		var icon_renderer = new Gtk.CellRendererPixbuf ();
		var method_renderer = new Gtk.CellRendererText ();

		name_column = new Gtk.TreeViewColumn ();
		name_column.title = _("Name");
		name_column.resizable = true;
		// NOTE: must pack_* before add_attribute
		name_column.pack_start (icon_renderer, false);
		name_column.pack_start (method_renderer, false);
		name_column.pack_start (name_renderer, true);
		name_column.add_attribute (icon_renderer, "pixbuf", Benchwell.Http.Columns.ICON);
		name_column.add_attribute (method_renderer, "text", Benchwell.Http.Columns.METHOD);
		name_column.add_attribute (name_renderer, "text", Benchwell.Http.Columns.TEXT);

		name_column.set_cell_data_func (method_renderer, (cell_layout, cell, tree_model, iter) => {
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
				cell.set_property ("visible", true);
			} else {
				cell.set_property ("visible", false);
			}
		});
		name_column.set_cell_data_func (icon_renderer, (cell_layout, cell, tree_model, iter) => {
			GLib.Value val;
			tree_model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val.get_object () as Benchwell.HttpItem;
			if (item == null) {
				return;
			}

			if (item.is_folder) {
				cell.set_property ("visible", true);
			} else {
				cell.set_property ("visible", false);
			}
		});

		treeview.append_column (name_column);
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

		edit_menu = new Benchwell.MenuItem(_("Rename"), "config");
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
		treeview.row_expanded.connect ((iter, path) => {
			if (selected_collection == null) {
				return;
			}
			GLib.Value val;
			store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val as Benchwell.HttpItem;

			Config.http_tree_state.insert (item.id, true);
		});

		treeview.row_collapsed.connect ((iter, path) => {
			GLib.Value val;
			store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val as Benchwell.HttpItem;

			Config.http_tree_state.remove (item.id);
		});

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

		treeview.drag_motion.connect ( (ctx, x, y, time) => {
			int cellx, celly;
			Gtk.TreePath path;
			treeview.get_path_at_pos (x,y, out path, null, out cellx, out celly);
			store.drop_path = path;

			return false;
		});

		Config.settings.changed["http-font"].connect (() => {
			override_font (Pango.FontDescription.from_string (Config.settings.get_string("http-font")));
		});

		Config.settings.changed["http-single-click-activate"].connect (() => {
			treeview.activate_on_single_click = Config.settings.get_boolean ("http-single-click-activate");
		});
	}

	public void on_resort () {
		//GLib.MainLoop loop = new GLib.MainLoop ();
		//resort.begin ((obj, res) => {
			//loop.quit ();
		//});
		//loop.run ();
		resort ();
	}

	public void resort () {
		var sort = 0;
		store.foreach ((model, path, iter) => {
			GLib.Value val;
			store.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val.get_object () as Benchwell.HttpItem;

			// PARENT CHANGE
			Gtk.TreeIter? parent_iter = null;
			var ok = store.iter_parent (out parent_iter, iter);
			if (ok) {
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
			sort++;

			// should stop iteration
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

		try {
			var new_item = selected_collection.clone_item (selected_item);

			iter = add_row (new_item, null, iter);
			var path = store.get_path (iter);
			name_renderer.editable = true;
			treeview.expand_to_path (path);
			treeview.set_cursor (path, name_column, true);
			item_activated (new_item, iter);
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (this, err.message);
		}
	}

	private void on_add_folder () {
		add_item (true);
	}

	private void on_add_item () {
		add_item ();
	}

	private void add_item(bool is_folder = false) {
		int cell_x, cell_y;
		Gtk.TreePath? path;
		Gtk.TreeIter? parent_iter = null;
		Gtk.TreeIter? sibling_iter = null;
		GLib.Value val;
		Benchwell.HttpItem? selected_item = null;
		var item = new Benchwell.HttpItem ();

		var ok = treeview.get_path_at_pos (cursor_x, cursor_y, out path, null, out cell_x, out cell_y);
		if (ok) {
			ok = store.get_iter (out parent_iter, path);
			if (ok) {
				store.get_value (parent_iter, Benchwell.Http.Columns.ITEM, out val);
				selected_item = val.get_object () as Benchwell.HttpItem;
			}
		}

		item.touch_without_save (() => {
			item.is_folder = is_folder;
			if (selected_item != null) {
				if (selected_item.is_folder) {
					item.parent_id = selected_item.id;
				 } else {
					sibling_iter = parent_iter;
					store.iter_parent (out parent_iter, parent_iter);
				}
			}

			try {
				selected_collection.add_item (item);
			} catch (ConfigError err) {
				Config.show_alert (this, err.message);
				return;
			}

		});

		var iter = add_row (item, parent_iter, sibling_iter);

		name_renderer.editable = true;
		path = store.get_path (iter);
		treeview.expand_to_path (path);
		treeview.set_cursor (path, name_column, true);
		if (item.is_folder)
			item_activated (item, iter);
		return;
	}

	private void on_delete_item () {
		Gtk.TreeIter iter;
		var item = get_selected_item (out iter);
		try {
			selected_collection.delete_item (item);
		} catch (ConfigError err) {
			Config.show_alert (this, err.message);
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
			Config.show_alert (this, err.message);
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

		item_activated (selected_item, iter);
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
				Config.show_alert (this, err.message);
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
				Config.show_alert (this, err.message);
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
			var expanded = Config.http_tree_state.get (item.id);
			if (expanded == null)
				expanded = false;

			if (item.is_folder) {
				build_tree (folder_parent, item.items);

				if (expanded) {
					treeview.expand_to_path (store.get_path(folder_parent));
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
				Config.show_alert (this, err.message);
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

public class Benchwell.HttpStore : Gtk.TreeStore, Gtk.TreeDragDest, Gtk.TreeModel {
	public Gtk.TreePath drop_path;
	public HttpStore.newv (GLib.Type[] types) {
		Object();
		set_column_types (types);
	}

	public bool row_drop_possible (Gtk.TreePath dest, Gtk.SelectionData selection_data) {
		var path = dest.copy ();

		// NOTE: this method is called for each index in the path.
		//       drop_path helps to know which was the actual drop path :shrug:
		if (drop_path.compare (path) != 0) {
			drop_path.append_index (0);
			if (drop_path.compare(path) != 0) {
				return false;
			}
		}

		if (path.get_depth () == 0) {
			return false;
		}

		// dropping between items. must make sure the parent is a folder
		var indices = path.get_indices ();
		if (indices[indices.length -1] == 0) {
			var ok = path.up ();
			if (!ok) {
				return false; //not never reach this point as depth is checked previously
			}
		}

		// we only care the drop area is a folder
		Gtk.TreeIter iter;
		var ok = get_iter (out iter, path);
		if (!ok) {
			return false;
		}

		GLib.Value val;
		get_value (iter, Benchwell.Http.Columns.ITEM, out val);

		var drop_item = val as Benchwell.HttpItem;
		return drop_item.is_folder;
	}
}
