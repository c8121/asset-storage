package thumbnails

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	ImageMagickBinPaths = []string{
		"/usr/bin/convert",
		"C:/Program Files/ImageMagick*/magick.exe",
	}

	ImageMagickBinPath = ""
)

// ImageMagickThumb executes ...
func ImageMagickThumb(input string, output string, width int, height int) error {

	binary := FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("FFmpeg not found (searching in %v)", FFmpegBinPaths)
	}

	var args []string

	args = append(args, input)

	if width > 0 && height > 0 {
		args = append(args, "-geometry")
		args = append(args, fmt.Sprintf("%d:%d", util.Iif(width > 0, width, -1), util.Iif(height > 0, height, -1)))
	}

	args = append(args, "-flatten")
	args = append(args, "-colorspace", "RGB")

	args = append(args, output)

	return run(binary, args...)
}

// FindImageMagickBin checks if one of FFmpegBinPaths exists
func FindImageMagickBin() string {

	if ImageMagickBinPath != "" {
		return ImageMagickBinPath
	}
	ImageMagickBinPath = util.FindFile(ImageMagickBinPaths)
	return ImageMagickBinPath
}
