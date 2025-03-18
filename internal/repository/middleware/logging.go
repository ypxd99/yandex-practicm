package middleware

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ypxd99/yandex-practicm/util"
)

var noWritten = -1

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
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

func (r *loggingResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

func (r *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if r.responseData.size < 0 {
		r.responseData.size = 0
	}
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *loggingResponseWriter) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

func (r *loggingResponseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (r *loggingResponseWriter) WriteHeaderNow() {
	if !r.Written() {
		r.responseData.size = 0
		r.ResponseWriter.WriteHeader(r.responseData.status)
	}
}

func (r *loggingResponseWriter) WriteString(s string) (n int, err error) {
	r.WriteHeaderNow()
	n, err = io.WriteString(r.ResponseWriter, s)
	r.responseData.size += n
	return
}

func (r *loggingResponseWriter) Status() int {
	return r.responseData.status
}

func (r *loggingResponseWriter) Size() int {
	return r.responseData.size
}

func (r *loggingResponseWriter) Written() bool {
	return r.responseData.size != noWritten
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
