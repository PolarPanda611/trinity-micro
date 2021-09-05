package config

// Field field interface .
// when you want to valid your field with
// config loader .
// please implement this interface
type Field interface {
	IsValid() bool
}
