package models

import (
	"errors"
	"fmt"
)

// Describe the attribut value of an Entity
type Value struct {
	Model
	IsNull    bool
	StringVal string
	FloatVal  float64
	IntVal    int
	BoolVal   bool

	// GORM relations
	EntityId   uint
	AttributId uint
	Attribut   *Attribut
}

var (
	ErrValueIsNull        = errors.New("you can't get the value from a null Value")
	ErrAskingForWrongType = errors.New("you can't get this type of value, the attribut type doesn't match")
)

// Create a new null value
func NewNullValue(attr *Attribut) (*Value, error) {
	val := new(Value)
	if !attr.IsNullable {
		return nil, fmt.Errorf("can't create new null with an non nullable attribut")
	}
	val.IsNull = true
	val.Attribut = attr
	return val, nil
}

// Create a new int value
func NewIntValue(attr *Attribut, i int) (*Value, error) {
	val := new(Value)
	if attr.ValueType != "int" {
		return nil, fmt.Errorf("can't create a new int value with a %s attribut", attr.ValueType)
	}
	val.IsNull = false
	val.IntVal = i
	val.Attribut = attr
	return val, nil
}

// Create a new bool value
func NewBoolValue(attr *Attribut, b bool) (*Value, error) {
	val := new(Value)
	if attr.ValueType != "bool" {
		return nil, fmt.Errorf("can't create a new bool value with a %s attribut", attr.ValueType)
	}
	val.IsNull = false
	val.BoolVal = b
	val.Attribut = attr
	return val, nil
}

// Create a new float value
func NewFloatValue(attr *Attribut, f float64) (*Value, error) {
	val := new(Value)
	if attr.ValueType != "float" {
		return nil, fmt.Errorf("can't create a new float value with a %s attribut", attr.ValueType)
	}
	val.IsNull = false
	val.FloatVal = f
	val.Attribut = attr
	return val, nil
}

// Create a new string value
func NewStringValue(attr *Attribut, s string) (*Value, error) {
	val := new(Value)
	if attr.ValueType != "string" {
		return nil, fmt.Errorf("can't create a new string value with a %s attribut", attr.ValueType)
	}
	val.IsNull = false
	val.StringVal = s
	val.Attribut = attr
	return val, nil
}

// Check if the Value is whole. eg, no fields are nil
func (v *Value) CheckWhole() error {
	if v.Attribut == nil {
		return fmt.Errorf("the Attribut pointer is nil in Value at %v", v)
	}
	return nil
}

// Return the string value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribut == nil, it panic
func (v *Value) GetStringVal() (string, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		return "", ErrValueIsNull
	}
	if v.Attribut.ValueType != "string" {
		return "", ErrAskingForWrongType
	}
	return v.StringVal, nil
}

// Return the float value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribut == nil, it panic
func (v *Value) GetFloatVal() (float64, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		return 0.0, ErrValueIsNull
	}
	if v.Attribut.ValueType != "float" {
		return 0.0, ErrAskingForWrongType
	}
	return v.FloatVal, nil
}

// Return the int value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribut == nil, it panic
func (v *Value) GetIntVal() (int, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		return 0, ErrValueIsNull
	}
	if v.Attribut.ValueType != "int" {
		return 0, ErrAskingForWrongType
	}
	return v.IntVal, nil
}

// Return the bool value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribut == nil, it panic
func (v *Value) GetBoolVal() (bool, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		return false, ErrValueIsNull
	}
	if v.Attribut.ValueType != "bool" {
		return false, ErrAskingForWrongType
	}
	return v.BoolVal, nil
}
