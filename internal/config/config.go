package config

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// Default values for testing
	AssetStorageBaseDir     = "/tmp/asset-storage"         // Base directory for all assets.
	AssetStorageTempDir     = "/tmp/asset-storage/tmp"     // Temporary directory. Should be on same drive as AssetStorageBaseDir
	AssetMetaDataBaseDir    = "/tmp/asset-metadata"        // Base directory for all meta-data of assets.
	AssetMetaDataDb         = "/tmp/asset-metadata.sqlite" // Data source name of database
	AssetCollectionsBaseDir = "/tmp/asset-collections"     // Base directory for collections.

	UseGzip = false //Note: Cannot be changed after storage was created!
	XorKey  []byte  //Note: Cannot be changed after storage was created!

	MaxMemFileSize int64 = 1000 * 1000 * 400

	SkipMetaDataIfExists = false

	SpaHttpRoot = filepath.Dir(os.Args[0]) + "/vue-ui" // Root directory to service SPA from

	cmdBaseDir              = flag.String("base", "", "Base directory for storage, meta-data, db...")
	cmdUseGzip              = flag.Bool("gzip", false, "Use GZIP compression")
	cmdSpaHttpRoot          = flag.String("spa", "", "HTTP root directory of SPA app")
	cmdMaxMemFileSize       = flag.Int64("maxmem", 0, "Max memory file size in bytes")
	cmdSkipMetaDataIfExists = flag.Bool("skip-meta", false, "Skip meta data update if file exist")

	cmdXorKey = flag.String("xor", "", "XOR Key for content obfusication")
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
	AssetCollectionsBaseDir = useBaseDir + "/asset-storage/collections"

	UseGzip = *cmdUseGzip
	if UseGzip {
		fmt.Printf("Using GZIP\n")
	}

	if *cmdMaxMemFileSize > 0 {
		MaxMemFileSize = *cmdMaxMemFileSize
		fmt.Printf("Max memory file size: %d\n", MaxMemFileSize)
	}

	SkipMetaDataIfExists = *cmdSkipMetaDataIfExists
	if SkipMetaDataIfExists {
		fmt.Printf("Will not update meta-data on existing files\n")
	}

	if *cmdSpaHttpRoot != "" {
		SpaHttpRoot = *cmdSpaHttpRoot
		fmt.Printf("Using SPA directory: %s\n", SpaHttpRoot)
	}

	if *cmdXorKey != "" {
		if len(*cmdXorKey) < 64 {
			sha := sha256.New()
			XorKey = fmt.Appendf(nil, "%x", sha.Sum([]byte(*cmdXorKey)))
		} else {
			XorKey = []byte(*cmdXorKey)
		}
		fmt.Printf("Xor obfusication enabled, key length: %d (%s)\n", len(XorKey), string(XorKey))
	}
}
