package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/faces"
	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
	"github.com/c8121/asset-storage/internal/storage"
)

func main() {

	command := flag.String("command", "indentify", "Command, either 'identify' or 'similarity'")
	threshold := flag.Float64("threshold", 0.45, "Minimun similariy threshold")

	config.LoadDefault()
	storage.CreateDirectories()

	fmt.Printf("Command: %s\n", *command)

	if strings.HasPrefix(*command, "i") {

		handler := func(path string) {
			hash := storage.HashFromStoragePath(path)
			faces, err := faces.GetFaces(hash)
			if err != nil {
				fmt.Printf("Cannot get faces from %s: %s\n", hash, err)
			} else {
				fmt.Printf("Found %d faces in %s\n", len(faces), hash)
			}
		}

		storage.Walk(handler)

	} else if strings.HasPrefix(*command, "s") {

		mdsqlite.Open()
		defer mdsqlite.Close()

		embeddings := faces.ReadEmbeddings(config.AssetFacesBaseDir)

		faces.CalculateSimilarity(embeddings, *threshold)
		fmt.Printf("Checked %d embeddings\n", len(embeddings))

	} else {
		fmt.Printf("Unknown command: %s\n", *command)
	}

}
