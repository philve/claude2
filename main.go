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
	err := r.Run(":8000")
	if err != nil {
		return
	}
}
