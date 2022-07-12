// This file contain function that are helpfull to test thing while adding new feature. It should never be added to a production release
package eav

import "gorm.io/gorm"

func PopulateDatabase(db *gorm.DB) error {
	//Describing a bird
	val1 := Value{Value: "blue"}
	val2 := Value{Value: "mésange"}
	val3 := Value{Value: "8cm"}

	attr1 := Attribut{Value: val1, Name: "color"}
	attr2 := Attribut{Value: val2, Name: "specie"}
	attr3 := Attribut{Value: val3, Name: "height"}

	mesange := Object{Identifier: "figndkfh", Attributs: []Attribut{attr1, attr2, attr3}}
	db.Create(&mesange)
	return nil
}
