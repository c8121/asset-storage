package restapi

import (
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {
	router.GET("/assets/:hash", GetAsset)

	router.GET("/assets/list", ListAssets)
	router.GET("/assets/list/:offset", ListAssets)
	router.GET("/assets/list/:offset/:count", ListAssets)
	router.GET("/assets/list/:offset/:count/:mimetype", ListAssets)

	router.GET("/assets/list/path/:path", ListAssets)
	router.GET("/assets/list/path/:path/:offset", ListAssets)
	router.GET("/assets/list/path/:path/:offset/:count", ListAssets)
	router.GET("/assets/list/path/:path/:offset/:count/:mimetype", ListAssets)

	router.GET("/assets/thumbnail/:hash", GetPreview)
	router.GET("/assets/metadata/:hash", GetMetaData)

	router.POST("/assets/upload", ReceiveUpload)
	router.POST("/assets/upload/add", AddUploadedFile)

	router.GET("/mimetypes/list", ListMimeTypes)
	router.GET("/pathitems/list", ListPathItems)
	router.GET("/pathitems/list/:parent", ListPathItems)
}
