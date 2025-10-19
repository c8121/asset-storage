package metadata_db

import (
	"fmt"
	"slices"
)

type ScoredId struct {
	Id    int64
	Score float32
}

type ScoredIdMap map[int64]float32

func (m ScoredIdMap) Add(id int64, score float32) {
	existingScore, ok := m[id]
	if ok {
		m[id] = score + existingScore
	} else {
		m[id] = score
	}
}

func (m ScoredIdMap) Reduce(remains ScoredIdMap) {
	fmt.Printf("Reduce %d by %d\n", len(m), len(remains))
	removeIds := make([]int64, len(m))
	for id, _ := range m {
		_, exists := remains[id]
		if !exists {
			removeIds = append(removeIds, id)
		}
	}
	for _, id := range removeIds {
		delete(m, id)
	}
	fmt.Printf("Reduced to %d\n", len(m))
}

func (m ScoredIdMap) Sort() []ScoredId {

	var items []ScoredId
	for k, v := range m {
		items = append(items, ScoredId{k, v})
	}
	slices.SortFunc(items, func(a, b ScoredId) int {
		return int((b.Score - a.Score) * 32000)
	})

	return items
}
