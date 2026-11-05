package rest_api

import (
	auth "github.com/c8121/go_vue_spa/auth/gin"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {
	router.GET("/assets/:hash", auth.UserAuthRequiredHandler(GetAsset))

	router.POST("/assets/list", auth.UserAuthRequiredHandler(ListAssets))

	router.GET("/assets/thumbnail/:hash", auth.UserAuthRequiredHandler(GetPreview))
	router.GET("/assets/cropped/:hash/:width/:x1/:y1/:x2/:y2", auth.UserAuthRequiredHandler(GetCoppedPreview))
	router.GET("/assets/metadata/:hash", auth.UserAuthRequiredHandler(GetMetaData))

	router.GET("/assets/filter/:filter/:hash", auth.UserAuthRequiredHandler(GetFiltered))

	router.POST("/assets/upload", auth.UserAuthRequiredHandler(ReceiveUpload))
	router.POST("/assets/upload/add", auth.UserAuthRequiredHandler(AddUploadedFile))

	router.GET("/collections/:uuid", auth.UserAuthRequiredHandler(GetCollection))

	router.POST("/collections/list", auth.UserAuthRequiredHandler(ListCollections))

	router.POST("/collections/add", auth.UserAuthRequiredHandler(AddCollection))
	router.POST("/collections/remove", auth.UserAuthRequiredHandler(RemoveCollection))

	router.GET("/mimetypes/list", auth.UserAuthRequiredHandler(ListMimeTypes))
	router.GET("/pathitems/list", auth.UserAuthRequiredHandler(ListPathItems))
	router.GET("/pathitems/list/:parent", auth.UserAuthRequiredHandler(ListPathItems))
	router.POST("/pathitem/", auth.UserAuthRequiredHandler(FindPathItem))

	router.GET("/notifications", auth.UserAuthRequiredHandler(GetNotifications))

	router.GET("/faces/:hash", auth.UserAuthRequiredHandler(GetFaces))
	router.GET("/faces/detect/:hash", auth.UserAuthRequiredHandler(DetectFaces))
}
