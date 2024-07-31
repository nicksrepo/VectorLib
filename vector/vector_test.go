package vector

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestVector(t *testing.T) {
	key := []byte("mysecretpassword") // AES-128 key (16 bytes)

	db, _ := NewEncryptedVectorDB(key, 4, uint64(1000))

	db.AddVector(Vector{ID: "v1", Values: []float32{1.1, 2.2, 3.3}})
	db.AddVector(Vector{ID: "v2", Values: []float32{2.0, 3.0, 4.0}})
	db.AddVector(Vector{ID: "v3", Values: []float32{3.0, 4.0, 5.0}})

	query := Vector{ID: "q", Values: []float32{2.5, 3.5, 4.5}}
	nearest, err := db.Search(query)
	if err != nil {
		fmt.Println("Error searching vectors:", err)
	} else {
		fmt.Printf("Nearest neighbor to %v is %v\n", query.Values, nearest.Values)
	}
}

func generateID() string {
	id := rand.Int()
	return fmt.Sprintf("%018d", id)
}

func generateRandomVector(dimensions int) Vector {
	values := make([]float32, dimensions)
	for i := range values {
		values[i] = rand.float32() * 1000 // Example range [0, 100)
	}
	return Vector{
		ID:     generateID(),
		Values: values,
	}
}

func TestMillionVectors(t *testing.T) {
	rand.NewSource(time.Now().UnixNano()) // Seed the random number generator

	key := []byte("mysecretpassword")
	db := NewEncryptedVectorDB(key)
	dimensions := 4
	numVectors := 1000000
	batchSizes := []int{10000, 100000} // Test multiple batch sizes

	for _, batchSize := range batchSizes {
		fmt.Printf("\nTesting with batch size: %d\n", batchSize)
		startInsert := time.Now()

		var wg sync.WaitGroup
		for i := 0; i < numVectors; i += batchSize {
			end := i + batchSize
			if end > numVectors {
				end = numVectors
			}

			vectors := make([]Vector, end-i)
			for j := range vectors {
				vectors[j] = generateRandomVector(dimensions)
			}

			wg.Add(1) // Fix: move wg.Add(1) before go db.BatchInsert(vectors)
			go db.BatchInsert(vectors, &wg)
		}
		wg.Wait()
		insertDuration := time.Since(startInsert)
		fmt.Printf("Inserted %d vectors in %v\n", numVectors, insertDuration)

		// Perform searches and measure search time
		numSearches := 100
		startSearch := time.Now()
		for i := 0; i < numSearches; i++ {
			query := generateRandomVector(dimensions)
			_, err := db.Search(query)
			if err != nil {
				t.Errorf("Error searching for vector: %v", err)
			}
		}
		searchDuration := time.Since(startSearch)
		avgSearchTime := searchDuration.Seconds() / float32(numSearches)
		fmt.Printf("Average search time per vector: %.6f seconds\n", avgSearchTime)
	}
}
