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

	metadata_db.SetDatabase(db)

	mimeTypeTest(t)
	assetTest(t)
}

func assetTest(t *testing.T) {

	var asset = &metadata_db.Asset{
		Hash: "test",
	}

	err := asset.Load()
	fmt.Println("asset:", asset)
	if err != metadata_db.ErrNotFound {
		t.Errorf("Asset found, but none created so far")
	}

	err = asset.Create(true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("asset:", asset)
	if asset.Id == 0 {
		t.Errorf("Asset not found")
	}

	var asset2 = &metadata_db.Asset{
		Hash: "test",
	}

	err = asset2.Load()
	util.PanicOnError(err, "Test failed")
	fmt.Println("asset:", asset2)
	if asset2.Id == 0 {
		t.Errorf("Asset not found")
	}
	if asset2.Id != asset.Id {
		t.Errorf("Wrong asset not found")
	}

	var asset3 = &metadata_db.Asset{
		Hash: "test2",
	}

	err = asset3.Save()
	util.PanicOnError(err, "Test failed")
	fmt.Println("asset:", asset3)
	if asset3.Id == 0 {
		t.Errorf("Asset not saved")
	}

	savedId := asset3.Id
	err = asset3.Save()
	util.PanicOnError(err, "Test failed")
	if asset3.Id != savedId {
		t.Errorf("Asset duplicated")
	}

	mimeType, err := metadata_db.DbGetMimeType("text/plain", true)
	util.PanicOnError(err, "Test failed")
	asset3.MimeType = mimeType.Id
	err = asset3.Save()
	util.PanicOnError(err, "Test failed")

	var asset4 = &metadata_db.Asset{Hash: "test2"}
	err = asset4.Load()
	util.PanicOnError(err, "Test failed")
	fmt.Println("asset:", asset4)
	if asset4.Id == 0 {
		t.Errorf("Asset not found")
	}
	if asset4.Id != asset3.Id {
		t.Errorf("Asset mismatch")
	}
}

func mimeTypeTest(t *testing.T) {

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
