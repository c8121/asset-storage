package filter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	shell_command "github.com/c8121/asset-storage/internal/shell-command"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type FFmpegFilter struct {
	DefaultWidth           string
	DefaultVideoPosition   string
	DefaultFileNamePattern string
	DefaultMimeType        string
}

func NewFFmpegFilter() *FFmpegFilter {
	f := &FFmpegFilter{}
	f.DefaultWidth = "400"
	f.DefaultVideoPosition = "00:00:01"
	f.DefaultFileNamePattern = "asset-thumb*.png"
	f.DefaultMimeType = "image/png"
	return f
}

func (f FFmpegFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "video/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	width, _ := strconv.Atoi(util.GetOrDefault(params, "width", f.DefaultWidth))
	videoPosition := util.GetOrDefault(params, "videoPosition", f.DefaultVideoPosition)
	tempFileNamePattern := util.GetOrDefault(params, "fileNamePattern", f.DefaultFileNamePattern)
	mimeType := util.GetOrDefault(params, "mimeType", f.DefaultMimeType)

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, tempFileNamePattern)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	err = shell_command.FFmpegThumb(in, out.Name(), width, -1, videoPosition)
	if err != nil {
		util.LogError(os.Remove(out.Name()))
		return nil, "", fmt.Errorf("Failed to create thumbnail: %w", err)
	}

	bytes, err := os.ReadFile(out.Name())
	if err != nil {
		util.LogError(os.Remove(out.Name()))
		return nil, "", fmt.Errorf("Failed to read thumbnail: %w", err)
	}

	util.LogError(os.Remove(out.Name()))
	return bytes, mimeType, nil
}
