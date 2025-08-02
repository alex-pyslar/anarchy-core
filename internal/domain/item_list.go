package domain

type ItemList struct {
	ID          int  `db:"id"`
	ObjectID    int  `db:"object_id"`
	Rarity      int  `db:"rarity"`
	IsStackable bool `db:"is_stackable"`
}
