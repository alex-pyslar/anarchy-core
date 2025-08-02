package domain

type Object struct {
	ID           int `db:"id"`
	ObjectListID int `db:"object_list_id"`
}
