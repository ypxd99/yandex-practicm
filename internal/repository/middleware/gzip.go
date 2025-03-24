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

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			handleGzipRequest(c)
		}

		wr := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            &bytes.Buffer{},
		}
		c.Writer = wr

		c.Next()

		contentType := c.Writer.Header().Get("Content-Type")
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if strings.Contains(contentType, "application/json") || 
			strings.Contains(contentType, "text/html") && 
			strings.Contains(acceptEncoding, "gzip") {
			handleGzipResponse(c, wr)
		} else {
			wr.ResponseWriter.Write(wr.buf.Bytes())
		}
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

func handleGzipResponse(c *gin.Context, wr *responseWriter) {
	gz := gzip.NewWriter(c.Writer)
	defer gz.Close()

	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Del("Content-Length")

	gz.Write(wr.buf.Bytes())
	gz.Flush()
}
