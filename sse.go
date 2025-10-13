package sse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var _ SSE = (*sse)(nil)

type SSE interface {
	Encode(msg *Message) error
	Err(err SSEData)
	Data(data SSEData)
	Done()
	SetRetry(retry uint32) SSE
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
		w http.ResponseWriter
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
func (s *sse) SetRetry(retry uint32) SSE {
	fmt.Fprintf(s.w, Retry, retry)
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

	buf := fmt.Appendf(nil, ID, msg.ID)
	buf = fmt.Appendf(buf, Data, msg.Data)
	buf = fmt.Appendf(buf, Event, msg.Event)
	s.w.Write(buf)
	s.flush()
	return nil
}

// Err
// send error messages
func (s *sse) Err(err SSEData) {
	msg := &Message{
		Event: SSEEventError,
		Data:  err,
	}
	s.Encode(msg)
}

// Data
// send normal messages
func (s *sse) Data(data SSEData) {
	msg := &Message{
		Event: SSEEventMessage,
		Data:  data,
	}
	s.Encode(msg)
}

// Done
// send the DONE message
func (s *sse) Done() {
	msg := &Message{
		Event: SSEEventMessage,
		Data:  Done,
	}
	s.Encode(msg)
}

// ============================ private methods ============================

func (s *sse) flush() {
	if flush, ok := s.w.(http.Flusher); ok {
		flush.Flush()
	}
}

func (s *sse) dft() {
	s.w.Header().Set("Connection", "keep-alive")
	s.w.Header().Set("Transfer-Encoding", "chunked")
	s.w.Header().Set("Content-Type", "text/event-stream")
	s.w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0, no-transform")
}
