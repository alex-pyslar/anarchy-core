package domain

type Inventory struct {
	ID       string `db:"id"`        // Идентификатор инвенторя
	EntityID string `db:"entity_id"` // Идентификатор энтити
	ItemID   string `db:"item_id"`   // Идетификатор предмета
}
