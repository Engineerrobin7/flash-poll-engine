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
		// BUG BOUNTY FIX: Respect X-Forwarded-For for Cloud Deployments (Render/Vercel)
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		// Clean the IP (remove ports or multiple proxy hops)
		if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
			ip = ip[:commaIndex]
		}
		if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
			ip = ip[:colonIndex]
		}
		ip = strings.TrimSpace(ip)

		mu.Lock()
		if c, exists := clients[ip]; exists {
			if time.Since(c.lastSeen) < 500*time.Millisecond {
				mu.Unlock()
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
