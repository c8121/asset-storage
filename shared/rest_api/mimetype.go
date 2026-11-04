package rest_api

import (
	"net/http"

	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

func ListMimeTypes(c *gin.Context) {

	items, err := metadata_db_entity.ListMimeTypes()
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to load mime-types", err)
		return
	}

	if len(items) > 0 {
		c.IndentedJSON(http.StatusOK, items)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}
