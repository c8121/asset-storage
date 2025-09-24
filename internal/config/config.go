package config

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	//Default values for testing
	AssetStorageBaseDir  = "/tmp/asset-storage"         // Base directory for all assets.
	AssetStorageTempDir  = "/tmp/asset-storage/tmp"     // Temporary directory. Should be on same drive as AssetStorageBaseDir
	AssetMetaDataBaseDir = "/tmp/asset-metadata"        // Base directory for all meta-data of assets.
	AssetMetaDataDb      = "/tmp/asset-metadata.sqlite" // Data source name of database

	SpaHttpRoot = filepath.Dir(os.Args[0]) + "/vue-ui" // Root directory to servce SPA from
)

// LoadDefault initializes configuration with defaults, paths in user-home
func LoadDefault() {

	userHome, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user home directory")
	}

	AssetStorageBaseDir = userHome + "/asset-storage/files"
	AssetStorageTempDir = userHome + "/asset-storage/tmp"
	AssetMetaDataBaseDir = userHome + "/asset-storage/meta"
	AssetMetaDataDb = userHome + "/asset-storage/db/asset-metadata.sqlite"

	//Development: When using "go run...", path must be set manually
	if strings.Contains(SpaHttpRoot, "go-build") {
		SpaHttpRoot = "/home/christianh/Workspace/go/asset-storage/vue-ui"
	}

}
