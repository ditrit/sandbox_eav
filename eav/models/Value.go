package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Describe the attribut value of an Entity
type Value struct {
	Model
	IsNull      bool
	StringVal   string
	FloatVal    float64
	IntVal      int
	BoolVal     bool
	RelationVal uint

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
	if attr.Required {
		return nil, fmt.Errorf("can't create new null value for a required attribut")
	}
	val.IsNull = true
	val.Attribut = attr
	return val, nil
}

// Create a new int value
func NewIntValue(attr *Attribut, i int) (*Value, error) {
	val := new(Value)
	if attr.ValueType != IntValueType {
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
	if attr.ValueType != BooleanValueType {
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
	if attr.ValueType != FloatValueType {
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
	if attr.ValueType != StringValueType {
		return nil, fmt.Errorf("can't create a new string value with a %s attribut", attr.ValueType)
	}
	val.IsNull = false
	val.StringVal = s
	val.Attribut = attr
	return val, nil
}

// Create a new relation value.
// If et is nil, then the function return an error
// If et is of the wrong types
func NewRelationValue(attr *Attribut, et *Entity) (*Value, error) {
	val := new(Value)
	if attr.ValueType != RelationValueType {
		return nil, fmt.Errorf("can't create a new relation value with a %s attribut", attr.ValueType)
	}
	if et == nil {
		return nil, fmt.Errorf("can't create a new relation with a nill entity pointer")
	}
	if et.EntityType.ID != attr.TargetEntityTypeId {
		return nil, fmt.Errorf("can't create a relation with an entity of wrong EntityType. (got the entityid=%d, expected=%d)", et.EntityType.ID, attr.TargetEntityTypeId)
	}
	val.IsNull = false
	val.RelationVal = et.ID
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
		if v.Attribut.Default {
			return v.Attribut.DefaultString, nil
		}
		return "", ErrValueIsNull
	}
	if v.Attribut.ValueType != StringValueType {
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
		if v.Attribut.Default {
			return v.Attribut.DefaultFloat, nil
		}
		return 0.0, ErrValueIsNull
	}
	if v.Attribut.ValueType != FloatValueType {
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
		if v.Attribut.Default {
			return v.Attribut.DefaultInt, nil
		}
		return 0, ErrValueIsNull
	}
	if v.Attribut.ValueType != IntValueType {
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
		if v.Attribut.Default {
			return v.Attribut.DefaultBool, nil
		}
		return false, ErrValueIsNull
	}
	if v.Attribut.ValueType != BooleanValueType {
		return false, ErrAskingForWrongType
	}
	return v.BoolVal, nil
}

// Return the Relation value as a *Entity
// If the Value is null, it return nil
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribut == nil, it panic
func (v *Value) GetRelationVal(db *gorm.DB) (*Entity, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	if v.Attribut.ValueType != RelationValueType {
		return nil, ErrAskingForWrongType
	}

	if v.IsNull {
		return nil, nil
	}
	var et Entity
	err = db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").First(&et, v.RelationVal).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

// Return the underlying value as an interface
func (v *Value) Value() interface{} {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	switch v.Attribut.ValueType {
	case StringValueType:
		s, err := v.GetStringVal()
		if err != nil {
			if errors.Is(err, ErrValueIsNull) {
				return nil
			}
			//else panic
			panic(err)
		}
		return s
	case IntValueType:
		i, err := v.GetIntVal()
		if err != nil {
			if errors.Is(err, ErrValueIsNull) {
				return nil
			}
			//else panic
			panic(err)
		}
		return i
	case FloatValueType:
		f, err := v.GetFloatVal()
		if err != nil {
			if errors.Is(err, ErrValueIsNull) {
				return nil
			}
			//else panic
			panic(err)
		}
		return f
	case BooleanValueType:
		b, err := v.GetFloatVal()
		if err != nil {
			if errors.Is(err, ErrValueIsNull) {
				return nil
			}
			//else panic
			panic(err)
		}
		return b
	case RelationValueType:
		if v.IsNull {
			return nil
		}
		return v.RelationVal
	default:
		panic(fmt.Errorf(
			"hmm this Attribut.ValueType does not exists (got=%s)",
			v.Attribut.ValueType,
		))
	}
}

// When Value isNull, it is impossible to build a Key/Value pair
var ErrCantBuildKVPairForNullValue = errors.New("can't build key/value pair from null value")

// Build a key/value pair to be included in a JSON
// If the value hold an int=8 with an attribut named "voila" then the string returned will be `"voila":8`
func (v *Value) BuildJsonKVPair() (string, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	var row string
	typ := v.Attribut.ValueType
	switch typ {
	case StringValueType:
		stringValue, err := v.GetStringVal()
		if err != nil {
			if err == ErrValueIsNull {
				return "", ErrCantBuildKVPairForNullValue
			}
			panic(err)
		}

		row = fmt.Sprintf("%q: %q", v.Attribut.Name, stringValue)
	case IntValueType:
		intValue, err := v.GetIntVal()
		if err != nil {
			if err == ErrValueIsNull {
				return "", ErrCantBuildKVPairForNullValue
			}
			panic(err)
		}
		row = fmt.Sprintf("%q: %d", v.Attribut.Name, intValue)
	case FloatValueType:
		floatValue, err := v.GetFloatVal()
		if err != nil {
			if err == ErrValueIsNull {
				return "", ErrCantBuildKVPairForNullValue
			}
			panic(err)
		}
		row = fmt.Sprintf("%q: %f", v.Attribut.Name, floatValue)
	case BooleanValueType:
		boolValue, err := v.GetBoolVal()
		if err != nil {
			if err == ErrValueIsNull {
				return "", ErrCantBuildKVPairForNullValue
			}
			panic(err)
		}
		row = fmt.Sprintf("%q: %t", v.Attribut.Name, boolValue)
	case RelationValueType:
		row = fmt.Sprintf("%q: %d", v.Attribut.Name, v.RelationVal)
	default:
		panic(fmt.Errorf("the type %q is not supported by the EAV (not implemented)", typ))

	}
	return row, nil
}
