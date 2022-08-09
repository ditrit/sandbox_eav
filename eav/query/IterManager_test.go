package query_test

import (
	"testing"

	"github.com/ditrit/sandbox_eav/eav/query"
)

func _instanciateDatasetWithIntegers() [][]int {
	l1 := []int{
		10001,
		10002,
		10003,
		10004,
		10005,
		10006,
		10007,
		10008,
		10009,
		10010,
		10011,
	}
	l2 := []int{
		20001,
		20002,
		20003,
		20004,
		20005,
		20006,
		20007,
		20008,
		20009,
	}

	l3 := []int{
		30001,
		30002,
		30003,
		30004,
		30005,
		30006,
		30007,
	}

	return [][]int{l1, l2, l3}
}
func TestNewIterManager(t *testing.T) {
	dataSet := _instanciateDatasetWithIntegers()
	rm := query.NewIterManager(dataSet)
	if rm == nil {
		t.Error("the newly created IterManager should not be nil")
	}
}

func TestNext(t *testing.T) {
	dataSet := _instanciateDatasetWithIntegers()
	rm := query.NewIterManager(dataSet)
	var cnt int = 0
	var limit int = 100 * 50 * 25 // +1 because the "next" function will run one last time to exit the loop
	for {
		cnt++
		ended := rm.Next()
		if ended {
			if cnt != limit {
				t.Errorf("Next should only count %v times,got %v", limit, cnt)
				return
			} else {
				return
			}
		}

	}
}

func TestGetElem(t *testing.T) {
	testCases := []struct {
		desc        string
		indexList   int
		indexElem   int
		expectedVal int
	}{
		{
			desc:        "Test on list 1",
			indexList:   0,
			indexElem:   1,
			expectedVal: 10002,
		},
		{
			desc:        "Test on list 2",
			indexList:   1,
			indexElem:   6,
			expectedVal: 20007,
		},
	}
	for _, tC := range testCases {
		dataSet := _instanciateDatasetWithIntegers()
		rm := query.NewIterManager(dataSet)
		t.Run(tC.desc, func(t *testing.T) {
			if res := rm.GetElem(tC.indexList, tC.indexElem); res != tC.expectedVal {
				t.Errorf("got %v, expected %v", res, tC.expectedVal)
			}
		})
	}
}

type Record struct {
	K string
	V interface{}
}

func TestGetCurrentRecords(t *testing.T) {
	testRecord11 := &Record{K: "10001", V: 10001}
	testRecord12 := &Record{K: "10002", V: 10002}
	testRecord13 := &Record{K: "10003", V: 10003}

	testRecord21 := &Record{K: "20001", V: 20001}
	testRecord22 := &Record{K: "20002", V: 20002}
	testRecord23 := &Record{K: "20003", V: 20003}

	testRecord31 := &Record{K: "30001", V: 30001}
	testRecord32 := &Record{K: "30002", V: 30002}
	l1 := []*Record{
		testRecord11,
		testRecord12,
		testRecord13,
	}

	l2 := []*Record{
		testRecord21,
		testRecord22,
		testRecord23,
	}

	l3 := []*Record{
		testRecord31,
		testRecord32,
	}
	dataSet := [][]*Record{l1, l2, l3}
	rm := query.NewIterManager(dataSet)
	tableExpectedResults := [][]*Record{
		{testRecord11, testRecord21, testRecord31},
		{testRecord11, testRecord21, testRecord32},
		{testRecord11, testRecord22, testRecord31},
		{testRecord11, testRecord22, testRecord32},
		{testRecord11, testRecord23, testRecord31},
		{testRecord11, testRecord23, testRecord32},
		{testRecord12, testRecord21, testRecord31},
		{testRecord12, testRecord21, testRecord32},
		{testRecord12, testRecord22, testRecord31},
		{testRecord12, testRecord22, testRecord32},
		{testRecord12, testRecord23, testRecord31},
		{testRecord12, testRecord23, testRecord32},
		{testRecord13, testRecord21, testRecord31},
		{testRecord13, testRecord21, testRecord32},
		{testRecord13, testRecord22, testRecord31},
		{testRecord13, testRecord22, testRecord32},
		{testRecord13, testRecord23, testRecord31},
		{testRecord13, testRecord23, testRecord32},
	}
	cnt := 0
	for {
		records := rm.GetSelectedElems()
		if !checkSliceIdentical(records, tableExpectedResults[cnt]) {
			t.Errorf("didn't get the expected Records '%v', got'%v'", tableExpectedResults[cnt], records)
			return
		}
		cnt++
		if rm.Next() {
			break
		}
	}
}

// helper
func checkSliceIdentical[T comparable](a1, a2 []T) bool {
	// check size
	if len(a1) != len(a2) {
		return false
	}
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}
