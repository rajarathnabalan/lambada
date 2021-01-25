package lambada

import (
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
