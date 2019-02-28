package rest_backend

import (
	gsql "database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/makeroo/taxi_scout/storage"
	"github.com/makeroo/taxi_scout/ts_errors"
	"net/http"
	"strconv"
)


type InvitationToken struct {
	Token string `json:"invitation"`
}

func (server *RestServer) Accounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		group, found := r.URL.Query()["group"]

		if !found || len(group) != 1 {
			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		groupId64, err := strconv.ParseInt(group[0], 10, 32)

		if err != nil {
			server.Logger.Debugw("can't parse group query parameter",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		groupId := int32(groupId64)

		userId, err := server.checkUserIdCookie(r)

		if err == http.ErrNoCookie {
			server.writeResponse(401, ts_errors.NotAuthorized, w)
			return
		} else if err != nil {
			server.Logger.Debugw("can't decode user cookie",
				"err", err)

			// TODO: unsure if it is corrrect 400 or 403
			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		accounts, err := server.Dao.QueryAccounts(groupId, userId)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)

			return

		case nil:
			server.writeResponse(200, accounts, w)

		default:
			server.Logger.Errorw("fetch accounts failed",
				"err", err)
			server.writeResponse(500, ts_errors.ServerError, w)
		}

	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		invitationToken := InvitationToken{}

		err := decoder.Decode(invitationToken)

		if err != nil {
			server.Logger.Debugw("InvitationToken decoding failed",
				"err", err)

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
				server.Logger.Debugw("both invitation processing and cookie decoding failed",
					"cookieErr", cookieErr)
			}

			switch t := err.(type) {
			case ts_errors.RestError:
				server.writeResponse(t.Code, t, w)
				return

			default:
				if t == gsql.ErrNoRows {
					server.writeResponse(404, ts_errors.NotFound, w)
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

	default:
		w.WriteHeader(405)
	}
}


func (server *RestServer) Account(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		vars := mux.Vars(r)

		var account *storage.Account
		var err error
		var id32 int32

		if vars["id"] == "me" {
			account, err = server.checkUserCookie(r)
			if err == nil {
				id32 = account.Id
			}
		} else {
			id, err := strconv.ParseInt(vars["id"], 10, 32)

			if err != nil {
				server.Logger.Errorw("illegal id parameter",
					"err", err)

				server.writeResponse(400, ts_errors.BadRequest, w)
				return
			}

			id32 = int32(id)
			account, err = server.checkUserCookie(r)
		}

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
		myId, err := server.checkUserIdCookie(r)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		decoder := json.NewDecoder(r.Body)

		account := storage.Account{}
		err = decoder.Decode(&account)

		if err != nil {
			server.Logger.Debugw("update account: illegal payload",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		if myId != account.Id {
			server.writeResponse(403, ts_errors.Forbidden, w)
			return
		}

		err = server.Dao.AccountUpdate(&account)

		if err != nil {
			server.Logger.Errorw("update account: update failed",
				"err", err)

			server.writeResponse(500, ts_errors.ServerError, w)
			return
		}

		// time.Sleep(20 * time.Second)
		server.writeResponse(200, nil, w)

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

		credentials := AccountAuthenticatePayload{}
		err := decoder.Decode(&credentials)

		if err != nil {
			server.Logger.Debugw("payload unmarshalling failed",
				"err", err,
			)
			server.writeResponse(400, ts_errors.BadRequest, w)

			return
		}

		accountId, err := server.Dao.AuthenticateAccount(credentials.Email, credentials.Pwd)

		if err != nil {
			server.Logger.Debugw("authentication failed",
				"err", err,
				"email", credentials.Email,
			)
			server.writeResponse(401, ts_errors.NotAuthorized, w)
			return
		}

		server.Logger.Infow("user authenticated",
			"email", credentials.Email)

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
		w.Header().Add("Content-Type", "application/json")

		userId, err := server.checkUserIdCookie(r)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		decoder := json.NewDecoder(r.Body)

		credentials := AccountPasswordPayload{}
		err = decoder.Decode(credentials)

		if err != nil {
			server.Logger.Debugw("json decoding failed",
				"err", err)

			w.WriteHeader(400)
			return
		}

		if userId != credentials.Id {
			server.writeResponse(403, ts_errors.Forbidden, w)
			return
		}

		err = server.Dao.UpdateAccountPassword(credentials.Id, credentials.OldPwd, credentials.NewPwd)

		if err != nil {
			server.Logger.Errorw("password update failed",
				"err", err)

			server.writeResponse(500, ts_errors.ServerError, w)
			return
		}

		server.writeResponse(200, map[string]int32{"id": credentials.Id}, w)

	default:
		w.WriteHeader(405)
	}
}

func (server *RestServer) AccountGroups(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		vars := mux.Vars(r)

		var err error
		var id32 int32

		id, err := strconv.ParseInt(vars["id"], 10, 32)

		if err != nil {
			server.Logger.Debugw("illegal id parameter",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		id32 = int32(id)

		uid, err := server.checkUserIdCookie(r)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		if id32 != uid {
			server.writeResponse(403, ts_errors.Forbidden, w)
			return
		}

		groups, err := server.Dao.AccountGroups(id32)

		if err == gsql.ErrNoRows {
			groups = []*storage.ScoutGroup{}
		} else if err != nil {
			server.Logger.Errorw("storage error:",
				"err", err)

			server.writeResponse(500, ts_errors.ServerError, w)
			return
		}

		server.writeResponse(200, groups, w)

	default:
		w.WriteHeader(405)
	}
}

func (server *RestServer) AccountScouts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		vars := mux.Vars(r)

		var err error
		var id32 int32

		id, err := strconv.ParseInt(vars["id"], 10, 32)

		if err != nil {
			server.Logger.Debugw("illegal id parameter",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		id32 = int32(id)

		uid, err := server.checkUserIdCookie(r)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		if id32 != uid {
			server.writeResponse(403, ts_errors.Forbidden, w)
			return
		}

		scouts, err := server.Dao.AccountScouts(id32)

		if err == gsql.ErrNoRows {
			scouts = []*storage.Scout{}
		} else if err != nil {
			server.Logger.Errorw("storage error:",
				"err", err)

			server.writeResponse(500, ts_errors.ServerError, w)
			return
		}

		server.writeResponse(200, scouts, w)

	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		vars := mux.Vars(r)

		var err error
		var id32 int32

		id, err := strconv.ParseInt(vars["id"], 10, 32)

		if err != nil {
			server.Logger.Debugw("illegal id parameter",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		id32 = int32(id)

		uid, err := server.checkUserIdCookie(r)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		if id32 != uid {
			server.writeResponse(403, ts_errors.Forbidden, w)
			return
		}

		decoder := json.NewDecoder(r.Body)

		scout := storage.Scout{}
		err = decoder.Decode(&scout)

		if err != nil {
			server.Logger.Debugw("update scout: illegal payload",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		scout.Id = -1

		scoutId, err := server.Dao.InsertOrUpdateScout(scout, uid)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)

		case nil:
			server.writeResponse(200, map[string]int32{
				"id": scoutId,
			}, w)

		default:
			server.Logger.Errorw("scout insert failed",
				"err", err)

			server.writeResponse(500, ts_errors.ServerError, w)
		}

	default:
		w.WriteHeader(405)
	}
}
