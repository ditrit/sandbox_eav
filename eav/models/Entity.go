package models

import (
	"fmt"

	"github.com/ditrit/sandbox_eav/utils"
)

// Describe an instance of an EntityType
type Entity struct {
	Model
	Fields []*Value

	// GORM relations
	EntityTypeId uint
	EntityType   *EntityType
}

// Encode the Entity to json
func (e *Entity) EncodeToJson() []byte {
	var pairs []string
	pairs = append(pairs,
		fmt.Sprintf("%q: %d", "id", e.ID),
	)
	pairs = append(pairs,
		fmt.Sprintf("%q: %s", "attrs", e.encodeAttributes()),
	)
	return []byte(utils.BuildJsonFromStrings(pairs))
}

// return the attribut in a json encoded string
func (e *Entity) encodeAttributes() string {
	var pairs []string
	var row string
	for _, f := range e.Fields {
		if f.IsNull {
			row = fmt.Sprintf("%q: %s", f.Attribut.Name, "null")
		} else {
			typ := f.Attribut.ValueType
			switch typ {
			case "string":
				stringValue, err := f.GetStringVal()
				if err != nil {
					panic(err)
				}
				row = fmt.Sprintf("%q: %q", f.Attribut.Name, stringValue)
			case "int":
				intValue, err := f.GetIntVal()
				if err != nil {
					panic(err)
				}
				row = fmt.Sprintf("%q: %d", f.Attribut.Name, intValue)
			case "float":
				floatValue, err := f.GetFloatVal()
				if err != nil {
					panic(err)
				}
				row = fmt.Sprintf("%q: %f", f.Attribut.Name, floatValue)
			case "bool":
				boolValue, err := f.GetBoolVal()
				if err != nil {
					panic(err)
				}
				row = fmt.Sprintf("%q: %t", f.Attribut.Name, boolValue)
			case "relation":
				row = fmt.Sprintf("%q: %d", f.Attribut.Name, f.RelationVal)
			default:
				panic(fmt.Errorf("the type %q is supported type by the EAV (not implemented)", typ))
			}
		}
		pairs = append(pairs, row)
	}
	return utils.BuildJsonFromStrings(pairs)
}
