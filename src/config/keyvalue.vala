[Flags]
public enum Benchwell.KeyValueTypes {
	STRING,
	MULTILINE,
	FILE
}

public interface Benchwell.KeyValueI : Object {
	public abstract int64  id                      { get; set; }
	public abstract string key                     { get; set; }
	public abstract string val                     { get; set; }
	public abstract bool   enabled                 { get; set; }
	public abstract Benchwell.KeyValueTypes kvtype { get; set; }
}

