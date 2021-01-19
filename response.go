package lambada

import (
	"bytes"
	"fmt"
	"net/http"
)

type responseWriter struct {
	header     http.Header
	body       bytes.Buffer
	statusCode int
	Binary     bool
}

func newResponseWriter() *responseWriter {
	return &responseWriter{
		header: http.Header{},
	}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if statusCode < 100 || statusCode >= 600 {
		panic(fmt.Errorf("Invalid status code %d", statusCode))
	}
	w.statusCode = statusCode
}

func (w *responseWriter) Body() []byte {
	return w.body.Bytes()
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

func SetBinary(w http.ResponseWriter) {
	if w, ok := w.(*responseWriter); ok {
		w.Binary = true
	}
}

func SetText(w http.ResponseWriter) {
	if w, ok := w.(*responseWriter); ok {
		w.Binary = false
	}
}
