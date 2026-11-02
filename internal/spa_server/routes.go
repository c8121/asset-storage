package spa_server

import (
	"io/fs"
	"os"

	"github.com/c8121/asset-storage/internal/config"
	apphttproot "github.com/c8121/asset-storage/shared/httproot"
	"github.com/c8121/go_vue_spa"
	"github.com/c8121/go_vue_spa/httproot"
	"github.com/gin-gonic/gin"
)

// CreateRoutes looks for all files/directories in config.SpaHttpRoot and creates static routes.
func CreateRoutes(router *gin.Engine) {

	//APP resources
	var rootFS fs.FS
	if config.SpaHttpRoot != "" {
		rootFS = os.DirFS(config.SpaHttpRoot)
	} else {
		rootFS = apphttproot.FS
	}

	resources := make([]go_vue_spa.Resource, 0)
	resources = append(resources, go_vue_spa.Resource{
		RootFs:     httproot.FS,
		SubDirName: httproot.VueUI,
	})
	resources = append(resources, go_vue_spa.Resource{
		RootFs:          rootFS,
		SubDirName:      apphttproot.AppDir,
		DefaultDocument: "asset-storage.html",
	})

	go_vue_spa.CreateRoutesFromFS(router, resources)
}
