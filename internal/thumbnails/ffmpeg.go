package thumbnails

import (
	"fmt"

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

	binary := FindFFmpegBin()
	if binary == "" {
		return fmt.Errorf("FFmpeg not found (searching in %v)", FFmpegBinPaths)
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
func FindFFmpegBin() string {

	if FFmpegBinPath != "" {
		return FFmpegBinPath
	}
	FFmpegBinPath = util.FindFile(FFmpegBinPaths)
	return FFmpegBinPath
}
