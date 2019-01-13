package rest_backend

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/makeroo/taxi_scout/storage"

	"github.com/gorilla/mux"
)

func (server *RestServer) Accounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		accounts, err := server.Dao.QueryAccounts()

		if err != nil {
			w.WriteHeader(500)
			return
		}

		res, err := json.Marshal(accounts)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(res)

	case http.MethodPost:
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

	default:
		w.WriteHeader(405)
	}
}

func (server *RestServer) Account(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			w.WriteHeader(400)
			return
		}

		accounts, err := server.Dao.QueryAccount(int32(id))

		if err != nil {
			w.WriteHeader(500)
			return
		}

		res, err := json.Marshal(accounts)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(res)

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
		decoder := json.NewDecoder(r.Body)

		credentials := AccountAuthenticatePayload{"", ""}
		err := decoder.Decode(&credentials)

		if err != nil {
			server.Logger.Errorw("payload unmarshalling failed",
				"err", err,
			)
			w.WriteHeader(400)
			return
		}

		accountId, err := server.Dao.AuthenticateAccount(credentials.Email, credentials.Pwd)

		if err != nil {
			server.Logger.Errorw("authentication failed",
				"err", err,
				"email", credentials.Email,
			)
			w.WriteHeader(401)
			return
		}

		bytes, err := json.Marshal(map[string]int32{"id": accountId})
		w.Write(bytes)

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

		bytes, err := json.Marshal(map[string]int32{"id": credentials.Id})
		w.Write(bytes)

	default:
		w.WriteHeader(405)
	}
}
