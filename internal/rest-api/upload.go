package restapi

import (
	"net/http"

	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

func ReceiveUpload(c *gin.Context) {

	//Read JSON request in chunks to temp file,
	//parse later
	w, err := storage.NewTempFileWriter()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer util.CloseOrLog(w)

	buf := make([]byte, 4096)
	for {
		n, err := c.Request.Body.Read(buf)
		if err != nil {
			break
		}

		_, err = w.Write(buf[:n])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	c.Data(http.StatusOK, "application/json", []byte("{\"message\": \"TBD\"}"))
}
