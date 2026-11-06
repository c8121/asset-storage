package metadata_db

import (
	"log/slog"

	"github.com/c8121/asset-storage/internal/access_control"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_vue_spa/util"
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

	slog.Info("Grant permission", "asset", assetId, "user", ctx.User)

	return nil
}

func (ctrl MetaDataDbAccessController) CanWritePath(path string, ctx *access_control.AccessContext) error {
	slog.Info("Grant permission", "path", path, "user", ctx.User)
	return nil
}

func InitAccessControl() {

	var ctrl = MetaDataDbAccessController{}
	access_control.SetAccessController(ctrl)

}
