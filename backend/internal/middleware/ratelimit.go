package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type client struct {
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
			ip = ip[:colonIndex]
		}

		mu.Lock()
		if c, exists := clients[ip]; exists {
			if time.Since(c.lastSeen) < 500*time.Millisecond {
				mu.Unlock()
				// Send JSON error instead of plain text
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]string{
						"code":    "RATE_LIMIT_EXCEEDED",
						"message": "TOO MANY SIGNALS. SLOW DOWN.",
					},
				})
				return
			}
		}
		clients[ip] = &client{lastSeen: time.Now()}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
