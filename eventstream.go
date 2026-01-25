package www

import (
	"net/http"
	"sync"
)

// Send Content-Type:text/event-stream.
// Note: Use NewEventStreamWriter instead of StartEventStream if you're going to be
// writing events from multiple threads.
func StartEventStream(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
}

// WriteEvent writes a single event to the event stream.
// Note: This is not thread safe. Use EventStreamWriter if you're going to be writing
// events from multiple threads.
func WriteEvent(w http.ResponseWriter, event string, data string) error {
	_, err := w.Write([]byte("event: " + event + "\n"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("data: " + data + "\n\n"))
	if err != nil {
		return err
	}
	// Flush the data immediately
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

// A thread-safe event stream writer
type EventStreamWriter struct {
	W    http.ResponseWriter
	lock sync.Mutex
}

// Start an event stream, and return a thread-safe writer
func NewEventStreamWriter(w http.ResponseWriter) *EventStreamWriter {
	StartEventStream(w)
	return &EventStreamWriter{W: w}
}

// WriteEvent writes a single event to the event stream in a thread-safe manner
func (w *EventStreamWriter) WriteEvent(event string, data string) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return WriteEvent(w.W, event, data)
}
