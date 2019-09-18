package middleware

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
	"time"
)

var limiter *rate.Limiter = DefaultLimiter()

func DefaultLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Minute), 10)
}

func SetLimiter(limit *rate.Limiter) {
	limiter = limit
}

//RateLimiter using global limiter
func RateLimiter() endpoint.Middleware {
	return ratelimit.NewErroringLimiter(limiter)
}

var limiters map[string]*rate.Limiter

//RateLimiterWithAuthToken create limiter by given key
func RateLimiterWithAuthToken(key string, limit *rate.Limiter) endpoint.Middleware {
	if nil == limit {
		limit = DefaultLimiter()
	}

	if _, ok := limiters[key]; !ok {
		limiters[key] = limit
	}

	return ratelimit.NewErroringLimiter(limiters[key])
}
