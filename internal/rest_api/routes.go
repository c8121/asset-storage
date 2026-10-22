package rest_api

import (
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {
	router.GET("/assets/:hash", AuthRequiredHandler(GetAsset))

	router.POST("/assets/list", AuthRequiredHandler(ListAssets))

	router.GET("/assets/thumbnail/:hash", AuthRequiredHandler(GetPreview))
	router.GET("/assets/cropped/:hash/:width/:x1/:y1/:x2/:y2", AuthRequiredHandler(GetCoppedPreview))
	router.GET("/assets/metadata/:hash", AuthRequiredHandler(GetMetaData))

	router.GET("/assets/filter/:filter/:hash", AuthRequiredHandler(GetFiltered))
	router.POST("/assets/filter/:filter/:hash", AuthRequiredHandler(GetFiltered))

	router.POST("/assets/upload", AuthRequiredHandler(ReceiveUpload))
	router.POST("/assets/upload/add", AuthRequiredHandler(AddUploadedFile))

	router.GET("/collections/:uuid", AuthRequiredHandler(GetCollection))

	router.POST("/collections/list", AuthRequiredHandler(ListCollections))

	router.POST("/collections/add", AuthRequiredHandler(AddCollection))
	router.POST("/collections/remove", AuthRequiredHandler(RemoveCollection))

	router.GET("/mimetypes/list", AuthRequiredHandler(ListMimeTypes))
	router.GET("/pathitems/list", AuthRequiredHandler(ListPathItems))
	router.GET("/pathitems/list/:parent", AuthRequiredHandler(ListPathItems))
	router.POST("/pathitem/", AuthRequiredHandler(FindPathItem))

	router.GET("/notifications", AuthRequiredHandler(GetNotifications))

	router.GET("/faces/:hash", AuthRequiredHandler(GetFaces))
	router.GET("/faces/detect/:hash", AuthRequiredHandler(DetectFaces))
}
