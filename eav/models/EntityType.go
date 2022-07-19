package models

// Describe an object type
type EntityType struct {
	Model
	Name string

	// GORM relations
	Attributs []*Attribut
}
