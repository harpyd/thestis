package correlationid

import "context"

type correlationCtx int

const correlationCtxKey correlationCtx = iota

func AssignToCtx(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationCtxKey, correlationID)
}

func FromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if correlationID, ok := ctx.Value(correlationCtxKey).(string); ok {
		return correlationID
	}

	return ""
}
