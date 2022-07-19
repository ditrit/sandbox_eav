package models

// Describe the attribut of a en EntityType
type Attribut struct {
	Model
	Name       string
	IsNullable bool
	Unique     bool
	// Required   bool // FIXME: to implement

	// Default values
	Default       bool // there is a default value
	DefaultInt    int
	DefaultBool   bool
	DefaultString string
	DefaultFloat  float64

	// the type the values of this attr are. Can be "int", "float", "string", "bool", "relationid"
	ValueType          string
	TargetEntityTypeId uint

	// GORM relations
	EntityTypeId uint
}
