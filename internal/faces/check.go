package faces

import (
	"github.com/c8121/go_vue_spa/notifications"
	"github.com/go-resty/resty/v2"
)

func CheckFaceRestServiceAvailable() bool {

	client := resty.New()
	_, err := client.R().
		Get(CheckFaceRestServiceEndpoint)
	if err != nil {
		notifications.AppNotifications.AddNotification("Face REST API is not available: " + GetFaceRestServiceEndpoint + " (" + err.Error() + ")")
		return false
	}

	return true
}
