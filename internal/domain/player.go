package domain

type Player struct {
	ID         int     `db:"id"`
	UserID     int     `db:"user_id"`
	X          float64 `db:"x"`
	Y          float64 `db:"y"`
	Z          float64 `db:"z"`
	PlayerName string  `db:"name"`
}
