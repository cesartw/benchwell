public class Benchwell.Http.Store : Gtk.TreeStore, Gtk.TreeDragDest, Gtk.TreeModel {
	public Store (GLib.Type[] types) {
		set_column_types ( types);
	}

	public bool row_drop_possible (Gtk.TreePath dest, Gtk.SelectionData selection_data) {
		var path = dest.copy ();

		var indices = dest.get_indices();
		var lastIndex = indices[indices.length - 1];

		if (lastIndex != 0) {
			return true;
		}

		Gtk.TreeIter iter;
		GLib.Value val;
		var ok = get_iter (out iter, path);
		if (ok) {
			return true;
		}

		ok = path.up ();
		if (!ok) {
			return false;
		}
		ok = get_iter (out iter, path);
		if (!ok) {
			return false;
		}

		get_value (iter, Benchwell.Http.Columns.ITEM, out val);
		var drop_item = val as Benchwell.HttpItem;
		return drop_item.is_folder;
	}
}
