package model

type User struct {
	ID       uint64
	Username string
	Password string
	Email    string
	Age      uint
	Gender   uint
	Orders   []Order
}
