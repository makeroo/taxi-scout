package rest

import (
	gsql "database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	tserrors "github.com/makeroo/taxi_scout/errors"
	"github.com/makeroo/taxi_scout/storage"
)

// AccountsRequest models /accounts POST request payload.
type AccountsRequest struct {
	Token string `json:"invitation"`
}

// Accounts method implement /accounts REST requests.
func (server *Server) Accounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		group, found := r.URL.Query()["group"]

		if !found || len(group) != 1 {
			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		groupID64, err := strconv.ParseInt(group[0], 10, 32)

		if err != nil {
			server.Logger.Debugw("can't parse group query parameter",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		groupID := int32(groupID64)

		userID, err := server.checkUserIDCookie(r)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:
			break

		default:
			server.Logger.Debugw("can't decode user cookie",
				"err", err)

			// TODO: unsure if it is corrrect 400 or 403
			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		accounts, err := server.Dao.QueryAccounts(groupID, userID)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)

			return

		case nil:
			server.writeResponse(200, accounts, w)

		default:
			server.Logger.Errorw("fetch accounts failed",
				"err", err)
			server.writeResponse(500, tserrors.ServerError, w)
		}

	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		invitationToken := AccountsRequest{}

		err := decoder.Decode(&invitationToken)

		if err != nil {
			server.Logger.Debugw("AccountsRequest decoding failed",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		userID, cookieErr := server.checkUserIDCookie(r)

		if cookieErr != nil {
			userID = storage.NoRequestingUser
		}

		account, found, scoutGroupID, joinedGroup, err := server.Dao.QueryInvitationToken(invitationToken.Token, userID)

		if err != nil {
			if err == tserrors.StolenToken {
				server.writeResponse(tserrors.StolenToken.Code, err, w)
				return
			}

			if cookieErr == nil {
				server.Logger.Debugw("ignoring invitation error, using authentication cookie",
					"err", err)

				server.writeResponse(200,
					map[string]interface{}{
						"id":            userID,
						"authenticated": false,
						"new_account":   false,
					},
					w)

				return
			}

			server.Logger.Debugw("both invitation processing and cookie decoding failed",
				"cookieErr", cookieErr)

			switch t := err.(type) {
			case tserrors.RestError:
				server.writeResponse(t.Code, t, w)
				return

			default:
				if t == gsql.ErrNoRows {
					server.writeResponse(404, tserrors.NotFound, w)
				} else {
					server.Logger.Errorw("account creation failed",
						"err", err,
					)

					server.writeResponse(500, tserrors.ServerError, w)
				}

				return
			}
		}

		if cookieErr != nil {
			server.setUserCookie(account.ID, w)
		}

		var resp = map[string]interface{}{
			"id":            account.ID,
			"authenticated": cookieErr != nil,
			"new_account":   !found,
		}

		if joinedGroup {
			resp["scout_group"] = scoutGroupID
		}

		server.writeResponse(200, resp, w)

	default:
		w.WriteHeader(405)
	}
}

// Account method implements /account REST requests.
func (server *Server) Account(w http.ResponseWriter, r *http.Request) {
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
				id32 = account.ID
			}
		} else {
			id, err := strconv.ParseInt(vars["id"], 10, 32)

			if err != nil {
				server.Logger.Errorw("illegal id parameter",
					"err", err)

				server.writeResponse(400, tserrors.BadRequest, w)
				return
			}

			id32 = int32(id)
			account, err = server.checkUserCookie(r)
		}

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			if err == gsql.ErrNoRows {
				server.writeResponse(401, tserrors.NotAuthorized, w)

			} else {
				server.Logger.Errorw("account query failed",
					"err", err,
				)

				server.writeResponse(500, tserrors.ServerError, w)
			}

			return
		}

		if id32 != account.ID {
			server.writeResponse(403, tserrors.Forbidden, w)
			return
		}

		server.writeResponse(200, account, w)

	case http.MethodPost:
		myID, err := server.checkUserIDCookie(r)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		decoder := json.NewDecoder(r.Body)

		account := storage.Account{}
		err = decoder.Decode(&account)

		if err != nil {
			server.Logger.Debugw("update account: illegal payload",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		if myID != account.ID {
			server.writeResponse(403, tserrors.Forbidden, w)
			return
		}

		err = server.Dao.AccountUpdate(&account)

		if err != nil {
			server.Logger.Errorw("update account: update failed",
				"err", err)

			server.writeResponse(500, tserrors.ServerError, w)
			return
		}

		// time.Sleep(20 * time.Second)
		server.writeResponse(200, nil, w)

	default:
		w.WriteHeader(405)
	}
}

// AccountAuthenticateRequest models /account/authenticate REST request payload.
type AccountAuthenticateRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

// AccountsAuthenticate implements /account/authenticate REST request.
func (server *Server) AccountsAuthenticate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		credentials := AccountAuthenticateRequest{}
		err := decoder.Decode(&credentials)

		if err != nil {
			server.Logger.Debugw("payload unmarshalling failed",
				"err", err,
			)
			server.writeResponse(400, tserrors.BadRequest, w)

			return
		}

		accountID, err := server.Dao.AuthenticateAccount(credentials.Email, credentials.Pwd)

		if err != nil {
			server.Logger.Debugw("authentication failed",
				"err", err,
				"email", credentials.Email,
			)
			server.writeResponse(401, tserrors.NotAuthorized, w)
			return
		}

		server.Logger.Infow("user authenticated",
			"email", credentials.Email)

		if encoded, err := server.Configuration.SecureCookies.Encode("_ts_u", accountID); err == nil {
			cookie := &http.Cookie{
				Name:     "_ts_u",
				Value:    encoded,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}

		server.writeResponse(200, map[string]int32{"id": accountID}, w)

	default:
		w.WriteHeader(405)
	}
}

// TODO: write docs
type AccountPasswordPayload struct {
	ID     int32  `json:"id"`
	OldPwd string `json:"old_pwd"`
	NewPwd string `json:"new_pwd"`
}

// TODO: write docs
func (server *Server) AccountPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		userID, err := server.checkUserIDCookie(r)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		decoder := json.NewDecoder(r.Body)

		credentials := AccountPasswordPayload{}
		err = decoder.Decode(&credentials)

		if err != nil {
			server.Logger.Debugw("json decoding failed",
				"err", err)

			w.WriteHeader(400)
			return
		}

		if userID != credentials.ID {
			server.writeResponse(403, tserrors.Forbidden, w)
			return
		}

		err = server.Dao.UpdateAccountPassword(credentials.ID, credentials.OldPwd, credentials.NewPwd)

		if err != nil {
			server.Logger.Errorw("password update failed",
				"err", err)

			server.writeResponse(500, tserrors.ServerError, w)
			return
		}

		server.writeResponse(200, map[string]int32{"id": credentials.ID}, w)

	default:
		w.WriteHeader(405)
	}
}

// AccountGroups implements /account/:id/groups REST request.
func (server *Server) AccountGroups(w http.ResponseWriter, r *http.Request) {
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

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		id32 = int32(id)

		uid, err := server.checkUserIDCookie(r)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		if id32 != uid {
			server.writeResponse(403, tserrors.Forbidden, w)
			return
		}

		groups, err := server.Dao.AccountGroups(id32)

		if err == gsql.ErrNoRows {
			groups = []*storage.ScoutGroup{}
		} else if err != nil {
			server.Logger.Errorw("storage error:",
				"err", err)

			server.writeResponse(500, tserrors.ServerError, w)
			return
		}

		server.writeResponse(200, groups, w)

	default:
		w.WriteHeader(405)
	}
}

// AccountScouts method implements /account/:id/group/:id/scouts REST requests.
func (server *Server) AccountScouts(w http.ResponseWriter, r *http.Request) {
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

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		id32 = int32(id)

		uid, err := server.checkUserIDCookie(r)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		if id32 != uid {
			server.writeResponse(403, tserrors.Forbidden, w)
			return
		}

		scouts, err := server.Dao.AccountScouts(id32)

		if err == gsql.ErrNoRows {
			scouts = []*storage.Scout{}
		} else if err != nil {
			server.Logger.Errorw("storage error:",
				"err", err)

			server.writeResponse(500, tserrors.ServerError, w)
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

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		id32 = int32(id)

		uid, err := server.checkUserIDCookie(r)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)
			return

		case nil:

		default:
			server.Logger.Debugw("unexpected cookie error",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		if id32 != uid {
			server.writeResponse(403, tserrors.Forbidden, w)
			return
		}

		decoder := json.NewDecoder(r.Body)

		scout := storage.Scout{}
		err = decoder.Decode(&scout)

		if err != nil {
			server.Logger.Debugw("update scout: illegal payload",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		scout.ID = -1

		scoutID, err := server.Dao.InsertOrUpdateScout(scout, uid)

		switch t := err.(type) {
		case tserrors.RestError:
			server.writeResponse(t.Code, t, w)

		case nil:
			server.writeResponse(200, map[string]int32{
				"id": scoutID,
			}, w)

		default:
			server.Logger.Errorw("scout insert failed",
				"err", err)

			server.writeResponse(500, tserrors.ServerError, w)
		}

	default:
		w.WriteHeader(405)
	}
}
