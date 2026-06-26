package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// sseHeartbeatInterval keeps idle connections alive through proxies and lets us
// notice a dead client (the write fails on the next heartbeat).
const sseHeartbeatInterval = 25 * time.Second

// handleEvents streams real-time events to the client over Server-Sent Events.
//
// This route is deliberately registered outside the logging/Timeout middleware
// group: a request Timeout would close the stream after 60s, and the logging
// wrapper can stop us from clearing the connection write deadline.
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	if s.eventBroker == nil {
		http.Error(w, "events not available", http.StatusServiceUnavailable)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// Disable response buffering in nginx/Traefik so events are delivered live.
	w.Header().Set("X-Accel-Buffering", "no")

	// The server's WriteTimeout would otherwise abort the stream; clear the
	// write deadline for this long-lived connection.
	rc := http.NewResponseController(w)
	if err := rc.SetWriteDeadline(time.Time{}); err != nil {
		s.log.Warn().Err(err).Msg("SSE: failed to clear write deadline; stream may drop")
	}

	ch, unsubscribe := s.eventBroker.Subscribe()
	defer unsubscribe()

	// Open the stream so the client's EventSource fires onopen.
	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	heartbeat := time.NewTicker(sseHeartbeatInterval)
	defer heartbeat.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-heartbeat.C:
			if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case ev := <-ch:
			data, err := json.Marshal(ev.Data)
			if err != nil {
				s.log.Error().Err(err).Str("event_type", ev.Type).Msg("SSE: failed to marshal event")
				continue
			}
			if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, data); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
