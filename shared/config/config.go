package config

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

type Config struct {
	AssetStorageConfigDir     string // Base directory for config files.
	AssetStorageBaseDir       string // Base directory for all assets.
	AssetStorageTempDir       string // Temporary directory. Should be on same drive as AssetStorageBaseDir
	AssetMetaDataBaseDir      string // Base directory for all metadata of assets.
	AssetMetaDataDb           string // Data source name of database
	AssetCollectionsBaseDir   string // Base directory for collections.
	AssetFulltextIndexBaseDir string // Base directory for fulltext index.

	FilePermissions fs.FileMode

	UseGzip bool   //Note: Cannot be changed after storage was created!
	XorKey  []byte //Note: Cannot be changed after storage was created!

	MaxMemFileSize int64

	SkipMetaDataIfExists bool
	CheckHashBeforeAdd   bool

	SpaHttpRoot string // Root directory to service SPA from. Using embedded FS if empty.

	ListenAddress string //Spa-/Rest-/Ssh-Server in format ip:port
	CertFile      string
	KeyFile       string

	PdfToTextCharCountToUseOcr int
}

var (
	AppConfig *Config = nil
)

func NewConfigFromFlags(errorHandling flag.ErrorHandling) *Config {

	flagSet := flag.NewFlagSet("asset-storage", errorHandling)

	cmdDataDir := flagSet.String("data", "", "Data directory for storage, meta-data, db...")
	cmdUseGzip := flagSet.Bool("gzip", false, "Use GZIP compression")
	cmdXorKey := flagSet.String("xor", "", "XOR Key for content obfuscation")
	cmdMaxMemFileSize := flagSet.Int64("maxmem", 0, "Max memory file size in bytes")
	cmdSpaHttpRoot := flagSet.String("spa", "", "HTTP root directory of SPA app")
	cmdSkipMetaDataIfExists := flagSet.Bool("skip-meta", false, "Skip meta data update if file exist")
	cmdCheckHashBeforeAdd := flagSet.Bool("check-hash", false, "Check hash before trying to add file.")
	cmdListen := flagSet.String("listen", "", "Listen Address (ip:port)")
	cmdCertFile := flagSet.String("cert", "", "Certificate file")
	cmdKeyFile := flagSet.String("key", "", "Key file")
	cmdCharCountPdfToTextOcr := flagSet.Int("pdf-ocr", 10, "Use OCR when converting PDF to Text, if result does not have more chars")
	cmdLogVerbose := flagSet.Bool("v", false, "Log verbose")

	flagSet.Parse(os.Args[1:])

	config := &Config{}

	useDataDir := *cmdDataDir
	if useDataDir == "" {
		if userHome, err := os.UserHomeDir(); err != nil {
			panic("failed to get user home directory")
		} else {
			useDataDir = filepath.Join(userHome, "asset-storage")
		}
	}

	validateDataDir(useDataDir)

	slog.Info("Init path", "data", useDataDir)
	SetDirectories(config, useDataDir)

	config.FilePermissions = 0744

	config.UseGzip = *cmdUseGzip
	if config.UseGzip {
		slog.Info("Using GZIP")
	}

	if *cmdMaxMemFileSize > 0 {
		config.MaxMemFileSize = *cmdMaxMemFileSize
		slog.Info("Max memory file size", "bytes", config.MaxMemFileSize)
	} else {
		config.MaxMemFileSize = 1000 * 1000 * 400
	}

	config.SkipMetaDataIfExists = *cmdSkipMetaDataIfExists
	if config.SkipMetaDataIfExists {
		slog.Info("Will not update meta-data on existing files")
	}

	config.CheckHashBeforeAdd = *cmdCheckHashBeforeAdd
	if config.CheckHashBeforeAdd {
		slog.Info("Checking hash before adding file (this is faster only if most if files already exist.")
	}

	if *cmdSpaHttpRoot != "" {
		config.SpaHttpRoot = *cmdSpaHttpRoot
		slog.Info("Init Path", "spa", config.SpaHttpRoot)
	}

	if *cmdXorKey != "" {
		if len(*cmdXorKey) < 64 {
			sha := sha256.New()
			config.XorKey = fmt.Appendf(nil, "%x", sha.Sum([]byte(*cmdXorKey)))
		} else {
			config.XorKey = []byte(*cmdXorKey)
		}
		slog.Info("Xor obfuscation enabled", "key length", len(config.XorKey))
	}

	if *cmdListen != "" {
		config.ListenAddress = *cmdListen
		slog.Info("Server", "bind", config.ListenAddress)
	} else {
		config.ListenAddress = "127.0.0.1:9999"
	}

	if *cmdCertFile != "" {
		config.CertFile = *cmdCertFile
		slog.Info("Init TLS", "cert", config.CertFile)
	}

	if *cmdKeyFile != "" {
		config.KeyFile = *cmdKeyFile
		slog.Info("Init TLS", "key", config.KeyFile)
	}

	if *cmdCharCountPdfToTextOcr >= 0 {
		config.PdfToTextCharCountToUseOcr = *cmdCharCountPdfToTextOcr
		slog.Info("Using OCR if for PDF to text conversion if non-OCR result", "char limit", config.PdfToTextCharCountToUseOcr)
	} else {
		config.PdfToTextCharCountToUseOcr = 10
	}

	if *cmdLogVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return config
}

func SetDirectories(config *Config, baseDir string) {
	config.AssetStorageConfigDir = filepath.Join(baseDir, "config")
	config.AssetStorageBaseDir = filepath.Join(baseDir, "files")
	config.AssetStorageTempDir = filepath.Join(baseDir, "tmp")
	config.AssetMetaDataBaseDir = filepath.Join(baseDir, "meta")
	config.AssetMetaDataDb = filepath.Join(baseDir, "db/asset-metadata.sqlite")
	config.AssetCollectionsBaseDir = filepath.Join(baseDir, "collections")
	config.AssetFulltextIndexBaseDir = filepath.Join(baseDir, "index")
}

// validateDataDir checks is given path is an asset-storage directory.
// A non-existing path is considered valid so it can be created later.
func validateDataDir(path string) {

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return
	} else if stat == nil || !stat.IsDir() {
		fmt.Printf("%s is not a directory\n", path)
		os.Exit(1)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	validNames := []string{"config", "files", "tmp", "meta", "db", "collections", "faces", "index", "bin"}

	for _, entry := range entries {
		if !slices.Contains(validNames, entry.Name()) {
			fmt.Printf("Invalid data dir: %s (invalid item \"%s\")\n", path, entry.Name())
			os.Exit(1)
		}
	}

}
