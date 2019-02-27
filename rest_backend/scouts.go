package rest_backend

import (
	"encoding/json"
	"github.com/makeroo/taxi_scout/storage"
	"github.com/makeroo/taxi_scout/ts_errors"
	"net/http"
)

func (server *RestServer) Scouts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

	default:
		w.WriteHeader(405)
	}
}

func (server *RestServer) Scout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		w.Header().Add("Content-Type", "application/json")

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

		scout := storage.Scout{}
		err = decoder.Decode(&scout)

		if err != nil {
			server.Logger.Debugw("update scout: illegal payload",
				"err", err)

			server.writeResponse(400, ts_errors.BadRequest, w)
			return
		}

		scoutId, err := server.Dao.UpdateScout(scout, myId)

		switch t := err.(type) {
		case ts_errors.RestError:
			server.writeResponse(t.Code, t, w)

		case nil:
			server.writeResponse(200, map[string]int32{
				"id": scoutId,
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
