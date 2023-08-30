package main

import (
	"claude2/handles"
	"claude2/initialize"
	"claude2/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	initialize.NewViper()
	r := gin.Default()
	r.Use(middleware.CORS)
	r.SetTrustedProxies(nil)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.OPTIONS("/v1/chat/completions", handles.OptionsHandler)
	r.POST("/v1/chat/completions", handles.ChatCompletionsHandler)

	var conf = &initialize.ServerConfig
	if conf.TlsCert != "" && conf.TlsKey != "" {
		r.RunTLS(conf.ListenHost, conf.TlsCert, conf.TlsKey)
	} else {
		err := r.Run(conf.ListenHost)
		if err != nil {
			return
		}
	}

}
