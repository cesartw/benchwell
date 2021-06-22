namespace Benchwell {
	namespace Http {
		enum CODES {
			OK,
			Created,
			Accepted,
			NoContent,
			BadRequest,
			Unauthorized,
			PaymentRequired,
			Forbidden,
			NotFound,
			MethodNotAllowed,
			RequestTimeout,
			Conflict,
			TooManyRequests,
			InternalServerError,
			NotImplemented,
			BadGateway,
			ServiceUnavailable,
			GatewayTimeout;

			public static Benchwell.Http.CODES? parse (uint code) {
				switch (code) {
					case 200:
						return OK;
					case 201:
						return Created;
					case 202:
						return Accepted;
					case 204:
						return NoContent;
					case 400:
						return BadRequest;
					case 401:
						return Unauthorized;
					case 402:
						return PaymentRequired;
					case 403:
						return Forbidden;
					case 404:
						return NotFound;
					case 405:
						return MethodNotAllowed;
					case 408:
						return RequestTimeout;
					case 409:
						return Conflict;
					case 429:
						return TooManyRequests;
					case 500:
						 return InternalServerError;
					case 501:
						 return NotImplemented;
					case 502:
						 return BadGateway;
					case 503:
						 return ServiceUnavailable;
					case 504:
						 return GatewayTimeout;
					default:
						 return null;
				}
			}

			public string to_string () {
				switch (this) {
					case OK:
						return "OK";
					case Created:
						return "Created";
					case Accepted:
						return "Accepted";
					case NoContent:
						return "No Content";
					case BadRequest:
						return "Bad Request";
					case Unauthorized:
						return "Unauthorized";
					case PaymentRequired:
						return "Payment Required";
					case Forbidden:
						return "Forbidden";
					case NotFound:
						return "NotFound";
					case MethodNotAllowed:
						return "Method Not Allowed";
					case RequestTimeout:
						return "Request Timeout";
					case Conflict:
						return "Conflict";
					case TooManyRequests:
						return "Too Many Requests";
					case InternalServerError:
						 return "Internal Server Error";
					case NotImplemented:
						 return "Not Implemented";
					case BadGateway:
						 return "Bad Gateway";
					case ServiceUnavailable:
						 return "Service Unavailable";
					case GatewayTimeout:
						 return "Gateway Timeout";
					default:
						 return "";
				}
			}
		}

		enum Columns {
			ICON,
			TEXT,
			METHOD,
			ITEM
		}

		enum Method {
			GET,
			POST,
			PUT,
			PATCH,
			DELETE,
			OPTIONS,
			HEAD;

			public string to_string () {
				switch (this) {
					case GET:
						return "GET";
					case POST:
						return  "POST";
					case PUT:
						return "PUT";
					case PATCH:
						return "PATCH";
					case DELETE:
						return "DELETE";
					case OPTIONS:
						return "OPTIONS";
					case HEAD:
						return "HEAD";
					default:
						assert_not_reached ();
				}
			}

			public static Method[] all () {
				return {
					GET,
					POST,
					PUT,
					PATCH,
					DELETE,
					OPTIONS,
					HEAD
				};
			}
		}
	}
}
