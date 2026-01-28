package filter_commands

import (
	"github.com/c8121/asset-storage/internal/util"
)

var (
	FFmpegBinPaths = []string{
		"/opt/ffmpeg*/bin/ffmpeg.exe",
	}

	FFmpegBinPath = ""
)

// FindFFmpegBin checks if one of FFmpegBinPaths exists
func FindFFmpegBin() string {

	if FFmpegBinPath != "" {
		return FFmpegBinPath
	}
	FFmpegBinPath = util.FindFile(FFmpegBinPaths)
	return FFmpegBinPath
}
