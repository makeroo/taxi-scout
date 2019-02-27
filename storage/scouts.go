package storage

type ScoutGroup struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
}

type Scout struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
	GroupId int32 `json:"group"`
}
