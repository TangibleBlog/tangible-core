package control

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"tangible-core/public/common"

	"github.com/gin-gonic/gin"
)

func handlerSecurity() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.GetHeader("Token"); token == "" {
			c.Abort()
			common.CodeResult(c, http.StatusUnauthorized)
		} else {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			var bodyJson map[string]interface{}
			err := json.Unmarshal(body, &bodyJson)
			if err != nil {
				c.Abort()
				common.CodeResult(c, http.StatusForbidden)
			}
			if value, ok := bodyJson["Sign"]; ok {
				if !validMAC(value.(string), token, []byte(Config.Password)) {
					c.Abort()
					common.CodeResult(c, http.StatusForbidden)
				}
			} else {
				c.Abort()
				common.CodeResult(c, http.StatusForbidden)
			}
		}
	}
}

func validMAC(message string, messageMAC string, key []byte) bool {
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(message))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	if expectedMAC == messageMAC {
		return true
	}
	return false
}
