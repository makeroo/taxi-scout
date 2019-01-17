package rest_backend

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (server *RestServer) Invitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
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
			return
		case err != nil:
			w.WriteHeader(500)
			server.Logger.Errorf("unexpected error: error=%v", err)
			return
		}

		var r map[string]interface{}

		if invitation != nil {
			r = map[string]interface{}{
				"type": "invitation",
				"invitation": invitation,
			}
		} else if account != nil {
			r = map[string]interface{}{
				"type": "account",
				"account": account,
			}
		}
		res, err := json.Marshal(r)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		_, err = w.Write(res)

		if err != nil {
			server.Logger.Errorf("unexpected error while writing response: error=%v", err)
		}

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

