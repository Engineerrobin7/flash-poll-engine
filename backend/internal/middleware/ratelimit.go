package middleware

import (
	"net/http"
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

		mu.Lock()
		if c, exists := clients[ip]; exists {
			// Limit to 1 request every 500ms
			if time.Since(c.lastSeen) < 500*time.Millisecond {
				mu.Unlock()
				http.Error(w, "TOO MANY SIGNALS. SLOW DOWN.", http.StatusTooManyRequests)
				return
			}
		}
		clients[ip] = &client{lastSeen: time.Now()}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
