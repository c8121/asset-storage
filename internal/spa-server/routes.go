package spaserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/gin-gonic/gin"
)

// CreateRoutes looks for all files/directories in config.SpaHttpRoot and creates static routes.
func CreateRoutes(router *gin.Engine) {

	if len(config.SpaHttpRoot) == 0 {
		panic("spaserver.HttpRoot not set")
	}

	dir, err := os.ReadDir(config.SpaHttpRoot)
	if err != nil {
		panic(fmt.Errorf("failed to read spaserver.HttpRoot: %s", config.SpaHttpRoot))
	}

	for _, e := range dir {
		if strings.HasPrefix(e.Name(), "index.") {
			router.StaticFile("/", filepath.Join(config.SpaHttpRoot, e.Name()))
		} else if !e.IsDir() {
			router.StaticFile("/"+e.Name(), filepath.Join(config.SpaHttpRoot, e.Name()))
		} else {
			router.Static("/"+e.Name(), filepath.Join(config.SpaHttpRoot, e.Name()))
		}
	}
}
