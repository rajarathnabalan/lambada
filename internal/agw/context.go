package agw

import "context"

type contextKeyType struct{}

var contextKey = contextKeyType{}

func WithRequest(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, contextKey, v)
}

func GetStageVariables(ctx context.Context) map[string]string {
	if res, ok := ctx.Value(contextKey).(map[string]string); ok {
		return res
	}
	return nil
}
