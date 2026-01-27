package filter

import (
	"fmt"
	"regexp"
	"strings"
)

type AssetFilter struct {
	Name      string //Internal, descriptive name
	Alias     string //Alias name which clients might use
	Filter    Filter
	MimeTypes []string //Asset Mime-Type RegEx, which the filter can handle
}

var (
	AvailableFilters = []AssetFilter{}
)

// GetFirstFilterByMimeType returns the first filter out of AvailableFilters
// which has a matching mime-type
func GetFirstFilterByMimeType(mimetype string) Filter {

	loadAvailableFilters()

	for _, f := range AvailableFilters {
		for _, m := range f.MimeTypes {
			match, _ := regexp.MatchString(m, mimetype)
			if match {
				fmt.Printf("Using filter '%s' for %s\n", f.Name, mimetype)
				return f.Filter
			}
		}
	}
	return nil
}

// GetFirstFilterByMimeType returns the first filter out of AvailableFilters
// which has a matching name or alias and mime-type
func GetFirstFilterByNameAndMimeType(name string, mimetype string) Filter {

	loadAvailableFilters()

	for _, f := range AvailableFilters {
		if !strings.EqualFold(name, f.Name) && !strings.EqualFold(name, f.Alias) {
			continue
		}
		for _, m := range f.MimeTypes {
			match, _ := regexp.MatchString(m, mimetype)
			if match {
				fmt.Printf("Using filter '%s' for %s\n", f.Name, mimetype)
				return f.Filter
			}
		}
	}
	return nil
}

func loadAvailableFilters() {

	//TODO make mime-type->filter mapping in configurable

	if len(AvailableFilters) > 0 {
		return
	}

	AvailableFilters = []AssetFilter{
		{
			Name:   "ImageMagick",
			Alias:  "Image",
			Filter: NewImageMagickFilter(),
			MimeTypes: []string{
				"(?i)^image/bmp$",
				"(?i)^image/tiff$",
				"(?i)^application/pdf$",
				"(?i)^image/.*xcf$",
				"(?i)photoshop",
			},
		},
		{
			Name:   "NativeImage",
			Filter: NewImageFilter(),
			MimeTypes: []string{
				"(?i)^image/",
			},
		},
		{
			Name:   "FFmpegFilter",
			Alias:  "VideoToImage",
			Filter: NewFFmpegFilter(),
			MimeTypes: []string{
				"(?i)^video/",
			},
		},
	}

}
