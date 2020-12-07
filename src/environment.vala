// modules: gtk+-3.0
// vapidirs: vapi

public class Benchwell.EnvironmentEditor : Gtk.Box {
	public Benchwell.Button btn_add;
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

		header_bar.pack_end (btn_add);

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

		Config.environments.foreach ((env) => {
			var panel = new Benchwell.EnvironmentPanel (env);
			panel.show ();
			stack.add_titled (panel, env.name, env.name);
		});

		paned.pack1 (switcher, false, true);
		paned.pack2 (stack, true, false);

		pack_start (header_bar, true, true, 0);
		pack_start (paned, false, false, 0);

		btn_add.clicked.connect (on_add_env);
	}

	private void on_add_env () {
		var env = new Benchwell.Environment ();
		env.name = "New environment";
		try {
			Config.save_environment (env);
		} catch (GLib.Error err) {
			stderr.printf (err.message);
		}

		var panel = new Benchwell.EnvironmentPanel (env);
		panel.show ();
		stack.add_titled (panel, env.name, env.name);
		stack.set_visible_child (panel);
	}
}

public class Benchwell.EnvironmentPanel : Gtk.Box {
	public Gtk.Entry  entry_name;
	public Benchwell.Environment environment { get; construct; }

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

		var vv = new Benchwell.KeyValues ();
		if (env.variables.length > 0) {
			vv.clear ();
			foreach (var v in env.variables) {
				vv.add (v);
			}
		}

		vbox.pack_start (entry_name, false, false, 5);
		vbox.pack_start (vv, true, true, 5);
		//vbox.pack_end (btn_box, false, false, 5);
		vv.show ();

		pack_start (vbox, true, true, 5);

		// signals

		vv.changed.connect (on_save);
	}

	private void on_save () {
		try {
			environment.save ();
		} catch (Error err) {
			stderr.printf (err.message);
		}

		if (Config.environment != null && environment.id == Config.environment.id) {
			Config.changed ();
		}
	}
}
