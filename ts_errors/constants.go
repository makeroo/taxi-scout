package ts_errors

import (
	"encoding/json"
	"errors"
)

type RestError struct {
	error
	Code int
}

func (error *RestError) Json () ([]byte, error) {
	return json.Marshal(map[string]string{
		"error": error.Error(),
	})
}

var BadRequest = RestError{errors.New("forbidden"), 400}
var NotAuthorized = RestError{errors.New("not_authorized"), 401}
var Forbidden = RestError{errors.New("forbidden"), 403}
var NotFound = RestError{errors.New("not_found"), 404}
var Expired = RestError{errors.New("expired"), 410}

var ServerError = RestError{errors.New("server_error"), 500}
