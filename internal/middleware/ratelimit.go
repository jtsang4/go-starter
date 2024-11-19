package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jtsang4/go-stater/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)
		if !l.Allow() {
			logger.Logger.Warn("rate limit exceeded",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
