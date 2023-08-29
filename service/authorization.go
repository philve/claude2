package service

import (
	"claude2/global"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetSessionKey(c *gin.Context) (sk string) {
	auth := c.Request.Header.Get("Authorization")
	hasPrefix := strings.HasPrefix(auth, "Bearer ")
	if hasPrefix && len(auth) > 7 && strings.HasPrefix(auth[7:], "sk-ant-sid") {
		sk = "sessionKey=" + auth[7:]
	}
	if sk == "" {
		sk = global.ServerConfig.Claude.GetSessionKey()
	}
	return sk
}
