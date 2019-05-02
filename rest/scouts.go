package rest

import (
	"encoding/json"
	"net/http"

	"github.com/makeroo/taxi_scout/storage"
	"github.com/makeroo/taxi_scout/ts_errors"
)

/*
func (server *Server) Scouts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

	default:
		w.WriteHeader(405)
	}
}
*/

// Scout implements /scout/:id REST API.
func (server *Server) Scout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		w.Header().Add("Content-Type", "application/json")

		// TODO: check that id in url and id in scout match
		myID, err := server.checkUserIDCookie(r)

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

		scout := storage.Scout{}
		err = decoder.Decode(&scout)

		if err != nil {
			server.Logger.Debugw("update scout: illegal payload",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		if scout.ID < 0 {
			server.Logger.Errorw("update scout without scout id")

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		scoutID, err := server.Dao.InsertOrUpdateScout(scout, myID)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)

		case nil:
			server.writeResponse(200, map[string]int32{
				"id": scoutID,
			}, w)

		default:
			server.Logger.Errorw("scout update failed",
				"err", err)

			server.writeResponse(500, ts_errors.ServerError, w)
		}

	default:
		w.WriteHeader(405)
	}
}
