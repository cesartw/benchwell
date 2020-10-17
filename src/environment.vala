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
	public Gtk.Entry name;
	public Gtk.Button btn_save;
	public Gtk.Button btn_remove;

	public EnvironmentPanel (Benchwell.Environment env) {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		name = new Gtk.Entry ();
		name.set_text (env.name);
		name.set_placeholder_text ("Name");
		name.show ();

		btn_save = new Gtk.Button.with_label (_("Save"));
		btn_save.show ();

		btn_remove = new Gtk.Button.with_label (_("Remove"));
		btn_remove.get_style_context ().add_class ("destructive-action");
		btn_remove.show ();

		var btn_box = new Gtk.ButtonBox (Gtk.Orientation.HORIZONTAL);
		btn_box.pack_end(btn_remove, false, false, 5);
		btn_box.pack_end(btn_save, false, false, 5);
		btn_box.show ();

		var vbox = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		vbox.show ();

		var vv = new Benchwell.KeyValues ();
		vv.clear ();
		foreach (var v in env.variables) {
			vv.add (v);
		}
		vbox.pack_start (name, false, false, 5);
		vbox.pack_start (vv, true, true, 5);
		vbox.pack_end (btn_box, false, false, 5);
		vv.show ();

		pack_start (vbox, true, true, 5);
	}
}
