package query

type IterManager[T any] struct {
	Indexes []int
	Data    [][]T
}

func NewIterManager[T any](data [][]T) *IterManager[T] {
	return &IterManager[T]{
		Indexes: make([]int, len(data)),
		Data:    data,
	}
}

func (it *IterManager[T]) Next() bool {
	for indexList := len(it.Data) - 1; indexList >= 0; indexList-- {

		it.Indexes[indexList]++
		if it.Indexes[0] == len(it.Data[0]) {
			return true
		}
		if it.Indexes[indexList] == len(it.Data[indexList]) {
			// si on arrive à la fin d'une des listes, alors on repars au début et on incrémente l'index de la liste précédente
			it.Indexes[indexList] = 0
			// fmt.Printf("DEBUG ITER MANAGER: currently selected=%v (indexes=%v)\n", indexList, rm.Indexes)
		} else {
			break
		}
	}
	return false
}

func (it *IterManager[T]) GetElem(tableIndex, elementIndex int) T {
	return it.Data[tableIndex][elementIndex]
}

// Return the current Records
func (it *IterManager[T]) GetSelectedElems() (returnSet []T) {
	for tableIndex, elementIndex := range it.Indexes {
		returnSet = append(returnSet, it.GetElem(tableIndex, elementIndex))
	}
	return returnSet
}
