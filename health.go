package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

const (
	readinessCacheTTL     = 30 * time.Second
	readinessCheckTimeout = 10 * time.Second
)

// healthzHandler is a liveness check. It only proves the process is serving requests and
// deliberately does not depend on the Atlas API.
func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		logger.Errorf("Failed to write response: %v", err)
	}
}

// readinessCheck reports whether the Atlas API is reachable. The result is cached so that
// frequent probes do not turn into frequent upstream requests, and the upstream call uses its
// own context so that a probe disconnecting does not cancel it.
type readinessCheck struct {
	check func(context.Context) error
	ttl   time.Duration

	mu        sync.Mutex
	checkedAt time.Time
	err       error
}

func newReadinessCheck(check func(context.Context) error, ttl time.Duration) *readinessCheck {
	return &readinessCheck{check: check, ttl: ttl}
}

func (r *readinessCheck) ready() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.checkedAt.IsZero() && time.Since(r.checkedAt) < r.ttl {
		return r.err
	}
	ctx, cancel := context.WithTimeout(context.Background(), readinessCheckTimeout)
	defer cancel()
	r.err = r.check(ctx)
	r.checkedAt = time.Now()
	if r.err != nil {
		logger.Errorf("Readiness check failed: %v", r.err)
	}
	return r.err
}

func (r *readinessCheck) handler(w http.ResponseWriter, _ *http.Request) {
	if err := r.ready(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, writeErr := w.Write([]byte("Atlas API not reachable")); writeErr != nil {
			logger.Errorf("Failed to write response: %v", writeErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		logger.Errorf("Failed to write response: %v", err)
	}
}
