package vector

import (
	"math"
	"sort"
)

type Vector struct {
	Root   *KDNode
	ID     string
	Values []float32
}

type VectorDB struct {
	Vectors []Vector
	index   Index
}

func NewVectorDB() *VectorDB {
	return &VectorDB{Vectors: []Vector{}}
}

func cosineSimilarity(v1, v2 []float32) float32 {
	dotProduct := 0.0
	normA := 0.0
	normB := 0.0
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		normA += v1[i] * v1[i]
		normB += v2[i] * v2[i]
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (db *VectorDB) Search(query []float32, topK int) []Vector {
	type Result struct {
		Vector     Vector
		Similarity float32
	}

	results := []Result{}
	for _, vector := range db.Vectors {
		similarity := cosineSimilarity(query, vector.Values)
		results = append(results, Result{Vector: vector, Similarity: similarity})
	}

	// Sort results by similarity in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Return the topK results
	topVectors := []Vector{}
	for i := 0; i < topK && i < len(results); i++ {
		topVectors = append(topVectors, results[i].Vector)
	}

	return topVectors
}
