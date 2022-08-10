package query

import "fmt"

// Integer is made up of all the int and float types
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~uint | ~uint8 | ~uint16 | ~uint32
}

func evalNumbers[T Number](ref, value T, operator string) (bool, error) {
	if !contains(operator, getNumericEvaluationOperators()) {
		return false, fmt.Errorf("operator not allowed for number evaluation")
	}
	if operator == EO_equals {
		return ref == value, nil
	} else if operator == EO_notEquals {
		return ref != value, nil
	} else if operator == EO_inferior {
		return ref < value, nil
	} else if operator == EO_inferiorOrEquals {
		return ref <= value, nil
	} else if operator == EO_superior {
		return ref > value, nil
	} else if operator == EO_superiorOrEquals {
		return ref >= value, nil
	} else {
		panic("This error should have been caught before. Operator is not supported")
	}
}
