// modules: gtk+-3.0
// vapidirs: vapi

public class Benchwell.EnvironmentEditor : Gtk.Paned {
	public Benchwell.Button btn_add;
	public Benchwell.Button btn_remove;
	public Benchwell.Button btn_clone;
	private Gtk.Stack stack;

	public EnvironmentEditor () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL
		);

		//btn_add = new Benchwell.Button ("white-add", Gtk.IconSize.SMALL_TOOLBAR);
		//btn_add.halign = Gtk.Align.START;
		//btn_add.show ();

		//btn_clone = new Benchwell.Button ("white-copy", Gtk.IconSize.SMALL_TOOLBAR);
		//btn_clone.halign = Gtk.Align.START;
		//btn_clone.show ();

		//btn_remove = new Benchwell.Button ("white-close", Gtk.IconSize.SMALL_TOOLBAR);
		//btn_remove.halign = Gtk.Align.END;
		//btn_remove.show ();

		//var btn_box = new Gtk.ButtonBox (Gtk.Orientation.HORIZONTAL);
		//btn_box.name = "EnvActions";
		//btn_box.get_style_context ().add_class("linked");
		//btn_box.layout_style = Gtk.ButtonBoxStyle.SPREAD;
		//btn_box.hexpand = false;
		//btn_box.homogeneous = false;
		//btn_box.spacing = 0;
		//btn_box.add (btn_add);
		//btn_box.add (btn_clone);
		//btn_box.add (btn_remove);
		//btn_box.height_request = 5;
		//btn_box.show ();

		var switcher = new Gtk.StackSwitcher ();
		switcher.orientation = Gtk.Orientation.VERTICAL;
		switcher.vexpand = true;
		switcher.hexpand = true;
		switcher.show ();

		stack = new Gtk.Stack ();
		stack.homogeneous = true;
		stack.vexpand = true;
		stack.hexpand = true;
		stack.transition_type = Gtk.StackTransitionType.SLIDE_UP;
		stack.show ();

		switcher.stack = stack;

		Config.environments.for_each ((item) => {
			var env = item as Environment;
			var panel = new Benchwell.EnvironmentPanel (env);
			panel.show ();
			stack.add_titled (panel, env.name, env.name);
			panel.entry_name.changed.connect (() => {
				stack.child_set_property(panel, "title", panel.entry_name.text);
				env.name = panel.entry_name.text;
			});
			return false;
		});

		//var env_list_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		//env_list_box.pack_start (switcher, true, true, 0);
		//env_list_box.pack_end (btn_box, false, false, 0);
		//env_list_box.show ();

		var sww = new Gtk.ScrolledWindow (null, null);
		sww.hscrollbar_policy = Gtk.PolicyType.NEVER;
		sww.show ();
		var sws = new Gtk.ScrolledWindow (null, null);
		sws.hscrollbar_policy = Gtk.PolicyType.NEVER;
		sws.show ();

		sww.add (switcher);
		sws.add (stack);

		pack1 (sww, false, true);
		pack2 (sws, true, false);

		//btn_add.clicked.connect (on_add_env);
		//btn_remove.clicked.connect (on_remove_env);
		//btn_clone.clicked.connect (on_clone);

		var menu = new Gtk.Menu ();
		var add_menu = new Benchwell.MenuItem (_("Add"), "add");
		add_menu.activate.connect (on_add_env);
		add_menu.show ();

		var remove_menu = new Benchwell.MenuItem (_("Remove"), "close");
		remove_menu.activate.connect (on_remove_env);
		remove_menu.show ();

		var clone_menu = new Benchwell.MenuItem (_("Copy"), "copy");
		clone_menu.activate.connect (on_clone);
		clone_menu.show ();

		menu.add (add_menu);
		menu.add (remove_menu);
		menu.add (clone_menu);

		switcher.button_press_event.connect ( (list, event) => {
			if (event.button != Gdk.BUTTON_SECONDARY) {
				return false;

			}
			menu.show ();
			menu.popup_at_pointer (event);

			return false;
		});

		Config.environments.added.connect ((env) => {
			env.name = @"Environment #$(Config.environments.length ())";
			var panel = new Benchwell.EnvironmentPanel (env);
			panel.show ();
			stack.add_titled (panel, env.name, env.name);
			stack.set_visible_child (panel);
			panel.entry_name.changed.connect (() => {
				stack.child_set_property(panel, "title", panel.entry_name.text);
				env.name = panel.entry_name.text;
			});
		});
	}

	private void on_add_env () {
		var env = new Benchwell.Environment();
		Config.environments.add (env);
	}

	private void on_clone () {
		try {
			var index = stack.get_children ().index (stack.get_visible_child ());
			Config.environments.at(index).clone ();
		} catch (ConfigError err) {
			Config.show_alert (this, err.message);
		}
	}

	private void on_remove_env () {
		var panel = stack.get_visible_child ();
		var index = stack.get_children ().index (panel);
		try {
			var env = Config.environments.at(index);
			env.remove ();
		} catch(ConfigError err) {
			Config.show_alert (this, err.message);
		}

		stack.remove (panel);
	}

	public void select_env (Benchwell.Environment env) {
		stack.set_visible_child_name (env.name);
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

		keyvalues = new Benchwell.KeyValues (Benchwell.KeyValueTypes.STRING | Benchwell.KeyValueTypes.MULTILINE);
		keyvalues.no_row_left.connect (on_row_added);
		if (env.variables.length > 0) {
			keyvalues.clear ();
			foreach (var v in env.variables) {
				keyvalues.add (v as Benchwell.KeyValueI);
			}
		}

		vbox.pack_start (entry_name, false, false, 5);
		vbox.pack_start (keyvalues, true, true, 5);
		keyvalues.show ();

		pack_start (vbox, true, true, 5);

		// signals
		entry_name.changed.connect (on_save);

		if (keyvalues.get_children ().length () == 0) {
			on_row_added ();
		}
	}

	private void on_row_added () {
		try {
			var kv = (Benchwell.KeyValueI) environment.add_variable ();
			keyvalues.add (kv);
		} catch (ConfigError err) {
			Config.show_alert (this, err.message);
		}
	}

	private void on_save () {
		environment.name = entry_name.text;

		//try {
			//environment.save ();
		//} catch (Error err) {
			//Config.show_alert (this, err.message);
		//}
	}
}
