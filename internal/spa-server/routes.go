package spaserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	HttpRoot = filepath.Dir(os.Args[0]) + "/vue-ui"
)

func CreateRoutes(router *gin.Engine) {

	if len(HttpRoot) == 0 {
		panic("spaserver.HttpRoot not set")
	}

	//Development: When using "go run...", path must be set manually
	if strings.Contains(HttpRoot, "go-build") {
		HttpRoot = "/home/christianh/Workspace/go/asset-storage/vue-ui"
	}

	dir, err := os.ReadDir(HttpRoot)
	if err != nil {
		panic(fmt.Errorf("failed to read spaserver.HttpRoot: %s", HttpRoot))
	}

	for _, e := range dir {
		if strings.HasPrefix(e.Name(), "index.") {
			router.StaticFile("/", filepath.Join(HttpRoot, e.Name()))
		} else if !e.IsDir() {
			router.StaticFile("/"+e.Name(), filepath.Join(HttpRoot, e.Name()))
		} else {
			router.Static("/"+e.Name(), filepath.Join(HttpRoot, e.Name()))
		}
	}
}
