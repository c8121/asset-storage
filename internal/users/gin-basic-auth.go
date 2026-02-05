package users

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequiredHandler(handlerFunc func(c *gin.Context)) func(c *gin.Context) {

	f := func(c *gin.Context) {
		if err := CheckBasicAuth(c); err != nil {
			fmt.Printf("Auth error: %s\n", err)
			RespondBasicAuthRequired(c)
			return
		}
		handlerFunc(c)
	}

	return f
}

func CheckBasicAuth(c *gin.Context) error {
	header := c.GetHeader("Authorization")
	if len(header) == 0 {
		return errors.New("Authentication required")
	}

	p := strings.Index(strings.ToLower(header), "basic")
	if p < 0 {
		return errors.New("Invalid authentication header (only basic supported)")
	}

	key := strings.TrimSpace(header[p+6:])
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		fmt.Printf("Basic auth: %s", err)
		return errors.New("Cannot decode authentication header")
	}
	key = string(decodedKey)

	p = strings.Index(key, ":")
	if p < 0 {
		fmt.Printf("Basic auth, invalid key: %s, header=%s\n", key, header)
		return errors.New("Invalid authentication header")
	}

	username := key[:p]
	password := key[p+1:]

	return Authenticate(username, []byte(password))
}

func RespondBasicAuthRequired(c *gin.Context) {
	c.Header("WWW-Authenticate", "Basic realm=\"Asset Storage\", charset=\"UTF-8\"")
	c.String(http.StatusUnauthorized, "Unauthorized")
}
