package metadata_db_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/c8121/asset-storage/internal/config"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
	"github.com/c8121/asset-storage/internal/util"
	_ "modernc.org/sqlite"
)

func TestDb(t *testing.T) {

	db, err := sql.Open("sqlite", "file::memory:")
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AssetMetaDataDb)
	defer util.CloseOrLog(db)

	metadata_db.SetDatabase(db)
	metadata_db_entity.AutoCreateEntities()

	mimeTypeTest(t)
	assetTest(t)
	pathTest(t, true)
	pathTest(t, false) //Test cache

	fmt.Println("OK")
}

func pathTest(t *testing.T, testDontCreateItem bool) {

	testPaths := []string{"", "test1", "test2/", "/test3/", "/test4", "/test5\\\\", "C5:/test/", "C6:\\test", "C7:/test/file.txt",
		"C8:\\test\\file.txt", "///test9///file.txt", "file:///C10:/home/test", "/", "\\"}
	expectLen := []int{1, 1, 1, 1, 1, 1, 2, 2, 3,
		3, 2, 3, 1, 1}

	for i, path := range testPaths {
		split := metadata_db_entity.SplitPath(path)
		fmt.Printf("%d: '%s' = %s\n", i, path, split)
		if len(split) != expectLen[i] {
			t.Errorf("Expected len is %d, but result len is %d: %s", expectLen[i], len(split), path)
		}

		if testDontCreateItem {
			pathItem, err := metadata_db_entity.GetPathItem(path, false)
			if err != metadata_db_entity.ErrNotFound {
				t.Errorf("Found path item, but should not be created %s: %v, %v", path, pathItem, err)
			}
		}

		pathItem, err := metadata_db_entity.GetPathItem(path, true)
		if err != nil {
			t.Errorf("Failed to get path item %s: %v", path, err)
		}
		if pathItem == nil {
			t.Errorf("Failed to get path item %s: %v", path, err)
		}
		fmt.Printf("    %v\n", pathItem)
	}
}

func assetTest(t *testing.T) {

	var asset = &metadata_db_entity.Asset{
		Hash: "test",
	}

	err := asset.Load()
	fmt.Println("asset:", asset)
	if err != metadata_db_entity.ErrNotFound {
		t.Errorf("Asset found, but none created so far")
	}

	err = asset.Get(true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("asset:", asset)
	if asset.Id == 0 {
		t.Errorf("Asset not found")
	}

	var asset2 = &metadata_db_entity.Asset{
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

	var asset3 = &metadata_db_entity.Asset{
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

	mimeType, err := metadata_db_entity.GetMimeType("text/plain", true)
	util.PanicOnError(err, "Test failed")
	asset3.MimeType = mimeType.Id
	err = asset3.Save()
	util.PanicOnError(err, "Test failed")

	var asset4 = &metadata_db_entity.Asset{Hash: "test2"}
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

	mimeType, err := metadata_db_entity.GetMimeType("text/plain", false)
	fmt.Println("mimeType:", mimeType)
	if err != metadata_db_entity.ErrNotFound {
		t.Errorf("MimeType found, but none created so far")
	}

	mimeType, err = metadata_db_entity.GetMimeType("text/plain", true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType)
	if mimeType == nil {
		t.Errorf("MimeType not found: %v", err)
	}

	mimeType2, err := metadata_db_entity.GetMimeType(" text/plain ", true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType2)
	if mimeType2.Id != mimeType.Id {
		t.Errorf("MimeType id missmatch: %v, %v", mimeType2, mimeType)
	}

	mimeType2, err = metadata_db_entity.GetMimeType(" text/plain;charset=UTF8 ", true)
	util.PanicOnError(err, "Test failed")
	fmt.Println("mimeType:", mimeType2)
	if mimeType2.Id != mimeType.Id {
		t.Errorf("MimeType id missmatch: %v, %v", mimeType2, mimeType)
	}
}
