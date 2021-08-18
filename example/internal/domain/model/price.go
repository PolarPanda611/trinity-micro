package model

type Price struct {
	ID     uint64
	ItemID uint64
	Item   *Item
	Price  int
}
