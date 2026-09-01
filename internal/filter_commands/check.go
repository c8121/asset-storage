package filter_commands

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/util"
)

func CheckAndNotify() {

	path := FindFFmpegBin()
	if path == "" {
		util.AppNotifications.AddNotification("FFmpeg is not installed. To generate video-thumbs and images, please install FFmpeg")
	} else {
		util.AppNotifications.AddNotification("Using FFmpeg from " + path)
	}

	path = FindImageMagickBin()
	if path == "" {
		util.AppNotifications.AddNotification("ImageMagick is not installed. To support more image formats, please install ImageMagick")
	} else {
		util.AppNotifications.AddNotification("Using ImageMagick from " + path)
	}

	for _, msg := range util.AppNotifications.Messages {
		fmt.Println(msg)
	}
}
