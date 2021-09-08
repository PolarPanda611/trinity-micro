package model

type Order struct {
	ID           uint64
	Code         string
	UserID       uint64
	User         *User
	OrderDetails []OrderDetail
	Total        string
}

type OrderDetail struct {
	ID       uint64
	Code     string
	PriceID  uint64
	Quantity int
	Total    string
}
