package thumbnails

import (
	"fmt"
	"strconv"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	ImageMagickBinPaths = []string{
		"/usr/bin/convert",
		"C:/Program Files/ImageMagick*/magick.exe",
	}

	ImageMagickBinPath = ""
)

// ImageMagickThumb executes ImageMagick for Image conversion ...
func ImageMagickThumb(input string, output string, width int, height int) error {

	binary := FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", ImageMagickBinPaths)
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

// ImageMagickThumbFromPdf executes ImageMagick for PDF to Image conversion ...
func ImageMagickThumbFromPdf(input string, output string, width int, height int) error {

	binary := FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", ImageMagickBinPaths)
	}

	var args []string

	args = append(args, input+"[0]")

	if width > 0 || height > 0 {
		args = append(args, "-thumbnail")
		args = append(args, fmt.Sprintf("%sx%s", util.Iif(width > 0, strconv.Itoa(width), ""), util.Iif(height > 0, strconv.Itoa(height), "")))
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
