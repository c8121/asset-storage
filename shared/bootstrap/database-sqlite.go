package bootstrap

import (
	"database/sql"
	"fmt"
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

// Open Connect to SQLite database file + init
func OpenDatabase() {

	dbDir := filepath.Dir(config.AppConfig.AssetMetaDataDb)
	apputil.CreateDirIfNotExists(dbDir, metadata.FilePermissions)

	url := "file:" + config.AppConfig.AssetMetaDataDb +
		"?_pragma=journal_mode(wal)" +
		"&_pragma=busy_timeout(2500)" +
		"&_pragma=synchronous(normal)" +
		"&_txlock=immediate"

	fmt.Printf("Open DB %s\n", config.AppConfig.AssetMetaDataDb)
	db, err := sql.Open("sqlite", url)
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AppConfig.AssetMetaDataDb)

	//Creation of face-embedding-table cannot be done with go_db_entities
	queryBuilder := sqlite.NewQueryBuilder()
	queryBuilder.CreateQueryBuilder = &metadata_db_sqlite.FaceEmbeddingCreateQueryBuilder{
		DefaultCreateQueryBuilder: queryBuilder.CreateQueryBuilder,
	}

	metadata_db_entity.DbContext = go_db_entities.NewContext(db, queryBuilder, sqlite.NewFieldMapping())

	//Finding face-embeddings cannot be done with go_db_entities
	metadata_db_entity.FacesDbContext = metadata_db_entity.NewContext(db, metadata_db_sqlite.FaceEmbeddingQueryBuilder{})

	metadata_db_entity.CheckCreateTables()

	fmt.Println("DB Ready")
}

// Close Disconnect from Database
func CloseDatabase() {
	fmt.Printf("Close DB %s\n", config.AppConfig.AssetMetaDataDb)
	util.CloseOrLog(metadata_db_entity.DbContext)
}
