package model

import "github.com/PolarPanda611/trinity-micro/core/dbx"

type Tenant struct {
	dbx.Model
	Name string
}
