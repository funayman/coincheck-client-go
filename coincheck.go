package coincheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type coinCheck struct {
	apiKey    string
	apiSecret string
}

func New(key, secret string) *coinCheck {
	return &coinCheck{key, secret}
}

func createSignature(nonce int64, url, secret string) string {
	message := fmt.Sprintf("%d%s%s", nonce, url, "")
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

func createRequest(method, url string, data url.Values, key, secret string) (req *http.Request) {
	nonce := time.Now().UnixNano()

	req, _ = http.NewRequest(method, url, bytes.NewBufferString(data.Encode()))

	req.Header.Add("ACCESS-KEY", key)
	req.Header.Add("ACCESS-NONCE", strconv.FormatInt(nonce, 10))
	req.Header.Add("ACCESS-SIGNATURE", createSignature(nonce, url, secret))
	return
}
