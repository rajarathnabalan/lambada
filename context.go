package lambada

import (
	"context"
	"net/http"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

func GetRequest(r *http.Request) *Request {
	if res, ok := r.Context().Value(contextKey).(*Request); ok {
		return res
	}
	return nil
}

func WithRequest(ctx context.Context, req *Request) context.Context {
	return context.WithValue(ctx, contextKey, req)
}
