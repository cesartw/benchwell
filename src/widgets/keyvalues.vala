public class Benchwell.KeyValues : Gtk.Box {
	public KeyValues () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);

		add (null);
	}

	public void add (KeyValueI? kvi) {
		var kv = new Benchwell.KeyValue ();
		kv.show ();

		if (kvi != null) {
			kv.key.set_text (kvi.key ());
			kv.val.set_text (kvi.val ());
			kv.enabled.set_active (kvi.enabled ());
		}

		kv.key.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (null);
		});
		kv.val.grab_focus.connect (() => {
			if (get_children ().index (kv) != get_children ().length () - 1) {
				return;
			}

			add (null);
		});
		kv.enabled.state_set.connect ((b) => {
			if (get_children ().index (kv) == get_children ().length () - 1) {
				add (null);
			}

			return false;
		});

		kv.remove_btn.clicked.connect( () => {
			remove(kv);
			if (get_children ().length () == 0) {
				add (null);
			}
		});

		pack_start (kv, false, false, 0);
	}

	public void clear () {
		get_children ().foreach ( (c) => {
			remove (c);
		});
	}
}

public class Benchwell.KeyValue : Gtk.Box {
	public Gtk.Switch enabled;
	public Gtk.Entry        key;
	public Gtk.Entry        val;
	public Benchwell.Button remove_btn;

	public KeyValue () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		key = new Gtk.Entry ();
		key.placeholder_text = _("Name");
		key.show ();

		val = new Gtk.Entry ();
		val.placeholder_text = _("Value");
		val.show ();

		remove_btn = new Benchwell.Button ("close", Gtk.IconSize.BUTTON);
		remove_btn.show ();

		enabled = new Gtk.Switch ();
		enabled.valign = Gtk.Align.CENTER;
		enabled.vexpand = false;
		enabled.set_active (true);
		enabled.show ();

		pack_start (key, true, true, 0);
		pack_start (val, true, true, 0);
		pack_end (remove_btn, false, false, 0);
		pack_end (enabled, false, false, 5);
	}
}

