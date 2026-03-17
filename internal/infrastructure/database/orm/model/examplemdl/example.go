package examplemdl

import "gorm.io/gorm"

type ExampleModel struct {
	gorm.Model
	Data string `gorm:"column:data;not null"`
}
