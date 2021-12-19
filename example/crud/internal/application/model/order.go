package model

import "github.com/PolarPanda611/trinity-micro/core/dbx"

type Order struct {
	dbx.Model
	Code         string
	UserID       uint64
	User         *User
	OrderDetails []OrderDetail
	Total        string
}

type OrderDetail struct {
	dbx.Model
	Code     string
	PriceID  uint64
	Quantity int
	Total    string
}
