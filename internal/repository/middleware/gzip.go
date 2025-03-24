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

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

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

func (w *responseWriter) handleGzipResponse() {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Del("Content-Length")

	gz := gzip.NewWriter(w.ResponseWriter)
	defer gz.Close()

	gz.Write(w.buf.Bytes())
}
