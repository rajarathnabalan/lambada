package lambada

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// responseWriter is an implementation of http.ResponseWriter which stores the data written to it internally.
// Trailers are not supported by this implementation.
type responseWriter struct {
	header       http.Header
	lockedHeader http.Header
	body         bytes.Buffer
	statusCode   int
	binary       bool
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
	if w.statusCode == 0 {
		// WriteHeader has not been called yet

		if statusCode < 100 || statusCode >= 600 {
			panic(fmt.Errorf("Invalid status code %d", statusCode))
		}
		w.statusCode = statusCode

		// Current headers are copied into lockedHeader, so further changed to the header map will not affect headers
		w.lockedHeader = w.header.Clone()
	}
}

func (w *responseWriter) finalize() {
	// Ensure the header has been written
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}

	body := w.body.Bytes()

	// Compute Content-Length
	w.lockedHeader.Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))

	// Compute Content-Type if not set
	if len(body) > 0 && w.lockedHeader.Get("Content-Type") == "" {
		w.lockedHeader.Set("Content-Type", http.DetectContentType(body))
	}
}

// SetBinary enforces binary mode for the given ResponseWriter.
// That is, the response will be encoded to Base64 when returned to API Gateway.
//
// If the passed ResponseWriter has not been provided by Lambada, this function has no effect.
func SetBinary(w http.ResponseWriter) {
	if w, ok := w.(*responseWriter); ok {
		w.binary = true
	}
}

// SetText enforces text mode for the given ResponseWriter.
// That is, the response will be be set to not be encoded to Base64 when returned to API Gateway.
//
// If the passed ResponseWriter has not been provided by Lambada, this function has no effect.
func SetText(w http.ResponseWriter) {
	if w, ok := w.(*responseWriter); ok {
		w.binary = false
	}
}
