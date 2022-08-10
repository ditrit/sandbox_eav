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
		fmt.Sprintf("%q: %q", "type", e.EntityType.Name),
	)
	pairs = append(pairs,
		fmt.Sprintf("%q: %s", "attrs", e.encodeAttributes()),
	)
	return []byte(utils.BuildJsonFromStrings(pairs))
}

// return the attribut in a json encoded string
func (e *Entity) encodeAttributes() string {
	var pairs []string
	for _, f := range e.Fields {
		if f.IsNull {
			continue
		}
		pair, _ := f.BuildJsonKVPair()
		pairs = append(pairs, pair)
	}
	return utils.BuildJsonFromStrings(pairs)
}

func (e *Entity) GetValue(attrName string) (interface{}, error) {
	var attrId uint = 0
	for _, a := range e.EntityType.Attributs {
		if a.Name == attrName {
			attrId = a.ID
			break
		}
	}
	if attrId == 0 {
		return nil, fmt.Errorf("attr not found: got=%s", attrName)
	}
	for _, v := range e.Fields {
		if v.AttributId == attrId {
			return v.Value(), nil
		}
	}
	return nil, fmt.Errorf("value not found")
}
