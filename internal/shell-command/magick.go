package shell_command

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
func ImageMagickThumb(inputFilePath string, outputFilePath string, width int, height int) error {

	binary := FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", ImageMagickBinPaths)
	}

	var args []string

	args = append(args, inputFilePath)

	if width > 0 || height > 0 {
		args = append(args, "-geometry")
		if width > 0 && height > 0 {
			args = append(args, fmt.Sprintf("%d:%d", width, height))
		} else if width > 0 {
			args = append(args, fmt.Sprintf("%d", width))
		} else {
			args = append(args, fmt.Sprintf("x%d", height))
		}
	}

	args = append(args, "-flatten")
	args = append(args, "-colorspace", "RGB")

	args = append(args, outputFilePath)

	return util.RunSilent(binary, args...)
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

	return util.RunSilent(binary, args...)
}

// FindImageMagickBin checks if one of FFmpegBinPaths exists
func FindImageMagickBin() string {

	if ImageMagickBinPath != "" {
		return ImageMagickBinPath
	}
	ImageMagickBinPath = util.FindFile(ImageMagickBinPaths)
	return ImageMagickBinPath
}
