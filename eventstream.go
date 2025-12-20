package www

import "net/http"

// Send Content-Type:text/event-stream
func StartEventStream(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
}

// WriteEvent writes a single event to the event stream
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
