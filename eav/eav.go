package eav

import (
	"fmt"

	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

const EavPrefix = "[EAV]"

func log(v ...any) {
	fmt.Println(EavPrefix + " " + fmt.Sprintln(v...))
}

func Init(db *gorm.DB) error {
	// Migrate the schema

	db.Migrator().DropTable(&models.EntityType{}, &models.Entity{}, &models.Attribut{}, &models.Value{})
	err := db.AutoMigrate(
		&models.EntityType{}, &models.Entity{}, &models.Attribut{}, &models.Value{},
	)
	if err != nil {
		log("migration failed error: ", err.Error())
		return err
	}
	return nil
}
