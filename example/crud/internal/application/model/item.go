package model

import "github.com/PolarPanda611/trinity-micro/core/dbx"

type Item struct {
	dbx.Model
	Code string
	Name string
}
