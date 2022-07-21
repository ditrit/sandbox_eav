package operations

import (
	"fmt"

	"github.com/ditrit/sandbox_eav/eav/models"
	"github.com/ditrit/sandbox_eav/utils"
	"gorm.io/gorm"
)

func UpdateEntity(db *gorm.DB, et *models.Entity, attrs map[string]interface{}) error {
	for _, a := range et.EntityType.Attributs {
		for _, value := range et.Fields {
			if a.ID == value.AttributId {
				for k, v := range attrs {
					if k == a.Name {
						switch t := v.(type) {
						case string:
							if a.ValueType != "string" {
								return fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
							}
							value.StringVal = v.(string)
						case float64:
							if a.ValueType != "int" && a.ValueType != "float" && a.ValueType != "relation" {
								return fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
							}
							if utils.IsAnInt(v.(float64)) {
								// is an int
								if a.ValueType == "relation" {
									value.RelationVal = uint(v.(float64))
								} else {
									value.IntVal = int(v.(float64))
								}
							} else {
								// is a float
								value.FloatVal = v.(float64)
							}

						case bool:
							if a.ValueType != "bool" {
								return fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
							}
							value.BoolVal = v.(bool)

						case nil:
							if a.Required {
								return fmt.Errorf("can't set a required variable to null (expected=%v got=%T)", a.ValueType, t)
							}
							value.IsNull = true
							value.IntVal = 0
							value.FloatVal = 0.0
							value.StringVal = ""
							value.BoolVal = false
							value.RelationVal = 0

						default:
							panic("mmmh you just discovered a new json type (https://go.dev/blog/json#generic-json-with-interface)")
						}
					}
				}
				value.Attribut = a
				db.Save(value)
			}
		}
	}
	return nil
}
