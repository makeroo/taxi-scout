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

var BadRequest = RestError{errors.New("bad_request"), 400}
var NotAuthorized = RestError{errors.New("not_authorized"), 401}
var Forbidden = RestError{errors.New("forbidden"), 403}
var NotFound = RestError{errors.New("not_found"), 404}
var Expired = RestError{errors.New("expired"), 410}

// Stoken token error occurs when processing an invitation token while a user
// has been already authenticated and the receiving invitation email does not
// match authenticated user's email.
var StolenToken = RestError{errors.New("stolen_token"), 403}

// An unexpected and unrecoverable error. It could be either a system error,
// eg. database unreachable, or a bug.
// TODO: in the future add indication whenever possible if it is worth retrying later or not.
var ServerError = RestError{errors.New("server_error"), 500}
