package filter_commands

import (
	"github.com/c8121/asset-storage/internal/util"
)

var (
	ImageMagickBinPaths = []string{
		"/usr/bin/convert",
		"C:/Program Files/ImageMagick*/magick.exe",
	}

	ImageMagickBinPath = ""
)

// FindImageMagickBin checks if one of FFmpegBinPaths exists
func FindImageMagickBin() string {

	if ImageMagickBinPath != "" {
		return ImageMagickBinPath
	}
	ImageMagickBinPath = util.FindFile(ImageMagickBinPaths)
	return ImageMagickBinPath
}
