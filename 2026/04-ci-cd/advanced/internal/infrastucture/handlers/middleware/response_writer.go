package middleware

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	started    time.Time
	statusCode int
	count      uint64
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		started:        time.Now(),
	}
}

func (r *ResponseWriter) StatusCode() int {
	return r.statusCode
}

func (r *ResponseWriter) StatusCodeStr() string {
	return strconv.Itoa(r.statusCode)
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	atomic.AddUint64(&r.count, uint64(n))
	return n, err
}

func (r *ResponseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseWriter) Count() uint64 {
	return atomic.LoadUint64(&r.count)
}
