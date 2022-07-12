package eav

import (
	"fmt"

	"gorm.io/gorm"
)

const EavPrefix = "[EAV]"

func log(v ...any) {
	fmt.Println(EavPrefix + " " + fmt.Sprintln(v))
}

func Init(db *gorm.DB) error {
	// Migrate the schema
	err := db.AutoMigrate(
		&Object{}, &Attribut{}, &Value{},
	)
	if err != nil {
		log("migration failed error: ", err.Error())
		return err
	}
	return nil
}
