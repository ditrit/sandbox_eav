package query

const (
	// Boolean operator AND
	BO_and string = "&&"

	// Boolean operator OR
	BO_or string = "||"

	// Evalution operator equals to
	EO_equals string = "=="

	// Evalution operator "not equals to"
	EO_notEquals string = "!="

	// Evalution operator superior or equals to
	EO_superiorOrEquals string = ">="

	// Evalution operator inferior or equals to
	EO_inferiorOrEquals string = "<="

	// Evalution operator superior to
	EO_superior string = ">"

	// Evalution operator inferior to
	EO_inferior string = "<"
)

// Get a list of the Numeric operators
func getNumericEvaluationOperators() []string {
	return []string{EO_equals, EO_inferior, EO_inferiorOrEquals, EO_notEquals, EO_superior, EO_superiorOrEquals}
}

// Get a list of the Evaluation operators
func getEvaluationOperators() []string {
	return []string{EO_equals, EO_inferior, EO_inferiorOrEquals, EO_notEquals, EO_superior, EO_superiorOrEquals, BO_and, BO_or}
}

// Get a list of the Booleans operators
func getBooleanOperators() []string {
	return []string{BO_and, BO_or}
}

// Get a list of all implemented operators
func getImplementedOperators() []string {
	return []string{EO_equals, EO_inferior, EO_inferiorOrEquals, EO_notEquals, EO_superior, EO_superiorOrEquals, BO_and, BO_or}
}
