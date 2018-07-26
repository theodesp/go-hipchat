package hipchat

const (
	defaultBaseUrl             = "https://api.hipchat.com/"
	userAgent                  = "go-hipchat"
	contentTypeApplicationJson = "application/json; charset=UTF-8"
	contentDispositionMetadata = `attachment; name="metadata"`
	contentDispositionFile     = `attachment; name="file"; filename="%v"`
	apiVersion2                = "v2"
)
