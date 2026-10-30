package spa_server

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/c8121/asset-storage/internal/rest_api"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/shared/httproot"
	"github.com/gin-gonic/gin"
)

const (
	DefaultDocument = "index.html"
)

// CreateRoutes looks for all files/directories in config.SpaHttpRoot and creates static routes.
func CreateRoutes(router *gin.Engine) {

	subDirs := []string{
		httproot.AppDir,
		httproot.VueUI,
	}

	httpFS := http.FS(httproot.FS)

	router.NoRoute(rest_api.AuthRequiredHandler(func(c *gin.Context) {

		path := c.Request.URL.Path

		for _, subDir := range subDirs {

			var fspath string
			if path == "/" {
				fspath = filepath.Join(subDir, DefaultDocument)
				file, err := httpFS.Open(fspath)
				if err == nil {
					defer util.CloseOrLog(file)
					stat, err := file.Stat()
					if err == nil {
						http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), file)
					}
				}
			} else {
				fspath = filepath.Join(subDir, strings.TrimPrefix(path, "/"))
				if fspath != "" {
					stat, err := fs.Stat(httproot.FS, fspath)
					if err == nil {
						if stat.Mode().IsRegular() {
							c.FileFromFS(fspath, httpFS)
							return
						}
					}
				}
			}
		}
	}))

}
