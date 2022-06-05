package correlationid

import "context"

type correlationCtx int

const correlationCtxKey correlationCtx = iota

func AssignToCtx(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationCtxKey, correlationID)
}

func FromCtx(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	correlationID, ok := ctx.Value(correlationCtxKey).(string)

	return correlationID, ok
}
