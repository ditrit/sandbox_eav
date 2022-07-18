package eav

import (
	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

func GetEntityType(db *gorm.DB, id uint) *models.EntityType {
	var ett models.EntityType
	db.Preload("Attributs").First(&ett, id)
	return &ett
}

func GetAllEntityType(db *gorm.DB, id uint) []*models.EntityType {
	var etts []*models.EntityType
	db.Preload("Attributs").Find(&etts)
	return etts
}

func GetEntity(db *gorm.DB, id uint) *models.Entity {
	var et models.Entity
	db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").First(&et)
	return &et
}

func GetAllEntity(db *gorm.DB, id uint) []*models.Entity {
	var ets []*models.Entity
	db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&ets)
	return ets
}
