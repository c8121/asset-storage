package filter

import (
	"fmt"
	"os"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/filter_commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type PandocOfficeToTextFilter struct {
}

func NewPandocOfficeToTextFilter() *PandocOfficeToTextFilter {
	f := &PandocOfficeToTextFilter{}
	return f
}

func (f PandocOfficeToTextFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	format := ""
	checkMeta := strings.ToLower(meta.MimeType)
	if strings.Contains(checkMeta, "word") {
		format = "docx"
	} else if strings.Contains(checkMeta, "opendocument.text") {
		format = "odt"
	}

	if format == "" {
		return nil, "", fmt.Errorf("cannot convert type: %s", meta.MimeType)
	}

	out, err := f.pandocOfficeToText(in, format)
	if err != nil {
		return nil, "", fmt.Errorf("failed to extract text: %w", err)
	}

	bytes, err := os.ReadFile(out)
	if err != nil {
		util.LogError(os.Remove(out))
		return nil, "", fmt.Errorf("failed to extract text: %w", err)
	}

	util.LogError(os.Remove(out))
	return bytes, "text/plain;charset=UTF-8", nil
}

// pandocOfficeToText executes pandoc to convert office documents to text
func (f PandocOfficeToTextFilter) pandocOfficeToText(input string, format string) (string, error) {

	binary := filter_commands.FindPandocBin()
	if binary == "" {
		return "", fmt.Errorf("pandoc not found (searching in %v)", filter_commands.PandocBinPaths)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, "pandoc-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	var args []string

	args = append(args, input)
	args = append(args, "-f", format)
	args = append(args, "-t", "plain")
	args = append(args, "-o", out.Name())

	err = util.RunSilent(binary, args...)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	return out.Name(), nil
}
