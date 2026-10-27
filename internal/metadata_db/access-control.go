package metadata_db

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
)

type MetaDataDbAccessController struct {
}

func (ctrl MetaDataDbAccessController) CanRead(assetHash string, ctx *access_control.AccessContext) error {

	asset, err := metadata_db_entity.LoadAsset(assetHash)
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

	fmt.Printf("Grant permission to asset %d, ctx=%v\n", assetId, ctx)

	return nil
}

func (ctrl MetaDataDbAccessController) CanWritePath(path string, ctx *access_control.AccessContext) error {
	fmt.Printf("Grant permission to path %s, ctx=%v\n", path, ctx)
	return nil
}

func InitAccessControl() {

	var ctrl = MetaDataDbAccessController{}
	access_control.SetAccessController(ctrl)

}
