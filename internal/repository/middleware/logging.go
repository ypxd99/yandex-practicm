package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ypxd99/yandex-practicm/util"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	gin.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware() gin.HandlerFunc {
	logger := util.GetLogger()

	return func(c *gin.Context) {
		start := time.Now()

		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		rw := &loggingResponseWriter{
			ResponseWriter: c.Writer,
			responseData:   &responseData{},
		}
		c.Writer = rw

		c.Next()

		entry := logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"uri":      c.Request.RequestURI,
			"duration": time.Since(start).String(),
			"status":   rw.responseData.status,
			"size":     rw.responseData.size,
		})

		entry.Info("HTTP request")
	}
}
