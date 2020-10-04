namespace Benchwell {
						              // orange, dark, grey, lightgrey, null-hl, pk-hl
	public const string[] on_white = {"#ff7305", "#202329", "#9d9d9c", "#dadada", "#9dfc62", "#ffaf72"};
	public const string[] on_dark  = {"#ff7305", "#575757", "#9d9d9c", "#ffffff", "#9dfc62", "#ffaf72"};
	public const string[] on_main  = {"#202329", "#ffffff", "#9d9d9c", "#ffffff", "#9dfc62", "#ffaf72"};
	public const string null_string = "<NULL>";

}

public enum Benchwell.ColorScheme {
	Light,
	Dark,
	Main
}


public enum Benchwell.Colors {
	Main,
	Black,
	Grey,
	LightGrey,
	NullHL,
	PkHL;

	public string to_string (Benchwell.ColorScheme scheme = Benchwell.ColorScheme.Dark) {
		string[] colors = Benchwell.on_white;
		switch (scheme) {
			case Benchwell.ColorScheme.Light:
				colors = Benchwell.on_white;
				break;
			case Benchwell.ColorScheme.Dark:
				colors = Benchwell.on_dark;
				break;
			case Benchwell.ColorScheme.Main:
				colors = Benchwell.on_main;
				break;
		}

		return colors[this];
	}
}
