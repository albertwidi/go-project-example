package session

import "context"

type contextKey string

const sessionContextKey contextKey = "session:context:key"

// WithSession return a context with appended session data
func WithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, sess)
}

// FromContext return a session instance from requestContext
func FromContext(ctx context.Context) *Session {
	sess, ok := ctx.Value(sessionContextKey).(*Session)
	if !ok {
		// do something when session not found
	}
	return sess
}
