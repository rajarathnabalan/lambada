package agw

import (
	"bytes"
	"fmt"
	"net/http"
)

type ResponseWriter struct {
	header     http.Header
	body       bytes.Buffer
	statusCode int
	Binary     bool
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(data)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if statusCode < 100 || statusCode >= 600 {
		panic(fmt.Errorf("Invalid status code %d", statusCode))
	}
	w.statusCode = statusCode
}

func (w *ResponseWriter) Body() []byte {
	return w.body.Bytes()
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}
