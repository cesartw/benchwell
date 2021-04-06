public class Benchwell.ImporterInsomnia : Benchwell.Importer, Object {
	public class Resource : Object {
		public string id {get;set;}
		public string resource_type {get;set;}
		public string parentId {get;set;}
		public string url {get;set;}
		public string name {get;set;}
		public string description {get;set;}
		public string method {get;set;}
		public string workspace_id {get;set;}
		public ResourceBody request_body {get;set;}
		public ResourceKV[] request_params {get;set;}
		public ResourceKV[] request_headers {get;set;}
		public Resource[] resources {get;set;}
		public ResourceKV[] variables {get;set;}
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
					break;
				case "workspace":
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

					break;
				case "environment":
					var data = node.get_object ().get_object_member ("data");
					ResourceKV[] vars = {};
					data.get_members ().foreach ((key) => {
						var val = data.get_string_member (key);
						vars += new ResourceKV () { name = key, value = val };
					});
					resource.variables = vars;
					break;
			}

			resources += resource;
		});

		var collection_map = new HashTable<string, Benchwell.HttpCollection> (str_hash, str_equal);
		var folder_map = new HashTable<string, Benchwell.HttpItem> (str_hash, str_equal);

		foreach (var environment in resources) {
			if (environment.resource_type != "environment") {
				continue;
			}

			var env = new Benchwell.Environment ();
			env.name = environment.name;
			Config.add_environment (env);
			foreach (var v in environment.variables) {
				env.add_variable (v.name, v.value);
			}

		}

		foreach (var workspace in resources) {
			if (workspace.resource_type != "workspace") {
				continue;
			}

			var collection = new Benchwell.HttpCollection ();
			collection.touch_without_save (() => {
				collection.name = workspace.name;
			});
			Config.add_http_collection (collection);
			collection_map.set (workspace.id, collection);


			// CREATE FOLDER TREE
			var created = true;
			while (created) {
				created = false;
				foreach (var request_group in resources) {
					if (request_group.resource_type != "request_group")  {
						continue;
					}

					if (folder_map.get (request_group.id) != null) {
						continue;
					}

					Benchwell.HttpItem? parent_item = null;
					if (request_group.parentId != workspace.id) {
						parent_item = folder_map.get (request_group.parentId);
						if (parent_item == null || parent_item.http_collection_id != collection.id) {
							continue;
						}
					}

					var item = new Benchwell.HttpItem ();
					item.touch_without_save (() => {
						item.http_collection_id = collection.id;
						item.parent_id = 0;
						if (parent_item != null)
							item.parent_id = parent_item.id;
						item.is_folder = true;
						item.name = request_group.name;
					});

					item.save ();
					created = true;
					folder_map.set (request_group.id, item);
				}
			}

			foreach (var request in resources) {
				if (request.resource_type != "request")  {
					continue;
				}

				var parent_folder = folder_map.get (request.parentId);
				if (request.id != workspace.id && (parent_folder == null || parent_folder.http_collection_id != collection.id)) {
					continue;
				}

				var item = new Benchwell.HttpItem ();
				item.touch_without_save (() => {
					item.http_collection_id = collection.id;
					item.parent_id = 0;
					if (parent_folder != null)
						item.parent_id = parent_folder.id;
					item.is_folder = false;
					item.name = request.name;
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

				created = true;
			}
		}
	}

	public Gtk.FileFilter get_file_filter () {
		var filter = new Gtk.FileFilter ();
		filter.set_filter_name ("*.json");
		filter.add_mime_type ("application/json");

		return filter;
	}

	//public static int main(string[] args) {
		//string text;
		//var ok = GLib.FileUtils.get_contents ("../Insomnia_2021-04-05.json", out text, null);
		//if (!ok) {
			//return 1;
		//}

		//try {
			//Benchwell.Config = new Benchwell._Config ();
		//} catch (Benchwell.ConfigError err) {
			//stderr.printf (err.message);
			//return 1;
		//}

		//var t = new Benchwell.InsomniaImporter ();
		//t.import (text);

		//return 0;
	//}
}
