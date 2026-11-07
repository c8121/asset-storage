package rest_api

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

// GetFiltered is a rest-api handler to filter/convert an asset
func GetFiltered(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid hash", nil)
		return
	}

	acc := GetAccessContext(c)
	ctx := access_control.GetAccessController()
	if ctx != nil && ctx.CanRead(hash, acc) != nil {
		util.AbortWithStatus(c, http.StatusForbidden, "forbidden", nil)
		return
	}

	filterParam := c.Param("filter")
	if len(filterParam) < 1 {
		util.AbortWithStatus(c, http.StatusNotFound, "no filter name given", nil)
		return
	}
	filterName, filterParams := parseFilterParam(filterParam)

	meta, err := metadata.LoadByHash(hash)
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid hash", err)
		return
	}

	if bytes, mimeType, err := filterAsset(hash, meta, filterName, filterParams); err == nil {
		c.Data(http.StatusOK, mimeType, bytes)
		return
	} else {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to get filtered asset", err)
	}
}

// filterAsset converts/filters/modified an asset an returns the filtered content
// Returns content, mimeType, error
func filterAsset(assetHash string, meta *metadata.JsonAssetMetaData, filterName string, filterParams map[string]string) ([]byte, string, error) {

	var f = filter.GetFirstFilterByNameAndMimeType(filterName, meta.MimeType)
	if f == nil {
		return nil, "", fmt.Errorf("no filter with name '%s' available for mime-type: %s", filterName, meta.MimeType)
	}

	return f.Apply(assetHash, meta, filterParams)

}

// parseFilterParam splits filter parameter: filterName:name|param=value|...
func parseFilterParam(s string) (string, map[string]string) {

	filterName := ""
	params := make(map[string]string)

	p := 0
	for p > -1 {
		nvp := ""
		e := strings.Index(s[p:], "|")
		if e > p {
			nvp = s[p:e]
			p = e + 1
		} else {
			nvp = s[p:]
			p = -1
		}

		k := ""
		v := ""
		var err error
		i := strings.Index(nvp, ":")
		if i > -1 {
			k, err = url.QueryUnescape(nvp[:i])
			v, err = url.QueryUnescape(nvp[i+1:])
		} else {
			k = nvp
			v = ""
		}
		if err == nil {
			if k == "filterName" {
				filterName = v
			} else {
				params[k] = v
			}
		} else {
			slog.Warn("Invalid filter string", "nvp", nvp, "error", err)
		}
	}

	slog.Info("Filter", "name", filterName, "params", params)
	return filterName, params
}
