# go-hipchat - WIP #

<a href="https://godoc.org/github.com/theodesp/go-hipchat">
<img src="https://godoc.org/github.com/theodesp/go-hipchat/hipchat?status.svg" alt="GoDoc">
</a>

<a href="https://opensource.org/licenses/BSD-3-Clause" rel="nofollow">
<img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="License"/>
</a>

<a href="https://travis-ci.org/theodesp/go-hipchat" rel="nofollow">
<img src="https://travis-ci.org/theodesp/go-hipchat.svg?branch=master" />
</a>

<a href="https://codecov.io/gh/theodesp/go-hipchat">
  <img src="https://codecov.io/gh/theodesp/go-hipchat/branch/master/graph/badge.svg" />
</a>

<a href="https://goreportcard.com/report/github.com/theodesp/go-hipchat">
  <img src="https://goreportcard.com/badge/github.com/theodesp/go-hipchat" />
</a>


go-hipchat is a Go client library for accessing the [Hipchat API](https://developer.atlassian.com/server/hipchat/about-the-hipchat-rest-api/).

go-hipchat requires Go version 1.8 or greater.


This is ***WIP!***


## Usage ##

```go
import "github.com/theodesp/go-hipchat/hipchat"
```

Construct a new HipChat client, then use the various services on the client to
access different parts of the HipChat API. For example:

```go

```

Some API methods have optional parameters that can be passed. For example:

```go

```

NOTE: Using the [context](https://godoc.org/context) package, one can easily
pass cancelation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

For more sample code snippets, head over to the
[examples](https://github.com/theodesp/go-hipchat/tree/master/examples) directory.

### Authentication ###

The go-hipchat library does not directly handle authentication. Instead, when
creating a new client, pass an `http.Client` that can handle authentication for
you. The easiest and recommended way to do this is using the [oauth2][]
library, but you can always use any other library that provides an
`http.Client`. If you have an OAuth2 access token (for example, a [personal
API token][]), you can use it with the oauth2 library using:

```go

```

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

See the [oauth2 docs][] for complete instructions on using that library.


### Rate Limiting ###

Hipchat imposes a rate limit on all API clients. 500 API requests per 5 minutes. 
Once you exceed the limit, calls will return HTTP status 429.

Learn more about HipChat rate limiting at
https://developer.atlassian.com/server/hipchat/hipchat-rest-api-rate-limits.

### Response Codes ###

https://developer.atlassian.com/server/hipchat/hipchat-rest-api-response-codes


For complete usage of go-hipchat, see the full [package docs][https://godoc.org/github.com/theodesp/go-hipchat/hipchat].

[HipChat API]: https://developer.atlassian.com/server/hipchat/about-the-hipchat-rest-api/
[oauth2]: https://github.com/golang/oauth2
[oauth2 docs]: https://godoc.org/golang.org/x/oauth2

## Roadmap ##

[Contributing](./CONTRIBUTING)

## Versioning ##

In general, go-hipchat follows [semver](https://semver.org/) as closely as we
can for tagging releases of the package. For self-contained libraries, the
application of semantic versioning is relatively straightforward and generally
understood. But because go-hipchat is a client library for the HipChat API, which
itself changes behavior, and because we are typically pretty aggressive about
implementing preview features of the HipChat API, we've adopted the following
versioning policy:

* We increment the **major version** with any incompatible change to
	non-preview functionality, including changes to the exported Go API surface
	or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to
	functionality, as well as any changes to preview functionality in the HipChat
	API.
* We increment the **patch version** with any backwards-compatible bug fixes.

Preview functionality may take the form of entire methods or simply additional
data returned from an otherwise non-preview method. Refer to the HipChat API
documentation for details on preview functionality.

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.