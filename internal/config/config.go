package config

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

var (
	// Default values for testing
	AssetStorageConfigDir   = "/tmp/asset-storage/config"                   // Base directory for config files.
	AssetStorageBaseDir     = "/tmp/asset-storage/files"                    // Base directory for all assets.
	AssetStorageTempDir     = "/tmp/asset-storage/tmp"                      // Temporary directory. Should be on same drive as AssetStorageBaseDir
	AssetMetaDataBaseDir    = "/tmp/asset-storage/meta"                     // Base directory for all meta-data of assets.
	AssetMetaDataDb         = "/tmp/asset-storage/db/asset-metadata.sqlite" // Data source name of database
	AssetCollectionsBaseDir = "/tmp/asset-collections"                      // Base directory for collections.
	AssetFacesBaseDir       = "/tmp/asset-storage/faces"                    // Base directory for all meta-data of assets.

	UseGzip = false //Note: Cannot be changed after storage was created!
	XorKey  []byte  //Note: Cannot be changed after storage was created!

	MaxMemFileSize int64 = 1000 * 1000 * 400

	SkipMetaDataIfExists = false
	CheckHashBeforeAdd   = false

	SpaHttpRoot = filepath.Dir(os.Args[0]) + "/../vue-ui" // Root directory to service SPA from

	ListenAddress = "127.0.0.1:9999" //Spa-/Rest-/Ssh-Server in format ip:port
	CertFile      = ""
	KeyFile       = ""

	cmdDataDir              = flag.String("data", "", "Data directory for storage, meta-data, db...")
	cmdUseGzip              = flag.Bool("gzip", false, "Use GZIP compression")
	cmdXorKey               = flag.String("xor", "", "XOR Key for content obfusication")
	cmdMaxMemFileSize       = flag.Int64("maxmem", 0, "Max memory file size in bytes")
	cmdSpaHttpRoot          = flag.String("spa", "", "HTTP root directory of SPA app")
	cmdSkipMetaDataIfExists = flag.Bool("skip-meta", false, "Skip meta data update if file exist")
	cmdCheckHashBeforeAdd   = flag.Bool("check-hash", false, "Check hash before trying to add file.")
	cmdListen               = flag.String("listen", "", "Listen Address (ip:port)")
	cmdCertFile             = flag.String("cert", "", "Certificate file")
	cmdKeyFile              = flag.String("key", "", "Key file")
)

// LoadDefault initializes configuration with defaults,
// paths either with base-directory from command-line argument "-base"
// or in user-home
func LoadDefault() {

	checkObsoleteCmdArgs()

	flag.Parse()

	useDataDir := *cmdDataDir
	if useDataDir == "" {
		if userHome, err := os.UserHomeDir(); err != nil {
			panic("failed to get user home directory")
		} else {
			useDataDir = filepath.Join(userHome, "asset-storage")
		}
	}

	validateDataDir(useDataDir)

	fmt.Printf("Using data directory: %s\n", useDataDir)

	AssetStorageConfigDir = filepath.Join(useDataDir, "config")
	AssetStorageBaseDir = filepath.Join(useDataDir, "files")
	AssetStorageTempDir = filepath.Join(useDataDir, "tmp")
	AssetMetaDataBaseDir = filepath.Join(useDataDir, "meta")
	AssetMetaDataDb = filepath.Join(useDataDir, "db/asset-metadata.sqlite")
	AssetCollectionsBaseDir = filepath.Join(useDataDir, "collections")
	AssetFacesBaseDir = filepath.Join(useDataDir, "faces")

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

	CheckHashBeforeAdd = *cmdCheckHashBeforeAdd
	if CheckHashBeforeAdd {
		fmt.Printf("Checking hash before adding file (this is faster only if most if files already exist.\n")
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
		fmt.Printf("Xor obfuscation enabled, key length: %d\n", len(XorKey))
	}

	if *cmdListen != "" {
		ListenAddress = *cmdListen
		fmt.Printf("Server address: %s\n", ListenAddress)
	}

	if *cmdCertFile != "" {
		CertFile = *cmdCertFile
		fmt.Printf("Certificate file: %s\n", CertFile)
	}

	if *cmdKeyFile != "" {
		KeyFile = *cmdKeyFile
		fmt.Printf("Key file: %s\n", KeyFile)
	}
}

// validateDataDir checks is given path is a asset-storage directory.
// A non-existing path is considered valid so it can be created later.
func validateDataDir(path string) {

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return
	}
	if stat == nil || !stat.IsDir() {
		fmt.Printf("%s is not a directory\n", path)
		os.Exit(1)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	validNames := []string{"config", "files", "tmp", "meta", "db", "collections", "faces", "bin", "vue-ui"}

	for _, entry := range entries {
		if !slices.Contains(validNames, entry.Name()) {
			fmt.Printf("Invalid data dir: %s (invalid item \"%s\")\n", path, entry.Name())
			os.Exit(1)
		}
	}

}

// checkObsoleteCmdArgs looks for command-line arguments which are not supported anymore.
func checkObsoleteCmdArgs() {

	obsoleteBase := flag.String("base", "", "")

	flag.Parse()

	if *obsoleteBase != "" {
		fmt.Println("Argument 'base' is obsolete. Use 'data' instead, pointing directly to the data directory")
		os.Exit(1)
	}
}
