package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ypxd99/yandex-practicm/util"
)

// responseData представляет структуру для хранения данных о HTTP ответе.
// Содержит статус код и размер ответа.
type responseData struct {
	// status HTTP статус код ответа
	status int
	// size размер ответа в байтах
	size int
}

// loggingResponseWriter представляет обертку над gin.ResponseWriter для логирования.
// Позволяет отслеживать статус код и размер ответа.
type loggingResponseWriter struct {
	gin.ResponseWriter
	responseData *responseData
}

// Write записывает данные в ответ и обновляет размер ответа.
// Реализует интерфейс io.Writer.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader устанавливает HTTP статус код ответа.
// Реализует интерфейс http.ResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LoggingMiddleware создает middleware для логирования HTTP запросов.
// Логирует метод, URI, длительность обработки, статус код и размер ответа.
// Возвращает gin.HandlerFunc.
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
