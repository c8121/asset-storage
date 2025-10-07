package thumbnails

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	FFmpegBinPaths = []string{
		"/opt/ffmpeg*/bin/ffmpeg.exe",
	}

	FFmpegBinPath      = ""
	ThumbVideoPosition = "00:00:01"
)

// FFmpegThumb executes ffmpeg
func FFmpegThumb(input string, output string, width int, height int) error {

	binary, err := FindFFmpegBin()
	if binary == "" || err != nil {
		return fmt.Errorf("FFmpeg not found (searching in %v): %v", FFmpegBinPaths, err)
	}

	var args []string

	if ThumbVideoPosition != "" {
		args = append(args, "-ss", ThumbVideoPosition)
	}

	args = append(args, "-y") //Overwrite
	args = append(args, "-i", input)

	args = append(args, "-vf")
	args = append(args, fmt.Sprintf("scale=%d:%d", util.Iif(width > 0, width, -1), util.Iif(height > 0, height, -1)))

	args = append(args, "-frames:v", "1")
	args = append(args, "-update", "true")

	args = append(args, output)

	return run(binary, args...)
}

// FindFFmpegBin checks if one of FFmpegBinPaths exists
func FindFFmpegBin() (string, error) {

	if FFmpegBinPath != "" {
		return FFmpegBinPath, nil
	}

	for _, path := range FFmpegBinPaths {
		l, err := filepath.Glob(path)
		if err != nil {
			continue
		}
		if len(l) > 0 {
			FFmpegBinPath = l[0]
			return FFmpegBinPath, nil
		}
	}

	return "", os.ErrNotExist
}
