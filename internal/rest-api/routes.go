package restapi

import (
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {
	router.GET("/assets/:hash", GetAsset)
	router.GET("/assets/list", ListAssets)
	router.GET("/assets/list/:offset", ListAssets)
	router.GET("/assets/list/:offset/:count", ListAssets)
	router.GET("/assets/thumbnail/:hash", GetPreview)
}
