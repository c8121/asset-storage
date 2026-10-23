package ssh_server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type SshHandler interface {
	GetUsername() string
	GetNewFiles() []SshFileInfo
}

type SshFileInfo struct {
	LocalPath string //Absolute path on local file system
	UserPath  string //Path as visible to the user (child of users virtual-root)
}

// AddFilesToArchive adds file to storage and creates meta-data
func AddFilesToArchive(h SshHandler) {

	tx, err := metadata_db_conn.BeginTransaction()
	if err != nil {
		util.LogError(err)
		tx = nil
	}

	files := h.GetNewFiles()
	for _, file := range files {

		fmt.Printf("Add file %s to archive\n", file.UserPath)

		stat, err := os.Stat(file.LocalPath)
		if err != nil {
			fmt.Printf("Cannot get file stat '%s'\n", file)
			continue
		}

		path := filepath.Dir(file.UserPath)
		acc := GetAccessContext(h)
		ctx := access_control.GetAccessController()
		if ctx != nil && ctx.CanWritePath(path, acc) != nil {
			fmt.Printf("Access dinied to path %s\n", path)
			continue
		}

		//Add file to storage
		infos, err := storage.AddFile(file.LocalPath)
		if err != nil {
			fmt.Printf("Error adding '%s': %s\n", file.UserPath, err)
			continue
		}

		for _, info := range infos {
			if info.IsNewFile || !config.SkipMetaDataIfExists {
				//Create/Update meta-data
				meta, err := metadata.AddMetaData(
					info.Hash,
					info.MimeType,
					filepath.Base(info.SourcePath),
					path,
					h.GetUsername(),
					stat.ModTime())
				if err != nil {
					fmt.Printf("Error adding meta-data '%s': %s\n", file.LocalPath, err)
					continue
				}

				//Create/Update meta-data-database
				if tx != nil {
					err = metadata_db_entity.AddAsset(tx, meta)
					if err != nil {
						fmt.Printf("Error adding meta-data to database '%s': %s\n", file, err)
					}
				}
			}
		}
	}

	if tx != nil {
		util.LogError(metadata_db_conn.CommitOrLog(tx))
	}

}
