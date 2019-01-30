package rest_backend

import (
	gsql "database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/makeroo/taxi_scout/sql"
	"github.com/makeroo/taxi_scout/storage"
	"github.com/makeroo/taxi_scout/ts_errors"
	"net/http"
	"strconv"
)

type InvitationToken struct {
	Token string `json:"invitation"`
}

func (server *RestServer) checkUserIdCookie (r *http.Request) (int32, error) {
	cookie, err := r.Cookie("_ts_u")

	if err == http.ErrNoCookie {
		return 0, ts_errors.NotAuthorized
	}

	if err != nil {
		server.Logger.Debugw("failed at reading cookies",
			"err", err)
		return 0, ts_errors.BadRequest
	}

	var userId int32

	err = server.Configuration.SecureCookies.Decode("_ts_u", cookie.Value, &userId)

	if err != nil {
		server.Logger.Debugw("cookie decoding failed",
			"err", err)

		err = ts_errors.BadRequest
	}

	return userId, err
}

func (server *RestServer) checkUserCookie (r *http.Request) (*storage.Account, error) {
	userId, err := server.checkUserIdCookie(r)

	if err != nil {
		return nil, err
	}

	return server.Dao.QueryAccount(userId)
}


func (server *RestServer) setUserCookie (accountId int32, w http.ResponseWriter) {
	if encoded, err := server.Configuration.SecureCookies.Encode("_ts_u", accountId); err == nil {
		cookie := &http.Cookie{
			Name:  "_ts_u",
			Value: encoded,
			Path:  "/",
			Secure: true,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
	} else {
		server.Logger.Errorf("failed encoding user cookie: error=%v", err)
	}
}

func (server *RestServer) writeResponse (statusCode int, payload interface{}, w http.ResponseWriter) {
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

func (server *RestServer) Accounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		group, found := r.URL.Query()["group"]

		if !found || len(group) != 1 {
			w.WriteHeader(400)
			return
		}

		groupId64, err := strconv.ParseInt(group[0], 10, 32)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		groupId := int32(groupId64)

		userId, err := server.checkUserIdCookie(r)

		if err == http.ErrNoCookie {
			server.writeResponse(401, ts_errors.NotAuthorized, w)
			return
		} else if err != nil {
			// TODO: in alternativa forbidden: mi ha dato un cookie ma non buono?
			// putroppo securecookie non pubblica "expired"
			w.WriteHeader(400)
			return
		}

		err = server.Dao.CheckPermission(userId, groupId, sql.PermissionMember)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)

			return

		case nil:

		default:
			server.Logger.Errorf("unexpected error: error=%v", err)
			server.writeResponse(500, ts_errors.ServerError, w)

			return
		}

		/*		if token, found := r.URL.Query()["invitation"]; found {
					if len(token) > 1 {
						server.writeResponse(400, "TODO errore troppi token", w)
						return
					}

					account, err := server.checkUserCookie(r)

					if err == nil {
						server.writeResponse(200, []*storage.Account{account}, w)
						return
					} else if err != http.ErrNoCookie {
						server.writeResponse(400, "TODO illegal cookie", w)
						return
					}

					account, err = server.Dao.QueryInvitationToken(token[0])

					switch {
					case err == sql.ErrNoRows:
						server.writeResponse(404, ts_errors.NotFound, w)

					case err != nil:
						server.Logger.Errorf("unexpected error: error=%v", err)
						server.writeResponse(500, ts_errors.ServerError, w)

					default:
						server.writeResponse(200, []*storage.Account{account}, w)
					}

					return
				}
		*/
		accounts, err := server.Dao.QueryAccounts(groupId)

		if err != nil {
			server.Logger.Errorf("unexpected error: error=%v", err)
			server.writeResponse(500, ts_errors.ServerError, w)
		} else {
			server.writeResponse(200, accounts, w)
		}

	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		invitationToken := new(InvitationToken)

		err := decoder.Decode(invitationToken)

		if err != nil {
			server.Logger.Errorw("InvitationToken decoding failed",
				"err", err,
				)
			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		account, found, err := server.Dao.QueryInvitationToken(invitationToken.Token)

		if err != nil {
			// ok, something went wrong
			// but if the user was authenticated then we return an account anyway
			// why? because we want that reloading the invitation after having used
			// it, thus being authenticated, bring the user to the home page.
			// this logic could be transferred to the client but imho it is clearer
			// to have it on the server although this huge comment seems to prove the contrary
			userId, cookieErr := server.checkUserIdCookie(r)

			if cookieErr == nil {
				server.Logger.Debugw("ignoring invitation error, using authentication cookie",
					"err", err)

				server.writeResponse(200,
					map[string]interface{}{
						"id": userId,
						"authenticated": false,
						"new_account": false,
					},
					w)

				return
			} else {
				server.Logger.Debugw("cookie decoding failed",
					"err", cookieErr)
			}

			switch t := err.(type) {
			case ts_errors.RestError:
				server.writeResponse(t.Code, t, w)
				return

			default:
				if t == gsql.ErrNoRows {
					server.writeResponse(404, ts_errors.NotFound.Error(), w)
				} else {
					server.Logger.Errorw("account creation failed",
						"err", err,
					)

					server.writeResponse(500, ts_errors.ServerError, w)
				}

				return
			}
		}

		server.setUserCookie(account.Id, w)

		server.writeResponse(200,
			map[string]interface{}{
				"id": account.Id,
				"authenticated": true,
				"new_account": !found,
			},
			w)
/*		account := storage.NewAccountWithCredentials()
		err := decoder.Decode(account)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		newId, err := server.Dao.InsertAccount(account)

		if err != nil {
			server.Logger.Errorw("account creation failed",
				"err", err,
				"name", account.Name,
				"email", account.Email,
			)
			w.WriteHeader(400)
			return
		}

		bytes, err := json.Marshal(map[string]int32{"id": newId})
		w.Write(bytes)
*/
	default:
		w.WriteHeader(405)
	}
}

func (server *RestServer) Account(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			server.Logger.Errorw("illegal id parameter",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		id32 := int32(id)
		account, err := server.checkUserCookie(r)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			if err == gsql.ErrNoRows {
				server.writeResponse(401, ts_errors.NotAuthorized, w)

			} else {
				server.Logger.Errorw("account query failed",
					"err", err,
				)

				server.writeResponse(500, ts_errors.ServerError, w)
			}

			return
		}

		if id32 != account.Id {
			server.writeResponse(403, ts_errors.Forbidden, w)
			return
		}

		server.writeResponse(200, account, w)

	case http.MethodPost:
	default:
		w.WriteHeader(405)
	}
}

type AccountAuthenticatePayload struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

func (server *RestServer) AccountsAuthenticate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		credentials := AccountAuthenticatePayload{"", ""}
		err := decoder.Decode(&credentials)

		if err != nil {
			server.Logger.Errorw("payload unmarshalling failed",
				"err", err,
			)
			server.writeResponse(400, ts_errors.BadRequest, w)

			return
		}

		accountId, err := server.Dao.AuthenticateAccount(credentials.Email, credentials.Pwd)

		if err != nil {
			server.Logger.Errorw("authentication failed",
				"err", err,
				"email", credentials.Email,
			)
			server.writeResponse(401, ts_errors.NotAuthorized, w)
			return
		}

		if encoded, err := server.Configuration.SecureCookies.Encode("_ts_u", accountId); err == nil {
			cookie := &http.Cookie{
				Name:  "_ts_u",
				Value: encoded,
				Path:  "/",
				Secure: true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}

		server.writeResponse(200, map[string]int32{"id": accountId}, w)

	default:
		w.WriteHeader(405)
	}
}

type AccountPasswordPayload struct {
	Id     int32  `json:"id"`
	OldPwd string `json:"old_pwd"`
	NewPwd string `json:"new_pwd"`
}

func (server *RestServer) AccountPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)

		credentials := AccountPasswordPayload{-1, "", ""}
		err := decoder.Decode(credentials)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		err = server.Dao.UpdateAccountPassword(credentials.Id, credentials.OldPwd, credentials.NewPwd)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		server.writeResponse(200, map[string]int32{"id": credentials.Id}, w)

	default:
		w.WriteHeader(405)
	}
}
