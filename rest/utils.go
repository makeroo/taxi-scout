package rest

import (
	"encoding/json"
	"github.com/makeroo/taxi_scout/ts_errors"
	"github.com/makeroo/taxi_scout/storage"
	"net/http"
)

func (server *Server) checkUserIDCookie (r *http.Request) (int32, error) {
	cookie, err := r.Cookie("_ts_u")

	if err == http.ErrNoCookie {
		return 0, ts_errors.NotAuthorized
	}

	if err != nil {
		server.Logger.Debugw("failed at reading cookies",
			"err", err)
		return 0, ts_errors.BadRequest
	}

	var userID int32

	err = server.Configuration.SecureCookies.Decode("_ts_u", cookie.Value, &userID)

	if err != nil {
		server.Logger.Debugw("cookie decoding failed",
			"err", err)

		err = ts_errors.BadRequest
	}

	return userID, err
}

func (server *Server) checkUserCookie (r *http.Request) (*storage.Account, error) {
	userID, err := server.checkUserIDCookie(r)

	if err != nil {
		return nil, err
	}

	return server.Dao.QueryAccount(userID)
}


func (server *Server) setUserCookie (accountID int32, w http.ResponseWriter) {
	if encoded, err := server.Configuration.SecureCookies.Encode("_ts_u", accountID); err == nil {
		cookie := &http.Cookie{
			Name:  "_ts_u",
			Value: encoded,
			Path:  "/",
			Secure: server.Configuration.HTTPSCookies,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
	} else {
		server.Logger.Errorf("failed encoding user cookie: error=%v", err)
	}
}

func (server *Server) writeResponse (statusCode int, payload interface{}, w http.ResponseWriter) {
	if val, ok := payload.(error); ok {
		payload = map[string]string{
			"error": val.Error(),
		}
	}

	res, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(res)

	if err != nil {
		server.Logger.Warnf("error while writing response: %v", err)
	}
}
