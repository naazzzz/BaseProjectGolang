package model

import "gorm.io/gorm"

type TestModel struct {
	gorm.Model
	Name string `json:"name"`
}
