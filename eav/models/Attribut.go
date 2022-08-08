package models

import "fmt"

const (
	RelationValueType string = "relation"
	BooleanValueType  string = "bool"
	StringValueType   string = "string"
	IntValueType      string = "int"
	FloatValueType    string = "float"
)

// Describe the attribut of a en EntityType
type Attribut struct {
	Model
	Name     string
	Unique   bool
	Required bool

	// Default values
	Default       bool // there is a default value
	DefaultInt    int
	DefaultBool   bool
	DefaultString string
	DefaultFloat  float64

	// the type the values of this attr are. Can be "int", "float", "string", "bool", "relation"
	ValueType          string
	TargetEntityTypeId uint // name of the EntityType

	// GORM relations
	EntityTypeId uint
}

func (a *Attribut) GetNewDefaultValue() (*Value, error) {
	switch a.ValueType {
	case StringValueType:
		v, err := NewStringValue(a, a.DefaultString)
		if err != nil {
			return nil, err
		}
		return v, nil
	case IntValueType:
		v, err := NewIntValue(a, a.DefaultInt)
		if err != nil {
			return nil, err
		}
		return v, nil
	case FloatValueType:
		v, err := NewFloatValue(a, a.DefaultFloat)
		if err != nil {
			return nil, err
		}
		return v, nil
	case BooleanValueType:
		v, err := NewBoolValue(a, a.DefaultBool)
		if err != nil {
			return nil, err
		}
		return v, nil
	case RelationValueType:
		return nil, fmt.Errorf("can't provide default value for relations")
	default:
		panic("hmmm we are not supposed to be here")
	}
}
