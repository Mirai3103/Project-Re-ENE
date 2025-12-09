package store

import (
	"context"
	"sort"
)

type MemoryWithScore struct {
	Memory
	Score float64
}

func (q *Queries) SimilarMemories(
	ctx context.Context,
	targetEmbedding []float32,
	limit int,
	threshold float64,
) ([]MemoryWithScore, error) {
	if limit <= 0 {
		return []MemoryWithScore{}, nil
	}

	embeddings, err := q.GetAllEmbeddings(ctx)
	if err != nil {
		return nil, err
	}

	type scored struct {
		id    string
		score float64
	}

	top := make([]scored, 0, limit)

	findWorst := func() int {
		worst := 0
		for i := 1; i < len(top); i++ {
			if top[i].score > top[worst].score {
				worst = i
			}
		}
		return worst
	}

	for _, e := range embeddings {
		vec := BytesToFloat32(e.Embedding)
		d := L2DistanceF32(vec, targetEmbedding)

		if d >= threshold {
			continue
		}

		if len(top) < limit {
			top = append(top, scored{id: e.ID, score: d})
		} else {
			worst := findWorst()
			if d < top[worst].score {
				top[worst] = scored{id: e.ID, score: d}
			}
		}
	}

	if len(top) == 0 {
		return []MemoryWithScore{}, nil
	}

	sort.Slice(top, func(i, j int) bool {
		return top[i].score < top[j].score
	})

	ids := make([]string, len(top))
	for i, s := range top {
		ids[i] = s.id
	}

	dbMemories, err := q.GetMemoriesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	memByID := make(map[string]Memory, len(dbMemories))
	for _, m := range dbMemories {
		memByID[m.ID] = m
	}

	memories := make([]MemoryWithScore, len(top))
	for i, s := range top {
		if m, ok := memByID[s.id]; ok {
			memories[i] = MemoryWithScore{
				Memory: m,
				Score:  s.score,
			}
		}
	}

	return memories, nil
}
