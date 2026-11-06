package bootstrap

import (
	"database/sql"
	"log/slog"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata_db/entity"
	metadata_db_sqlite "github.com/c8121/asset-storage/internal/metadata_db/sqlite"
	apputil "github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_db_entities/sqlite"
	"github.com/c8121/go_vue_spa/util"
	_ "gosqlite.org" // Bundles modernc + sqlite-vec
)

var (
	database *sql.DB
)

// Open Connect to SQLite database file + init
func OpenDatabase() {

	dbDir := filepath.Dir(config.AppConfig.AssetMetaDataDb)
	err := apputil.CreateDirIfNotExists(dbDir, metadata.FilePermissions)
	if err != nil {
		panic(err)
	}

	url := "file:" + config.AppConfig.AssetMetaDataDb +
		"?_pragma=journal_mode(wal)" +
		"&_pragma=busy_timeout(2500)" +
		"&_pragma=synchronous(normal)" +
		"&_txlock=immediate"

	slog.Info("Open DB", "database", config.AppConfig.AssetMetaDataDb)

	database, err = sql.Open("sqlite", url)
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AppConfig.AssetMetaDataDb)

	//Creation of face-embedding-table cannot be done with go_db_entities
	queryBuilder := sqlite.NewQueryBuilder()
	queryBuilder.CreateQueryBuilder = &metadata_db_sqlite.FaceEmbeddingCreateQueryBuilder{
		DefaultCreateQueryBuilder: queryBuilder.CreateQueryBuilder,
	}

	metadata_db_entity.DbContext = go_db_entities.NewContext(database, queryBuilder, sqlite.NewFieldMapping())

	//Finding face-embeddings cannot be done with go_db_entities
	metadata_db_entity.FacesDbContext = metadata_db_entity.NewContext(database, metadata_db_sqlite.FaceEmbeddingQueryBuilder{})

	metadata_db_entity.CheckCreateTables()

	slog.Info("DB Ready")
}

func GetDatabase() *sql.DB {
	return database
}

// Close Disconnect from Database
func CloseDatabase() {
	slog.Info("Close DB", "database", config.AppConfig.AssetMetaDataDb)
	util.CloseOrLog(database)
}
