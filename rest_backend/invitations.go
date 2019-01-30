package rest_backend

import (
	"net/http"
)

func (server *RestServer) Invitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	/*
	case http.MethodGet:
		result := map[string]interface{}{}

		if cookie, err := r.Cookie("_ts_u"); err == nil && cookie != nil {
			var userId int32

			if err := server.Configuration.SecureCookies.Decode("_ts_u", cookie.Value, &userId); err == nil {
				account, err := server.Dao.QueryAccount(userId)

				if err == nil {
					result["type"] = "account"
					result["account"] = account

				} else {
					server.Logger.Warnf("cant retrieve account: userIdFromCookie=%v, error=%v", userId, err)
				}
			} else {
				server.Logger.Warnf("_ts_u cookie decoding failed: value=%v", cookie.Value)
			}
		} else {
			server.Logger.Debug("_ts_u cookie not found")
		}

		if len(result) == 0 {
			vars := mux.Vars(r)

			token, found := vars["token"]

			if !found {
				w.WriteHeader(400)
				return
			}

			invitation, account, err := server.Dao.QueryInvitationToken(token)

			switch {
			case err == sql.ErrNoRows:
				w.WriteHeader(404)

				result["error"] = NotFound

			case err != nil:
				w.WriteHeader(500)
				server.Logger.Errorf("unexpected error: error=%v", err)

				result["error"] = ServerError

			default:

				if invitation != nil {
					if invitation.Expires.Before(time.Now()) {
						w.WriteHeader(410)

						result["error"] = Expired
					} else {
						result["type"] = "invitation"
						result["invitation"] = invitation
					}
				} else if account != nil {
					result["type"] = "account"
					result["account"] = account
				}
			}
		}

		res, err := json.Marshal(result)

		if err != nil {
			server.Logger.Errorf("json marshalling failed: result=%v, error=%v", r, err)

			w.WriteHeader(500)
			return
		}

		_, err = w.Write(res)

		if err != nil {
			server.Logger.Errorf("unexpected error while writing response: error=%v", err)
		}
*/
/*	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)

		account := storage.NewAccountWithCredentials()
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

