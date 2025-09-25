package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// Default values for testing
	AssetStorageBaseDir  = "/tmp/asset-storage"         // Base directory for all assets.
	AssetStorageTempDir  = "/tmp/asset-storage/tmp"     // Temporary directory. Should be on same drive as AssetStorageBaseDir
	AssetMetaDataBaseDir = "/tmp/asset-metadata"        // Base directory for all meta-data of assets.
	AssetMetaDataDb      = "/tmp/asset-metadata.sqlite" // Data source name of database

	UseGzip = false //Note: Cannot be changed after storage was created!

	SpaHttpRoot = filepath.Dir(os.Args[0]) + "/vue-ui" // Root directory to service SPA from

	cmdBaseDir     = flag.String("base", "", "Base directory for storage, meta-data, db...")
	cmdUseGzip     = flag.Bool("gzip", false, "Use GZIP compression")
	cmdSpaHttpRoot = flag.String("spa", "", "HTTP root directory of SPA app")
)

// LoadDefault initializes configuration with defaults,
// paths either with base-directory from command-line argument "-base"
// or in user-home
func LoadDefault() {

	flag.Parse()

	useBaseDir := *cmdBaseDir
	if useBaseDir == "" {
		if userHome, err := os.UserHomeDir(); err != nil {
			panic("failed to get user home directory")
		} else {
			useBaseDir = userHome
		}
	}

	fmt.Printf("Using base directory: %s\n", useBaseDir)

	AssetStorageBaseDir = useBaseDir + "/asset-storage/files"
	AssetStorageTempDir = useBaseDir + "/asset-storage/tmp"
	AssetMetaDataBaseDir = useBaseDir + "/asset-storage/meta"
	AssetMetaDataDb = useBaseDir + "/asset-storage/db/asset-metadata.sqlite"

	UseGzip = *cmdUseGzip
	if UseGzip {
		fmt.Printf("Using GZIP\n")
	}

	useSpaHttpRoot := *cmdSpaHttpRoot
	if useSpaHttpRoot != "" {
		SpaHttpRoot = useSpaHttpRoot
		fmt.Printf("Using SPA directory: %s\n", SpaHttpRoot)
	}
}
