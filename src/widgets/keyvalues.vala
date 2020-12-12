public class Benchwell.KeyValues : Gtk.Box {
	public signal void changed ();
	public signal Benchwell.KeyValueI row_added();

	public KeyValues () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);
	}

	public void get_kvs (out string[] keys, out string[] values) {
		string[] ks = {};
		string[] vs = {};
		get_children ().foreach ( (child) => {
			var kv = child as Benchwell.KeyValue;
			var key = kv.entry_key.get_text ();
			var val = kv.entry_val.get_text ();
			if (key == "" || val == "") {
				return;
			}

			ks += key;
			vs += val;
		});

		keys = ks;
		values = vs;
	}

	public void add (Benchwell.KeyValueI kvi) {
		var kv = new Benchwell.KeyValue (kvi);
		kv.show ();
		kv.keyvalue = (Benchwell.KeyValueI) kvi;

		kv.entry_key.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (row_added ());
		});

		kv.entry_val.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (row_added ());
		});

		kv.switch_enabled.state_set.connect ((b) => {
			if (get_children ().index (kv) == get_children ().length () - 1) {
				add (row_added ());
			}

			return false;
		});

		kv.btn_remove.clicked.connect( () => {
			remove(kv);
			if (get_children ().length () == 0) {
				add (row_added ());
			}
		});

		pack_start (kv, false, false, 0);

		kv.changed.connect (() => { changed ();});
	}

	public void clear () {
		get_children ().foreach ( (c) => {
			remove (c);
		});
	}
}

public class Benchwell.KeyValue : Gtk.Box {
	public Gtk.Switch switch_enabled;
	public Gtk.Entry        entry_key;
	public Gtk.Entry        entry_val;
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
			} else {
			}
			switch_enabled.state = true;
			enabled_update = true;
		}
	}
	private Benchwell.KeyValueI _keyvalue;
	private bool enabled_update = true;

	public signal void changed ();

	public KeyValue (Benchwell.KeyValueI? kv) {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		entry_key = new Gtk.Entry ();
		entry_key.placeholder_text = _("Name");
		entry_key.show ();

		entry_val = new Gtk.Entry ();
		entry_val.placeholder_text = _("Value");
		entry_val.show ();

		btn_remove = new Benchwell.Button ("close", Gtk.IconSize.BUTTON);
		btn_remove.show ();

		switch_enabled = new Gtk.Switch ();
		switch_enabled.valign = Gtk.Align.CENTER;
		switch_enabled.vexpand = false;
		switch_enabled.set_active (true);
		switch_enabled.show ();

		pack_start (entry_key, true, true, 0);
		pack_start (entry_val, true, true, 0);
		pack_end (btn_remove, false, false, 0);
		pack_end (switch_enabled, false, false, 5);

		entry_key.changed.connect (on_change);
		entry_val.changed.connect (on_change);
		switch_enabled.state_set.connect (on_change_state);

		keyvalue = kv;
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
		keyvalue.enabled = switch_enabled.active;
		changed ();
	}

	private bool on_change_state (bool state) {
		if ( !enabled_update ) {
			return true;
		}
		if (keyvalue == null) {
			return true;
		}

		keyvalue.enabled = switch_enabled.active;
		changed ();

		return true;
	}
}

