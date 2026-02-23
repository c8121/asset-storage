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
		face := existingFaces.getFace(idx)
		if face == nil {
			return nil, fmt.Errorf("index out of bounds")
		}
		fileName := filepath.Join(facesDir, strconv.Itoa(face.Index)+FaceImageExtension)
		return os.ReadFile(fileName)
	}

	faces, err := GetFaces(sourceHash)
	if err != nil {
		return nil, err
	}
	face := faces.getFace(idx)
	if face == nil {
		return nil, fmt.Errorf("index out of bounds")
	}
	fileName := filepath.Join(facesDir, strconv.Itoa(face.Index)+FaceImageExtension)
	return os.ReadFile(fileName)
}

// GetFaces finds faces in image, returns ["name",...]
func GetFaces(sourceHash string) (*Faces, error) {

	existingFaces, err := loadFaces(sourceHash)
	if err == nil {
		return existingFaces, nil
	}
	fmt.Printf("Faces not found, will use REST-Service: %s\n", err)

	faces, err := restExecExtractFaces(sourceHash)
	if err != nil {
		return nil, err
	}

	saveExtractedFaces(sourceHash, faces)

	return faces, nil
}

// Find face by index
func (faces *Faces) getFace(idx int) *Face {
	for _, face := range faces.Faces {
		if face.Index == idx {
			return &face
		}
	}
	return nil
}

// loadFaces loads previously created faces, if exists. Returns ["name",...]
func loadFaces(sourceHash string) (*Faces, error) {

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

	faces := &Faces{
		Faces: make([]Face, 0),
	}

	for _, e := range entries {
		if strings.HasSuffix(e.Name(), FaceImageExtension) {
			idx, err := strconv.Atoi(e.Name()[0 : len(e.Name())-len(FaceImageExtension)])
			if err != nil {
				fmt.Printf("Failed to parse index: %s\n", err)
				continue
			}
			face := &Face{
				Index:     idx,
				Embedding: nil,
				Image:     "",
			}
			faces.Faces = append(faces.Faces, *face)
		}
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

// Save extracted faces to disk
func saveExtractedFaces(sourceHash string, faces *Faces) {

	if len(faces.Faces) == 0 {
		return
	}

	facesDir := getFacesDir(sourceHash, true)

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
	}
}

// Create dir: AssetFacesBaseDir/Hash[:2]/Hash[2:]
func getFacesDir(sourceHash string, create bool) string {

	dir := filepath.Join(config.AssetFacesBaseDir, sourceHash[0:2], sourceHash[2:])
	if create {
		util.CreateDirIfNotExists(dir, FilePermissions)
	}
	return dir
}
