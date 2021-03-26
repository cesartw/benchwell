public class Benchwell.SettingsPanel : Gtk.Box {
	private Gtk.Notebook notebook;
	private Benchwell.EnvironmentEditor env_editor;
	private Benchwell.SettingsEditor editor_settings;
	private Gtk.Switch dark_switch;

	public SettingsPanel () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			vexpand: true,
			hexpand: true
		);

		var header_bar = new Gtk.HeaderBar ();
		header_bar.title = _("Settings");
		header_bar.show ();

		dark_switch = new Gtk.Switch ();
		dark_switch.show();
		var dark_icon = new Gtk.Image.from_icon_name ("weather-clear-night", Gtk.IconSize.SMALL_TOOLBAR);
		dark_icon.show ();


		notebook = new Gtk.Notebook ();
		notebook.tab_pos = Gtk.PositionType.TOP;
		notebook.show ();

		env_editor = new Benchwell.EnvironmentEditor ();
		env_editor.show ();

		editor_settings = new Benchwell.SettingsEditor ();
		editor_settings.show ();

		notebook.append_page (env_editor, new Gtk.Label (_("Environments")));
		notebook.append_page (editor_settings, new Gtk.Label (_("Editor")));

		header_bar.pack_end (dark_switch);
		header_bar.pack_end (dark_icon);
		pack_start (header_bar, false, false, 0);
		pack_start (notebook, true, true, 0);

		dark_switch.state = Config.settings.get_boolean ("dark-mode");
		dark_switch.state_set.connect ((state) => {
			Config.settings.set_boolean ("dark-mode", state);
			Gtk.Settings.get_default ().gtk_application_prefer_dark_theme = state;
			return false;
		});
	}

	public void select_env (Benchwell.Environment env) {
		env_editor.select_env (env);
	}
}

public class Benchwell.SettingsEditor : Gtk.Box {
	private Gtk.SourceStyleSchemeManager stylemanager;

	public SettingsEditor () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			spacing: 5
		);

		// LOOK & FEEL
		var laf_label = new Gtk.Label ("");
		laf_label.set_markup (_("<b>Look &amp; Feel</b>"));
		laf_label.show ();

		var laf_frame = new Gtk.Frame (null);
		laf_frame.set_label_widget (laf_label);
		laf_frame.shadow_type = Gtk.ShadowType.NONE;
		laf_frame.show ();

		var laf_theme_combo = new Gtk.ComboBoxText ();
		stylemanager = Gtk.SourceStyleSchemeManager.get_default ();
		foreach (var id in stylemanager.scheme_ids) {
			laf_theme_combo.append (id, id);
		}
		laf_theme_combo.set_active_id (Config.settings.get_string ("editor-theme"));
		laf_theme_combo.show ();

		var laf_grid = new Gtk.Grid () {
			orientation = Gtk.Orientation.VERTICAL,
			margin_top = margin_bottom = 5
		};
		var laf_theme_label = new Gtk.Label (_("Theme"));
		laf_theme_label.show ();

		laf_grid.attach (new Gtk.Separator (Gtk.Orientation.HORIZONTAL), 0, 0, 1, 1);
		laf_grid.attach (laf_theme_label, 1, 0, 1, 1);
		laf_grid.attach (laf_theme_combo, 2, 0, 1, 1);
		laf_grid.show_all ();

		laf_frame.add (laf_grid);


		laf_theme_combo.changed.connect (() => {
			Config.settings.set_string ("editor-theme", laf_theme_combo.get_active_id ());
		});
		//});
		//////////////
	}
}
