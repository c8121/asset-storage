package restapi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

type (
	AddUploadedFileRequest struct {
		TempName string
		Name     string
		Owner    string
		FileTime time.Time
	}
)

// ReceiveUpload is a rest-api handler to receive binary data.
// Adding the file requires a second request: AddUploadedFile(...)
func ReceiveUpload(c *gin.Context) {

	//Read binary file, save as temp-file
	w, err := storage.NewTempFileWriter()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer util.CloseOrLog(w)

	buf := make([]byte, 4096)
	for {
		n, err := c.Request.Body.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		util.PanicOnIoError(err, "Failed to read body")

		_, err = w.Write(buf[:n])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	c.Data(http.StatusOK, "application/json", []byte("{\"tempName\": \""+filepath.Base(w.Name())+"\"}"))
}

// AddUploadedFile is a rest-api handler to add one file previously uploaded (see ReceiveUpload)
func AddUploadedFile(c *gin.Context) {

	var req AddUploadedFileRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	path := filepath.Join(config.AssetStorageTempDir, req.TempName)

	//Add file to storage
	infos, err := storage.AddFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	for _, info := range infos {
		if info.IsNewFile || !config.SkipMetaDataIfExists {

			//Create/Update meta-data
			meta, err := metadata.AddMetaData(
				info.Hash,
				info.MimeType,
				req.Name,
				"",
				req.Owner,
				req.FileTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			//Create/Update meta-data-database
			err = metadata_db.AddMetaData(meta)
			if err != nil {
				fmt.Printf("Error adding meta-data to database '%s': %s\n", path, err)
			}
		}
	}

	util.LogError(os.Remove(path))

	c.JSON(http.StatusOK, req)
}
