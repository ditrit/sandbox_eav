package query

import "fmt"

func allowedOperatorsForStringEvaluations() []string {
	return []string{EO_equals, EO_notEquals}
}

func evalStrings(ref, value, operator string) (bool, error) {
	if !contains(operator, allowedOperatorsForStringEvaluations()) {
		return false, fmt.Errorf("operator not allowed for strings evaluation")
	}
	if operator == EO_equals {
		return ref == value, nil
	} else if operator == EO_notEquals {
		return ref != value, nil
	} else {
		panic("This error should have been caught before. Operator is not supported")
	}
}
