public class Benchwell.KeyValues : Gtk.Box {
	public Benchwell.KeyValueTypes supported_types { get; construct; }
	private Gtk.SizeGroup sgkey;
	private Gtk.SizeGroup sgval;
	private Gtk.SizeGroup sgbtn;
	public signal void changed ();
	public signal Benchwell.KeyValueI row_wanted ();
	public signal void row_added (Benchwell.KeyValueI kvi);
	public signal void row_removed (Benchwell.KeyValueI kvi);

	public KeyValues (Benchwell.KeyValueTypes types = Benchwell.KeyValueTypes.FILE) {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5,
			supported_types: types
		);

		sgkey = new Gtk.SizeGroup (Gtk.SizeGroupMode.BOTH);
		sgval = new Gtk.SizeGroup (Gtk.SizeGroupMode.BOTH);
		sgbtn = new Gtk.SizeGroup (Gtk.SizeGroupMode.HORIZONTAL);

		get_style_context ().add_class ("keyvalues");
	}

	public void get_kvs (out string[] keys, out string[] values) {
		string[] ks = {};
		string[] vs = {};

		get_children ().foreach ( (child) => {
			var kv = child as Benchwell.KeyValue;
			var key = kv.entry_key.get_text ();
			var val = kv.entry_val.get_text ();
			if (key == "" || !kv.switch_enabled.state) {
				return;
			}

			ks += key;
			vs += val;
		});

		keys = ks;
		values = vs;
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
			row_removed (kv.keyvalue);
			remove(kv);
			if (get_children ().length () == 0) {
				add (row_wanted ());
			}
		});

		pack_start (kv, false, false, 0);

		kv.changed.connect (() => { changed ();});

		row_added (kvi);
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

		add (row_wanted ());
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
	public Benchwell.KeyValueTypes keyvaluetype { get; set; }
	public int supported_types { get; construct; }

	private Gtk.Menu type_menu;
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
			spacing: 5,
			supported_types: types
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

		pack_start (entry_key, true, true, 0);
		pack_end (btn_remove, false, false, 0);
		pack_end (switch_enabled, false, false, 5);

		var string_enabled = (supported_types ^ Benchwell.KeyValueTypes.STRING) != 0;
		var multi_enabled = (supported_types ^ Benchwell.KeyValueTypes.MULTILINE) != 0;
		var file_enabled = (supported_types ^ Benchwell.KeyValueTypes.FILE) != 0;

		Gtk.MenuItem[] menu_opts = {};

		if (string_enabled) {
			entry_val = new Gtk.Entry ();
			entry_val.placeholder_text = _("Value");
			pack_start (entry_val, true, true, 0);
			entry_val.changed.connect (on_change);
			var m = new Gtk.MenuItem.with_label (_("Text"));
			m.show ();
			m.activate.connect (() => {
				kv.kvtype = Benchwell.KeyValueTypes.STRING;
			});
			menu_opts += m;
			sgval.add_widget (entry_val);
		}

		if (multi_enabled) {
			multi_edit_btn = new Gtk.Button.with_label (_("Edit"));
			pack_start (multi_edit_btn, true, true, 0);
			var m = new Gtk.MenuItem.with_label (_("Multiline"));
			m.show ();
			m.activate.connect (() => {
				kv.kvtype = Benchwell.KeyValueTypes.MULTILINE;
			});
			menu_opts += m;
			multi_edit_btn.clicked.connect (on_multi_edit);
			sgval.add_widget (multi_edit_btn);
		}

		if (file_enabled) {
			select_file_btn = new Gtk.Button.with_label (_("Select"));
			pack_start (select_file_btn, true, true, 0);
			var m = new Gtk.MenuItem.with_label (_("File"));
			m.show ();
			m.activate.connect (() => {
				kv.kvtype = Benchwell.KeyValueTypes.FILE;
			});
			menu_opts += m;
			sgval.add_widget (select_file_btn);
		}

		if (menu_opts.length > 1) {
			type_menu = new Gtk.Menu ();
			var type_button = new Benchwell.Button ("config", Gtk.IconSize.BUTTON);
			type_button.show ();
			sgbtn.add_widget (type_button);

			for (var i = 0; i < menu_opts.length; i++) {
				type_menu.add (menu_opts[i]);
			}

			type_button.clicked.connect (() => {
				type_menu.popup_at_widget (type_button, Gdk.Gravity.SOUTH_EAST, Gdk.Gravity.NORTH_EAST, null);
			});

			pack_start (type_button, false, false, 0);
		}

		keyvalue = kv;

		keyvalue.notify["kvtype"].connect (on_kvtype_changed);
		entry_key.changed.connect (on_change);
		switch_enabled.state_set.connect (on_change_state);

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
		}

		if (show_file) {
			select_file_btn.show ();
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
		}
		dialog.destroy ();
	}
}
