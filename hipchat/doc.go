/*
go-hipchat is a Go client library for accessing the https://developer.atlassian.com/server/hipchat/about-the-hipchat-rest-api/.

go-hipchat requires Go version 1.8 or greater.

Usage:


Construct a new HipChat client, then use the various services on the client to
access different parts of the HipChat API. For example:


Some API methods have optional parameters that can be passed. For example:

NOTE: Using the https://godoc.org/context package, one can easily
pass cancelation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

For more sample code snippets, head over to the
https://github.com/theodesp/go-hipchat/tree/master/examples directory.


Authentication:

The go-hipchat library does not directly handle authentication. Instead, when
creating a new client, pass an `http.Client` that can handle authentication for
you. The easiest and recommended way to do this is using the golang.org/x/oauth2
library, but you can always use any other library that provides an
`http.Client`. If you have an OAuth2 access token (for example, a
https://developer.atlassian.com/server/hipchat/hipchat-rest-api-access-tokens/, you can use it with the oauth2 library using:


Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

See the oauth2 docs for complete instructions on using that library.


Rate Limiting:

Hipchat imposes a rate limit on all API clients. 500 API requests per 5 minutes.
Once you exceed the limit, calls will return HTTP status 429.

Learn more about HipChat rate limiting at
https://developer.atlassian.com/server/hipchat/hipchat-rest-api-rate-limits.

Response Codes:

https://developer.atlassian.com/server/hipchat/hipchat-rest-api-response-codes

Pagination:

*/
package hipchat
