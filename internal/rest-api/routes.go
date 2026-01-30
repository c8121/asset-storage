package restapi

import (
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {
	router.GET("/assets/:hash", GetAsset)

	router.POST("/assets/list", ListAssets)

	router.GET("/assets/thumbnail/:hash", GetPreview)
	router.GET("/assets/metadata/:hash", GetMetaData)

	router.GET("/assets/filter/:filter/:hash", GetFiltered)
	router.POST("/assets/filter/:filter/:hash", GetFiltered)

	router.POST("/assets/upload", ReceiveUpload)
	router.POST("/assets/upload/add", AddUploadedFile)

	router.GET("/collections/:hash", GetCollection)

	router.POST("/collections/add", AddCollection)

	router.GET("/mimetypes/list", ListMimeTypes)
	router.GET("/pathitems/list", ListPathItems)
	router.GET("/pathitems/list/:parent", ListPathItems)
}
