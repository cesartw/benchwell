public errordomain Benchwell.ImportError {
	BASE
}

public interface Benchwell.Importer : Object {
	public abstract Gtk.FileFilter get_file_filter ();
	public abstract void import (InputStream source) throws Benchwell.ImportError, Benchwell.ConfigError;
}
