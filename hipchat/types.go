package hipchat

// A struct representing the API Version
type ApiVersionOptions struct {
	V2 string
}

// A struct representing an HTTP Method type
type HttpMethodOptions struct {
	POST    string
	GET     string
	PUT     string
	DELETE  string
	OPTIONS string
}

// A struct representing a Content Type header
type ContentTypesOptions struct {
	ApplicationJson string
}

var (
	ApiVersions = ApiVersionOptions{
		V2: "v2",
	}

	HttpMethods = HttpMethodOptions{
		POST:    "POST",
		GET:     "GET",
		PUT:     "PUT",
		DELETE:  "DELETE",
		OPTIONS: "OPTIONS",
	}

	ContentTypes = ContentTypesOptions{
		ApplicationJson: "application/json",
	}
)
