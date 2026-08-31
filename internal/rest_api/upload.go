package rest_api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

type (
	AddUploadedFileRequest struct {
		TempName string
		Name     string
		Path     string
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

	var list = make([]metadata.JsonAssetMetaData, 0)

	for _, info := range infos {
		if info.IsNewFile || !config.SkipMetaDataIfExists {

			addedFileName := filepath.Base(info.SourcePath)
			if addedFileName == req.TempName {
				addedFileName = req.Name
			} else {
				//Added file was unpackable and contained multiple other files
				fmt.Printf("Unpacked %s from %s\n", addedFileName, req.Name)
			}

			//Create/Update meta-data
			meta, err := metadata.AddMetaData(
				info.Hash,
				info.MimeType,
				addedFileName,
				req.Path,
				req.Owner,
				req.FileTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			list = append(list, *meta)

			//Create/Update meta-data-database
			err = metadata_db_entity.AddMetaData(meta)
			if err != nil {
				fmt.Printf("Error adding meta-data to database '%s': %s\n", path, err)
			}
		}
	}

	util.LogError(os.Remove(path))

	c.JSON(http.StatusOK, list)
}
