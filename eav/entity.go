package eav

import (
	"fmt"

	"github.com/ditrit/sandbox_eav/eav/models"
	"github.com/ditrit/sandbox_eav/utils"
	"gorm.io/gorm"
)

func GetEntityType(db *gorm.DB, id uint) (*models.EntityType, error) {
	var ett models.EntityType
	err := db.Preload("Attributs").First(&ett, id).Error
	if err != nil {
		return nil, err
	}
	return &ett, nil
}

func GetEntityTypeByName(db *gorm.DB, name string) (*models.EntityType, error) {
	var ett models.EntityType
	err := db.Preload("Attributs").First(&ett, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &ett, nil
}

func GetAllEntityType(db *gorm.DB) []*models.EntityType {
	var etts []*models.EntityType
	db.Preload("Attributs").Find(&etts)
	return etts
}

func GetEntity(db *gorm.DB, id uint) (*models.Entity, error) {
	var et models.Entity
	err := db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").First(&et, id).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

// Return all Entities for the ett EntityType
func GetEntities(db *gorm.DB, ett *models.EntityType) []*models.Entity {
	var ets []*models.Entity
	db.Where("entity_type_id = ?", ett.ID).Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&ets)
	return ets
}

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
					if a.ValueType != "int" && a.ValueType != "float" {
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
					if !a.IsNullable {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					value = models.Value{IsNull: true}

				default:
					panic("mmmh you just discovered a new json type (https://go.dev/blog/json#generic-json-with-interface)")
				}
			}
		}
		if !a.IsNullable && !present {
			return nil, fmt.Errorf("field %q is missing and can't be null", a.Name)
		}
		if !present {
			value = models.Value{IsNull: true}
		}
		value.Attribut = a
		et.Fields = append(et.Fields, &value)
	}
	et.EntityType = ett
	return &et, db.Create(&et).Error
}

func ModifyEntity(db *gorm.DB, ett *models.EntityType, et *models.Entity, attrs map[string]interface{}) error {

	return nil
}
