package vector

import (
	"encoding/json"
	"math"
)

type KDNode struct {
	Vector        Vector
	Left          *KDNode
	Right         *KDNode
	Dimension     int
	EncryptedData string
}

func insertKDNode(root *KDNode, vector Vector, encryptedData string, depth int) *KDNode {
	if root == nil {
		return &KDNode{
			Vector:        vector,
			EncryptedData: encryptedData,
			Left:          nil,
			Right:         nil,
			Dimension:     depth % len(vector.Values),
		}
	}

	dim := root.Dimension
	if vector.Values[dim] < root.Vector.Values[dim] {
		root.Left = insertKDNode(root.Left, vector, encryptedData, depth+1)
	} else {
		root.Right = insertKDNode(root.Right, vector, encryptedData, depth+1)
	}

	return root
}

func distance(v1, v2 []float64) float64 {
	sum := 0.0
	for i := range v1 {
		diff := v1[i] - v2[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func nearestNeighborKD(root *KDNode, query Vector, depth int, best *KDNode, bestDist *float64, key []byte) (*KDNode, float64, error) {
	if root == nil {
		return best, *bestDist, nil
	}

	decryptedData, err := decrypt(root.EncryptedData, key)
	if err != nil {
		return nil, 0, err
	}

	var values []float64
	err = json.Unmarshal(decryptedData, &values)
	if err != nil {
		return nil, 0, err
	}

	d := distance(query.Values, values)
	if d < *bestDist {
		*bestDist = d
		best = root
	}

	dim := root.Dimension
	diff := query.Values[dim] - values[dim]

	var first, second *KDNode
	if diff < 0 {
		first, second = root.Left, root.Right
	} else {
		first, second = root.Right, root.Left
	}

	best, *bestDist, err = nearestNeighborKD(first, query, depth+1, best, bestDist, key)
	if err != nil {
		return nil, 0, err
	}

	if math.Abs(diff) < *bestDist {
		best, *bestDist, err = nearestNeighborKD(second, query, depth+1, best, bestDist, key)
		if err != nil {
			return nil, 0, err
		}
	}

	return best, *bestDist, nil
}
