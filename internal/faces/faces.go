package faces

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/go-resty/resty/v2"
)

var (
	GetFaceRestServiceEndpoint = "http://localhost:8000/extract"
	FaceImageExtension         = ".jpg"
	FaceEmbeddingExtension     = ".json"
)

type (
	Embedding []float64

	Face struct {
		Index     int
		Embedding Embedding
		Image     string
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

// GetFaceImage finds face in image, returns face image
func GetFaceImage(sourceHash string, idx int) ([]byte, error) {

	facesDir := getFacesDir(sourceHash, false)

	existingFaces, err := loadFaces(sourceHash)
	if err == nil {
		if idx >= len(existingFaces) {
			return nil, fmt.Errorf("index out of bounds")
		}
		fileName := filepath.Join(facesDir, existingFaces[idx]+FaceImageExtension)
		return os.ReadFile(fileName)
	}

	faces, err := GetFaces(sourceHash)
	if err != nil {
		return nil, err
	}
	if idx >= len(faces) {
		return nil, fmt.Errorf("index out of bounds")
	}
	fileName := filepath.Join(facesDir, faces[idx]+FaceImageExtension)
	return os.ReadFile(fileName)
}

// GetFaces finds faces in image, returns ["name",...]
func GetFaces(sourceHash string) ([]string, error) {

	existingFaces, err := loadFaces(sourceHash)
	if err == nil {
		return existingFaces, nil
	}
	fmt.Printf("Faces not found, will use REST-Service: %s\n", err)

	faces, err := restExecGetFaces(sourceHash)
	if err != nil {
		return nil, err
	}

	facesDir := getFacesDir(sourceHash, true)

	list := make([]string, 0)

	for idx, face := range faces.Faces {

		dec, err := base64.StdEncoding.DecodeString(face.Image)
		if err != nil {
			fmt.Printf("Failed to decode image\n")
			continue
		}

		imageFileName := filepath.Join(facesDir, strconv.Itoa(idx)+FaceImageExtension)
		err = os.WriteFile(imageFileName, dec, FilePermissions)
		if err != nil {
			fmt.Printf("Failed to write: %s\n", err)
		}

		embeddingFileName := filepath.Join(facesDir, strconv.Itoa(idx)+FaceEmbeddingExtension)
		embeddingJson, err := json.Marshal(face.Embedding)
		if err != nil {
			fmt.Printf("Failed to create json: %s\n", err)
		} else {
			err = os.WriteFile(embeddingFileName, embeddingJson, FilePermissions)
			if err != nil {
				fmt.Printf("Failed to write: %s\n", err)
			}
		}

		list = append(list, strconv.Itoa(idx))
	}

	return list, nil
}

// loadFaces loads previously created faces, if exists. Returns ["name",...]
func loadFaces(sourceHash string) ([]string, error) {

	facesDir := getFacesDir(sourceHash, false)
	stat, err := os.Stat(facesDir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("not a directory: %s\n", facesDir)
	}

	entries, err := os.ReadDir(facesDir)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)

	for _, e := range entries {
		if strings.HasSuffix(e.Name(), FaceImageExtension) {
			list = append(list, e.Name()[0:len(e.Name())-len(FaceImageExtension)])
		}
	}

	return list, nil
}

// restExecGetFaces calls REST-Service, see services/insightface/service.py
func restExecGetFaces(sourceHash string) (*Faces, error) {

	meta, err := metadata.LoadByHash(sourceHash)
	if err != nil {
		return nil, err
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

// Create dir: AssetFacesBaseDir/Hash[:2]/Hash[2:]
func getFacesDir(sourceHash string, create bool) string {

	dir := filepath.Join(config.AssetFacesBaseDir, sourceHash[0:2], sourceHash[2:])
	if create {
		util.CreateDirIfNotExists(dir, FilePermissions)
	}
	return dir
}
