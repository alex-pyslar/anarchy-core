package domain

type Item struct {
	ID         int `db:"id"` // Идетификатор предмета
	ObjectID   int `db:"object_id"`
	ItemListID int `db:"item_list_id"`
}
