// This file contain function that are helpfull to test thing while adding new feature. It should never be added to a production release
package eav

import (
	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

func PopulateDatabase(db *gorm.DB) error {
	// Defining a bird
	colorAttr := &models.Attribut{Name: "color", ValueType: "string", IsNullable: false}
	specieAttr := &models.Attribut{Name: "specie", ValueType: "string", IsNullable: false}
	heightAttr := &models.Attribut{Name: "height", ValueType: "int", IsNullable: true}
	weightAttr := &models.Attribut{Name: "weight", ValueType: "float", IsNullable: true}

	BirdType := &models.EntityType{
		Name: "bird",
	}
	BirdType.Attributs = append(
		BirdType.Attributs, colorAttr, specieAttr, heightAttr, weightAttr,
	)
	val1, err := models.NewStringValue(colorAttr, "blue")
	if err != nil {
		panic(err)
	}
	val2, err := models.NewStringValue(specieAttr, "mésange")
	if err != nil {
		panic(err)
	}
	val3, err := models.NewIntValue(heightAttr, 8)
	if err != nil {
		panic(err)
	}
	val4, err := models.NewNullValue(weightAttr)
	if err != nil {
		panic(err)
	}

	mesange := &models.Entity{EntityType: BirdType}
	mesange.Fields = append(mesange.Fields, val1, val2, val3, val4)

	db.Create(mesange)
	log("Finished populating the database")

	return nil
}
