package filter_commands

import (
	"fmt"

	"github.com/c8121/go_vue_spa/notifications"
)

func CheckAndNotify() {

	path := FindFFmpegBin()
	if path == "" {
		notifications.AppNotifications.AddNotification("FFmpeg is not installed. To generate video-thumbs and images, please install FFmpeg")
	} else {
		notifications.AppNotifications.AddNotification("Using FFmpeg from " + path)
	}

	path = FindImageMagickBin()
	if path == "" {
		notifications.AppNotifications.AddNotification("ImageMagick is not installed. To support more image formats, please install ImageMagick")
	} else {
		notifications.AppNotifications.AddNotification("Using ImageMagick from " + path)
	}

	path = FindTesseractBin()
	if path == "" {
		notifications.AppNotifications.AddNotification("Tesseract-OCR is not installed. To support image to text conversion, please install Tesseract-OCR")
	} else {
		notifications.AppNotifications.AddNotification("Using Tesseract from " + path)
	}

	path = FindPopplerPdfBin()
	if path == "" {
		notifications.AppNotifications.AddNotification("Poppler pdftotext is not installed. To support pdf to text conversion, please install pdftotext")
	} else {
		notifications.AppNotifications.AddNotification("Using pdftotext from " + path)
	}

	path = FindPandocBin()
	if path == "" {
		notifications.AppNotifications.AddNotification("Pandoc is not installed. To support doc to text conversion, please install pandoc")
	} else {
		notifications.AppNotifications.AddNotification("Using pandoc from " + path)
	}

	for _, msg := range notifications.AppNotifications.Messages {
		fmt.Println(msg)
	}
}
