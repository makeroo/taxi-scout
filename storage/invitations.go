package storage

import "time"

type Invitation struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Expires time.Time `json:"expires"`
	ScoutGroup *ScoutGroup `json:"scout_group"`
}
