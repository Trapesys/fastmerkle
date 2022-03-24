package fastmerkle

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	errEmptyDataSet = errors.New("empty data set provided")
)

// GenerateMerkleTree generates a Merkle tree based on the input data
func GenerateMerkleTree(inputData [][]byte) (*MerkleTree, error) {
	// Check if the data set is valid
	if !isValidDataSet(inputData) {
		return nil, errEmptyDataSet
	}

	// Create the worker pool and put them on standby
	workerPool := newWorkerPool(runtime.NumCPU())
	workerPool.startWorkerPool()
	defer workerPool.close()

	// Generate the leaves of the Merkle tree
	nodes, leafErr := generateLeaves(inputData, workerPool)
	if leafErr != nil {
		return nil, fmt.Errorf(
			"unable to generate leaf nodes, %w",
			leafErr,
		)
	}

	// While the root is not derived, create new hashing jobs
	// for the worker pool
	for len(nodes) > 1 {
		// A hashing job is just hashing two subsequent
		// siblings in the tree. Since the tree is a perfect
		// Merkle tree, the node array will always be a power of 2
		for i := 0; i < len(nodes); i += 2 {
			workerPool.addJob(&workerJob{
				storeIndex: i,
				sourceData: [][]byte{
					nodes[i].Hash,
					nodes[i+1].Hash,
				},
			})
		}

		// The Merkle tree is being built from bottom to top,
		// so each level has exactly 1/2 fewer nodes
		// than the previous level (property of perfect binary trees).
		// Therefore, for N nodes on a single tree level, only N/2 results can be expected
		for i := 0; i < len(nodes)/2; i++ {
			result := workerPool.getResult()
			if result.error != nil {
				return nil, fmt.Errorf(
					"unable to perform hashing, %w",
					result.error,
				)
			}

			// Create a placeholder for the parent node
			parent := &Node{
				// Save the hashing data of the 2 children
				Hash: result.hashData,
				// Save a reference to the left child
				Left: nodes[result.storeIndex],
				// Save a reference to the right child
				Right: nodes[result.storeIndex+1],
			}

			// Save the parent reference with the children
			nodes[result.storeIndex].Parent = parent
			nodes[result.storeIndex+1].Parent = parent

			// Overwrite the left child's slot in the array,
			// since it's no longer needed. The right child
			// is also not needed anymore in the original array,
			// and will be overwritten later
			nodes[result.storeIndex] = parent
		}

		// Now that results are gathered for the level,
		// the array can be shifted and shrunk
		shiftAndShrinkArray(&nodes)
	}

	return &MerkleTree{
		Root: nodes[0],
	}, nil
}

// isValidDataSet makes sure the input set has elements
func isValidDataSet(inputData [][]byte) bool {
	return inputData != nil && len(inputData) > 0
}

// shiftAndShrinkArray shifts every other node to the
// beginning of the array, and discards half of it (shrinks it).
// Due to the way results are being stored (index of left child),
// and the fact that the Merkle tree is a perfect binary tree,
// it can be guaranteed that the results are on every other index in the node level array
func shiftAndShrinkArray(nodes *[]*Node) {
	// Put the results in the first half of the array.
	// One counter keeps track of the next slot to place the value (moves by 1) (saveIndx)
	// and the other keeps track of which element should be stored (moves by 2) (resultIndx)
	initialLevelSize := len(*nodes)
	saveIndx := 0
	for resultIndx := 0; resultIndx < initialLevelSize; resultIndx += 2 {
		(*nodes)[saveIndx] = (*nodes)[resultIndx]
		saveIndx++
	}

	// Wipe the other half of the array, since
	// all useful and needed results are in the first half
	*nodes = (*nodes)[:initialLevelSize/2]
}

// generateLeaves generates the initial (leaf) level of the Merkle tree.
// The leaf level needs to be a power of 2, since the Merkle tree is considered
// to be a perfect binary tree
func generateLeaves(inputData [][]byte, wp *workerPool) ([]*Node, error) {
	inputDataSize := len(inputData)
	leafLevelSize := inputDataSize

	if inputDataSize > 1 {
		// Find the nearest power of 2 for the base leaf level size
		leafLevelSize = nearestPowerOf2(len(inputData))
	}

	leaves := make([]*Node, leafLevelSize)

	// Create the initial job set for the leaf nodes,
	// where each job is a single leaf node to be processed
	for i := 0; i < inputDataSize; i++ {
		wp.addJob(&workerJob{
			storeIndex: i,
			sourceData: [][]byte{
				inputData[i],
			},
		})
	}

	// Grab the results from the worker pool
	for i := 0; i < inputDataSize; i++ {
		result := wp.getResult()
		if result.error != nil {
			return nil, fmt.Errorf(
				"unable to perform hashing, %w",
				result.error,
			)
		}

		// Save the leaf nodes
		leaves[i] = &Node{
			Hash:   result.hashData,
			Left:   nil,
			Right:  nil,
			Parent: nil,
		}
	}

	// Since there is a possibility an expansion of the leaves array
	// took place to make it a power of 2, the last element in the original leaf set
	// needs to be duplicated to fill out the remaining (expanded) slots
	lastNode := leaves[inputDataSize-1]
	for i := inputDataSize; i < len(leaves); i++ {
		leaves[i] = lastNode.duplicate()
	}

	return leaves, nil
}

// nearestPowerOf2 returns the nearest power of 2 to
// the provided number.
// Courtesy of BitTwiddlingHacks:
// https://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func nearestPowerOf2(num int) int {
	num--
	num |= num >> 1
	num |= num >> 2
	num |= num >> 4
	num |= num >> 8
	num |= num >> 16
	num++

	return num
}
