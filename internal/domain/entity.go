package domain

type Entity struct {
	ID           int     `db:"id"`
	ObjectID     int     `db:"object_id"`
	EntityListID int     `db:"entity_list_id"`
	Health       float64 `db:"health"`
	X            float64 `db:"x"`
	Y            float64 `db:"y"`
	Z            float64 `db:"z"`
}
