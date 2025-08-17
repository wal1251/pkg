package mw

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

// AddIP создает новый rateLimiter и добавляет его в мапу ips с IP-адресом в виде ключа.
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter возвращает rateLimiter для предоставленного IP-адреса, если он уже иммется.
// В противном случае он вызывает AddIP, чтобы добавить IP-адрес на мапу.
func (i *IPRateLimiter) GetLimiter(requestIP string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[requestIP]
	i.mu.RUnlock()

	if !exists {
		i.mu.Lock()
		defer i.mu.Unlock()
		limiter, exists = i.ips[requestIP]
		if !exists {
			limiter = rate.NewLimiter(i.r, i.b)
			i.ips[requestIP] = limiter
		}
	}

	return limiter
}

func RateLimiter(limit int, per time.Duration) func(http.Handler) http.Handler {
	var (
		rateLimiters  = make(map[string]*IPRateLimiter)
		rateLimiterMu = sync.RWMutex{}
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestPath := r.URL.Path

			rateLimiterMu.RLock()
			ipRateLimiter, exists := rateLimiters[requestPath]
			rateLimiterMu.RUnlock()

			if !exists {
				rateLimiterMu.Lock()
				defer rateLimiterMu.Unlock()

				ipRateLimiter, exists = rateLimiters[requestPath]
				if !exists {
					ipRateLimiter = NewIPRateLimiter(rate.Every(per/time.Duration(limit)), limit)
					rateLimiters[requestPath] = ipRateLimiter
				}
			}
			// Получаем IP-адрес клиента
			clientIP := getIP(r)

			// Получаем rateLimiter для IP-адреса
			limiter := ipRateLimiter.GetLimiter(clientIP)
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// getIp извлекает IP-адрес клиента из запроса.
func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
