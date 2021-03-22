public enum Benchwell.KEYVALUETYPES {
	String,
	File,
	Multiline
}

public interface Benchwell.KeyValueI : Object {
	public abstract int64  id                      { get; set; }
	public abstract string key                     { get; set; }
	public abstract string val                     { get; set; }
	public abstract bool   enabled                 { get; set; }
	public abstract Benchwell.KEYVALUETYPES kvtype { get; set; }
}

