package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Default values for testing
	AssetStorageBaseDir  = "/tmp/asset-storage"         // Base directory for all assets.
	AssetStorageTempDir  = "/tmp/asset-storage/tmp"     // Temporary directory. Should be on same drive as AssetStorageBaseDir
	AssetMetaDataBaseDir = "/tmp/asset-metadata"        // Base directory for all meta-data of assets.
	AssetMetaDataDb      = "/tmp/asset-metadata.sqlite" // Data source name of database

	SpaHttpRoot = filepath.Dir(os.Args[0]) + "/vue-ui" // Root directory to service SPA from

	baseDir = flag.String("base", "", "Base directory for storage, meta-data, db...")
)

// LoadDefault initializes configuration with defaults,
// paths either with base-directory from command-line argument "-base"
// or in user-home
func LoadDefault() {

	flag.Parse()

	useBaseDir := *baseDir
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

	//Development: When using "go run...", path must be set manually
	if strings.Contains(SpaHttpRoot, "go-build") {
		SpaHttpRoot = "/home/christianh/Workspace/go/asset-storage/vue-ui"
	}
	//Testing Build: Adjust path if in "bin" dir
	if strings.Contains(SpaHttpRoot, "/bin/") || strings.Contains(SpaHttpRoot, "\\bin/") {
		SpaHttpRoot = filepath.Dir(os.Args[0]) + "/../vue-ui"
	}
}
