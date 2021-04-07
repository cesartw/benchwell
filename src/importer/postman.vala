public class Benchwell.ImporterPostman : Object {
	public class Item {
		public string name { get; set; }
	}

	public class Request : Object {
		public string method { get; set; }
	}

	public class RequestAuth : Object, Json.Serializable {
		public string auth_type { get; set; }
		public HashTable <string, string>[] bearer { get; set; }

		public unowned ParamSpec? find_property (string property_name) {
			if (property_name == "type")
				return ((ObjectClass) get_type ().class_ref ()).find_property ("auth_type");
			return ((ObjectClass) get_type ().class_ref ()).find_property (property_name);
		}
	}

	public class Url : Object {
		public string raw { get; set; }
		public string protocol { get; set; }
		public string[] host { get; set; }
		public string[] path { get; set; }
		public KV[] query { get; set; }
	}

	public class KV : Object {
		public string key { get; set; }
		public string @value { get; set; }
	}
}
