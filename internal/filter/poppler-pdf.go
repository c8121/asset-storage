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

type PopplerPdfToTextFilter struct {
}

func NewPopplerPdfToTextFilter() *PopplerPdfToTextFilter {
	f := &PopplerPdfToTextFilter{}
	return f
}

func (f PopplerPdfToTextFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "application/pdf") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	out, err := f.pdfPopplerPdfToText(in)
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

// pdfPopplerPdfToText executes pdftotext for PDF to Text conversion ...
func (f PopplerPdfToTextFilter) pdfPopplerPdfToText(input string) (string, error) {

	binary := filter_commands.FindPopplerPdfBin()
	if binary == "" {
		return "", fmt.Errorf("pdftotext not found (searching in %v)", filter_commands.TesseractBinPaths)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, "poppler-pdftotext-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	var args []string

	args = append(args, input)
	args = append(args, out.Name())

	err = util.RunSilent(binary, args...)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	return out.Name(), nil
}
