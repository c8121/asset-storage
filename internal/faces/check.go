package faces

import (
	"github.com/c8121/asset-storage/internal/util"
	"github.com/go-resty/resty/v2"
)

func CheckFaceRestServiceAvailable() {

	client := resty.New()
	_, err := client.R().
		Get(CheckFaceRestServiceEndpoint)
	if err != nil {
		util.AppNotifications.AddNotification("Face REST API is not available: " + GetFaceRestServiceEndpoint + " (" + err.Error() + ")")
	}

}
