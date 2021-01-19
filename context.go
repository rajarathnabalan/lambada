package lambada

import "context"

type contextKeyType struct{}

var contextKey = contextKeyType{}

func GetRequest(ctx context.Context) *Request {
	if res, ok := ctx.Value(contextKey).(*Request); ok {
		return res
	}
	return nil
}

func withRequest(ctx context.Context, req *Request) context.Context {
	return context.WithValue(ctx, contextKey, req)
}
