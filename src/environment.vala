// modules: gtk+-3.0
// vapidirs: vapi

public class Benchwell.EnvironmentEditor : Gtk.Box {
	public Benchwell.Button btn_add;
	public Benchwell.Button btn_remove;
	private Gtk.Stack stack;

	public EnvironmentEditor () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 0
		);

		var header_bar = new Gtk.HeaderBar ();
		header_bar.show ();
		header_bar.title = _("Environments");

		btn_add = new Benchwell.Button ("white-add", Gtk.IconSize.BUTTON);
		btn_add.show ();
		btn_add.get_style_context ().add_class ("suggested-action");

		btn_remove = new Benchwell.Button ("white-close", Gtk.IconSize.BUTTON);
		btn_remove.show ();
		btn_remove.get_style_context ().add_class ("destructive-action");

		header_bar.pack_end (btn_add);
		header_bar.pack_start (btn_remove);

		var paned = new Gtk.Paned (Gtk.Orientation.HORIZONTAL);
		paned.show ();

		var switcher = new Gtk.StackSwitcher ();
		switcher.orientation = Gtk.Orientation.VERTICAL;
		switcher.vexpand = true;
		switcher.hexpand = true;
		switcher.show ();

		stack = new Gtk.Stack ();
		stack.show ();
		stack.homogeneous = true;
		stack.vexpand = true;
		stack.hexpand = true;

		switcher.stack = stack;

		for (var i = 0; i < Config.environments.length; i++) {
			var env = Config.environments[i];
			var panel = new Benchwell.EnvironmentPanel (env);
			panel.show ();
			stack.add_titled (panel, env.name, env.name);
			panel.entry_name.changed.connect (() => {
				stack.child_set_property(panel, "title", panel.entry_name.text);
			});
		}

		Config.environment_added.connect ((env) => {
			var panel = new Benchwell.EnvironmentPanel (env);
			panel.show ();
			stack.add_titled (panel, env.name, env.name);
			stack.set_visible_child (panel);
		});

		paned.pack1 (switcher, false, true);
		paned.pack2 (stack, true, false);

		pack_start (header_bar, true, true, 0);
		pack_start (paned, false, false, 0);

		btn_add.clicked.connect (on_add_env);
		btn_remove.clicked.connect (on_remove_env);
	}

	private void on_add_env () {
		try {
			Config.add_environment ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}
	}

	private void on_remove_env () {
		var panel = stack.get_visible_child ();
		var index = stack.get_children ().index (panel);
		try {
			Config.environments[index].remove ();
		} catch(ConfigError err) {
			stderr.printf (err.message);
		}

		stack.remove (panel);
	}
}

public class Benchwell.EnvironmentPanel : Gtk.Box {
	public Gtk.Entry  entry_name;
	public Benchwell.Environment environment { get; construct; }
	public Benchwell.KeyValues keyvalues;

	public EnvironmentPanel (Benchwell.Environment env) {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5,
			environment: env
		);

		entry_name = new Gtk.Entry ();
		entry_name.set_text (env.name);
		entry_name.set_placeholder_text ("Name");
		entry_name.show ();

		var vbox = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		vbox.show ();

		keyvalues = new Benchwell.KeyValues ();
		keyvalues.row_added.connect (on_row_added);
		if (env.variables.length > 0) {
			keyvalues.clear ();
			foreach (var v in env.variables) {
				keyvalues.add ((Benchwell.KeyValueI) v);
			}
		}

		vbox.pack_start (entry_name, false, false, 5);
		vbox.pack_start (keyvalues, true, true, 5);
		keyvalues.show ();

		pack_start (vbox, true, true, 5);

		// signals
		keyvalues.changed.connect (on_save);
		entry_name.changed.connect (on_save);

		if (keyvalues.get_children ().length () == 0) {
			keyvalues.add (on_row_added ());
		}
	}

	private Benchwell.KeyValueI on_row_added () {
		KeyValueI kv = null;
		try {
			kv = (Benchwell.KeyValueI) environment.add_variable ();
		} catch (ConfigError err) {
			stderr.printf (err.message);
		}

		return kv;
	}

	private void on_save () {
		environment.name = entry_name.text;

		try {
			environment.save ();
		} catch (Error err) {
			stderr.printf (err.message);
		}
	}
}
