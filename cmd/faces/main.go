package main

import (
	"fmt"
	"math"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/storage"
)

func main() {

	config.LoadDefault()
	storage.CreateDirectories()

	embeddings := faces.ReadEmbeddings(config.AssetFacesBaseDir)
	fmt.Printf("Found %d embeddings\n", len(embeddings))

	CalculateSimilarity(embeddings, 0.45)

}

func CalculateSimilarity(embeddings map[string]faces.Embedding, threshold float64) {

	len := len(embeddings)
	paths := make([]string, len)
	embeds := make([]faces.Embedding, len)

	i := 0
	for path, emb := range embeddings {
		paths[i] = path
		embeds[i] = emb
		i++
	}

	for i := range len {
		for j := i + 1; j < len; j++ {
			sim := CosineSimilarity(embeds[i], embeds[j])
			if sim >= threshold {
				fmt.Printf("Similarity %s <-> %s = %f\n", paths[i], paths[j], sim)
			}
		}
	}
}

func CosineSimilarity(a, b faces.Embedding) float64 {
	var dot, normA, normB float64

	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
