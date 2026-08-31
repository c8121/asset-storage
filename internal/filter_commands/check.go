package filter_commands

import "fmt"

func CheckAndNotify() {
	path := FindFFmpegBin()
	if path == "" {
		fmt.Println("FFmpeg is not installed. To generate video-thumbs and images, please install FFmpeg")
	} else {
		fmt.Printf("Using FFmpeg from %s\n", path)
	}

	path = FindImageMagickBin()
	if path == "" {
		fmt.Println("ImageMagick is not installed. To support more image formats, please install ImageMagick")
	} else {
		fmt.Printf("Using ImageMagick from %s\n", path)
	}
}
