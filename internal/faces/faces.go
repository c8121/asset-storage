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
	Embedding []float32

	Face struct {
		Index     int
		Embedding Embedding
		Image     []int
	}

	Faces struct {
		Faces []Face
	}
)

const (
	FilePermissions = 0744
)

// Init creates required directories
func Init() {
	util.CreateDirIfNotExists(config.AssetFacesBaseDir, FilePermissions)
}

// GetFaces finds faces in image, returns ["name",...]
func GetFaces(sourceHash string) (*Faces, error) {

	faces, err := restExecExtractFaces(sourceHash)
	if err != nil {
		return nil, err
	}

	return faces, nil
}

// restExecExtractFaces calls REST-Service, see services/insightface/service.py
func restExecExtractFaces(sourceHash string) (*Faces, error) {

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

	result := &Faces{}
	err = json.Unmarshal(response.Body(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
