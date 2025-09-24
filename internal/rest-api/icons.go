package restapi

import (
	"embed"
	"strings"
)

var (
	//go:embed icons/file-regular.svg
	DefaultFileThumbnail         []byte
	DefaultFileThumbnailMimeType = "image/svg+xml"

	ByMimeTypeThumbnailExt      = ".svg"
	ByMimeTypeThumbnailMimeType = "image/svg+xml"

	//go:embed icons/*.svg
	ByMimeTypeThumbnails embed.FS
)

// loadIconByMimeType returns an icon matching the mime-type: icons/type[-subtype].svg
func loadIconByMimeType(mimeType string) ([]byte, string, error) {

	thumb := strings.Replace(strings.ToLower(mimeType), "/", "-", 1) + ByMimeTypeThumbnailExt
	if bytes, err := ByMimeTypeThumbnails.ReadFile("icons/" + thumb); err == nil {
		return bytes, DefaultFileThumbnailMimeType, nil
	}

	thumb = strings.ToLower(mimeType)
	if p := strings.Index(thumb, "/"); p > 1 {
		thumb = thumb[:p] + ByMimeTypeThumbnailExt
	} else {
		thumb += ByMimeTypeThumbnailExt
	}

	if bytes, err := ByMimeTypeThumbnails.ReadFile("icons/" + thumb); err == nil {
		return bytes, DefaultFileThumbnailMimeType, nil
	} else {
		return nil, "", err
	}
}

// loadDefaultIcon returns the default icon & mime-type
func loadDefaultIcon() ([]byte, string, error) {
	return DefaultFileThumbnail, DefaultFileThumbnailMimeType, nil
}
