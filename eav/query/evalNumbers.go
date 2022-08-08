package query

import "fmt"

func evalNumbers(ref, value float64, operator string) (bool, error) {
	if !ContainsOperator(operator, getNumericEvaluationOperators()) {
		return false, fmt.Errorf("operator not allowed for float evaluation")
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
