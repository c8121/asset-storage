package faces

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
	"github.com/c8121/asset-storage/internal/util"
)

func CalculateSimilarity(embeddings map[string]Embedding, threshold float64) {

	cnt := len(embeddings)
	paths := make([]string, cnt)
	embeds := make([]Embedding, cnt)

	i := 0
	for path, emb := range embeddings {
		paths[i] = path
		embeds[i] = emb
		i++
	}

	for i := range cnt {
		for j := i + 1; j < cnt; j++ {
			sim := CosineSimilarity(embeds[i], embeds[j])
			if sim >= threshold {
				fmt.Printf("Similarity %s <-> %s = %f\n", paths[i], paths[j], sim)
				p := strings.Index(paths[i], "/")
				hashA := paths[i][:p]
				faceA, _ := strconv.Atoi(paths[i][p+1:])
				hashB := paths[j][:p]
				faceB, _ := strconv.Atoi(paths[j][p+1:])

				err := metadata_db_entity.AddFaceSimilarity(hashA, faceA, hashB, faceB)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func CosineSimilarity(a, b Embedding) float64 {
	var dot, normA, normB float64

	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func ReadEmbeddings(baseDir string) map[string]Embedding {

	embedding := make(map[string]Embedding, 0)

	doReadEmbeddings(baseDir, embedding)

	return embedding

}

func doReadEmbeddings(dir string, into map[string]Embedding) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read dir: %s\n", err)
		return
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		stat, err := os.Stat(path)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}

		if stat.IsDir() {
			doReadEmbeddings(path, into)
		} else if strings.HasSuffix(e.Name(), FaceEmbeddingExtension) {

			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}

			embedding := Embedding{}
			err = json.Unmarshal(data, &embedding)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}

			key := keyFromFacePath(path)

			into[key] = embedding
		}
	}
}

// Exctract hash from .../hash[:2]/hash[2:]/faceIdx.ext
func keyFromFacePath(path string) string {
	items := util.SplitPath(path)
	l := len(items)
	if l < 3 {
		return "invalid-path"
	}

	faceIdx := items[l-1]
	p := strings.Index(faceIdx, ".")
	if p > -1 {
		faceIdx = faceIdx[:p]
	}

	return items[l-3] + items[l-2] + "/" + faceIdx
}
