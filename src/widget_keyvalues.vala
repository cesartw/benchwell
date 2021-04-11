public class Benchwell.KeyValues : Gtk.Box {
	public Benchwell.KeyValueTypes _supported_types;
	public Benchwell.KeyValueTypes supported_types {
		get {
			return _supported_types;
		}
		set {
			_supported_types = value;
			get_children ().foreach ((w) => {
				var kv = w as Benchwell.KeyValue;
				kv.supported_types = _supported_types;
			});
		}
	}
	private Gtk.SizeGroup sgkey;
	private Gtk.SizeGroup sgval;
	private Gtk.SizeGroup sgbtn;

	public signal void changed ();
	public signal void no_row_left ();
	public signal void row_added (Benchwell.KeyValue kv);
	public signal void row_removed (Benchwell.KeyValue kv);

	public KeyValues (Benchwell.KeyValueTypes types = Benchwell.KeyValueTypes.STRING) {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);

		supported_types= types;
		sgkey = new Gtk.SizeGroup (Gtk.SizeGroupMode.BOTH);
		sgval = new Gtk.SizeGroup (Gtk.SizeGroupMode.BOTH);
		sgbtn = new Gtk.SizeGroup (Gtk.SizeGroupMode.HORIZONTAL);

		get_style_context ().add_class ("keyvalues");
	}

	public Benchwell.KeyValueI[] get_kvs () {
		Benchwell.KeyValueI[] items = null;

		get_children ().foreach ( (child) => {
			var kv = child as Benchwell.KeyValue;
			if (kv.entry_key.get_text () == "" || !kv.switch_enabled.state) {
				return;
			}

			items += kv.keyvalue;
		});

		return items;
	}

	public new void add (Benchwell.KeyValueI kvi) {
		var kv = new Benchwell.KeyValue (kvi, sgkey, sgval, sgbtn, supported_types);
		kv.show ();

		kv.entry_key.grab_focus.connect (() => { add_blank (kv); });

		kv.entry_val.grab_focus.connect (() => { add_blank (kv); });

		kv.switch_enabled.state_set.connect ((b) => {
			add_blank (kv);
			return false;
		});

		kv.btn_remove.clicked.connect( () => {
			row_removed (kv);
			remove(kv);
			if (get_children ().length () == 0) {
				//add (row_wanted ());
				no_row_left ();
			}
		});

		pack_start (kv, false, false, 0);

		kv.changed.connect (() => { changed ();});

		row_added (kv);
	}

	public void clear () {
		get_children ().foreach ( (c) => {
			remove (c);
		});
	}

	private void add_blank (KeyValue kv) {
		var last = get_children ().nth_data (get_children ().length () - 1) as KeyValue;
		if (last != null) {
			if ((last.entry_key.text == null || last.entry_key.text == "") &&
				(last.entry_val.text == null || last.entry_val.text == "")) {
				return;
			}
		}

		get_focus_child ();
		if (get_children ().index (kv) != get_children ().length () - 1) {
			return;
		}

		no_row_left ();
	}
}

public class Benchwell.KeyValue : Gtk.Box {
	public Gtk.Switch switch_enabled;
	public Gtk.Entry        entry_key;
	public Gtk.Entry        entry_val;
	public Gtk.Button       multi_edit_btn;
	public Gtk.Button       select_file_btn;
	public Benchwell.Button btn_remove;
	public Benchwell.KeyValueI keyvalue {
		get { return _keyvalue; }
		set {
			enabled_update = false;
			_keyvalue = value;
			if (_keyvalue != null) {
				if (_keyvalue.key != null)
					entry_key.text = _keyvalue.key;
				if (_keyvalue.val != null)
					entry_val.text = _keyvalue.val;
				switch_enabled.state = _keyvalue.enabled;
			}
			switch_enabled.state = true;
			enabled_update = true;
		}
	}
	public Benchwell.KeyValueTypes _supported_types;
	public Benchwell.KeyValueTypes supported_types {
		get {
			return _supported_types;
		}
		set {
			_supported_types = value;

			if (!(Benchwell.KeyValueTypes.STRING in _supported_types) ) {
				text_menu_opt.hide ();
			} else {
				text_menu_opt.show ();
			}

			if (!(Benchwell.KeyValueTypes.MULTILINE in _supported_types)) {
				multi_menu_opt.hide ();
			} else {
				multi_menu_opt.show ();
			}

			if (!(Benchwell.KeyValueTypes.FILE in _supported_types)) {
				file_menu_opt.hide ();
			} else {
				file_menu_opt.show ();
			}
		}
	}

	private Gtk.Menu type_menu;
	private Gtk.MenuItem text_menu_opt;
	private Gtk.MenuItem multi_menu_opt;
	private Gtk.MenuItem file_menu_opt;
	private Benchwell.KeyValueI _keyvalue;
	private bool enabled_update = true;

	public signal void changed ();

	public KeyValue (
		Benchwell.KeyValueI kv,
		Gtk.SizeGroup sgkey,
		Gtk.SizeGroup sgval,
		Gtk.SizeGroup sgbtn,
		Benchwell.KeyValueTypes types = Benchwell.KeyValueTypes.STRING
	) {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		get_style_context ().add_class ("keyvalue");

		entry_key = new Gtk.Entry ();
		entry_key.placeholder_text = _("Name");
		entry_key.show ();
		sgkey.add_widget (entry_key);

		btn_remove = new Benchwell.Button ("close", Gtk.IconSize.BUTTON);
		btn_remove.show ();
		sgbtn.add_widget (btn_remove);

		switch_enabled = new Gtk.Switch ();
		switch_enabled.valign = Gtk.Align.CENTER;
		switch_enabled.vexpand = false;
		switch_enabled.state = true;
		switch_enabled.show ();
		sgbtn.add_widget (switch_enabled);

		entry_val = new Gtk.Entry ();
		entry_val.placeholder_text = _("Value");

		text_menu_opt = new Gtk.MenuItem.with_label (_("Text"));
		sgval.add_widget (entry_val);

		multi_edit_btn = new Gtk.Button.with_label (_("Edit"));
		multi_menu_opt = new Gtk.MenuItem.with_label (_("Multiline"));
		sgval.add_widget (multi_edit_btn);

		select_file_btn = new Gtk.Button.with_label (_("Select"));
		file_menu_opt = new Gtk.MenuItem.with_label (_("File"));
		sgval.add_widget (select_file_btn);

		var type_button = new Benchwell.Button ("config", Gtk.IconSize.BUTTON);
		type_button.show ();
		sgbtn.add_widget (type_button);

		type_menu = new Gtk.Menu ();
		type_menu.add (text_menu_opt);
		type_menu.add (multi_menu_opt);
		type_menu.add (file_menu_opt);


		pack_start (entry_key, true, true, 0);
		pack_start (entry_val, true, true, 0);
		pack_start (multi_edit_btn, true, true, 0);
		pack_start (select_file_btn, true, true, 0);
		pack_end (btn_remove, false, false, 0);
		pack_end (switch_enabled, false, false, 5);
		pack_start (type_button, false, false, 0);

		keyvalue = kv;

		type_button.clicked.connect (() => {
			type_menu.popup_at_widget (type_button, Gdk.Gravity.SOUTH_EAST, Gdk.Gravity.NORTH_EAST, null);
		});
		keyvalue.notify["kvtype"].connect (on_kvtype_changed);
		entry_key.changed.connect (on_change);
		switch_enabled.state_set.connect (on_change_state);
		entry_val.changed.connect (on_change);

		text_menu_opt.activate.connect (() => {
			kv.kvtype = Benchwell.KeyValueTypes.STRING;
		});
		file_menu_opt.activate.connect (() => {
			kv.kvtype = Benchwell.KeyValueTypes.FILE;
		});
		select_file_btn.clicked.connect (on_select_file);
		multi_menu_opt.activate.connect (() => {
			kv.kvtype = Benchwell.KeyValueTypes.MULTILINE;
		});
		multi_edit_btn.clicked.connect (on_multi_edit);

		supported_types = types;
		on_kvtype_changed ();
	}

	private void on_kvtype_changed () {
		var show_string = keyvalue.kvtype == Benchwell.KeyValueTypes.STRING;
		var show_multi = keyvalue.kvtype == Benchwell.KeyValueTypes.MULTILINE;
		var show_file = keyvalue.kvtype == Benchwell.KeyValueTypes.FILE;

		if (entry_val != null) {
			entry_val.hide ();
		}

		if (multi_edit_btn != null) {
			multi_edit_btn.hide ();
		}

		if (select_file_btn != null) {
			select_file_btn.hide ();
		}

		if (show_string) {
			entry_val.show ();
		}

		if (show_multi) {
			multi_edit_btn.show ();
			if (keyvalue.val != null && keyvalue.val != "") {
				multi_edit_btn.tooltip_text = keyvalue.val;
			}
		}

		if (show_file) {
			select_file_btn.show ();
			if (keyvalue.val != null && keyvalue.val != "") {
				select_file_btn.label = _("Selected");
				select_file_btn.tooltip_text = keyvalue.val;
			}
		}
	}

	private void on_change () {
		if ( !enabled_update ) {
			return;
		}
		if (keyvalue == null) {
			return;
		}

		keyvalue.key = entry_key.text;
		keyvalue.val = entry_val.text;
		keyvalue.enabled = switch_enabled.state;
		changed ();
	}

	private bool on_change_state (bool state) {
		if ( !enabled_update ) {
			return false;
		}
		if (keyvalue == null) {
			return false;
		}

		keyvalue.enabled = switch_enabled.state;
		changed ();

		return false;
	}

	private void on_multi_edit () {
		var w = get_toplevel () as Gtk.Window;
		var dialog = new Gtk.Dialog.with_buttons (keyvalue.key, w,
								Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
								_("Ok"), Gtk.ResponseType.OK,
								_("Cancel"), Gtk.ResponseType.CANCEL);
		dialog.set_default_size (400, 400);

		var sv = new Benchwell.SourceView ();
		sv.show ();
		sv.get_buffer ().set_text (keyvalue.val);

		var sw = new Gtk.ScrolledWindow (null, null);
		sw.add (sv);
		sw.show ();

		dialog.get_content_area ().add (sw);

		var result = dialog.run ();
		if (result == Gtk.ResponseType.OK) {
			keyvalue.val = sv.get_buffer ().text;
			multi_edit_btn.tooltip_text = keyvalue.val;
			if (entry_val != null) {
				entry_val.text = keyvalue.val;
			}
		}
		dialog.destroy ();
	}

	private void on_select_file () {
		var w = get_toplevel () as Gtk.Window;
		var dialog = new Gtk.FileChooserDialog (_("Select file"), w,
											 Gtk.FileChooserAction.OPEN,
											_("Select"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		var resp = (Gtk.ResponseType) dialog.run ();

		if (resp == Gtk.ResponseType.OK) {
			keyvalue.val = dialog.get_filename ();
			select_file_btn.label = _("Selected");
			select_file_btn.tooltip_text = keyvalue.val;
			if (entry_val != null) {
				entry_val.text = keyvalue.val;
			}
		}

		dialog.destroy ();
	}
}
