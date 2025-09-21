package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/storage"
)

func main() {
	flag.Parse()
	files := flag.Args()
	if files == nil || len(files) == 0 {
		fmt.Printf("You must specify at least one file: %s file [file...]\n", filepath.Base(os.Args[0]))
	}

	for _, file := range files {
		_, _, err := storage.AddFile(file)
		if err != nil {
			fmt.Printf("Error adding '%s': %s\n", file, err)
		}
	}
}
