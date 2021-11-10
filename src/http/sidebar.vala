public class Benchwell.Http.SideBar : Gtk.Box {
	public Benchwell.ApplicationWindow    window { get; construct; }
	public Gtk.TreeView treeview;
	public Benchwell.Http.Store store;
	public Gtk.TreeModelFilter filter_store;
	public CollectionsComboBox collections_combo;
	public Gtk.SearchEntry search_entry;

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
	private bool auto_expanding = false;

	public SideBar (Benchwell.ApplicationWindow window) {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			window: window
		);

		if (Config.settings.http_font != "") {
			Config.set_font (this, Pango.FontDescription.from_string (Config.settings.http_font));
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
		treeview.activate_on_single_click = Config.settings.http_single_click_activate;

		search_entry = new Gtk.SearchEntry ();

		// VISIBILITY
		// SORT
		// ICON
		// TEXT
		// METHOD
		// ITEM
		store = new Benchwell.Http.Store ({
			GLib.Type.BOOLEAN,
			GLib.Type.INT64,
			GLib.Type.OBJECT,
			GLib.Type.STRING,
			GLib.Type.STRING,
			GLib.Type.OBJECT
		});
		store.set_sort_column_id (Benchwell.Http.Columns.SORT, Gtk.SortType.ASCENDING);
		store.set_default_sort_func((iter1,iter2) => {return 0;});
		filter_store = new Gtk.TreeModelFilter (store, null);
		filter_store.set_visible_column (Benchwell.Http.Columns.VISIBILITY);

		name_renderer = new Gtk.CellRendererText ();
		//var icon_renderer = new Gtk.CellRendererPixbuf ();
		var method_renderer = new Gtk.CellRendererText ();

		name_column = new Gtk.TreeViewColumn ();
		name_column.title = _("Name");
		name_column.resizable = true;
		// NOTE: must pack_* before add_attribute
		//name_column.pack_start (icon_renderer, false);
		name_column.pack_start (method_renderer, false);
		name_column.pack_start (name_renderer, true);
		//name_column.add_attribute (icon_renderer, "pixbuf", Benchwell.Http.Columns.ICON);
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
				var color = Benchwell.HighlightColors.parse (item.method);
				if (color == null) {
					return;
				}
				cell.set_property ("markup", @"<span foreground=\"$color\">$(item.method)</span>");
				cell.set_property ("visible", true);
			} else {
				cell.set_property ("visible", false);
			}
		});
		//name_column.set_cell_data_func (icon_renderer, (cell_layout, cell, tree_model, iter) => {
			//GLib.Value val;
			//tree_model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			//var item = val.get_object () as Benchwell.HttpItem;
			//if (item == null) {
				//return;
			//}

			//if (item.is_folder) {
				//cell.set_property ("visible", true);
			//} else {
				//cell.set_property ("visible", false);
			//}
		//});

		treeview.append_column (name_column);
		treeview.set_model (store);
		treeview.show ();

		var treeview_sw = new Gtk.ScrolledWindow (null, null);
		treeview_sw.add (treeview);
		treeview_sw.show ();
		///////////

		collections_combo = new CollectionsComboBox (window);
		collections_combo.show ();

		menu = new Gtk.Menu ();

		add_folder_menu = new Benchwell.MenuItem (_("New folder"), "directory");
		add_folder_menu.show ();

		add_request_menu = new Benchwell.MenuItem (_("New request"), "add");
		add_request_menu.show ();

		clone_request_menu = new Benchwell.MenuItem (_("Clone request"), "copy");
		clone_request_menu.show ();

		delete_menu = new Benchwell.MenuItem (_("Delete"), "close");
		delete_menu.show ();

		edit_menu = new Benchwell.MenuItem (_("Rename"), "config");
		edit_menu.show ();

		menu.add (add_request_menu);
		menu.add (add_folder_menu);
		menu.add (clone_request_menu);
		menu.add (edit_menu);
		menu.add (delete_menu);

		pack_start (collections_combo, false, false, 0);
		pack_start (search_entry, false, false, 5);
		pack_start (treeview_sw, true, true, 0);

		// SIGNALS

		search_entry.key_press_event.connect ((widget, event) => {
			if (event.keyval != Gdk.Key.Escape) {
				return false;
			}
			//treeview.set_model (filter_store);

			clear_search ();

			search_entry.get_buffer().delete_text (0 , -1);
			treeview.grab_focus ();

			return false;
		});

		search_entry.search_changed.connect (() => {
			var term = search_entry.get_buffer ().get_text ();
			if (term == "") {
				clear_search ();
				return;
			}

			treeview.freeze_child_notify ();
			store.set_sort_column_id(-1, Gtk.SortType.ASCENDING);
			store.foreach ((model, path, iter) => {
				GLib.Value val;
				int score = 0;

				store.get_value (iter, Columns.TEXT, out val);
				var matched = Utils.fuzzy_match(term, val.get_string (), out score);
				if (!matched) {
					score = 9999;
				}
				//store.set_value (iter, Columns.VISIBILITY, true);
				store.set_value (iter, Columns.SORT, score);

				return false;
			});
			store.set_sort_column_id(Columns.SORT, Gtk.SortType.ASCENDING);
			treeview.thaw_child_notify ();

			auto_expanding = true;
			treeview.expand_all ();
			auto_expanding = false;
			//filter_store.refilter ();
		});

		search_entry.focus_out_event.connect (() => {
			if (search_entry.get_buffer ().get_text () == "") {
				search_entry.hide ();
				treeview.grab_focus ();
				clear_search ();
			}

			return Gdk.EVENT_PROPAGATE;
		});

		collections_combo.changed.connect (on_collection_selected);
		treeview.row_activated.connect (on_load_item);
		treeview.row_expanded.connect ((iter, path) => {
			if (auto_expanding) { return; }
			if (selected_collection == null) {
				return;
			}
			GLib.Value val;
			treeview.model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val as Benchwell.HttpItem;

			Config.http_tree_state.insert (item.id.to_string (), true);
			Config.save_http_tree_state ();
		});

		treeview.row_collapsed.connect ((iter, path) => {
			if (auto_expanding) { return; }
			GLib.Value val;
			treeview.model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val as Benchwell.HttpItem;

			Config.http_tree_state.remove (item.id.to_string ());
			Config.save_http_tree_state ();
		});

		treeview.key_press_event.connect ((widget, event) => {
			if (event.keyval == Gdk.Key.F2) {
				Gtk.TreeIter iter;
				treeview.get_selection ().get_selected (null, out iter);

				var path = treeview.model.get_path (iter);
				name_renderer.editable = true;
				treeview.set_cursor (path, name_column, true);
				return true;
			}

			if (event.keyval != Gdk.Key.f) {
				return false;
			}

			if (event.state != Gdk.ModifierType.CONTROL_MASK) {
				return false;
			}

			search_entry.show ();
			search_entry.grab_focus ();
			//treeview.set_model (filter_store);

			return true;
		});

		selected_collection = Config.get_selected_http_collection ();
		if (selected_collection != null) {
			collections_combo.set_active_id (selected_collection.id.to_string ());
		}

		if (selected_collection != null)
			selected_collection.item_added.connect ((item) => {
				Gtk.TreeIter? iter = null;
				if (item.parent_id == 0) {
					store.append (out iter, null);
				} else {
					store.foreach((model, path, parent_iter) => {
						GLib.Value val;
						store.get_value (parent_iter,  Benchwell.Http.Columns.ITEM, out val);
						var iter_item = val as Benchwell.HttpItem;
						if (iter_item.id != item.parent_id) {
							return false;
						}

						store.append (out iter, parent_iter);
						return true;
					});
				}

				if (iter == null) {
					return;
				}

				store.set_value (iter, Benchwell.Http.Columns.ITEM, item);
				store.set_value (iter, Benchwell.Http.Columns.VISIBILITY, true);
				store.set_value (iter, Benchwell.Http.Columns.METHOD, item.method);
				store.set_value (iter, Benchwell.Http.Columns.TEXT, item.name);
			});

		// edit hack
		edit_menu.activate.connect ( () => {
			Gtk.TreeIter iter;
			treeview.get_selection ().get_selected (null, out iter);

			var path = treeview.model.get_path (iter);
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

		Config.settings.changed["http-font"].connect (() => {
			Config.set_font (this, Pango.FontDescription.from_string (Config.settings.http_font));
		});

		Config.settings.changed["http-single-click-activate"].connect (() => {
			treeview.activate_on_single_click = Config.settings.http_single_click_activate;
		});

		if (Config.settings.http_collection_id != 0) {
			on_collection_selected ();
		}
	}

	private void clear_search () {
		treeview.freeze_child_notify ();
		store.set_sort_column_id(-1, Gtk.SortType.ASCENDING);
		store.foreach ((model, path, iter) => {
			GLib.Value val;
			store.get_value (iter, Columns.ITEM, out val);
			var item = val.get_object () as HttpItem;
			store.set_value (iter, Columns.SORT, item.sort);

			var expanded = Config.http_tree_state.get (item.id.to_string ());
			if (expanded == null)
				expanded = false;

			var selected_path = store.get_path (iter);
			if (expanded)
				treeview.expand_row (selected_path, false);
			else
				treeview.collapse_row (selected_path);

			return false;
		});
		store.set_sort_column_id(Columns.SORT, Gtk.SortType.ASCENDING);
		treeview.thaw_child_notify ();
	}

	public void on_resort () {
		resort ();
	}

	public void resort () {
		var sort = 0;
		treeview.model.foreach ((model, path, iter) => {
			GLib.Value val;
			treeview.model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
			var item = val.get_object () as Benchwell.HttpItem;

			// PARENT CHANGE
			Gtk.TreeIter? parent_iter = null;
			var ok = treeview.model.iter_parent (out parent_iter, iter);
			if (ok) {
				treeview.model.get_value (parent_iter, Benchwell.Http.Columns.ITEM, out val);
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
		var ok = treeview.get_selection ().get_selected (null, out iter);
		if (!ok)
			return null;

		GLib.Value val;
		treeview.model.get_value (iter, Benchwell.Http.Columns.ITEM, out val);
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
			var path = treeview.model.get_path (iter);
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
			ok = treeview.model.get_iter (out parent_iter, path);
			if (ok) {
				treeview.model.get_value (parent_iter, Benchwell.Http.Columns.ITEM, out val);
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
					treeview.model.iter_parent (out parent_iter, parent_iter);
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
		path = treeview.model.get_path (iter);
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
		item_removed (item);
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
			selected_item.save ();
		} catch (ConfigError err) {
			Config.show_alert (this, err.message);
			return;
		}

		//Gtk.TreeIter child_iter;
		//filter_store.convert_iter_to_child_iter (out child_iter, iter);
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
			return;
		}

		var collection_id = int64.parse (collections_combo.get_active_id ());

		Config.settings.http_collection_id = collection_id;
		var collection = Config.get_selected_http_collection ();

		load_collection (collection);
		selected_collection = collection;
	}

	// private void on_add_new_collection () {
	// 	var dialog = new Gtk.Dialog.with_buttons (_("Add collection"), window,
	// 								Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
	// 								_("Save"), Gtk.ResponseType.OK,
	// 								_("Cancel"), Gtk.ResponseType.CANCEL);
	// 	dialog.set_default_size (250, 130);

	// 	var label = new Gtk.Label (_("Enter collection name"));
	// 	label.show ();

	// 	var entry = new Gtk.Entry ();
	// 	entry.show ();

	// 	var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 10);
	// 	box.show ();

	// 	box.pack_start (label, true, true, 0);
	// 	box.pack_start (entry, true, true, 0);

	// 	dialog.get_content_area ().add (box);

	// 	var resp = (Gtk.ResponseType) dialog.run ();
	// 	var name = entry.get_text ();
	// 	dialog.destroy ();

	// 	if (resp != Gtk.ResponseType.OK) {
	// 		return;
	// 	}

	// 	try {
	// 		var collection = new Benchwell.HttpCollection ();
	// 		collection.name = name;
	// 		Config.add_http_collection (collection);
	// 		//collections_combo.append (collection.id.to_string (), collection.name);
	// 		collections_combo.set_active_id (collection.id.to_string ());
	// 	} catch (Benchwell.ConfigError err) {
	// 		Config.show_alert (this, err.message);
	// 	}
	// }

	private void load_collection (Benchwell.HttpCollection collection) {
		store.clear ();
		build_tree (null, collection.items, true);
	}

	private void build_tree (Gtk.TreeIter? parent, Benchwell.HttpItem[] items, bool is_parent_expanded) {
		var selected_item_id = Config.settings.http_item_id;
		foreach (var item in items) {
			var folder_parent = add_row (item, parent, null);
			var expanded = Config.http_tree_state.get (item.id.to_string ());
			if (expanded == null)
				expanded = false;

			if (item.is_folder) {
				build_tree (folder_parent, item.items, expanded && is_parent_expanded);

				if (expanded && is_parent_expanded) {
					treeview.expand_to_path (store.get_path(folder_parent));
				}
			}
			if (item.id == selected_item_id) {
				Timeout.add (0, () => { // so tree has finished building
					treeview.get_selection ().select_iter (folder_parent);
					treeview.row_activated (store.get_path (folder_parent), treeview.get_column(0));
					return false;
				});
			}
		}
	}

	private Gtk.TreeIter add_row (Benchwell.HttpItem item, Gtk.TreeIter? parent = null, Gtk.TreeIter? sibling = null) {
		Gtk.TreeIter iter;
		store.insert_before (out iter, parent, sibling);

		if (item.is_folder) {
			//try {
				//var px = Gtk.IconTheme.get_default ().load_icon ("bw-directory", Gtk.IconSize.BUTTON, Gtk.IconLookupFlags.FORCE_SIZE);
				//px.scale_simple ( Gtk.IconSize.BUTTON, Gtk.IconSize.BUTTON, Gdk.InterpType.NEAREST);
				//store.set_value (iter, Benchwell.Http.Columns.ICON, px);
			//} catch (GLib.Error err) {
				//Config.show_alert (this, err.message);
			//}
		} else {
			store.set_value (iter, Benchwell.Http.Columns.METHOD, item.method);
		}

		store.set_value (iter, Benchwell.Http.Columns.VISIBILITY, true);
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
