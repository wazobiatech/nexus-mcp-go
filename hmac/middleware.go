package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	headerSignature = "x-signature"
	headerTimestamp = "x-timestamp"
	maxAgeSeconds   = 300
)

// Middleware returns an HTTP middleware that validates HMAC-SHA256 signatures.
func Middleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sig := r.Header.Get(headerSignature)
			ts := r.Header.Get(headerTimestamp)

			if sig == "" || ts == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			timestamp, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			now := time.Now().Unix()
			if now-timestamp > maxAgeSeconds || timestamp-now > maxAgeSeconds {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			expected := sign(r.Method, r.URL.RequestURI(), secret, ts)
			if !hmac.Equal([]byte(sig), []byte(expected)) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
