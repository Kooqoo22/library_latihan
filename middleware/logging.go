package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func Logging() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":     method,
			"path":       path,
			"status":     status,
			"latency":    latency.String(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("request completed")
	}
}
