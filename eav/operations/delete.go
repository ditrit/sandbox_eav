package operations

import (
	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

// Delete an entity
func DeleteEntity(db *gorm.DB, et *models.Entity) error {
	for _, v := range et.Fields {
		err := db.Delete(v).Error
		if err != nil {
			return err
		}
	}
	return db.Delete(et).Error
}
