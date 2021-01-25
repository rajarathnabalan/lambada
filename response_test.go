package lambada

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWriter(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	w := newResponseWriter()
	require.NotNil(w)
	require.NotNil(w.header)
	assert.False(w.binary)

	w.Header().Set("X-Test", "value")

	_, err := w.Write(([]byte)("Hello, World!"))
	assert.NoError(err)
	assert.Equal(http.StatusOK, w.statusCode)

	// This must have no effect on lockedHeaders
	w.Header().Set("X-Test-2", "value2")

	w.finalize()
	assert.Equal("13", w.lockedHeader.Get("Content-Length"))
	assert.NotEqual("", w.lockedHeader.Get("Content-Type"))
	assert.Equal("value", w.lockedHeader.Get("X-Test"))
	assert.Equal("", w.lockedHeader.Get("X-Test-2"))
}

func TestResponseWriterFinalize(t *testing.T) {
	assert := assert.New(t)

	w := newResponseWriter()

	// Should trigger a WriteHeader
	w.finalize()
	assert.Equal(http.StatusOK, w.statusCode)
}

func TestResponseWriterInvalidStatus(t *testing.T) {
	assert := assert.New(t)

	w := newResponseWriter()
	for value := range []int{0, 10, 99, 625, 999} {
		assert.Panics(func() {
			w.WriteHeader(value)
		})
	}
}

func TestResponseWriterBinary(t *testing.T) {
	assert := assert.New(t)

	w := newResponseWriter()
	assert.False(w.binary)

	SetBinary(w)
	assert.True(w.binary)

	SetText(w)
	assert.False(w.binary)
}

func TestResponseWriterBinary2(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	assert.NotPanics(func() {
		SetBinary(w)
	})
	assert.NotPanics(func() {
		SetText(w)
	})
}
