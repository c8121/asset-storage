package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

func GetFaceImage(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	idx, err := strconv.Atoi(c.Param("idx"))
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid index")))
		return
	}

	face, err := faces.GetFaceImage(hash, idx)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, err))
		return
	}

	c.Data(http.StatusOK, "image/jpeg", face)
}

func GetFaces(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	faces, err := faces.GetFaces(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	json, err := json.Marshal(faces)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	c.Data(http.StatusOK, "application/json", json)
}
