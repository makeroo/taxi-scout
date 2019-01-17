package storage

type ScoutGroup struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
}

func NewScoutGroup() *ScoutGroup {
	return &ScoutGroup{0, ""}
}
