package domain

type EntityList struct {
	ID           int     `db:"id"`
	ObjectListID int     `db:"object_list_id"`
	Damage       float64 `db:"damage"`
	Speed        float64 `db:"speed"`
	Cooldown     float64 `db:"cooldown"`
	DamageRadius float64 `db:"damage_radius"`
	IsAngry      bool    `db:"is_angry"`
	VisualRadius float64 `db:"visual_radius"`
	MaxHealth    float64 `db:"max_health"`
	Model        string  `db:"model"`
	Spawn        string  `db:"spawn"`
	IsOpen       bool    `db:"is_open"`
	IsSpawning   bool    `db:"is_spawning"`
	IsPickUp     bool    `db:"is_pick_up"`
}
