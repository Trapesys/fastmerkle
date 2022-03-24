package fastmerkle

import (
	"crypto/rand"
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
		randomData = append(randomData, data)
	}

	return randomData
}
