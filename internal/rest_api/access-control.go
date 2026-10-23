package rest_api

import (
	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/users"
	"github.com/gin-gonic/gin"
)

func GetAccessContext(c *gin.Context) *access_control.AccessContext {
	user, authenticated := c.Get("user")
	if !authenticated {
		return nil
	}
	return &access_control.AccessContext{User: user.(*users.AuthenticatedUser)}
}
