package rest_api

import (
	"net/http"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, util.AppNotifications.Messages)
}
