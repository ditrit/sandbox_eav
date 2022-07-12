package eav

import "gorm.io/gorm"

type Object struct {
	gorm.Model
	Identifier string
	Attributs  []Attribut
}

type Attribut struct {
	gorm.Model
	ValueId  int
	ObjectId int
	Name     string
	Value    Value
}

type Value struct {
	gorm.Model
	Value string
}
