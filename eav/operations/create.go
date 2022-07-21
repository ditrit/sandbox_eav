package operations

import (
	"fmt"

	"github.com/ditrit/sandbox_eav/eav/models"
	"github.com/ditrit/sandbox_eav/utils"
	"gorm.io/gorm"
)

// Create a brand new entity
func CreateEntity(db *gorm.DB, ett *models.EntityType, attrs map[string]interface{}) (*models.Entity, error) {
	var et models.Entity
	for _, a := range ett.Attributs {
		present := false
		var value models.Value
		for k, v := range attrs {
			if k == a.Name {
				present = true
				switch t := v.(type) {
				case string:
					if a.ValueType != "string" {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					value = models.Value{StringVal: v.(string)}

				case float64:
					if a.ValueType != "int" && a.ValueType != "float" && a.ValueType != "relation" {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					if utils.IsAnInt(v.(float64)) {
						// is an int
						if a.ValueType == "relation" {
							value = models.Value{RelationVal: uint(v.(float64))}
						} else {
							value = models.Value{IntVal: int(v.(float64))}
						}
					} else {
						// is a float
						value = models.Value{FloatVal: v.(float64)}
					}

				case bool:
					if a.ValueType != "bool" {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					value = models.Value{BoolVal: v.(bool)}

				case nil:
					if a.Required {
						return nil, fmt.Errorf("can't have a null field with a required attribut")
					}
					value = models.Value{IsNull: true}

				default:
					panic("mmmh you just discovered a new json type (https://go.dev/blog/json#generic-json-with-interface)")
				}
			}
		}
		if a.Required && !present {
			return nil, fmt.Errorf("field %q is missing and is required", a.Name)
		}
		if !present {
			if a.Default {
				v, err := a.GetNewDefaultValue()
				if err != nil {
					return nil, err
				} else {
					value = *v
				}
			} else {
				value = models.Value{IsNull: true}
			}
		}
		value.Attribut = a
		et.Fields = append(et.Fields, &value)
	}
	et.EntityType = ett
	return &et, db.Create(&et).Error
}
