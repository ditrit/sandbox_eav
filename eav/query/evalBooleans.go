package query

import "fmt"

func allowedOperatorsForBooleanEvaluations() []string {
	return []string{EO_equals, EO_notEquals, BO_and, BO_or}
}

func evalBooleans(ref bool, value bool, operator string) (bool, error) {
	if !ContainsOperator(operator, allowedOperatorsForBooleanEvaluations()) {
		return false, fmt.Errorf("operator not allowed for boolean evaluation")
	}
	if operator == EO_equals {
		return ref == value, nil
	} else if operator == EO_notEquals {
		return ref != value, nil
	} else if operator == BO_and {
		return ref && value, nil
	} else if operator == BO_or {
		return ref || value, nil
	} else {
		panic("This error should have been caught before. Operator is not supported")
	}
}
