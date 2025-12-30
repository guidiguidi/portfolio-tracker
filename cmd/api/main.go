package main

import (
    "github.com/gin-gonic/gin"
    "internal/config"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(err)
	}
	
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"logger": cfg.Logger.Level,
			"version": cfg.App.Version,
		})
	})
	r.Run()
}
