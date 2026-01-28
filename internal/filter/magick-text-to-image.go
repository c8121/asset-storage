package filter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	filter_commands "github.com/c8121/asset-storage/internal/filter-commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type ImageMagickTextToImageFilter struct {
	DefaultWidth           int
	FontSizeFactor         float64
	MaxLines               int
	MaxLineLength          int
	DefaultFileNamePattern string
	DefaultMimeType        string
}

func NewImageMagickTextToImageFilter() *ImageMagickTextToImageFilter {
	f := &ImageMagickTextToImageFilter{}
	f.DefaultWidth = 400
	f.FontSizeFactor = 0.05
	f.MaxLines = 50
	f.MaxLineLength = 160
	f.DefaultFileNamePattern = "asset-thumb*.png"
	f.DefaultMimeType = "image/png"
	return f
}

func (f ImageMagickTextToImageFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "text/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	width, _ := strconv.Atoi(util.GetOrDefault(params, "width", strconv.Itoa(f.DefaultWidth)))
	height, _ := strconv.Atoi(util.GetOrDefault(params, "height", "0"))
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

	err = f.imageMagickTextToImage(in, out.Name(), width, height)
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

// imageMagickTextToImage executes ImageMagick for Text to Image conversion ...
func (f ImageMagickTextToImageFilter) imageMagickTextToImage(input string, output string, width int, height int) error {

	binary := filter_commands.FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", filter_commands.ImageMagickBinPaths)
	}

	text, err := headText(input, f.MaxLines, f.MaxLineLength)
	if err != nil {
		return fmt.Errorf("failed to get head of textfile: %w", err)
	}

	var pointsize int
	var args []string

	args = append(args, "-size")
	if width > 0 && height > 0 {
		args = append(args, fmt.Sprintf("%dx%d", width, height))
		pointsize = int(float64(width) * f.FontSizeFactor)
	} else if width > 0 {
		args = append(args, fmt.Sprintf("%dx%d", width, width))
		pointsize = int(float64(width) * f.FontSizeFactor)
	} else if height > 0 {
		args = append(args, fmt.Sprintf("%dx%d", height, height))
		pointsize = int(float64(height) * f.FontSizeFactor)
	} else {
		args = append(args, fmt.Sprintf("%dx%d", f.DefaultWidth, f.DefaultWidth))
		pointsize = int(float64(f.DefaultWidth) * f.FontSizeFactor)
	}

	args = append(args, "xc:white")
	args = append(args, "-font", "Courier")
	args = append(args, "-pointsize", fmt.Sprintf("%d", pointsize))
	args = append(args, "-fill", "black")
	args = append(args, "-annotate", fmt.Sprintf("+%d+%d", pointsize*2, pointsize*2))
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
