package constant

type ContextKey string

var (
	ContextKeyUserID       ContextKey = "user_id"
	ContextKeyRequestID    ContextKey = "request_id"
	ContextKeyIncomingTime ContextKey = "incoming_time"
)
