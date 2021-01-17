function call(string_length, type, addon_set) {
	if (string_length) {
		string_legnth = 5;
	}
	if (!addon_set) {
		addon_set = "";
	}

	var set = "";
	switch (type) {
		case "number":
			//(Math.floor (Math.random () * string_length)+"").padStart (string_length);
			set = '0123456789';
			break;
		case "alpha":
			set = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
			break;
		case "letter":
			set = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
			break;
		case "upper-letter":
			set = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
			break;
		case "lower-letter":
			set = 'abcdefghijklmnopqrstuvwxyz';
			break;
		default:
			set = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
			break;
	}

	set += addon_set;

    var str = '';
    for (var i = 0; i < string_length; i++) {
        str += set.charAt(Math.floor(Math.random() * set.length));
    }

    return str;
}
