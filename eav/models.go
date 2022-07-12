package eav

import "gorm.io/gorm"

type EntityType struct {
	gorm.Model
	Name      string
	Attributs []Attribut
	Instances []Entity
}

type Entity struct {
	gorm.Model
	Identifier   string
	Fields       []Value
	EntityTypeId uint
}

type Attribut struct {
	gorm.Model
	EntityTypeId uint
	Name         string
	Values       []Value
}

type Value struct {
	gorm.Model
	Value      string
	EntityId   uint
	AttributId uint
}
