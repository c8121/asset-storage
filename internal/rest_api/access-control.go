package rest_api

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/users"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

type AccessContext struct {
	User *users.AuthenticatedUser
}

func GetAccessContext(c *gin.Context) *AccessContext {
	user, authenticated := c.Get("user")
	if !authenticated {
		return nil
	}
	return &AccessContext{User: user.(*users.AuthenticatedUser)}
}

func CanRead(assetId int64, ctx *AccessContext) bool {

	if ctx == nil {
		return true
	}
	if ctx.User == nil {
		return false
	}

	fmt.Printf("Checking permission for asset %d, ctx=%s\n", assetId, ctx.User)

	return true
}

func CanReadHash(hash string, ctx *AccessContext) bool {

	asset, err := metadata_db_entity.LoadAsset(metadata_db_conn.GetDatabase(), hash)
	if err != nil || asset == nil {
		util.LogError(err)
		return false
	}

	if !CanRead(asset.Id, ctx) {
		fmt.Printf("Access denied to %s\n", hash)
		return false
	}

	return true
}
