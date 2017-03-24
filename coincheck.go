//Package coincheck provides minimal interaction with the CoinCheck API (https://coincheck.com/documents/exchange/api)
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

//CoinCheck is a client for the CoinCheck Api
//apiKey and apiSecret are required for non-public endpoints (e.g. Orders, Account, etc)
type CoinCheck struct {
	apiKey     string
	apiSecret  string
	httpClient http.Client
}

//New returns a new instance of the CoicCheck client with the specified API key & secret
func New(key, secret string) *CoinCheck {
	return &CoinCheck{key, secret, http.Client{}}
}

func createSignature(nonce int64, url, secret string) string {
	message := fmt.Sprintf("%d%s%s", nonce, url, "")
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

func createRequest(method, url string, data url.Values, key, secret string) (req *http.Request, err error) {
	nonce := time.Now().UnixNano()

	req, err = http.NewRequest(method, url, bytes.NewBufferString(data.Encode()))
	req.Header.Add("ACCESS-KEY", key)
	req.Header.Add("ACCESS-NONCE", strconv.FormatInt(nonce, 10))
	req.Header.Add("ACCESS-SIGNATURE", createSignature(nonce, url, secret))

	return
}

//Public Api

//Ticker holds all information provided by the Ticker endpoint
//Timestamp is converted into a time.Time type rather than a Unix Timestamp
//Raw houses the original JSON request data
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
	type Alias Ticker
	tmp := &struct {
		Timestamp int64  `json:"timestamp"`
		Volume    string `json:"volume"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err = json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	v, err := strconv.ParseFloat(tmp.Volume, 64)
	if err != nil {
		return err
	}

	t.Timestamp = time.Unix(tmp.Timestamp, 0)
	t.Volume = v
	t.Raw = b
	return nil
}

//Ticker checks the latest information
//More information can be found at https://coincheck.com/documents/exchange/api#ticker
func (client *CoinCheck) Ticker() (t Ticker, err error) {
	url := "https://coincheck.com/api/ticker"
	req, err := createRequest("GET", url, nil, client.apiKey, client.apiSecret)
	if err != nil {
		return t, err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return t, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}
