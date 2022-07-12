// This file contain function that are helpfull to test thing while adding new feature. It should never be added to a production release
package eav

import "gorm.io/gorm"

func PopulateDatabase(db *gorm.DB) error {
	//Defining a bird
	colorAttr := Attribut{Name: "color"}
	specieAttr := Attribut{Name: "specie"}
	heightAttr := Attribut{Name: "height"}

	val1 := Value{Value: "blue"}
	val2 := Value{Value: "mésange"}
	val3 := Value{Value: "8cm"}

	colorAttr.Values = append(colorAttr.Values, val1)
	specieAttr.Values = append(specieAttr.Values, val2)
	heightAttr.Values = append(heightAttr.Values, val3)

	BirdType := EntityType{
		Name: "Bird",
	}
	BirdType.Attributs = append(
		BirdType.Attributs, colorAttr, specieAttr, heightAttr,
	)
	mesange := Entity{Identifier: "bird1"}
	mesange.Fields = append(mesange.Fields, val1, val2, val3)
	BirdType.Instances = append(BirdType.Instances, mesange)
	db.Create(&BirdType)

	return nil
}
