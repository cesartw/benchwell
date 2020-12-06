public class Benchwell.KeyValues : Gtk.Box {

	public signal void changed ();

	public KeyValues () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);

		add (null);
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

	public void add (Benchwell.KeyValueI? kvi) {
		var kv = new Benchwell.KeyValue ();
		kv.show ();
		kv.keyvalue = (Benchwell.KeyValueI) kvi;

		kv.entry_key.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (null);
		});
		kv.entry_val.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (null);
		});
		kv.switch_enabled.state_set.connect ((b) => {
			if (get_children ().index (kv) == get_children ().length () - 1) {
				add (null);
			}

			return false;
		});

		kv.btn_remove.clicked.connect( () => {
			remove(kv);
			if (get_children ().length () == 0) {
				add (null);
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
			entry_key.text = _keyvalue.key ();
			entry_val.text = _keyvalue.val ();
			switch_enabled.state = _keyvalue.enabled ();
			enabled_update = true;
		}
	}
	private Benchwell.KeyValueI _keyvalue;
	private bool enabled_update = true;

	public signal void changed ();

	public KeyValue () {
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
	}

	private void on_change () {
		if ( !enabled_update ) {
			return;
		}

		keyvalue.set_key (entry_key.text);
		keyvalue.set_val (entry_val.text);
		keyvalue.set_enabled (switch_enabled.active);
		changed ();
	}

	private bool on_change_state (bool state) {
		if ( !enabled_update ) {
			return true;
		}

		keyvalue.set_enabled (switch_enabled.active);
		changed ();

		return true;
	}
}

