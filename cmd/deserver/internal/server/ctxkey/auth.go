package ctxkey

import "context"

type usernameCtxKey struct{}

var usernameKey usernameCtxKey

// WithUsername adds a user's username into the context.
func WithUsername(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, usernameKey, s)
}

func UsernameFromCtx(ctx context.Context) string {
	str, ok := ctx.Value(usernameKey).(string)
	if !ok {
		return ""
	}
	return str
}
