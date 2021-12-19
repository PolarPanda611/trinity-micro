package model

import "github.com/PolarPanda611/trinity-micro/core/dbx"

type Price struct {
	dbx.Model
	ItemID uint64
	Item   *Item
	Price  int
}
