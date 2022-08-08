package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ditrit/sandbox_eav/eav/models"
	"github.com/ditrit/sandbox_eav/eav/operations"
	"github.com/ditrit/sandbox_eav/utils"
	"gorm.io/gorm"
)

// Represent a "sql" like query
// run the query with .Run(db *gorm.DB)
type Query struct {
	// The attributs that are selected to be returned
	Attrs []string `json:"attrs"` // ["bird.color", "bird.weight"

	// The table we will run the query on
	Table string `json:"table"` // "bird"

	// The condition we will evaluate (similar to "WHERE ..." in SQL)
	Condition EvaluationComposite `json:"condition"`
}

// A node in the condition tree
type EvaluationComposite struct {
	// an operator, can only be a Boolean operator (please check constants.go)
	Operator   string                `json:"operator"`
	Evaluation Evaluation            `json:"comparaison"`
	Composites []EvaluationComposite `json:"conditions"`
}

type Evaluation struct {
	// an operator, can be any operator as long as it's applicable in the considered case
	Operator string `json:"operator"`
	// The expressions to evaluate
	Expre1 Expression `json:"expre1"`
	Expre2 Expression `json:"expre2"`
}

type Expression struct {
	Type  string      `json:"type"`  // Either "ref" or "value"
	Value interface{} `json:"value"` // Either a ref (ex: "bird.color") or a value of type (float, int, string, null)
}

func (q *Query) Run(db *gorm.DB) ([]byte, error) {
	var b strings.Builder

	// used when returning an error
	var emptyBytesSlice []byte = []byte("")

	ett, err := getTable(db, q.Table)
	if err != nil {
		return emptyBytesSlice, err
	}

	entities, err := operations.GetEntities(db, ett)
	if err != nil {
		return emptyBytesSlice, err
	}
	var elems []string
	for _, et := range entities {
		r, err := q.Condition.Eval(et)
		if err != nil {
			return emptyBytesSlice, err
		}
		if r {
			var fieldsToOutput []string
			for _, atName := range q.Attrs {
				v, err := getValue(et, atName)
				if err != nil {
					return emptyBytesSlice, err
				}
				pair, err := v.BuildJsonKVPair()
				if err != nil {
					return emptyBytesSlice, err
				}

				fieldsToOutput = append(fieldsToOutput, pair)
			}
			elems = append(elems, utils.BuildJsonFromStrings(fieldsToOutput))
		}
	}
	b.WriteString(utils.BuildJsonListFromStrings(elems))
	return []byte(b.String()), nil
}

// return the EntityType that matches the name
func getTable(db *gorm.DB, name string) (*models.EntityType, error) {

	ett, err := operations.GetEntityTypeByName(db, name)
	if err != nil {
		return nil, err
	}
	return ett, nil
}

func (ec *EvaluationComposite) Eval(et *models.Entity) (bool, error) {
	ec.Operator = strings.TrimSpace(ec.Operator)
	if ec.Operator == "" {
		if ec.Composites == nil {
			// we are suposedly on a terminal EvaluationComposite
			// ec.Composites is suposedly nulll
			fmt.Println("We now will evaluate the following", ec.Evaluation)
			return ec.Evaluation.Eval(et)
		}
		return false, ErrOperatorEmptyWithNoConditions

	}
	if !ContainsOperator(ec.Operator, getBooleanOperators()) {
		return false, ErrBooleanOperatorNotImplemented
	}

	// EVALUATION for AND OPERATOR
	if ec.Operator == BO_and {
		var results []bool
		for _, ev := range ec.Composites {
			res, err := ev.Eval(et)
			if err != nil {
				return false, err
			}
			results = append(results, res)
		}
		return allTrue(results), nil
	}
	// EVALUATION for OR OPERATOR
	for _, ev := range ec.Composites {
		res, err := ev.Eval(et)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

func (ev *Evaluation) Eval(et *models.Entity) (bool, error) {
	ev.Operator = strings.TrimSpace(ev.Operator)
	if ev.Operator == "" {
		return false, ErrOperatorEmpty
	}
	if !ContainsOperator(ev.Operator, getEvaluationOperators()) {
		return false, ErrEvalutationOperatorNotImplemented
	}
	var result bool
	if ev.Expre1.Type == "ref" && ev.Expre2.Type == "value" {
		value, err := getValue(et, ev.Expre1.Value.(string))
		if err != nil {
			return false, err
		}
		result, err = evalValueVsRef(value, ev.Expre2.Value, ev.Operator)
		if err != nil {
			return false, err
		}
	} else if ev.Expre1.Type == "value" && ev.Expre2.Type == "ref" {
		value, err := getValue(et, ev.Expre2.Value.(string))
		if err != nil {
			return false, err
		}
		result, err = evalValueVsRef(value, ev.Expre1.Value, ev.Operator)
		if err != nil {
			return false, err
		}
	} else if ev.Expre1.Type == "ref" && ev.Expre2.Type == "ref" {
		value1, err := getValue(et, ev.Expre1.Value.(string))
		if err != nil {
			return false, err
		}
		value2, err := getValue(et, ev.Expre2.Value.(string))
		if err != nil {
			return false, err
		}
		result, err = evalRefVsRef(value1, value2, ev.Operator)
		if err != nil {
			return false, err
		}
	} else if ev.Expre1.Type == "value" || ev.Expre2.Type == "value" {
		return false, ErrCantEvaluateValueVsValue
	} else {
		panic(
			fmt.Errorf("expression type does not exist. got %q and %q", ev.Expre1.Type, ev.Expre2.Type),
		)
	}

	return result, nil
}

func evalRefVsRef(refV1 *models.Value, refV2 *models.Value, operator string) (bool, error) {
	if refV1.Attribut.ValueType == refV2.Attribut.ValueType {
		// Same type
		switch refV1.Attribut.ValueType {
		case models.StringValueType:
			strVal1, err := refV1.GetStringVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			strVal2, err := refV2.GetStringVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			return evalStrings(strVal1, strVal2, operator)

		case models.BooleanValueType:
			boolVal1, err := refV1.GetBoolVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			boolVal2, err := refV2.GetBoolVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			return evalBooleans(boolVal1, boolVal2, operator)
		case models.IntValueType:
			intVal1, err := refV1.GetIntVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			intVal2, err := refV2.GetIntVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			return evalNumbers(float64(intVal1), float64(intVal2), operator)
		case models.FloatValueType:
			floatVal1, err := refV1.GetFloatVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			floatVal2, err := refV2.GetFloatVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
			return evalNumbers(floatVal1, floatVal2, operator)

		case models.RelationValueType:
			panic("NOT IMPLEMENTED")

		default:
			panic("mmh should not be there")
		}
	} else if (refV1.Attribut.ValueType == models.FloatValueType && refV2.Attribut.ValueType == models.IntValueType) || (refV1.Attribut.ValueType == models.IntValueType && refV2.Attribut.ValueType == models.FloatValueType) {
		// number comparaison
		var float1 float64
		var float2 float64
		var err error
		if refV1.Attribut.ValueType == models.FloatValueType {
			float1, err = refV1.GetFloatVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
		} else if refV1.Attribut.ValueType == models.IntValueType {
			int1, err := refV1.GetIntVal()
			float1 = float64(int1)
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
		}

		if refV2.Attribut.ValueType == models.FloatValueType {
			float1, err = refV2.GetFloatVal()
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
		} else if refV2.Attribut.ValueType == models.IntValueType {
			int2, err := refV2.GetFloatVal()
			float2 = float64(int2)
			if err != nil {
				if errors.Is(err, models.ErrValueIsNull) {
					return false, ErrCantEvaluateNullRefAgainstAnythingNotNull
				}
				panic(err)
			}
		}
		return evalNumbers(float1, float2, operator)
	}

	return false, nil
}

// Evaluate a models.Value against a value.
func evalValueVsRef(refValue *models.Value, value interface{}, operator string) (bool, error) {
	switch value := value.(type) {
	case string:
		strVal, err := refValue.GetStringVal()
		if err != nil {
			panic(err)
		}
		return evalStrings(strVal, value, operator)

	case float64:
		var floatVal float64
		var err error
		if refValue.Attribut.ValueType == models.IntValueType {
			intVal, err := refValue.GetIntVal()
			if err != nil {
				panic(err)
			}
			floatVal = float64(intVal)
		} else {
			floatVal, err = refValue.GetFloatVal()
			if err != nil {
				panic(err)
			}
		}
		return evalNumbers(floatVal, value, operator)

	case bool:
		boolVal, err := refValue.GetBoolVal()
		if err != nil {
			panic(err)
		}
		return evalBooleans(boolVal, value, operator)
	case nil:
		return evalNulls(refValue.IsNull, true, operator)

	default:
		return false, fmt.Errorf("this json type (%T) is not available with this server, please use one type that is supported by golang (https://go.dev/blog/json#generic-json-with-interface)", value)

	}
}

func ContainsOperator(op string, ops []string) bool {
	for _, v := range ops {
		if v == op {
			return true
		}
	}
	return false
}

func allTrue(lis []bool) bool {
	for _, b := range lis {
		if !b {
			return false
		}
	}
	return true
}

func getValue(et *models.Entity, ref string) (*models.Value, error) {
	parts := strings.SplitN(ref, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("ref value is not a valid string: got=%s, wanted something like table.valuename", ref)
	}
	var attrId uint = 0
	for _, a := range et.EntityType.Attributs {
		if a.Name == parts[1] {
			attrId = a.ID
			break
		}
	}
	if attrId == 0 {
		return nil, fmt.Errorf("attr not found: got=%s", parts[1])
	}
	for _, v := range et.Fields {
		if v.AttributId == attrId {
			return v, nil
		}
	}
	return nil, fmt.Errorf("value not found: got=%s", parts[1])

}
