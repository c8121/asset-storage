package rest_api

import (
	"net/http"

	"github.com/c8121/go_vue_spa/notifications"
	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, notifications.AppNotifications.Messages)
}
