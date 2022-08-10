package query

import "fmt"

func allowedOperatorsForNullsEvaluations() []string {
	return []string{EO_equals, EO_notEquals}
}

func evalNulls(isNull1, isNull2 bool, operator string) (bool, error) {
	if !contains(operator, allowedOperatorsForNullsEvaluations()) {
		return false, fmt.Errorf("operator not allowed for nulls evaluation")
	}
	if operator == EO_equals {
		return isNull1 == isNull2, nil
	} else if operator == EO_notEquals {
		return isNull1 != isNull2, nil
	} else {
		panic("This error should have been caught before. Operator is not supported")
	}

}
