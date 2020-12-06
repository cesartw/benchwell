// modules: gtk+-3.0
// vapidirs: vapi

public class Benchwell.EnvironmentEditor : Gtk.Box {
	public Benchwell.Button btn_add;

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

		var stack = new Gtk.Stack ();
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
	}
}

public class Benchwell.EnvironmentPanel : Gtk.Box {
	public Gtk.Entry entry_name;
	public Gtk.Button btn_save;
	public Gtk.Button btn_remove;
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
		vv.clear ();
		foreach (var v in env.variables) {
			vv.add (v);
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
