package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Base for GORM, see https://gorm.io/docs/models.html#Declaring-Models
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// Describe an object type
type EntityType struct {
	Model
	Name string

	// GORM relations
	Attributs []*Attribut
}

// Describe the attribut of a en EntityType
type Attribut struct {
	Model
	Name      string
	CanBeNull bool

	// the type the values of this attr are. Can be "int", "float", "string", "bool"
	ValueType string

	// GORM relations
	EntityTypeId uint
}

// Describe an instance of an EntityType
type Entity struct {
	Model
	Fields []*Value

	// GORM relations
	EntityTypeId uint
	EntityType   *EntityType
}

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

var ErrValueIsNull = errors.New("You can't get the value from a null Value")

// Check if the Value is whole. eg, no fields are nil
func (v *Value) CheckWhole() error {
	if v.Attribut == nil {
		return fmt.Errorf("The Attribut pointer is nil in Value at %v", v)
	}
	return nil
}

func (v *Value) GetStringVal() string {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		panic(ErrValueIsNull)
	}
	return v.StringVal
}

func (v *Value) GetFloatVal() float64 {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		panic(ErrValueIsNull)
	}
	return v.FloatVal
}

func (v *Value) GetIntVal() int {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		panic(ErrValueIsNull)
	}
	return v.IntVal
}

func (v *Value) GetBoolVal() bool {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		panic(ErrValueIsNull)
	}
	return v.BoolVal
}

func (e *Entity) EncodeToJson() []byte {
	var pairs []string
	pairs = append(pairs,
		fmt.Sprintf("%q: %d", "id", e.ID),
	)
	pairs = append(pairs,
		fmt.Sprintf("%q: %s", "attrs", e.encodeAttributes()),
	)
	return []byte(buildJson(pairs))
}

func (e *Entity) encodeAttributes() string {
	var pairs []string
	var row string
	for _, f := range e.Fields {
		if f.IsNull {
			continue
		}
		typ := f.Attribut.ValueType
		switch typ {
		case "string":
			row = fmt.Sprintf("%q: %q", f.Attribut.Name, f.GetStringVal())
		case "int":
			row = fmt.Sprintf("%q: %d", f.Attribut.Name, f.GetIntVal())
		case "float":
			row = fmt.Sprintf("%q: %f", f.Attribut.Name, f.GetFloatVal())
		case "bool":
			row = fmt.Sprintf("%q: %t", f.Attribut.Name, f.GetBoolVal())
		default:
			panic(fmt.Errorf("the type %q is not a json type", typ))
		}
		pairs = append(pairs, row)
	}
	return buildJson(pairs)
}

func buildJson(pairs []string) string {
	var b strings.Builder
	// starting the json string
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ","))
	//ending the json string
	b.WriteString("}")
	return b.String()
}
