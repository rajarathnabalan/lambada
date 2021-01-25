package lambada

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	assert := assert.New(t)
	req := new(Request)

	ctx := WithRequest(context.TODO(), req)
	assert.Same(req, GetRequest(httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)))

	assert.Nil(GetRequest(httptest.NewRequest(http.MethodGet, "/", nil)))
}
