package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guidiguidi/portfolio-tracker/internal/assets"
)

func NewRouter(log *slog.Logger, assetsHandler *assets.Handler) *gin.Engine {
	r := gin.New()

	// Middlewares
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware(log))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	assetsGroup := r.Group("/assets")
	assetsGroup.GET("", assetsHandler.ListAssets)
	assetsGroup.POST("", assetsHandler.CreateAsset)
	assetsGroup.GET("/:id", assetsHandler.GetAsset)

	return r
}

// LoggerMiddleware logs requests using the structured logger.
func LoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)

		log.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}
