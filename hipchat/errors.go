package hipchat

import "errors"

var invalidSetApiVersion = errors.New("set_api_version: apiVersion string parameter is prefixed with a forward slash (/)")
