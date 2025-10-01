package metadata_db_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/c8121/asset-storage/internal/config"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	"github.com/c8121/asset-storage/internal/util"
	_ "modernc.org/sqlite"
)

func TestDb(t *testing.T) {

	db, err := sql.Open("sqlite", "file::memory:")
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AssetMetaDataDb)
	defer db.Close()

	metadata_db.DB = db

	metadata_db.DbInitMimeType()

	mimeType, err := metadata_db.DbGetMimeType("text/plain", false)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType)
	if mimeType != nil {
		t.Errorf("MimeType found, but none created so far")
	}

	mimeType, err = metadata_db.DbGetMimeType("text/plain", true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType)
	if mimeType == nil {
		t.Errorf("MimeType not found: %v", err)
	}

	mimeType2, err := metadata_db.DbGetMimeType(" text/plain ", true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType2)
	if mimeType2.Id != mimeType.Id {
		t.Errorf("MimeType id missmatch: %v, %v", mimeType2, mimeType)
	}

	mimeType2, err = metadata_db.DbGetMimeType(" text/plain;charset=UTF8 ", true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType2)
	if mimeType2.Id != mimeType.Id {
		t.Errorf("MimeType id missmatch: %v, %v", mimeType2, mimeType)
	}
}
