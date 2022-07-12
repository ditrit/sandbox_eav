package main

import (
	"log"

	"github.com/ditrit/sandbox_eav/eav"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// were using sqlite now since it's easy to use but we will move to cockroachDb as soon as possible
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connected to database")
	eav.Init(db)
	log.Println("Automigrate finished")
	eav.PopulateDatabase(db)
	log.Println("Finished populating the database")
}
