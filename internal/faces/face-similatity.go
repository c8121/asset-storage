package faces

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

			into[path] = embedding
		}
	}
}
