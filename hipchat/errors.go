package hipchat

import "errors"

var invalidSetApiVersion = errors.New("set_api_version: apiVersion string parameter is prefixed with a forward slash (/)")
var emptyParam = errors.New("empty_param: required parameter is empty")
var invalidFileUpload = errors.New("file_upload: the file to upload can't be a directory")
