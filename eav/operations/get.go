package operations

import (
	"github.com/ditrit/sandbox_eav/eav/models"
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
