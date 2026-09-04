package faces

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/go-resty/resty/v2"
)

var (
	CheckFaceRestServiceEndpoint = "http://localhost:8000/status"
	GetFaceRestServiceEndpoint   = "http://localhost:8000/extract"
)

type (
	RestApiFace struct {
		Index     int
		Embedding []float32
		Image     []int
	}

	RestApiFacesResponse struct {
		Faces []RestApiFace
	}
)

const (
	FilePermissions = 0744
)

// init creates required directories
func init() {
	util.CreateDirIfNotExists(config.AssetFacesBaseDir, FilePermissions)
}

// GetFaces finds faces in image, returns ["name",...]
func GetFaces(sourceHash string) (*[]RestApiFace, error) {

	faces, err := restExecExtractFaces(sourceHash)
	if err != nil {
		return nil, err
	}

	return faces, nil
}

// restExecExtractFaces calls REST-Service, see services/insightface/service.py
func restExecExtractFaces(sourceHash string) (*[]RestApiFace, error) {

	meta, err := metadata.LoadByHash(sourceHash)
	if err != nil {
		return nil, err
	}

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "image/") {
		return nil, fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	reader, err := storage.Open(sourceHash)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(reader)

	client := resty.New()
	response, err := client.R().
		SetFileReader("file", meta.Origins[0].Name, reader).
		Post(GetFaceRestServiceEndpoint)
	if err != nil {
		return nil, err
	}

	result := &RestApiFacesResponse{}
	err = json.Unmarshal(response.Body(), result)
	if err != nil {
		return nil, err
	}

	return &result.Faces, nil
}
