[Compact]
public struct Benchwell.Token
{
	uint start;
	uint end;
	string val;
	int type;

	public string to_string () {
		return @"$(val):$(type):$(start):$(end)";
	}
}

public enum Benchwell.TokenType {
	STRING,
	NUMBER,
	FLOAT,
	KEYWORD,
	IDENTIFIER
}

public interface Benchwell.TokenTypeLexer : Object {
	public abstract bool check (char c);
	public abstract bool consume (out Benchwell.Token token, out uint taken, uint position, owned string input, char espace_char);
}

public class Benchwell.Lexer : Object {
	//public bool     double_quoted_string  { get; construct; default = true; }
	//public bool     single_quoted_string  { get; construct; default = true; }
	//public bool     bare_string           { get; construct; default = true; }
	//public string   bare_string_separator { get; construct; default = " "; }
	//public bool     numbers               { get; construct; default = true; }
	//public bool     float_numbers         { get; construct; default = true; }
	public char     escape_char           { get; construct; default = '\\'; }
	//public string[] keywords              { get; construct; }
	//public string   identifier            { get; construct; default = "@acbcdefghijklmnopqrstuvxwyzACBCDEFGHIJKLMNOPQRSTUVXWYZ0123456789_"; }

	private Benchwell.TokenTypeLexer[] token_types;

	private uint position { get; set; }

	public Lexer (Benchwell.TokenTypeLexer[] tt) {
		token_types = tt;
	}

	public Benchwell.Token[]? parse (string input) {
		position = 0;
		Benchwell.Token[] tokens = {};

		while (true) {
			if (position >= input.length - 1) {
				return tokens;
			}

			var c = input.to_utf8 ()[position];
			bool found = false;
			foreach (var tt in token_types) {
				if (!tt.check (c)) {
					continue;
				}

				Benchwell.Token token;
				uint taken = 0;
				if (!tt.consume (out token, out taken, position, input, escape_char)) {
					continue;
				}

				token.val = input[token.start:token.end+1];
				found = true;
				position += taken;
				tokens += token;
				break;
			}
			if (!found) {
				break;
			}
		}

		return tokens;
	}

	public static int main(string[] args) {
		var lx = new Lexer ({
			new Benchwell.TemplateStart (),
			new Benchwell.TemplateEnd (),
			new Benchwell.BareString ()
		});
		var tokens = lx.parse ("http://{{token}}/desk/api/v1/tickets.json");

		foreach (var t in tokens) {
			print (@"=====$(t.start):$(t.end):$(t.val)\n");
		}

		return 0;
	}
}

public class Benchwell.TemplateStart : Benchwell.TokenTypeLexer, Object {
	public bool check (char c) {
		return c == '{';
	}

	public bool consume (out Benchwell.Token token, out uint taken, uint position, owned string input, char espace_char) {
		taken = 0;
		var chars = input.to_utf8 ();

		if (input.length < 2){
			return false;
		}

		if (chars[position] != '{' || chars[position+1] != '{') {
			return false;
		}

		taken = 2;
		token = Benchwell.Token () {
			start = position,
			end = position + taken - 1,
			type = Benchwell.TokenType.STRING
		};
		return true;
	}
}

public class Benchwell.TemplateEnd : Benchwell.TokenTypeLexer, Object {
	public bool check (char c) {
		return c == '}';
	}

	public bool consume (out Benchwell.Token token, out uint taken, uint position, owned string input, char espace_char) {
		taken = 0;
		var chars = input.to_utf8 ();

		if (input.length < 2){
			return false;
		}

		if (chars[position] != '}' || chars[position+1] != '}') {
			return false;
		}

		taken = 2;
		token = Benchwell.Token () {
			start = position,
			end = position+taken - 1,
			type = Benchwell.TokenType.STRING
		};
		return true;
	}
}

public class Benchwell.BareString : Benchwell.TokenTypeLexer, Object {
	public bool check (char c) {
		return c != '{';
	}

	public bool consume (out Benchwell.Token token, out uint taken, uint position, owned string input, char espace_char) {
		taken = 0;
		var chars = input.to_utf8 ();
		while (position < input.length && chars[position] != '{' && chars[position] != '}') {
			taken++;
			position++;
		}

		token = Benchwell.Token () {
			start = position-taken,
			end = position - 1,
			type = Benchwell.TokenType.STRING
		};

		return true;
	}
}
