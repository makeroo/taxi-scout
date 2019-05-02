package storage

// ScoutGroup collects informations about a scout group.
type ScoutGroup struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Scout collects informations about a scout.
// A scout belongs to a scout group.
type Scout struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	GroupID int32  `json:"group"`
}
