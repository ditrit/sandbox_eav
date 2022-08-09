// This file contain function that are helpfull to test things while adding new feature. It should never be added to a production release
package main

import (
	"fmt"

	"github.com/ditrit/sandbox_eav/eav/models"
	"gorm.io/gorm"
)

func PopulateDatabase(db *gorm.DB) error {
	//defining
	HumanType := &models.EntityType{
		Name: "human",
	}
	nameAttr := &models.Attribut{Name: "name", ValueType: "string", Required: true}
	HumanType.Attributs = append(
		HumanType.Attributs, nameAttr,
	)
	bobName, err := models.NewStringValue(nameAttr, "bob")
	if err != nil {
		panic(err)
	}
	bob := &models.Entity{EntityType: HumanType}
	bob.Fields = append(bob.Fields, bobName)
	db.Create(bob)

	jeanName, err := models.NewStringValue(nameAttr, "jean")
	if err != nil {
		panic(err)
	}
	jean := &models.Entity{EntityType: HumanType}
	jean.Fields = append(jean.Fields, jeanName)

	db.Create(jean)

	// Defining a bird
	colorAttr := &models.Attribut{Name: "color", ValueType: "string", Required: true}
	specieAttr := &models.Attribut{Name: "specie", ValueType: "string", Required: true}
	heightAttr := &models.Attribut{Name: "height", ValueType: "int", Default: true, DefaultInt: 12, Required: false}
	weightAttr := &models.Attribut{Name: "weight", ValueType: "float", Default: true, DefaultFloat: 12.500, Required: false}
	ownerAttr := &models.Attribut{Name: "owner", ValueType: "relation", Required: false, TargetEntityTypeId: HumanType.ID}

	BirdType := &models.EntityType{
		Name: "bird",
	}
	BirdType.Attributs = append(
		BirdType.Attributs, colorAttr, specieAttr, heightAttr, weightAttr, ownerAttr,
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

	val5, err := models.NewRelationValue(ownerAttr, bob)
	if err != nil {
		panic(err)
	}

	mesange := &models.Entity{EntityType: BirdType}
	mesange.Fields = append(mesange.Fields, val1, val2, val3, val4, val5)

	db.Create(mesange)
	fmt.Println("Finished populating the database")

	return nil
}
