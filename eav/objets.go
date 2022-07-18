package eav

import (
	"errors"
	"fmt"
	"reflect"

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
	err := db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").First(&et).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

func GetAllEntity(db *gorm.DB) []*models.Entity {
	var ets []*models.Entity
	db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&ets)
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
				switch v.(type) {
				case string:
					if a.ValueType != "string" {
						panic(fmt.Errorf("types dont match (expected=%v got=%v)", a.ValueType, reflect.ValueOf(v).Type().String()))
					}
					value = models.Value{StringVal: v.(string)}

				case float64:
					if a.ValueType != "int" && a.ValueType != "float" {
						panic(fmt.Errorf("types dont match (expected=%v got=%v)", a.ValueType, reflect.ValueOf(v).Type().String()))
					}
					if utils.IsAnInt(v.(float64)) {
						// is an int
						value = models.Value{IntVal: int(v.(float64))}
					} else {
						// is a float
						value = models.Value{FloatVal: v.(float64)}
					}

				case bool:
					if a.ValueType != "bool" {
						panic(fmt.Errorf("types dont match (expected=%v got=%v)", a.ValueType, reflect.ValueOf(v).Type().String()))
					}
					value = models.Value{BoolVal: v.(bool)}

				case nil:
					return nil, errors.New("null type not supported")

				default:
					panic("mmmh you just discovered a new json type (https://go.dev/blog/json#generic-json-with-interface)")
				}
			}
		}
		if !a.CanBeNull && !present {
			return nil, fmt.Errorf("field %q can't be null and is missing", a.Name)
		}
		value.Attribut = a
		et.Fields = append(et.Fields, &value)
	}
	return &et, db.Create(&et).Error
}
