package fastmerkle

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkGenerateMerkleTree5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMerkleTree(
			generateRandomData(5),
		)
	}
}

func BenchmarkGenerateMerkleTree50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMerkleTree(
			generateRandomData(50),
		)
	}
}

func BenchmarkGenerateMerkleTree500(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMerkleTree(
			generateRandomData(500),
		)
	}
}

func BenchmarkGenerateMerkleTree1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMerkleTree(
			generateRandomData(1000),
		)
	}
}

func BenchmarkGenerateMerkleTree10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMerkleTree(
			generateRandomData(10000),
		)
	}
}

func BenchmarkGenerateMerkleTree1000000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMerkleTree(
			generateRandomData(1000000),
		)
	}
}

// generateRandomData generates random byte data
func generateRandomData(count int) [][]byte {
	randomData := make([][]byte, count)

	for i := 0; i < count; i++ {
		data := make([]byte, 32)
		_, _ = rand.Read(data)
		randomData[i] = data
	}

	return randomData
}

// generateInputSet generates an input set from the
// fixed random data
func generateInputSet(count int) [][]byte {
	if count < 1 {
		return nil
	}

	inputSet := [][]byte{
		[]byte("Lazar"),
		[]byte("Vuksan"),
		[]byte("Dusan"),
		[]byte("Aleksa"),
		[]byte("Yoshiki"),
		[]byte("Milos"),
		[]byte("Zeljko"),
	}

	if count > len(inputSet) {
		count = len(inputSet)
	}

	return inputSet[:count]
}

// getHexBytes converts an input string to bytes
func getHexBytes(t *testing.T, input string) []byte {
	t.Helper()

	hexBytes, err := hex.DecodeString(input)
	if err != nil {
		t.Fatalf("Unable to decode hex, %v", err)
	}

	return hexBytes
}

// TestGenerateMerkleTree test Merkle tree generation logic
func TestGenerateMerkleTree(t *testing.T) {
	testTable := []struct {
		name          string
		inputElements [][]byte
		expectedRoot  []byte
		expectedError error
	}{
		{
			"no input data provided",
			nil,
			nil,
			errEmptyDataSet,
		},
		{
			"single element input data",
			generateInputSet(1),
			getHexBytes(
				t,
				"6c4cc993464af6cca9101c82d9a5733d6b8453726834f7fd9b1e7a6104915065",
			),
			nil,
		},
		{
			"input data set is power of 2",
			generateInputSet(2),
			getHexBytes(
				t,
				"2997f58b4810eb8d4e779f69e51ab80dc85d1a962a5036d02b21f485e1557c35",
			),
			nil,
		},
		{
			"input data set is not a power of 2 (even)",
			generateInputSet(6),
			getHexBytes(
				t,
				"de0aa53f3e453e031fc92844aa1845f29b69909f6123fed4b047bc253174a497",
			),
			nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Generate the Merkle tree
			merkleTree, genErr := GenerateMerkleTree(testCase.inputElements)

			// Check the errors
			if testCase.expectedError != nil {
				// Make sure the error is correct
				assert.ErrorIs(t, genErr, testCase.expectedError)

				// Make sure the Merkle tree has not been generated
				assert.Nil(t, merkleTree)
			} else {
				// Make sure no error occurred
				assert.NoError(t, genErr)

				if merkleTree == nil || merkleTree.Root == nil {
					t.Fatalf("Merkle tree is not initialized")
				}

				// Make sure the Merkle roots match
				assert.Equal(t, testCase.expectedRoot, merkleTree.GetRootHash())
			}
		})
	}
}
