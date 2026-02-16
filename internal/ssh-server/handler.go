package ssh_server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	"github.com/c8121/asset-storage/internal/storage"
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

	files := h.GetNewFiles()
	for _, file := range files {

		fmt.Printf("Adding file %s to archive\n", file)
		fmt.Printf("Adding file %s to archive\n", filepath.Base(file.LocalPath))
		fmt.Printf("Adding file %s to archive\n", filepath.Dir(file.UserPath))

		stat, err := os.Stat(file.LocalPath)
		if err != nil {
			fmt.Printf("Cannot get file stat '%s'\n", file)
			continue
		}

		//Add file to storage
		infos, err := storage.AddFile(file.LocalPath)
		if err != nil {
			fmt.Printf("Error adding '%s': %s\n", file, err)
			continue
		}

		for _, info := range infos {
			if info.IsNewFile || !config.SkipMetaDataIfExists {
				//Create/Update meta-data
				meta, err := metadata.AddMetaData(
					info.Hash,
					info.MimeType,
					filepath.Base(info.SourcePath),
					filepath.Dir(file.UserPath),
					h.GetUsername(),
					stat.ModTime())
				if err != nil {
					fmt.Printf("Error adding meta-data '%s': %s\n", file, err)
					continue
				}

				//Create/Update meta-data-database
				err = metadata_db.AddMetaData(meta)
				if err != nil {
					fmt.Printf("Error adding meta-data to database '%s': %s\n", file, err)
				}
			}
		}
	}

}
