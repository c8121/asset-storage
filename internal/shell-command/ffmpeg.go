package shell_command

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
func FFmpegThumb(inputFilePath string, outputFilePath string, width int, height int) error {

	binary := FindFFmpegBin()
	if binary == "" {
		return fmt.Errorf("FFmpeg not found (searching in %v)", FFmpegBinPaths)
	}

	var args []string

	if ThumbVideoPosition != "" {
		args = append(args, "-ss", ThumbVideoPosition)
	}

	args = append(args, "-y") //Overwrite
	args = append(args, "-i", inputFilePath)

	args = append(args, "-vf")
	args = append(args, fmt.Sprintf("scale=%d:%d", util.Iif(width > 0, width, -1), util.Iif(height > 0, height, -1)))

	args = append(args, "-frames:v", "1")
	args = append(args, "-update", "true")

	args = append(args, outputFilePath)

	return util.RunSilent(binary, args...)
}

// FindFFmpegBin checks if one of FFmpegBinPaths exists
func FindFFmpegBin() string {

	if FFmpegBinPath != "" {
		return FFmpegBinPath
	}
	FFmpegBinPath = util.FindFile(FFmpegBinPaths)
	return FFmpegBinPath
}
