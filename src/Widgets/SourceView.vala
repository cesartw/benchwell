public class Benchwell.SourceView : Gtk.SourceView {
	public SourceView (string lang = "sql", string style = "oblivion") {
		Object (
			show_line_numbers: false,
			show_right_margin: true,
			hexpand: true,
			vexpand: true
		);

		var buffer = (Gtk.SourceBuffer) get_buffer ();
		var lm = Gtk.SourceLanguageManager.get_default ();
		buffer.set_language (lm.get_language (lang));

		var sm = Gtk.SourceStyleSchemeManager.get_default ();
		buffer.set_style_scheme (sm.get_scheme (style));
	}
}
