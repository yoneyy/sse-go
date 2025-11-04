package sse

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var _ SSE = (*sse)(nil)

type SSE interface {
	Encode(msg *Message) error
	Err(err SSEData) SSE
	Data(data SSEData) SSE
	Done() SSE
	SetRetry(retry time.Duration) SSE
	SetHeader(key, value string) SSE
	SetWriter(w http.ResponseWriter) SSE
}

type (
	SSEData  string
	SSEEvent string

	Message struct {
		ID    string
		Data  SSEData
		Event SSEEvent
	}

	sse struct {
		mu sync.Mutex
		w  http.ResponseWriter
	}
)

// SSE
// create a sse server
// @author yoneyy (y.tianyuan)
func NewSSE(w http.ResponseWriter) SSE {
	sse := &sse{w: w}
	sse.dft()
	return sse
}

// SetHeader
// set header or override the default writer header
func (s *sse) SetHeader(key, value string) SSE {
	s.w.Header().Set(key, value)
	return s
}

// SetWriter
// override the default writer
func (s *sse) SetWriter(w http.ResponseWriter) SSE {
	s.w = w
	return s
}

// SetRetry
func (s *sse) SetRetry(retry time.Duration) SSE {
	fmt.Fprintf(s.w, Retry, retry.Milliseconds())
	s.flush()
	return s
}

// Encode
// encode sse messages
func (s *sse) Encode(msg *Message) error {
	if msg == nil {
		return ErrMessageRequired
	}

	if strings.TrimSpace(msg.ID) == "" {
		msg.ID = uuid.NewString()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	buf := fmt.Appendf(nil, ID, msg.ID)
	buf = fmt.Appendf(buf, Event, msg.Event)
	buf = fmt.Appendf(buf, Data, msg.Data)
	s.w.Write(buf)
	s.flush()
	return nil
}

// Err
// send error messages
func (s *sse) Err(err SSEData) SSE {
	msg := &Message{
		Event: SSEEventError,
		Data:  err,
	}
	s.Encode(msg)
	return s
}

// Data
// send normal messages
func (s *sse) Data(data SSEData) SSE {
	msg := &Message{
		Event: SSEEventMessage,
		Data:  data,
	}
	s.Encode(msg)
	return s
}

// Done
// send the DONE message
func (s *sse) Done() SSE {
	msg := &Message{
		Event: SSEEventMessage,
		Data:  Done,
	}
	s.Encode(msg)
	return s
}

// ============================ private methods ============================

func (s *sse) flush(b ...[]byte) {
	if len(b) > 0 && s.w != nil {
		s.w.Write(b[0])
	}

	if flush, ok := s.w.(http.Flusher); ok && flush != nil {
		flush.Flush()
	}
}

func (s *sse) dft() {
	s.w.Header().Set("Connection", "keep-alive")
	s.w.Header().Set("Transfer-Encoding", "chunked")
	s.w.Header().Set("Content-Type", "text/event-stream")
	s.w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0, no-transform")
}
