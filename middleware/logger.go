package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

func SetupLogger() {
	// Output logs to a file
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Out = os.Stdout
	}

	// Use JSON formatter for structured logs
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process request

		duration := time.Since(start)
		log.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"latency": duration.String(),
			"agent":   c.Request.UserAgent(),
			"error":   c.Errors.String(),
		}).Info("Request details")
	}
}
