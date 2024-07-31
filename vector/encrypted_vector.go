package vector

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/oligo/hnswgo"
)

type EncryptedVector struct {
	ID            string
	EncryptedData string
}

type EncryptedVectorDB struct {
	hnswIndex *hnswgo.HnswIndex
	Index     map[uint64]*Vector // Use uint64 labels from HNSW
	Key       []byte
	mutex     sync.RWMutex
}

func NewEncryptedVectorDB(key []byte, dimensions int, maxElements uint64) (*EncryptedVectorDB, error) {
	db := &EncryptedVectorDB{
		Key:   key,
		Index: make(map[uint64]*Vector), // Initialize the map
	}

	db.hnswIndex = hnswgo.New(
		dimensions,  // Dimensions of the vectors
		32,          // M - max number of connections per node
		500,         // efConstruction - quality vs. speed tradeoff
		12345,       // Random seed
		maxElements, // Max number of elements (estimate)
		hnswgo.L2,   // Space type (L2 for Euclidean distance)
		false,       // Use heuristic to find best M (optional)
	)

	return db, nil
}

func (db *EncryptedVectorDB) BatchInsert(vectors []Vector, wg *sync.WaitGroup) error {
	defer wg.Done()

	vectorData := make([][]float32, len(vectors))
	for i, vector := range vectors {
		vectorData[i] = []float32(vector.Values)
	}

	// Add vectors to the HNSW index
	err := db.hnswIndex.AddPoints(vectorData, nil, 1, false)
	if err != nil {
		return fmt.Errorf("failed to add points to HNSW index: %v", err)
	}

	// Store vectors and their labels in the index map
	db.mutex.Lock()
	for i, label := range labels {
		db.Index[vectors[i].ID] = label
	}
	db.mutex.Unlock()

	return nil
}

func (db *EncryptedVectorDB) AddVector(vector Vector) error {
	// ... (same encryption logic as before) ...

	db.mutex.Lock()
	label, err := db.hnswIndex.AddPoint(vector.Values, uint64(len(db.Index))+1)
	if err != nil {
		db.mutex.Unlock()
		return fmt.Errorf("failed to add point to HNSW index: %v", err)
	}
	db.Index[label] = &vector
	db.mutex.Unlock()

	return nil
}

func (db *EncryptedVectorDB) GetVector(label uint64) (*Vector, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	vector, ok := db.Index[label]
	if !ok {
		return nil, fmt.Errorf("vector not found for label %d", label)
	}

	// ... (decryption logic if needed) ...

	return vector, nil
}

func (db *EncryptedVectorDB) BatchInsert(vectors []Vector, wg *sync.WaitGroup) error {
	defer wg.Done()

	vectorData := make([][]float32, len(vectors))
	for i, vector := range vectors {
		vectorData[i] = vector.Values
	}

	// Add vectors to the HNSW index
	labels, err := db.hnswIndex.AddPoints(vectorData, nil, 1, false)
	if err != nil {
		return fmt.Errorf("failed to add points to HNSW index: %v", err)
	}

	// Store vectors and their labels in the index map
	db.mutex.Lock()
	for i, label := range labels {
		db.Index[label] = &vectors[i]
	}
	db.mutex.Unlock()

	return nil
}

func (db *EncryptedVectorDB) Search(query [][]float32) (*Vector, error) {
	db.hnswIndex.SetEf(10) // Adjust ef for search (experiment for best results)

	// Search for the nearest neighbor
	result, err := db.hnswIndex.SearchKNN(query, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to search HNSW index: %v", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no neighbors found")
	}

	// Get the nearest neighbor's label and retrieve the vector
	nearestLabel := result[0][0]
	db.mutex.RLock()
	nearestVector, found := db.Index[nearestLabel]
	db.mutex.RUnlock()
	if !found {
		return nil, fmt.Errorf("vector not found for label %d", nearestLabel)
	}

	return nearestVector, nil
}

func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(ciphertextHex string, key []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func parallelEncrypt(vectors []Vector, key []byte, db *EncryptedVectorDB, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, vector := range vectors {
		plaintext, err := json.Marshal(vector.Values)
		if err != nil {
			fmt.Println("Error marshaling vector:", err)
			continue
		}

		data, err := encrypt(plaintext, key)
		if err != nil {
			fmt.Println("Error encrypting vector:", err)
			continue
		}

		db.mutex.Lock()
		db.Root = insertKDNode(db.Root, vector, data, 0)
		db.Index[vector.ID] = db.Root
		db.mutex.Unlock()
	}
}

func (db *EncryptedVectorDB) rebuildTree() {
	// Collect all vectors from the index
	vectors := make([]Vector, 0, len(db.Index))
	for _, node := range db.Index {
		vectors = append(vectors, node.Vector)
	}

	// Rebuild the tree from scratch
	db.Root = buildKDTree(vectors, 0, db.Key)

	// Update the index to point to the new nodes
	db.Index = make(Index)
	buildIndex(db.Root, db.Index)

	db.numInsertionsSinceRebuild = 0
}

// buildKDTree recursively builds a balanced k-d tree from the given vectors
func buildKDTree(vectors []Vector, depth int, key []byte) *KDNode {
	if len(vectors) == 0 {
		return nil
	}

	// Determine the splitting dimension
	dim := depth % len(vectors[0].Values)

	// Sort the vectors based on the splitting dimension
	sort.Slice(vectors, func(i, j int) bool {
		return vectors[i].Values[dim] < vectors[j].Values[dim]
	})

	// Find the median vector
	medianIndex := len(vectors) / 2
	medianVector := vectors[medianIndex]

	// Encrypt the median vector's values
	plaintext, _ := json.Marshal(medianVector.Values) // Error handling omitted for brevity
	encryptedData, _ := encrypt(plaintext, key)       // Error handling omitted for brevity

	// Create a new node
	newNode := &KDNode{
		Vector:        medianVector,
		EncryptedData: encryptedData,
		Left:          nil,
		Right:         nil,
	}

	// Recursively build left and right subtrees
	newNode.Left = buildKDTree(vectors[:medianIndex], depth+1, key)
	newNode.Right = buildKDTree(vectors[medianIndex+1:], depth+1, key)

	return newNode
}

// buildIndex recursively builds the index mapping vector IDs to their corresponding nodes
func buildIndex(node *KDNode, index Index) {
	if node == nil {
		return
	}

	index[node.Vector.ID] = node
	buildIndex(node.Left, index)
	buildIndex(node.Right, index)
}
