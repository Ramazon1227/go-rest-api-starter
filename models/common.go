package models

type PrimaryKey struct {
	Id string `json:"id"`
}

type QueryParam struct {
	Limit  int32
	Page   int32
	Offset int32
	Search string
	Role   string
}