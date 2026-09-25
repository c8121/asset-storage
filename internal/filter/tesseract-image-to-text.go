package filter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/filter_commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type TesseractImageToTextFilter struct {
	DefaultLanguage string
}

func NewTesseractImageToTextFilter() *TesseractImageToTextFilter {
	f := &TesseractImageToTextFilter{}
	f.DefaultLanguage = ""
	return f
}

func (f TesseractImageToTextFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "image/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	lang := util.GetOrDefault(params, "lang", f.DefaultLanguage)

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	out, err := f.imageTesseractImageToText(in, lang)
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

// imageTesseractImageToText executes Tesseract for Image to Text conversion ...
func (f TesseractImageToTextFilter) imageTesseractImageToText(input string, lang string) (string, error) {

	binary := filter_commands.FindTesseractBin()
	if binary == "" {
		return "", fmt.Errorf("tesseract not found (searching in %v)", filter_commands.TesseractBinPaths)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, "tesseract*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	util.CloseOrLog(out)
	util.LogError(os.Remove(out.Name()))

	var args []string

	args = append(args, input)
	args = append(args, out.Name())

	if lang != "" {
		args = append(args, "-l", lang)
	}

	var env []string
	env = append(env, "OMP_THREAD_LIMIT=1")

	err = util.RunSilentWithEnv(binary, env, args...)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	entries, err := os.ReadDir(config.AssetStorageTempDir)
	if err != nil {
		return "", fmt.Errorf("failed to read dir: %w", err)
	}

	outFileName := filepath.Base(out.Name())
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), outFileName) {
			return filepath.Join(config.AssetStorageTempDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("output file not found")
}
