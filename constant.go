package sse

const (
	SSEEventError   SSEEvent = "error"
	SSEEventMessage SSEEvent = "message"
)

const (
	ID    = "id: %s\n"
	Event = "event: %s\n"
	Data  = "data: %s\n\n"
	Done  = "[DONE]"
)
