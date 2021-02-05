public errordomain Benchwell.ImportError {
	BASE
}

public interface Benchwell.Importer : Object {
	public abstract void import (string source) throws Benchwell.ImportError;
}

public class Benchwell.InsomniaImpoter : Benchwell.Importer, Object {

	public class Resource : Object {
		public string id {get;set;}
		public string resource_type {get;set;}
		public string parentId {get;set;}
		public string url {get;set;}
		public string name {get;set;}
		public string description {get;set;}
		public string method {get;set;}
		public ResourceBody request_body {get;set;}
		public ResourceKV[] request_params {get;set;}
		public ResourceKV[] request_headers {get;set;}
		public Resource[] resources {get;set;}
	}

	public class ResourceBody : Object {
		public string mimeType {get;set;}
		public string text {get;set;}
	}

	public class ResourceKV : Object {
		public string name {get;set;}
		public string value {get;set;}
		public bool disabled {get;set;}

		public string to_string () {
			return @"$name:$value";
		}
	}
	public class ResourceAuth : Object {
		public bool disabled {get;set;}
		public string username {get;set;}
		public string password {get;set;}
		public string token {get;set;}
	}

	public void import (string source) throws Benchwell.ImportError {
		//var collection = new Benchwell.HttpCollection ();

		Json.Parser parser = new Json.Parser ();
		parser.load_from_data (source, source.length);

		Json.Object root_object = parser.get_root ().get_object ();

		assert (root_object.get_int_member ( "__export_format") == 4);

		var resource_array = root_object.get_array_member ("resources");

		Resource[] resources = {};
		resource_array.foreach_element ((array, index, node) => {
			var resource = Json.gobject_deserialize (typeof (Resource), node) as Resource;
			assert (resource != null);

			resource.resource_type = node.get_object ().get_string_member ("_type");
			resource.id = node.get_object ().get_string_member ("_id");
			switch (resource.resource_type) {
				case "request_group":
					print (@" FOLDER: $(resource.name)\n");
					break;
				case "workspace":
					print (@" COLLECTION $(resource.name)\n");
					break;
				case "request":
					var body_node = node.get_object ().get_member ("body");
					resource.request_body = Json.gobject_deserialize (typeof (ResourceBody), body_node) as ResourceBody;
					assert (resource.request_body != null);
					if (resource.request_body.mimeType == null) {
						resource.request_body.mimeType = "none";
					}

					ResourceKV[] headers = {};
					node.get_object ().get_array_member ("headers").foreach_element ((harray, index, hnode) => {
						var header = Json.gobject_deserialize (typeof (ResourceKV), hnode) as ResourceKV;
						assert(header != null);
						headers += header;
					});
					ResourceKV[] parameters = {};
					node.get_object ().get_array_member ("parameters").foreach_element ((harray, index, hnode) => {
						var parameter = Json.gobject_deserialize (typeof (ResourceKV), hnode) as ResourceKV;
						assert(parameter != null);
						parameters += parameter;
					});

					var authtype = node.get_object ().get_object_member ("authentication").get_string_member_with_default ("type", "");
					switch (authtype) {
						case "basic":
							var auth = Json.gobject_deserialize (typeof (ResourceAuth), node.get_object ().get_member ("authentication")) as ResourceAuth;
							var header = new ResourceKV ();
							header.name = "Authorization";
							header.value = @"Basic {{ base64 '$(auth.username):$(auth.password)'}}";
							header.disabled = auth.disabled;
							headers += header;
							break;
						case "bearer":
							var auth = Json.gobject_deserialize (typeof (ResourceAuth), node.get_object ().get_member ("authentication")) as ResourceAuth;
							var header = new ResourceKV ();
							header.name = "Authorization";
							header.value = @"Bearer $(auth.token)";
							header.disabled = auth.disabled;
							headers += header;
							break;
					}

					resource.request_headers = headers;
					resource.request_params = parameters;

					print (@" REQUEST $(resource.method) $(resource.url)\n");
					print (@"         mime: $(resource.request_body.mimeType)\n");
					foreach (var kv in resource.request_headers) {
						print (@"       header: $(kv)\n");
					}
					foreach (var kv in resource.request_params) {
						print (@"        param: $(kv)\n");
					}

					break;
			}

			resources += resource;
		});

		Benchwell.HttpCollection[] collections = {};
		Benchwell.HttpItem[] dirs = {};
		Benchwell.HttpItem[] items = {};

		var collection_map = new HashTable<string,int64> (str_hash, str_equal);
		var folder_map = new HashTable<string,int64> (str_hash, str_equal);
		var item_map = new HashTable<string,int64> (str_hash, str_equal);

		foreach (var workspace in resources) {
			if (workspace.resource_type != "workspace") {
				continue;
			}

			var collection = new Benchwell.HttpCollection ();
			collection.touch_without_save (() => {
				collection.name = workspace.name;
			});
			collection.save ();
			collection_map.set (workspace.id, collection.id);

			collections += collection;
		}

		foreach (var request_group in resources) {
			if (request_group.resource_type != "request_group") {
				continue;
			}

			var item = new Benchwell.HttpItem ();
			item.touch_without_save (() => {
				item.http_collection_id = 0;
				item.parent_id = 0;
				item.is_folder = true;
				item.name = request_group.name;
			});

			item.save ();
			folder_map.set (request_group.id, item.id);
		}

		foreach (var request in resources) {
			if (request.resource_type != "request") {
				continue;
			}

			var item = new Benchwell.HttpItem ();
			item.touch_without_save (() => {
				item.is_folder = false;
				item.name = request.name;
				item.parent_id = 0;
				item.method = request.method;
				item.url = request.url;
				item.body = request.request_body.text;
				item.mime = request.request_body.mimeType;
				item.description = request.description;
			});
			item.save_all ();

			foreach (var header in request.request_headers)
				item.add_header (header.name, header.value);
			foreach (var param in request.request_params)
				item.add_param (param.name, param.value);

			item_map.set (request.id, item.id);
		}

		foreach (var resource in resources) {
			switch (resource.resource_type) {
				case "request_group":
					var folder_id = folder_map.get (resource.id);
					Benchwell.HttpItem folder;
					foreach (var item in dirs) {
						if (item.id == folder_id) {
							folder = item;
							break;
						}
					}


					if (resource.parentId.has_prefix ("wrk_")) {
						folder.http_collection_id = collection_map.get (resource.parentId);
						update_tree (folder, resource.id, folder_map, item_map, ref resources);
					}

					break;
				case "request":
					break;
			}
		}

		//return collection;
	}

	private void update_tree (Benchwell.HttpItem item,
							  string insomnia_id,
							  HashTable<string,int64?> folder_map,
							  HashTable<string,int64?> item_map,
							  ref Resource[] resources) {
		foreach (var resource in resources) {
			if (resource.parentId == insomnia_id) {

			}
		}
	}

	public static int main(string[] args) {
		string text;
		var ok = GLib.FileUtils.get_contents ("../Insomnia_2021-01-30.json", out text, null);
		if (!ok) {
			return 1;
		}

		var t = new Benchwell.InsomniaImpoter ();
		t.import (text);

		return 0;
	}
}
