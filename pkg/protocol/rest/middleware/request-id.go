package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

// Key to use when setting thee request ID
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context
const RequestIDKey ctxKeyRequestID = 0

var (
	// prefix is const prefix for request ID
	prefix string

	// reqID is counter for request ID
	reqID uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// RequestID is a middleware that injects a request ID into the context of
// each request. A Request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies thi go
// process,and where the last number is an atomically increment request counterfunc AddRequestID(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		myid := atomic.AddUint64(&reqID, 1)
//		ctx := r.Context()
//		ctx = context.WithValue(ctx, RequestIDKey, fmt.Sprintf("%s-%06d", prefix, myid))
//		h.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
func AddRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myid := atomic.AddUint64(&reqID, 1)
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDKey, fmt.Sprintf("%s-%06d", prefix, myid))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetReqID returns a request ID from the given context if one is present .
// Returns the empty string if a request ID cannot be found
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqId,ok := ctx.Value(RequestIDKey).(string); ok {
		return reqId
	}
	return ""
}


