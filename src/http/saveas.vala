public class Benchwell.Http.SaveAsSelector : Gtk.Grid {
	public Benchwell.ApplicationWindow   window { get; construct; }
	private Gtk.Entry name_entry;
	private Benchwell.Http.FolderSelector folder_selector;
	public weak Benchwell.HttpCollection collection {
		get {
			return folder_selector.collection;
		}
	}
	public weak Benchwell.HttpItem folder {
		get {
			return folder_selector.folder;
		}
	}

	public signal void changed ();

	public SaveAsSelector (Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			row_spacing: 5,
			column_spacing: 5
		);

		folder_selector = new Benchwell.Http.FolderSelector (Config.get_selected_http_collection (), window, 5, Gtk.Orientation.HORIZONTAL);
		folder_selector.show ();

		name_entry = new Gtk.Entry ();
		name_entry.show ();

		attach (folder_selector, 0, 0, 1, 1);
		attach (name_entry, 0, 1, 1, 1);

		folder_selector.changed.connect (()=>{changed ();});
		name_entry.changed.connect (()=>{changed ();});
	}

	public string get_name () {
		return name_entry.text;
	}
}

