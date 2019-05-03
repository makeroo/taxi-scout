package rest

import (
	"encoding/json"
	"net/http"

	tserrors "github.com/makeroo/taxi_scout/errors"
)

// InvitationsRequest models /invitations request payload.
type InvitationsRequest struct {
	Email string `json:"email"`
}

// Invitations implements /invitations REST request.
func (server *Server) Invitations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		invitation := InvitationsRequest{}
		err := decoder.Decode(&invitation)

		if err != nil {
			server.Logger.Debugw("invitation: illegal payload",
				"err", err)

			server.writeResponse(400, tserrors.BadRequest, w)
			return
		}

		account, err := server.checkUserCookie(r)

		switch err.(type) {
		case tserrors.RestError:
			break

		case nil:
			if account.Email == invitation.Email {
				server.writeResponse(200, map[string]interface{}{
					"authenticated": true,
				}, w)
			} else {
				server.writeResponse(403, tserrors.Forbidden, w)
			}
			return

		default:
			server.Logger.Errorw("unexpected dao error",
				"err", err)

			server.writeResponse(500, tserrors.ServerError, w)
			return
		}

		newInvitation, err := server.Dao.CreateInvitationForExistingMember(invitation.Email)

		if err != nil {
			server.writeResponse(500, tserrors.BadRequest, w)
			return
		}

		server.Logger.Infow("invitation created",
			"invitation", newInvitation)

		// TODO: send email

		server.writeResponse(200, map[string]interface{}{
			"exipires": newInvitation.Expires,
		}, w)

	default:
		w.WriteHeader(405)
	}
}
