package lambada

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSingleValueHeaders(t *testing.T) {
	cases := []struct {
		h   http.Header
		res map[string]string
	}{
		{
			h: http.Header{
				"Content-Type": []string{"application/json"},
				"X-Empty":      nil,
				"X-Empty-2":    []string{},
			},
			res: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			h: http.Header{
				"Content-Type": []string{"application/json", "text/plain"},
			},
			res: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)

			res := toSingleValueHeaders(c.h)
			assert.Equal(c.res, res)
		})
	}
}

func TestFromSingleValueHeaders(t *testing.T) {
	cases := []struct {
		h   map[string]string
		res http.Header
	}{
		{
			h: map[string]string{
				"content-type": "application/json",
			},
			res: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)

			res := fromSingleValueHeaders(c.h)
			assert.Equal(c.res, res)
		})
	}
}

func TestCanonicalizeHeaders(t *testing.T) {
	assert := assert.New(t)

	h := http.Header{
		"content-type": []string{"application/json"},
		"accept":       []string{"*/*"},
		"cookie":       []string{"cookie=yummy"},
	}

	res := canonicalizeHeader(h)
	for k, v := range h {
		vv, ok := res[http.CanonicalHeaderKey(k)]
		assert.Truef(ok, "Missing header %s (from %s)", http.CanonicalHeaderKey(k), k)
		if ok {
			assert.Equalf(v, vv, "Invalid header value %s (from %s)", http.CanonicalHeaderKey(k), k)
		}
	}
}
