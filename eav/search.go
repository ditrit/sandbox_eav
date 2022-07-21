package eav

import (
	"strconv"

	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

func GetEntitiesWithParams(db *gorm.DB, ett *models.EntityType, params map[string]string) []*models.Entity {
	var ets []*models.Entity
	db.Where("entity_type_id = ?", ett.ID).Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&ets)
	var resultSet []*models.Entity
	var keep bool
	for _, et := range ets {
		keep = true
		for _, value := range et.Fields {
			for k, v := range params {
				if k == value.Attribut.Name {
					switch value.Attribut.ValueType {
					case "string":
						if v != value.StringVal {
							keep = false
						}
					case "int":
						intVal, err := strconv.Atoi(v)
						if err != nil {
							break
						}
						if intVal != value.IntVal {
							keep = false
						}
					case "float":
						floatVal, err := strconv.ParseFloat(v, 64)
						if err != nil {
							break
						}
						if floatVal != value.FloatVal {
							keep = false
						}
					case "relation":
						relVal, err := strconv.Atoi(v)
						if err != nil {
							break
						}
						if relVal != int(value.RelationVal) {
							keep = false
						}
					case "bool":
						boolVal, err := strconv.ParseBool(v)
						if err != nil {
							break
						}
						if boolVal != value.BoolVal {
							keep = false
						}
					}
				}
			}
		}
		if keep {
			resultSet = append(resultSet, et)
		}
	}
	return resultSet
}
