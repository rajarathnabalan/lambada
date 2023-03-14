package lambada

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBodyToBytes(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		body      string
		isBase64  bool
		expectErr bool
		expected  []byte
	}{
		{
			body:      "Hello, World!",
			isBase64:  false,
			expectErr: false,
			expected:  []byte("Hello, World!"),
		},
		{
			body:      "SGVsbG8sIFdvcmxkIQ==",
			isBase64:  true,
			expectErr: false,
			expected:  []byte("Hello, World!"),
		},
		{
			body:      "This is invalid",
			isBase64:  true,
			expectErr: true,
		},
	}

	for i, c := range cases {
		data, err := bodyToBytes(c.body, c.isBase64)
		if c.expectErr {
			assert.Errorf(err, "Case %d", i)
		} else {
			assert.NoErrorf(err, "Case %d", i)
			if err == nil {
				assert.Equalf(c.expected, data, "Case %d", i)
			}
		}
	}
}

func TestBytesToBody(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		bytes    []byte
		isBase64 bool
		expected string
	}{
		{
			bytes:    []byte("Hello, World!"),
			isBase64: false,
			expected: "Hello, World!",
		},
		{
			bytes:    []byte("Hello, World!"),
			isBase64: true,
			expected: "SGVsbG8sIFdvcmxkIQ==",
		},
	}

	for i, c := range cases {
		str := bytesToBody(c.bytes, c.isBase64)
		assert.Equalf(c.expected, str, "Case %d", i)
	}
}

func TestIsBinary(t *testing.T) {
	cases := []struct {
		contentType     string
		contentEncoding string
		status          binaryStatus
	}{
		{
			contentType:     "application/json",
			contentEncoding: "",
			status:          bsText,
		},
		{
			contentType:     "text/csv",
			contentEncoding: "",
			status:          bsText,
		},
		{
			contentType:     "application/pdf",
			contentEncoding: "",
			status:          bsBinary,
		},
		{
			contentType:     "text/vnd.a",
			contentEncoding: "",
			status:          bsText,
		},
		{
			contentType:     "video/H264",
			contentEncoding: "",
			status:          bsBinary,
		},
		{
			contentType:     "application/pskc+xml",
			contentEncoding: "",
			status:          bsText,
		},
		{
			contentType:     "text/plain",
			contentEncoding: "gzip",
			status:          bsBinary,
		},
		{
			contentType:     "unknown/unknown",
			contentEncoding: "",
			status:          bsUnknown,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d_%s_%s", i, c.contentType, c.contentEncoding), func(t *testing.T) {
			assert := assert.New(t)
			status := isBinary(c.contentType, c.contentEncoding)
			assert.Equal(c.status, status)
		})
	}
}
