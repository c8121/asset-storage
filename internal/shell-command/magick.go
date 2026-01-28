package shell_command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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

// ImageMagickThumbFromPdf executes ImageMagick for Text to Image conversion ...
func ImageMagickThumbFromTxt(input string, output string, width int, height int) error {

	binary := FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", ImageMagickBinPaths)
	}

	text, err := headText(input, 50, 160)
	if err != nil {
		return fmt.Errorf("failed to get head of textfile: %w", err)
	}

	var pointsize int

	var args []string

	args = append(args, "-size")
	if width > 0 && height > 0 {
		args = append(args, fmt.Sprintf("%dx%d", width, height))
		pointsize = int(float64(width) * 0.05)
	} else if width > 0 {
		args = append(args, fmt.Sprintf("%dx%d", width, width))
		pointsize = int(float64(width) * 0.05)
	} else if height > 0 {
		args = append(args, fmt.Sprintf("%dx%d", height, height))
		pointsize = int(float64(height) * 0.05)
	} else {
		args = append(args, fmt.Sprintf("%dx%d", 400, 400))
		pointsize = int(float64(400) * 0.05)
	}

	args = append(args, "xc:white")
	args = append(args, "-font", "Courier")
	args = append(args, "-pointsize", fmt.Sprintf("%d", pointsize))
	args = append(args, "-fill", "black")
	args = append(args, "-annotate", "+15+15")
	args = append(args, text)

	args = append(args, output)

	return util.RunSilent(binary, args...)
}

func headText(input string, numLines int, maxLineLen int) (string, error) {

	var err error

	f, err := os.Open(input)
	if err != nil {
		return "", err
	}
	defer util.CloseOrLog(f)

	var buf strings.Builder

	reader := bufio.NewReaderSize(f, maxLineLen)
	for n := 0; err == nil && n < numLines; n++ {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return "", err
		}
		buf.Write(line)
		buf.WriteRune('\n')
	}
	return buf.String(), nil
}

// FindImageMagickBin checks if one of FFmpegBinPaths exists
func FindImageMagickBin() string {

	if ImageMagickBinPath != "" {
		return ImageMagickBinPath
	}
	ImageMagickBinPath = util.FindFile(ImageMagickBinPaths)
	return ImageMagickBinPath
}
