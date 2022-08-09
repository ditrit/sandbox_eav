package operations

import (
	"errors"
	"fmt"

	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

var (
	ErrIdDontMatchEntityType = errors.New("this object doesn't belong to this type")
)

// Get EntityType by id (uint)
func GetEntityType(db *gorm.DB, id uint) (*models.EntityType, error) {
	var ett models.EntityType
	err := db.Preload("Attributs").First(&ett, id).Error
	if err != nil {
		return nil, err
	}
	return &ett, nil
}

// Get EntityType by name (string)
func GetEntityTypeByName(db *gorm.DB, name string) (*models.EntityType, error) {
	var ett models.EntityType
	err := db.Preload("Attributs").First(&ett, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(" EntityType named %q not found", name)
		}
	}
	return &ett, nil
}

// Get all the types defined in the schema
func GetAllEntityType(db *gorm.DB) ([]*models.EntityType, error) {
	var etts []*models.EntityType
	err := db.Preload("Attributs").Find(&etts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("0 EntityType found")
		}
	}
	return etts, nil
}

func GetEntity(db *gorm.DB, ett *models.EntityType, id uint) (*models.Entity, error) {
	var et models.Entity
	err := db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").First(&et, id).Error
	if err != nil {
		return nil, err
	}
	if ett.ID != et.EntityTypeId {
		return nil, ErrIdDontMatchEntityType
	}
	return &et, nil
}

// Return all Entities for the ett EntityType
func GetEntities(db *gorm.DB, ett *models.EntityType) ([]*models.Entity, error) {
	var ets []*models.Entity
	tx := db.Where("entity_type_id = ?", ett.ID).Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&ets)
	if tx.Error != nil {
		return make([]*models.Entity, 0), tx.Error
	}

	return ets, nil
}
