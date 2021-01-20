namespace Benchwell {
	public const string[] on_white = {
		"#ff7305",
		"#202329",
		"#9d9d9c",
		"#dadada",
		"#9dfc62",
		"#575756",
		"#2bb6c7", // gdk.NewRGBA(43.0, 182.0, 199.0, 1),
		"#e68700", // gdk.NewRGBA(230.0, 135.0, 0.0, 1),
		"#de8f00", // gdk.NewRGBA(222.0, 143.0, 0.0, 1),
		"#ff0000", // gdk.NewRGBA(255.0, 0, 0.0, 1),
		"#27d000", // gdk.NewRGBA(39.0, 208.0, 0.0, 1),
		"#19a800", // gdk.NewRGBA(25.0, 168.0, 0.0, 1),
		"#298700" // gdk.NewRGBA(41.0, 135.0, 0.0, 1),
	};

	public const string[] on_dark  = {
		"#ff7305",
		"#575757",
		"#9d9d9c",
		"#ffffff",
		"#9dfc62",
		"#575756",
		"#2bb6c7", // gdk.NewRGBA(43.0, 182.0, 199.0, 1),
		"#e68700", // gdk.NewRGBA(230.0, 135.0, 0.0, 1),
		"#de8f00", // gdk.NewRGBA(222.0, 143.0, 0.0, 1),
		"#ff0000", // gdk.NewRGBA(255.0, 0, 0.0, 1),
		"#27d000", // gdk.NewRGBA(39.0, 208.0, 0.0, 1),
		"#19a800", // gdk.NewRGBA(25.0, 168.0, 0.0, 1),
		"#298700" // gdk.NewRGBA(41.0, 135.0, 0.0, 1),
	};
	public const string[] on_main  = {
		"#202329",
		"#ffffff",
		"#9d9d9c",
		"#ffffff",
		"#9dfc62",
		"#575756"
	};
	public const string null_string = "<NULL>";


	public const string[] Methods = {
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
		"HEAD",
	};
}

public enum Benchwell.ColorScheme {
	LIGHT,
	DARK,
	MAIN
}


public enum Benchwell.Colors {
	MAIN,
	BLACK,
	GREY,
	LIGHTGREY,
	NULLHL,
	PKHL,
	POST,
	PATCH,
	PUT,
	DELETE,
	GET,
	HEAD,
	OPTIONS;

	public string to_string (Benchwell.ColorScheme scheme = Benchwell.ColorScheme.DARK) {
		string[] colors = Benchwell.on_white;
		switch (scheme) {
			case Benchwell.ColorScheme.LIGHT:
				colors = Benchwell.on_white;
				break;
			case Benchwell.ColorScheme.DARK:
				colors = Benchwell.on_dark;
				break;
			case Benchwell.ColorScheme.MAIN:
				colors = Benchwell.on_main;
				break;
		}

		return colors[this];
	}

	public static Benchwell.Colors? parse(string? s) {
		if  (s == null)
			return GET;

		// NOTE: "".casefold is not supported as `case` option
		switch (s.casefold ()) {
			case "POST", "post":
				return POST;
			case "PATCH", "patch":
				return PATCH;
			case "PUT", "put":
				return PUT;
			case "DELETE", "delete":
				return DELETE;
			case "GET", "get":
				return GET;
			case "HEAD", "head":
				return HEAD;
			case "OPTIONS", "options":
				return OPTIONS;
		}

		return null;
	}
}

public enum Benchwell.Settings {
	WINDOW_SIZE_W,
	WINDOW_SIZE_H,
	WINDOW_POS_X,
	WINDOW_POS_Y,
	ENVIRONMENT_ID,
	HTTP_COLLECTION_ID,
	HTTP_ITEM_ID;

	public string to_string () {
		switch (this) {
			case WINDOW_SIZE_W:
				return "window-size-w";
			case WINDOW_SIZE_H:
				return "window-size-h";
			case WINDOW_POS_X:
				return "window-pos-x";
			case WINDOW_POS_Y:
				return "window-pos-y";
			case ENVIRONMENT_ID:
				return "environment-id";
			case HTTP_COLLECTION_ID:
				return "http-collection-id";
			case HTTP_ITEM_ID:
				return "http-item-id";
			default:
				return "";
		}
	}
}
