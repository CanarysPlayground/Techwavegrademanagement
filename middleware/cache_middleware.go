package middleware

import (
	"context"
	"net/http"
)

// CacheStatus represents cache hit/miss status
type CacheStatus string

const (
	CacheHit  CacheStatus = "HIT"
	CacheMiss CacheStatus = "MISS"
	CacheSkip CacheStatus = "SKIP"
)

// cacheContextKey is the key for storing cache status in request context
type cacheContextKey struct{}

// SetCacheStatus sets the cache status in the request context
func SetCacheStatus(r *http.Request, status CacheStatus) *http.Request {
	ctx := context.WithValue(r.Context(), cacheContextKey{}, status)
	return r.WithContext(ctx)
}

// GetCacheStatus retrieves the cache status from request context
func GetCacheStatus(r *http.Request) CacheStatus {
	if status, ok := r.Context().Value(cacheContextKey{}).(CacheStatus); ok {
		return status
	}
	return CacheSkip
}

// CacheStatusMiddleware adds X-Cache-Status header to responses
func CacheStatusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap response writer to capture status
		wrapped := &cacheStatusWriter{
			ResponseWriter: w,
			request:        r,
		}
		
		next.ServeHTTP(wrapped, r)
		
		// Add cache status header after handler completes
		status := GetCacheStatus(r)
		w.Header().Set("X-Cache-Status", string(status))
	})
}

// cacheStatusWriter wraps http.ResponseWriter
type cacheStatusWriter struct {
	http.ResponseWriter
	request *http.Request
}

func (w *cacheStatusWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
