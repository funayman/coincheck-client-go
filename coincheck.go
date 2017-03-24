package coincheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CoinCheck struct {
	apiKey     string
	apiSecret  string
	httpClient http.Client
}

func New(key, secret string) *CoinCheck {
	return &CoinCheck{key, secret, http.Client{}}
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

//Public Api
type Ticker struct {
	Last      int       `json:"last"`
	Bid       int       `json:"bid"`
	Ask       int       `json:"ask"`
	High      int       `json:"high"`
	Low       int       `json:"low"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
	Raw       []byte
}

func (t *Ticker) UnmarshalJSON(b []byte) (err error) {
	type alias Ticker
	tmp := &struct {
		Timestamp int64 `json:"timestamp"`
		*Ticker
	}{
		Ticker: (*Ticker)(t),
	}

	if err = json.Unmarshal(b, tmp); err != nil {
		return err
	}

	t.Timestamp = time.Unix(tmp.Timestamp, 0)
	t.Raw = b
	return nil
}

func (client *CoinCheck) Ticker() (t Ticker, err error) {
	url := "https://coincheck.com/api/ticker"
	req := createRequest("GET", url, nil, client.apiKey, client.apiSecret)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return t, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}
