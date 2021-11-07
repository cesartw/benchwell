public class Benchwell.SettingsPanel : Gtk.Box {
	public weak Benchwell.ApplicationWindow window { get; construct; }
	private Gtk.Notebook notebook;
	private Benchwell.EnvironmentEditor env_editor;
	private Benchwell.EditorSettings editor_settings;
	private Benchwell.HttpSettings http_settings;
	private Benchwell.PomodoroSettings pomodoro_settings;
	private Benchwell.About about;
	private Gtk.Switch dark_switch;

	public SettingsPanel (Benchwell.ApplicationWindow w) {
		Object (
			window: w,
			orientation: Gtk.Orientation.VERTICAL,
			vexpand: true,
			hexpand: true
		);

		set_name ("settings-panel");
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

		editor_settings = new Benchwell.EditorSettings ();
		editor_settings.show ();

		http_settings = new Benchwell.HttpSettings (w);
		http_settings.show ();

		pomodoro_settings = new Benchwell.PomodoroSettings ();
		pomodoro_settings.show ();

		about = new Benchwell.About (w);
		about.show ();

		env_editor.get_style_context ().add_class ("bw-spacing");
		editor_settings.get_style_context ().add_class ("bw-spacing");
		http_settings.get_style_context ().add_class ("bw-spacing");
		pomodoro_settings.get_style_context ().add_class ("bw-spacing");
		about.get_style_context ().add_class ("bw-spacing");

		notebook.append_page (env_editor, new Gtk.Label (_("Environments")));
		notebook.append_page (editor_settings, new Gtk.Label (_("Editor")));
		notebook.append_page (http_settings, new Gtk.Label (_("HTTP")));
		notebook.append_page (pomodoro_settings, new Gtk.Label (_("Pomodoro")));
		notebook.append_page (about, new Gtk.Label (_("About")));

		header_bar.pack_end (dark_switch);
		header_bar.pack_end (dark_icon);
		pack_start (header_bar, false, false, 0);
		pack_start (notebook, true, true, 0);

		dark_switch.state = Config.settings.dark_mode;
		dark_switch.state_set.connect ((state) => {
			Config.settings.dark_mode = state;
			Gtk.Settings.get_default ().gtk_application_prefer_dark_theme = state;
			return false;
		});
	}

	public void select_env (Benchwell.Environment env) {
		env_editor.select_env (env);
	}
}

public class Benchwell.EditorSettings : Gtk.Grid {
	public EditorSettings () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);

		var look_and_feel_col = new Benchwell.SettingColumn (_("Look & Feel"));
		look_and_feel_col.show ();

		var editing_col = new Benchwell.SettingColumn (_("Editing"));
		editing_col.show ();

		// EDITOR THEME
		var laf_theme_combo = new Gtk.ComboBoxText ();
		foreach (var id in Gtk.SourceStyleSchemeManager.get_default ().scheme_ids) {
			laf_theme_combo.append (id, id);
		}
		laf_theme_combo.set_active_id (Config.settings.editor_theme);
		laf_theme_combo.show ();

		var laf_theme_label = new Gtk.Label (_("Theme"));
		laf_theme_label.show ();
		look_and_feel_col.add (laf_theme_label, laf_theme_combo);
		///////////////

		// FONT
		var laf_font_label = new Gtk.Label (_("Font"));
		laf_font_label.show ();

		var laf_font_btn = new Gtk.FontButton ();
		if (Config.settings.editor_font != "") {
			laf_font_btn.font = Config.settings.editor_font;
		}
		laf_font_btn.show ();
		look_and_feel_col.add (laf_font_label, laf_font_btn);
		///////

		// EDITING
		// TAB WIDTH
		var editing_tabwidth_label = new Gtk.Label (_("Tab width"));
		editing_tabwidth_label.show ();

		var editing_tabwidth_spin = new Gtk.SpinButton.with_range (2, 8, 2);
		editing_tabwidth_spin.value = (double) Config.settings.editor_tab_width;
		editing_tabwidth_spin.show ();
		editing_col.add (editing_tabwidth_label, editing_tabwidth_spin);
		////////////

		// SHOW LINE NUMBER
		var editing_ln_label = new Gtk.Label (_("Show line number"));
		editing_ln_label.show ();

		var editing_ln_sw = new Gtk.Switch ();
		editing_ln_sw.state = Config.settings.editor_line_number;
		editing_ln_sw.halign = Gtk.Align.START;
		editing_ln_sw.show ();
		editing_col.add (editing_ln_label, editing_ln_sw);
		///////////////////

		// HIGHLIGHT CURRENT LINE
		var editing_hl_label = new Gtk.Label (_("Highlight current line"));
		editing_hl_label.show ();

		var editing_hl_sw = new Gtk.Switch ();
		editing_hl_sw.state = Config.settings.editor_highlight_line;
		editing_hl_sw.halign = Gtk.Align.START;
		editing_hl_sw.show ();
		editing_col.add (editing_hl_label, editing_hl_sw);
		/////////////////////////

		// SPACES FOR TABS
		var editing_notabs_label = new Gtk.Label (_("Insert spaces instead of tabs"));
		editing_notabs_label.show ();

		var editing_notabs_sw = new Gtk.Switch ();
		editing_notabs_sw.state = Config.settings.editor_no_tabs;
		editing_notabs_sw.halign = Gtk.Align.START;
		editing_notabs_sw.show ();
		editing_col.add (editing_notabs_label, editing_notabs_sw);
		//////////////////

		attach (look_and_feel_col, 0, 0, 1, 1);
		attach (editing_col, 1, 0, 1, 1);

		laf_theme_combo.changed.connect (() => {
			Config.settings.editor_theme = laf_theme_combo.get_active_id ();
		});

		laf_font_btn.font_set.connect (() => {
			Config.settings.editor_font = laf_font_btn.font;
		});

		editing_tabwidth_spin.changed.connect (() => {
			Config.settings.editor_tab_width = (int)editing_tabwidth_spin.value;
		});

		editing_ln_sw.state_set.connect ((state) => {
			Config.settings.editor_line_number = state;
			return false;
		});

		editing_hl_sw.state_set.connect ((state) => {
			Config.settings.editor_highlight_line = state;
			return false;
		});

		editing_notabs_sw.state_set.connect ((state) => {
			Config.settings.editor_no_tabs = state;
			return false;
		});
	}
}

public class Benchwell.SettingColumn : Gtk.Grid {
	private int current_row = 1;

	public class SettingColumn (string title) {
		Object (
			column_spacing: 5,
			row_spacing: 5
		);

		var lbl = new Gtk.Label(title) {
			valign = Gtk.Align.START,
			halign = Gtk.Align.START
		};
		lbl.show ();
		attach (lbl, 0, 0, 2, 1);
	}

	public new void add (Gtk.Widget? lbl, Gtk.Widget? w) {
		if (lbl != null) {
			lbl.valign = Gtk.Align.START;
			lbl.halign = Gtk.Align.START;
			attach (lbl, 1, current_row, 1, 1);
		}

		if (w != null)
			attach (w, 2, current_row, 1, 1);

		if (w != null || lbl != null)
			current_row += 1;
	}
}

public class Benchwell.HttpSettings : Gtk.Grid {
	public weak Benchwell.ApplicationWindow window { get; construct; }
	private Benchwell.ImporterInsomnia importer_insomnia;

	public HttpSettings (Benchwell.ApplicationWindow w) {
		Object (
			window: w,
			row_spacing: 5,
			column_spacing: 5
		);

		importer_insomnia = new Benchwell.ImporterInsomnia ();

		// LOOK & FEEL
		var look_and_feel_col = new Benchwell.SettingColumn (_("Look & Feel"));
		look_and_feel_col.show ();

		var other_col = new Benchwell.SettingColumn (_("Requests"));
		other_col.show ();

		var import_col = new Benchwell.SettingColumn (_("Import"));
		import_col.show ();


		// FONT
		var laf_font_label = new Gtk.Label (_("Font"));
		laf_font_label.show ();

		var laf_font_btn = new Gtk.FontButton ();
		if (Config.settings.http_font != "") {
			laf_font_btn.font = Config.settings.http_font;
		}
		laf_font_btn.show ();
		look_and_feel_col.add (laf_font_label, laf_font_btn);
		///////

		// ACTIVATE ON SINGLE CLICK
		var other_single_click_label = new Gtk.Label (_("Open on single click"));
		other_single_click_label.show ();

		var other_single_click_sw = new Gtk.Switch ();
		other_single_click_sw.state = Config.settings.http_single_click_activate;
		other_single_click_sw.show ();
		other_col.add (other_single_click_label, other_single_click_sw);
		//////////////////

		// FOLLOW REDIRECT
		var follow_redirect_click_label = new Gtk.Label (_("Follow redirect"));
		follow_redirect_click_label.show ();

		var follow_redirect_click_sw = new Gtk.Switch ();
		follow_redirect_click_sw.state = Config.settings.http_follow_redirect;
		follow_redirect_click_sw.show ();
		other_col.add (follow_redirect_click_label, follow_redirect_click_sw);
		//////////////////

		// IMPORTERS
		var other_import_type = new Gtk.ComboBoxText ();
		other_import_type.append ("insomnia", "Insomnia v4(JSON)");
		other_import_type.set_active_id ("insomnia");
		other_import_type.show ();

		var import_dialog = new Gtk.FileChooserDialog (_("Select Insomnia v4(JSON) file"), window,
											 Gtk.FileChooserAction.OPEN,
											_("Import"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		import_dialog.add_filter (importer_insomnia.get_file_filter ());

		var other_import_file_btn = new Gtk.FileChooserButton.with_dialog (import_dialog);
		other_import_file_btn.show ();

		var other_import_spinner = new Gtk.Spinner ();
		other_import_spinner.show ();
		import_col.add (other_import_type, other_import_file_btn);
		import_col.add (other_import_spinner, null);
		////////////

		attach (look_and_feel_col, 0, 0, 1, 1);
		attach (other_col, 1, 0, 1, 1);
		attach (import_col, 1, 1, 1, 1);
		//attach (laf_label, 0, 0, 2, 1);
		//attach (laf_font_label, 1, 1, 1, 1);
		//attach (laf_font_btn, 2, 1, 1, 1);

		//attach (other_label, 3, 0, 2, 1);
		//attach (other_single_click_label, 4, 1, 1, 1);
		//attach (other_single_click_sw, 5, 1, 1, 1);

		//attach (follow_redirect_click_label, 4, 2, 1, 1);
		//attach (follow_redirect_click_sw, 5, 2, 1, 1);

		//attach (other_import_label, 3, 3, 2, 1);
		//attach (other_import_type, 4, 4, 1, 1);
		//attach (other_import_file_btn, 5, 4, 1, 1);
		//attach (other_import_spinner, 6, 4, 1, 1);

		laf_font_btn.font_set.connect (() => {
			Config.settings.http_font = laf_font_btn.font;
		});

		other_single_click_sw.state_set.connect ((state) => {
			Config.settings.http_single_click_activate = state;
			return false;
		});

		follow_redirect_click_sw.state_set.connect ((state) => {
			Config.settings.http_follow_redirect = state;
			return false;
		});

		other_import_file_btn.file_set.connect (() => {
			Benchwell.Importer importer = null;
			switch (other_import_type.get_active_id ()) {
				case "insomnia":
					importer = importer_insomnia;
					break;
			}

			if (importer == null)
				return;

			other_import_spinner.start ();
			other_import_file_btn.sensitive = false;
			other_import_type.sensitive = false;
			import.begin (importer, other_import_file_btn.get_file (), (obj, res) => {
				var ok = import.end (res);
				if (ok)
					Config.show_alert (this, "Done", Gtk.MessageType.INFO, true, 5000);
				other_import_file_btn.unselect_all ();
				other_import_spinner.stop ();
				other_import_file_btn.sensitive = true;
				other_import_type.sensitive = true;
			});
		});
	}

	private async bool import (Benchwell.Importer importer, File file) {
		SourceFunc callback = import.callback;

		bool ok = false;
		ThreadFunc<bool> run = () => {
			try {
				importer.import (file.read (null));
				ok = true;
			} catch (Benchwell.ImportError err) {
				Config.show_alert (this, err.message);
			} catch (GLib.Error err) {
				Config.show_alert (this, err.message);
			}

			Idle.add((owned) callback);
			return true;
		};

		new Thread<bool>("benchwell-http-import", run);
		yield;
		return ok;
	}
}

public class Benchwell.PomodoroSettings : Gtk.Grid {
	public PomodoroSettings () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);

		var timing_col = new Benchwell.SettingColumn (_("Timings"));
		timing_col.show ();

		// WORK DURATION
		var duration_label = new Gtk.Label (_("Work session duration"));
		duration_label.show ();

		var duration_spin = new Gtk.SpinButton.with_range (20, 60, 5);
		duration_spin.value = (double)Config.settings.pomodoro_duration / 60;
		duration_spin.show ();
		timing_col.add (duration_label, duration_spin);
		////////////////

		// BREAK DURATION
		var break_duration_label = new Gtk.Label (_("Short break duration"));
		break_duration_label.show ();

		var break_duration_spin = new Gtk.SpinButton.with_range (5, 10, 1);
		break_duration_spin.value = (double)Config.settings.pomodoro_break_duration / 60;
		break_duration_spin.show ();
		timing_col.add (break_duration_label, break_duration_spin);
		/////////////////

		// LONG BREAK DURATION
		var long_break_duration_label = new Gtk.Label (_("Long break duration"));
		long_break_duration_label.show ();

		var long_break_duration_spin = new Gtk.SpinButton.with_range (10, 20, 1);
		long_break_duration_spin.value = (double)Config.settings.pomodoro_long_break_duration / 60;
		long_break_duration_spin.show ();
		timing_col.add (long_break_duration_label, long_break_duration_spin);
		/////////////////

		// SET COUNT
		var set_count_label = new Gtk.Label (_("Sets before long break"));
		set_count_label.show ();

		var set_count_spin = new Gtk.SpinButton.with_range (10, 20, 1);
		set_count_spin.value = (double)Config.settings.pomodoro_set_count;
		set_count_spin.show ();
		timing_col.add (set_count_label, set_count_spin);
		/////////////////

		attach (timing_col, 0, 0, 2, 1);

		duration_spin.changed.connect (() => {
			Config.settings.pomodoro_duration = (int)duration_spin.value*60;
		});

		break_duration_spin.changed.connect (() => {
			Config.settings.pomodoro_break_duration = (int)break_duration_spin.value*60;
		});

		long_break_duration_spin.changed.connect (() => {
			Config.settings.pomodoro_long_break_duration = (int)long_break_duration_spin.value*60;
		});

		set_count_spin.changed.connect (() => {
			Config.settings.pomodoro_set_count = (int)set_count_spin.value;
		});
	}
}

public class Benchwell.About : Gtk.Grid {
	public weak Benchwell.ApplicationWindow window { get; construct; }

	public About (Benchwell.ApplicationWindow w) {
		Object (
			window: w,
			row_spacing: 5,
			column_spacing: 5,
			halign: Gtk.Align.CENTER
		);

		var logo = new Gtk.Image.from_icon_name ("io.benchwell", Gtk.IconSize.DIALOG);
		logo.show ();

		attach (logo, 0, 0, 1, 1);
	}
}
