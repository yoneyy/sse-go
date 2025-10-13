package sse

import (
	"net/http"
	"testing"
)

func TestNewSSE(t *testing.T) {
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sse := NewSSE(w)
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Err("internal error")
		sse.Done()
	})
}

func TestSetWriter(t *testing.T) {
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sse := NewSSE(w)

		w.Header().Set("Cache-Control", "no-cache")
		sse.SetWriter(w)
	})
}

func TestSetHeader(t *testing.T) {
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sse := NewSSE(w)
		sse.SetHeader("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0, no-transform")
	})
}
