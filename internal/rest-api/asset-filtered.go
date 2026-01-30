package restapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// GetFiltered is a rest-api handler to filter/convert an asset
func GetFiltered(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	filterName := c.Param("filter")
	if len(filterName) < 1 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("no filter name given")))
		return
	}

	filterParamsReader := c.Request.Body
	b, err := io.ReadAll(filterParamsReader)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to read request body")))
		return
	}
	filterParams := paramsToMap(string(b))

	meta, err := metadata.LoadByHash(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash (not found)")))
		return
	}

	if bytes, mimeType, err := filterAsset(hash, meta, filterName, filterParams); err == nil {
		c.Data(http.StatusOK, mimeType, bytes)
		return
	} else {
		util.LogError(c.AbortWithError(http.StatusNotFound, err))
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

// paramsToMap splits query parameters: name=value&...
func paramsToMap(s string) map[string]string {

	m := make(map[string]string)

	p := 0
	for p > -1 {
		nvp := ""
		e := strings.Index(s[p:], "&")
		if e > -1 {
			nvp = s[p:e]
			p = e + 1
		} else {
			nvp = s[p:]
			p = -1
		}

		k := ""
		v := ""
		var err error
		i := strings.Index(nvp, "=")
		if i > -1 {
			k, err = url.QueryUnescape(nvp[:i])
			v, err = url.QueryUnescape(nvp[i+1:])
		} else {
			k = nvp
			v = ""
		}
		if err == nil {
			//fmt.Printf("%s %s: %s = '%s'\n", s, nvp, k, v)
			m[k] = v
		} else {
			util.LogError(fmt.Errorf("failed to decode parameter: %s, %w\n", nvp, err))
		}
	}

	return m
}
