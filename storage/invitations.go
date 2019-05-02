package storage

import "time"

// Invitation is an invite to join a Scout Group.
type Invitation struct {
	Token      string      `json:"token"`
	Email      string      `json:"email"`
	Expires    time.Time   `json:"expires"`
	ScoutGroup *ScoutGroup `json:"scout_group"`
}
