package webhookx

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type limiterStore struct {
	mu      sync.Mutex
	entries map[string]*limiterEntry
	limit   rate.Limit
	burst   int
}

func newLimiterStore(limit rate.Limit, burst int) *limiterStore {
	return &limiterStore{
		entries: make(map[string]*limiterEntry),
		limit:   limit,
		burst:   burst,
	}
}

func (s *limiterStore) allow(clientID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.entries[clientID]
	if !ok {
		entry = &limiterEntry{limiter: rate.NewLimiter(s.limit, s.burst)}
		s.entries[clientID] = entry
	}
	entry.lastSeen = time.Now()

	if len(s.entries) > 10000 {
		cutoff := time.Now().Add(-30 * time.Minute)
		for key, val := range s.entries {
			if val.lastSeen.Before(cutoff) {
				delete(s.entries, key)
			}
		}
	}

	return entry.limiter.Allow()
}
