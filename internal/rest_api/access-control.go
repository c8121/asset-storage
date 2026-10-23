package rest_api

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/users"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

type MetaDataDbAccessController struct {
}

func (ctrl MetaDataDbAccessController) CanRead(assetHash string, ctx *access_control.AccessContext) error {

	asset, err := metadata_db_entity.LoadAsset(metadata_db_conn.GetDatabase(), assetHash)
	if err != nil || asset == nil {
		util.LogError(err)
		return access_control.ErrAccessDenied
	}

	return ctrl.CanReadId(asset.Id, ctx)
}

func (ctrl MetaDataDbAccessController) CanReadId(assetId int64, ctx *access_control.AccessContext) error {
	if ctx == nil {
		return nil
	}
	if ctx.User == nil {
		return access_control.ErrAccessDenied
	}

	fmt.Printf("Checking permission for asset %d, ctx=%s\n", assetId, ctx.User)

	return nil
}

func InitAccessControl() {

	var ctrl = MetaDataDbAccessController{}
	access_control.SetAccessController(ctrl)

}

func GetAccessContext(c *gin.Context) *access_control.AccessContext {
	user, authenticated := c.Get("user")
	if !authenticated {
		return nil
	}
	return &access_control.AccessContext{User: user.(*users.AuthenticatedUser)}
}
