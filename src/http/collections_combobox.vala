public class Benchwell.Http.FolderSelector : Gtk.Box {
	private CollectionsComboBox collections_combobox;
	private FoldersComboBox folders_combobox;
	public weak Benchwell.HttpCollection collection {
		get {
			return collections_combobox.collection;
		}
	}
	public weak Benchwell.HttpItem folder {
		get {
			return folders_combobox.folder;
		}
	}

	public signal void changed ();

	public FolderSelector (Benchwell.HttpCollection? collection,Benchwell.ApplicationWindow window, int spacing, Gtk.Orientation orientation) {
		Object(spacing:spacing, orientation: orientation);

		collections_combobox = new CollectionsComboBox (window);
		collections_combobox.show ();
		folders_combobox = new FoldersComboBox (window);
		folders_combobox.show ();

		pack_start (collections_combobox, true, true, 0);
		pack_start (folders_combobox, true, true, 0);

		collections_combobox.changed.connect (() => {
			on_collection_selected ();
			changed ();
		});
		folders_combobox.changed.connect (() => { changed(); });
		if (collection != null) {
			collections_combobox.collection = collection;
			on_collection_selected ();
		}
	}

	private void on_collection_selected () {
		folders_combobox.set_collection (collections_combobox.collection);
	}
}

public class Benchwell.Http.CollectionsComboBox : Gtk.ComboBoxText {
	public Benchwell.ApplicationWindow    window { get; construct; }
	private weak Benchwell.HttpCollection _collection;
	public Benchwell.HttpCollection? collection {
		get {
			return _collection;
		}
		set {
			_collection = value;
			try {
				Config.load_http_items (collection);
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
				return;
			}
			set_active_id (value.id.to_string ());
		}
	}

	public CollectionsComboBox (Benchwell.ApplicationWindow window) {
		Object (
			window: window,
			name: "HttpCollectionSelect"
		);

		append ("new", _("Add collection"));
		foreach (var collection in Config.http_collections) {
			append (collection.id.to_string (), collection.name);
		}

		changed.connect (on_collection_selected);
		Config.http_collection_added.connect ((collection) => {
			append (collection.id.to_string (), collection.name);
		});
		set_active_id (Config.settings.http_collection_id.to_string ());
	}

	private void on_collection_selected () {
		if (get_active_id () == "new") {
			on_add_new_collection ();
			return;
		}

		var collection_id = int64.parse (get_active_id ());

		foreach(var c in Config.http_collections) {
			if (c.id != collection_id){
				continue;
			}

			collection = c;

			// TODO: make this eager loading optional. have no use case yet
			try {
				Config.load_http_items (collection);
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
				return;
			}

			break;
		}
	}

	private void on_add_new_collection () {
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
			var collection = new Benchwell.HttpCollection ();
			collection.name = name;
			Config.add_http_collection (collection);
			//append (collection.id.to_string (), collection.name);
			set_active_id (collection.id.to_string ());
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (this, err.message);
		}
	}
}

public class Benchwell.Http.FoldersComboBox : Gtk.ComboBoxText {
	public Benchwell.ApplicationWindow    window { get; construct; }
	public weak Benchwell.HttpCollection collection;
	public weak Benchwell.HttpItem? folder;

	public FoldersComboBox (Benchwell.ApplicationWindow window) {
		Object (
			window: window,
			name: "HttpCollectionSelect"
		);

		append ("new", _("Add new folder"));

		changed.connect (on_folder_selected);
	}

	public void set_collection (Benchwell.HttpCollection collection) {
		this.collection = collection;
		append_items_name_with_parent (null, collection.items);
	}

	private void append_items_name_with_parent (string? parent, HttpItem[] items) {
		foreach(var item in items) {
			if (!item.is_folder)
				continue;

			var name = item.name;
			if (parent != null)
				name = @"$parent -> $name";

			append (item.id.to_string (), name);
			append_items_name_with_parent (name, item.items);
		}
	}

	private void on_folder_selected () {
		if (get_active_id () != "new")
			return;

		on_add_new_folder ();
		var id = int64.parse (get_active_id ());
		foreach (var item in collection.items) {
			if (item.id != id)
				continue;

			folder = item;
			break;
		}
	}

	private void on_add_new_folder () {
		var dialog = new Gtk.Dialog.with_buttons (_("Add folder"), window,
									Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
									_("Save"), Gtk.ResponseType.OK,
									_("Cancel"), Gtk.ResponseType.CANCEL);
		dialog.set_default_size (250, 130);

		var label = new Gtk.Label (_("Enter folder name"));
		label.show ();

		var entry = new Gtk.Entry ();
		entry.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 10);
		box.show ();

		box.pack_start (label, true, true, 0);
		box.pack_start (entry, true, true, 0);

		dialog.get_content_area ().add (box);

		var resp = (Gtk.ResponseType) dialog.run ();
		if (resp != Gtk.ResponseType.OK) {
			dialog.destroy ();
			return;
		}

		var name = entry.get_text ();
		dialog.destroy ();

		try {
			var f = new Benchwell.HttpItem () {
				is_folder = true,
				http_collection_id = collection.id,
				name = name
			};
			folder = f;
			folder.save ();
			append (folder.id.to_string (), folder.name);
			set_active_id (folder.id.to_string ());
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (this, err.message);
		}
	}
}
