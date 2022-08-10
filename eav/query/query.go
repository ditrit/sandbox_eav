package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ditrit/sandbox_eav/eav/models"
	"github.com/ditrit/sandbox_eav/eav/operations"
	"gorm.io/gorm"
)

var (
	debugLogger = log.New(os.Stdout, "QUERY: ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
)

// Represent a "sql" like query
// run the query with .Run(db *gorm.DB)
type Query struct {
	// The attributs that are selected to be returned
	Attrs []string `json:"attrs"` // ["bird.color", "bird.weight"

	// The tables we will run the query on
	Tables []string `json:"tables"` // ["bird", "human"]

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
	// used when returning an error
	var emptyResponse []byte = []byte("[]")

	if len(q.Tables) == 0 {
		return emptyResponse, fmt.Errorf("table names are needed to query the database")
	}
	// Retrieve EntityTypes.
	etts, err := getEntityTypes(db, q.Tables...)
	if err != nil {
		return emptyResponse, err
	}
	// Retrieve Entities from with the EntityType aforementioned
	var data = make([][]*models.Entity, len(etts))
	for i, ett := range etts {
		entities, err := operations.GetEntities(db, ett)
		if err != nil {
			return emptyResponse, err
		}
		if len(entities) == 0 {
			// if there is no entities to filter then we return an empty response
			return []byte("[]"), nil
		}
		data[i] = entities
	}
	var b strings.Builder
	b.WriteString("[")
	// Make an IterManager holding the entities
	entityManager := NewIterManager(data)
	var selectedResults []string
	for {
		// Get selected entities
		selectedEntities := entityManager.GetSelectedElems()

		// get Records from the selected Entities
		rcs := buildRecordSliceFromEntities(selectedEntities)
		fmt.Println("Records", rcs)
		// pass thought the Condition system
		r, err := q.Condition.Eval(rcs)
		if err != nil {
			return emptyResponse, err
		}
		// if the condition tree is validated, then add to the returnResultSet
		if r {
			resultFields, err := getResultFields(rcs, q.Attrs)
			if err != nil {
				return emptyResponse, err
			}
			byt, err := json.Marshal(resultFields)
			if err != nil {
				return emptyResponse, fmt.Errorf("error while marshalling the response data")
			}

			selectedResults = append(selectedResults, string(byt))
		}

		if entityManager.Next() {
			break
		}
	}
	b.WriteString(strings.Join(selectedResults, ","))
	b.WriteString("]")
	return []byte(b.String()), nil
}

func getResultFields(rcs RecordMap, attrs []string) (map[string]any, error) {
	resurnSet := make(map[string]any, len(attrs))
	for _, att := range attrs {
		val, ok := rcs[att]
		if !ok {
			return resurnSet, fmt.Errorf("the attr %q does not exist", att)
		}
		resurnSet[att] = val.Value()
	}
	return resurnSet, nil
}

type RecordMap map[string]*models.Value

func buildRecordSliceFromEntities(ets []*models.Entity) RecordMap {
	var rcs RecordMap = RecordMap{}
	for _, et := range ets {
		rcs[et.EntityType.Name+".id"] = &models.Value{
			IntVal:   int(et.ID),
			IsNull:   false,
			EntityId: et.ID,
			Attribut: &models.Attribut{
				Name:      "id",
				ValueType: models.IntValueType,
			},
		}
		for _, f := range et.Fields {
			key := fmt.Sprintf("%s.%s", et.EntityType.Name, f.Attribut.Name)
			rcs[key] = f
		}
	}
	return rcs
}

// return the EntityTypes that matches the names
func getEntityTypes(db *gorm.DB, names ...string) ([]*models.EntityType, error) {
	var etts []*models.EntityType
	for _, name := range names {
		ett, err := operations.GetEntityTypeByName(db, name)
		if err != nil {
			return nil, err
		}
		etts = append(etts, ett)
	}
	return etts, nil
}

func (ec *EvaluationComposite) Eval(rcs RecordMap) (bool, error) {
	ec.Operator = strings.TrimSpace(ec.Operator)
	if ec.Operator == "" {
		if ec.Composites != nil {
			return false, ErrOperatorEmptyWithNoConditions
		}
		// we are suposedly on a terminal EvaluationComposite
		// ec.Composites is suposedly null
		fmt.Println("We now will evaluate the following", ec.Evaluation)
		r, err := ec.Evaluation.Eval(rcs)
		fmt.Println("RESULT: ", r)
		return r, err
	}
	if !contains(ec.Operator, getBooleanOperators()) {
		return false, ErrBooleanOperatorNotImplemented
	}

	// EVALUATION for AND OPERATOR
	if ec.Operator == BO_and {
		return andEval(ec.Composites, rcs)
	}
	// EVALUATION for OR OPERATOR
	return orEval(ec.Composites, rcs)
}

// Evaluate the Conditions and apply an AND operator on the result
func andEval(evcs []EvaluationComposite, rcs RecordMap) (bool, error) {
	for _, ev := range evcs {
		res, err := ev.Eval(rcs)
		if err != nil {
			return false, err
		}
		if !res {
			return false, nil
		}
	}
	return true, nil
}

// Evaluate the Conditions and apply an OR operator on the result
func orEval(evcs []EvaluationComposite, rcs RecordMap) (bool, error) {
	for _, ev := range evcs {
		res, err := ev.Eval(rcs)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

func (ev *Evaluation) Eval(rcs RecordMap) (result bool, err error) {
	ev.Operator = strings.TrimSpace(ev.Operator)
	if ev.Operator == "" {
		return false, ErrOperatorEmpty
	}
	if !contains(ev.Operator, getEvaluationOperators()) {
		return false, ErrEvalutationOperatorNotImplemented
	}

	if ev.Expre1.Type == "ref" && ev.Expre2.Type == "value" {
		tableDotAttr, ok := ev.Expre1.Value.(string)
		if !ok {
			return false, fmt.Errorf("can't cast comparation.expre1.value to a string. (got=%v)", ev.Expre1.Value)
		}
		val, ok := rcs[tableDotAttr]
		if !ok {
			return false, fmt.Errorf("%q not found", ev.Expre1.Value)
		}
		result, err = evalValueVsRef(val, ev.Expre2.Value, ev.Operator)
		if err != nil {
			return false, err
		}

	} else if ev.Expre1.Type == "value" && ev.Expre2.Type == "ref" {
		tableDotAttr, ok := ev.Expre2.Value.(string)
		if !ok {
			return false, fmt.Errorf("can't cast comparation.expre2.value to a string. (got=%v)", ev.Expre2.Type)
		}
		val, ok := rcs[tableDotAttr]
		if !ok {
			return false, fmt.Errorf("%q not found", ev.Expre2.Value)
		}
		result, err = evalValueVsRef(val, ev.Expre1.Value, ev.Operator)
		if err != nil {
			return false, err
		}

	} else if ev.Expre1.Type == "ref" && ev.Expre2.Type == "ref" {
		expre1ValueString, ok := ev.Expre1.Value.(string)
		if !ok {
			return false, fmt.Errorf("can't cast comparation.expre1.value to a string. (got=%v)", ev.Expre1.Value)
		}
		expre2ValueString, ok := ev.Expre2.Value.(string)
		if !ok {
			return false, fmt.Errorf("can't cast comparation.expre2.value to a string. (got=%v)", ev.Expre1.Value)
		}
		val1, ok := rcs[expre1ValueString]
		if !ok {
			return false, fmt.Errorf("%q not found", ev.Expre2.Value)
		}
		val2, ok := rcs[expre2ValueString]
		if !ok {
			return false, fmt.Errorf("%q not found", ev.Expre2.Value)
		}

		result, err = evalRefVsRef(val1, val2, ev.Operator)
		if err != nil {
			return false, err
		}

	} else if ev.Expre1.Type == "value" || ev.Expre2.Type == "value" {
		return false, ErrCantEvaluateValueVsValue
	} else {
		return false, fmt.Errorf("expression type does not exist. got %q and %q", ev.Expre1.Type, ev.Expre2.Type)
	}

	return result, nil
}

func jsontruc(v any) string {
	r, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(r)
}

func evalRefVsRef(refV1 *models.Value, refV2 *models.Value, operator string) (bool, error) {
	debugLogger.Printf("evalRefVsRef: ref %v and ref %v", jsontruc(refV1), jsontruc(refV2))

	val1 := refV1.Value()
	val2 := refV1.Value()
	if val1 == nil || val2 == nil {
		return false, nil
	}
	if refV1.Attribut.ValueType == refV2.Attribut.ValueType {
		// Same type
		switch refV1.Attribut.ValueType {
		case models.StringValueType:
			return evalStrings(val1.(string), val2.(string), operator)
		case models.BooleanValueType:
			return evalBooleans(val1.(bool), val2.(bool), operator)
		case models.IntValueType:
			return evalNumbers(val1.(int), val2.(int), operator)
		case models.FloatValueType:
			return evalNumbers(val1.(float64), val2.(float64), operator)
		case models.RelationValueType:
			return evalNumbers(val1.(uint), val2.(uint), operator) // FIXME: get only equality OP
		default:
			panic("mmh should not be there")
		}
	} else if (refV1.Attribut.ValueType == models.FloatValueType || refV1.Attribut.ValueType == models.IntValueType || refV1.Attribut.ValueType == models.RelationValueType) && (refV2.Attribut.ValueType == models.FloatValueType || refV2.Attribut.ValueType == models.IntValueType || refV2.Attribut.ValueType == models.RelationValueType) {

		float1, err := refToNumVal(refV1)
		if err != nil {
			return false, err
		}
		float2, err := refToNumVal(refV2)
		if err != nil {
			return false, err
		}

		return evalNumbers(float1, float2, operator)
	}

	return false, nil
}

func refToNumVal(ref *models.Value) (floatVal float64, err error) {
	switch ref.Attribut.ValueType {
	case models.FloatValueType:
		floatVal, err = ref.GetFloatVal()
		if err != nil {
			if errors.Is(err, models.ErrValueIsNull) {
				return 0.0, ErrCantEvaluateNullRefAgainstAnythingNotNull
			}
			panic(err)
		}
	case models.IntValueType:
		int1, err := ref.GetIntVal()
		floatVal = float64(int1)
		if err != nil {
			if errors.Is(err, models.ErrValueIsNull) {
				return 0.0, ErrCantEvaluateNullRefAgainstAnythingNotNull
			}
			panic(err)
		}
	case models.RelationValueType:
		if ref.IsNull {
			return 0.0, fmt.Errorf("can't evaluate a null reference to anything")
		}
		floatVal = float64(ref.RelationVal) // FIXME: bruh this is a terrible implementation be it will do the job for a poc.
	default:
		panic("we are not supposed to end up here")
	}
	return floatVal, nil
}

// Evaluate a models.Value against a value.
func evalValueVsRef(refValue *models.Value, value interface{}, operator string) (bool, error) {
	debugLogger.Printf("evalValueVsRef: ref %v and value %v", refValue, value)
	switch value := value.(type) {
	case string:
		val := refValue.Value()
		if val == nil {
			return false, nil
		}
		return evalStrings(val.(string), value, operator)
	case float64:
		val := refValue.Value()
		if val == nil {
			return false, nil
		}
		return evalNumbers(val.(float64), value, operator)

	case bool:
		val := refValue.Value()
		if val == nil {
			return false, nil
		}
		return evalBooleans(val.(bool), value, operator)
	case nil:
		return evalNulls(refValue.IsNull, true, operator)

	default:
		return false, fmt.Errorf("this json type (%T) is not available with this server, please use one type that is supported by golang (https://go.dev/blog/json#generic-json-with-interface)", value)

	}
}

func contains[T comparable](op T, ops []T) bool {
	for _, v := range ops {
		if v == op {
			return true
		}
	}
	return false
}
