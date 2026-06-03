package hmac

import (
	ghmac "crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// SignRequest generates an HMAC-SHA256 signature for an HTTP request.
// It returns the hex-encoded signature and the Unix timestamp used.
func SignRequest(method, path, secret string) (signature, timestamp string) {
	ts := time.Now().Unix()
	timestamp = strconv.FormatInt(ts, 10)
	signature = sign(method, path, secret, timestamp)
	return
}

// SignRequestWithTimestamp generates a signature using a fixed timestamp.
// This is exposed for testing against contract vectors.
func SignRequestWithTimestamp(method, path, secret, timestamp string) string {
	return sign(method, path, secret, timestamp)
}

func sign(method, path, secret, timestamp string) string {
	payload := strings.ToUpper(method) + path + timestamp
	mac := ghmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
