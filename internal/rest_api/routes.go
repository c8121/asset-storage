package rest_api

import (
	"github.com/c8121/asset-storage/internal/users"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {
	router.GET("/assets/:hash", users.AuthRequiredHandler(GetAsset))

	router.POST("/assets/list", users.AuthRequiredHandler(ListAssets))

	router.GET("/assets/thumbnail/:hash", users.AuthRequiredHandler(GetPreview))
	router.GET("/assets/cropped/:hash/:width/:x1/:y1/:x2/:y2", users.AuthRequiredHandler(GetCoppedPreview))
	router.GET("/assets/metadata/:hash", users.AuthRequiredHandler(GetMetaData))

	router.GET("/assets/filter/:filter/:hash", users.AuthRequiredHandler(GetFiltered))
	router.POST("/assets/filter/:filter/:hash", users.AuthRequiredHandler(GetFiltered))

	router.POST("/assets/upload", users.AuthRequiredHandler(ReceiveUpload))
	router.POST("/assets/upload/add", users.AuthRequiredHandler(AddUploadedFile))

	router.GET("/collections/:hash", users.AuthRequiredHandler(GetCollection))

	router.POST("/collections/list", users.AuthRequiredHandler(ListCollections))

	router.POST("/collections/add", users.AuthRequiredHandler(AddCollection))

	router.GET("/mimetypes/list", users.AuthRequiredHandler(ListMimeTypes))
	router.GET("/pathitems/list", users.AuthRequiredHandler(ListPathItems))
	router.GET("/pathitems/list/:parent", users.AuthRequiredHandler(ListPathItems))
	router.POST("/pathitem/", users.AuthRequiredHandler(FindPathItem))

	router.GET("/notifications", users.AuthRequiredHandler(GetNotifications))

	router.GET("/faces/:hash", users.AuthRequiredHandler(GetFaces))
	router.GET("/faces/detect/:hash", users.AuthRequiredHandler(DetectFaces))
}
