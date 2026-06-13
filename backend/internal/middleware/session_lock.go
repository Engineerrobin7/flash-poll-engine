package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

// SessionStore tracks IP -> map[PollID]bool
type SessionStore struct {
	mu    sync.RWMutex
	votes map[string]map[string]bool
}

var GlobalSessionStore = &SessionStore{
	votes: make(map[string]map[string]bool),
}

func (s *SessionStore) HasVoted(ip, pollID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if userVotes, exists := s.votes[ip]; exists {
		return userVotes[pollID]
	}
	return false
}

func (s *SessionStore) RecordVote(ip, pollID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.votes[ip]; !exists {
		s.votes[ip] = make(map[string]bool)
	}
	s.votes[ip][pollID] = true
}

// Helper to get clean IP
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
		ip = ip[:commaIndex]
	}
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex]
	}
	return strings.TrimSpace(ip)
}
