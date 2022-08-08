package query

import (
	"errors"
	"fmt"
)

var (
	ErrOperatorEmpty                 = errors.New("operator shouldn't be empty")
	ErrOperatorNotImplemented        = fmt.Errorf("operator not implemented, please choose one from this list %v", getImplementedOperators())
	ErrBooleanOperatorNotImplemented = fmt.Errorf("boolean operator not implemented, please choose one from this list %v", getBooleanOperators())

	ErrEvalutationOperatorNotImplemented = fmt.Errorf("conditional operator not implemented, please choose one from this list %v",
		getEvaluationOperators())
	ErrOperatorEmptyWithNoConditions = errors.New("operator empty while not providing comparaison")

	ErrCantEvaluateValueVsValue                  = errors.New("can't evaluate a value against another value. (ref vs value and ref vs ref are accepted)")
	ErrCantEvaluateNullRefAgainstAnythingNotNull = errors.New("can't evaluate a null ref against another value/ref that is not null")
)
