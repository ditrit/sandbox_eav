// This file contain function that are helpfull to test thing while adding new feature. It should never be added to a production release
package eav

import (
	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

func PopulateDatabase(db *gorm.DB) error {
	// Defining a bird
	colorAttr := &models.Attribut{Name: "color", ValueType: "string"}
	specieAttr := &models.Attribut{Name: "specie", ValueType: "string"}
	heightAttr := &models.Attribut{Name: "height", ValueType: "int"}

	BirdType := &models.EntityType{
		Name: "bird",
	}
	BirdType.Attributs = append(
		BirdType.Attributs, colorAttr, specieAttr, heightAttr,
	)
	val1 := &models.Value{StringVal: "blue", Attribut: colorAttr}
	val2 := &models.Value{StringVal: "mésange", Attribut: specieAttr}
	val3 := &models.Value{IntVal: 8, Attribut: heightAttr}

	mesange := &models.Entity{EntityType: BirdType}
	mesange.Fields = append(mesange.Fields, val1, val2, val3)

	db.Create(mesange)
	log("Finished populating the database")

	return nil
}
