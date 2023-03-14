package lambada

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/morelj/httptools/header"
)

// ResponseWriter is an implementation of http.ResponseWriter which stores the data written to it internally.
// Trailers are not supported by this implementation.
//
// You usually access ResponseWriter through the http.ResponseWriter interface.
// If you need to access the underlying ResponseWriter use:
//
//	// w is an http.ResponseWriter
//	rw, ok := w.(*lambada.ResponseWriter)
type ResponseWriter struct {
	outputMode            OutputMode
	header                http.Header
	lockedHeader          http.Header
	body                  bytes.Buffer
	statusCode            int
	binary                bool
	ignoreBinaryDetection bool
}

func newResponseWriter(outputMode OutputMode, binary bool) *ResponseWriter {
	return &ResponseWriter{
		outputMode: outputMode,
		header:     http.Header{},
		binary:     binary,
	}
}

// Header returns the response's header set.
// Once the WriteHeader has been called, Header returns a new copy of the response's headers, preventing them from
// being modified.
func (w *ResponseWriter) Header() http.Header {
	if w.statusCode == 0 {
		return w.header
	}
	return w.lockedHeader.Clone()
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(data)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		// WriteHeader has not been called yet

		if statusCode < 100 || statusCode >= 600 {
			panic(fmt.Errorf("invalid status code %d", statusCode))
		}
		w.statusCode = statusCode

		// Current headers are copied into lockedHeader, so further changed to the header map will not affect headers
		w.lockedHeader = w.header.Clone()
	}
}

func (w *ResponseWriter) finalize() {
	// Ensure the header has been written
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}

	body := w.body.Bytes()

	// Compute Content-Length
	w.lockedHeader.Set(header.ContentLength, strconv.FormatInt(int64(len(body)), 10))

	if w.outputMode >= AutoContentType {
		// Compute Content-Type if not set
		if len(body) > 0 && w.lockedHeader.Get(header.ContentType) == "" {
			w.lockedHeader.Set(header.ContentType, http.DetectContentType(body))
		}

		if w.outputMode >= Automatic && !w.ignoreBinaryDetection {
			// Detect if output is binary
			// We don't change the mode if we can't determine if the output is binary or not
			switch isBinary(w.lockedHeader.Get(header.ContentType), w.lockedHeader.Get(header.ContentEncoding)) {
			case bsBinary:
				w.SetBinary(true)
			case bsText:
				w.SetBinary(false)
			}
		}
	}
}

// StatusCode returns w's current status code.
// If WriteHeaders() has not been called yet, returns 200.
func (w *ResponseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

// Body returns the current body's byte.
// If nothing has been written, Body returns nil.
// The returned slice is valid until the next call to Write.
func (w *ResponseWriter) Body() []byte {
	return w.body.Bytes()
}

// SetBinary sets whether or not the binary mode should be enabled or not.
// When binary mode is enabled, the response is encoded to Base64 before being returned to API Gateway.
// The mode set through this function is forced - meaning that automatic binary detection will be skipped.
// Use AllowBinaryDetection() to revert this behavior.
func (w *ResponseWriter) SetBinary(binary bool) {
	w.binary = binary
	w.ignoreBinaryDetection = true
}

// AllowBinaryDetection allows binary detection to happen on w.
// Automatic binary detection will only happen when the OutputMode is set to Automatic.
func (w *ResponseWriter) AllowBinaryDetection() {
	w.ignoreBinaryDetection = false
}

// SetOutputMode sets the output mode for w.
func (w *ResponseWriter) SetOutputMode(outputMode OutputMode) {
	w.outputMode = outputMode
}

// SetBinary enforces binary mode for the given ResponseWriter.
// That is, the response will be encoded to Base64 when returned to API Gateway.
//
// If the passed ResponseWriter has not been provided by Lambada, this function has no effect.
func SetBinary(w http.ResponseWriter) {
	if w, ok := w.(*ResponseWriter); ok {
		w.SetBinary(true)
	}
}

// SetText enforces text mode for the given ResponseWriter.
// That is, the response will be be set to not be encoded to Base64 when returned to API Gateway.
//
// If the passed ResponseWriter has not been provided by Lambada, this function has no effect.
func SetText(w http.ResponseWriter) {
	if w, ok := w.(*ResponseWriter); ok {
		w.SetBinary(false)
	}
}

// SetOutputMode sets the given output mode to the given ResponseWriter.
func SetOutputMode(w http.ResponseWriter, outputMode OutputMode) {
	if w, ok := w.(*ResponseWriter); ok {
		w.SetOutputMode(outputMode)
	}
}
