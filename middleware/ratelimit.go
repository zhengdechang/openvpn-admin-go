package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipBucket 记录每个 IP 的请求桶
type ipBucket struct {
	count    int
	resetAt  time.Time
	mu       sync.Mutex
}

var (
	buckets   = make(map[string]*ipBucket)
	bucketsMu sync.Mutex
)

// getOrCreateBucket 获取或创建 IP 对应的令牌桶
func getOrCreateBucket(ip string, window time.Duration) *ipBucket {
	bucketsMu.Lock()
	defer bucketsMu.Unlock()

	b, ok := buckets[ip]
	if !ok || time.Now().After(b.resetAt) {
		b = &ipBucket{resetAt: time.Now().Add(window)}
		buckets[ip] = b
	}
	return b
}

// RateLimit 对指定路由限制每个 IP 在 window 时间内最多 maxReqs 次请求
func RateLimit(maxReqs int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		b := getOrCreateBucket(ip, window)

		b.mu.Lock()
		if time.Now().After(b.resetAt) {
			b.count = 0
			b.resetAt = time.Now().Add(window)
		}
		b.count++
		count := b.count
		b.mu.Unlock()

		if count > maxReqs {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
