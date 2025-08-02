package domain

type ObjectList struct {
	ID          int    `db:"id"`
	name        string `db:"name"`
	image       string `db:"image"`
	description string `db:"description"`
}
