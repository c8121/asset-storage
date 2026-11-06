package rest_api

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_vue_spa/util"
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
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to create writer", err)
		return
	}
	defer util.CloseOrLog(w)

	buf := make([]byte, 4096)
	for {
		n, err := c.Request.Body.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			util.AboutWithStatus(c, http.StatusInternalServerError, "failed to read", err)
			return
		}

		_, err = w.Write(buf[:n])
		if err != nil {
			util.AboutWithStatus(c, http.StatusInternalServerError, "failed to write", err)
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
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to read request", err)
		return
	}

	acc := GetAccessContext(c)
	if acc != nil && acc.User != nil && len(acc.User.Username) > 0 {
		req.Owner = acc.User.Username
	}
	ctx := access_control.GetAccessController()
	if ctx != nil && ctx.CanWritePath(req.Path, acc) != nil {
		util.AboutWithStatus(c, http.StatusForbidden, "forbidden", nil)
		return
	}

	if strings.Contains(req.TempName, "..") || strings.ContainsAny(req.TempName, "/\\") {
		util.AboutWithStatus(c, http.StatusBadRequest, "invalid request", nil)
		return
	}

	path := filepath.Join(config.AppConfig.AssetStorageTempDir, req.TempName)

	//Add file to storage
	infos, err := storage.AddFile(path)
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to add file", err)
		return
	}

	var list = make([]metadata.JsonAssetMetaData, 0)

	for _, info := range infos {
		if info.IsNewFile || !config.AppConfig.SkipMetaDataIfExists {

			addedFileName := filepath.Base(info.SourcePath)
			if addedFileName == req.TempName {
				addedFileName = req.Name
			} else {
				//Added file was unpackable and contained multiple other files
				slog.Info("Unpacked", "file", addedFileName, "request", req.Name)
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
				util.AboutWithStatus(c, http.StatusInternalServerError, "failed to add meta-data", err)
				return
			}

			list = append(list, *meta)

			//Create/Update meta-data-database
			err = metadata_db_entity.AddAsset(meta)
			if err != nil {
				slog.Warn("Error adding meta-data to database", "path", path, "error", err)
			}

		}
	}

	util.LogError(os.Remove(path))

	c.JSON(http.StatusOK, list)
}
