package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// responseWriter представляет обертку над gin.ResponseWriter для сжатия ответов.
// Позволяет буферизировать ответ перед его отправкой.
type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

// Write записывает данные в буфер.
// Реализует интерфейс io.Writer.
func (w *responseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

// WriteHeader устанавливает HTTP статус код ответа.
// Реализует интерфейс http.ResponseWriter.
func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// GzipMiddleware создает middleware для сжатия HTTP запросов и ответов.
// Поддерживает сжатие запросов с Content-Encoding: gzip.
// Сжимает ответы, если клиент поддерживает gzip и тип контента подходит.
// Возвращает gin.HandlerFunc.
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			handleGzipRequest(c)
		}

		wr := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            &bytes.Buffer{},
		}
		c.Writer = wr

		c.Next()

		acceptsGzip := strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip")
		contentType := c.Writer.Header().Get("Content-Type")
		if (strings.Contains(contentType, "application/json") ||
			strings.Contains(contentType, "text/html")) &&
			acceptsGzip {
			wr.handleGzipResponse()
			return
		}

		wr.ResponseWriter.Write(wr.buf.Bytes())
	}
}

// handleGzipRequest обрабатывает входящий запрос, сжатый с помощью gzip.
// Распаковывает тело запроса и обновляет его содержимое.
// В случае ошибки прерывает обработку запроса.
func handleGzipRequest(c *gin.Context) {
	gz, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("invalid gzip body"))
		return
	}
	defer gz.Close()

	body, err := io.ReadAll(gz)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("failed to read gzip body"))
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	c.Request.ContentLength = int64(len(body))
}

// handleGzipResponse сжимает ответ с помощью gzip.
// Устанавливает соответствующие заголовки и записывает сжатые данные.
func (w *responseWriter) handleGzipResponse() {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Del("Content-Length")

	gz := gzip.NewWriter(w.ResponseWriter)
	defer gz.Close()

	gz.Write(w.buf.Bytes())
}
