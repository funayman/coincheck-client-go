//Package coincheck provides minimal interaction with the CoinCheck API (https://coincheck.com/documents/exchange/api)
package coincheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/funayman/coincheck-client-go/errors"
)

const (
	BUY  = "buy"
	SELL = "sell"

	BTC_JPY = "btc_jpy"
	ETH_JPY = "eth_jpy"
	ETH_BTC = "eth_btc"
	ETC_JPY = "etc_jpy"
	ETC_BTC = "etc_btc"
	DAO_JPY = "dao_jpy"
	LSK_JPY = "lsk_jpy"
	LSK_BTC = "lsk_btc"
	FCT_JPY = "fct_jpy"
	FTC_BTC = "fct_btc"
	XMR_JPY = "xmr_jpy"
	XMR_BTC = "xmr_btc"
	REP_JPY = "rep_jpy"
	REP_BTC = "rep_btc"
	XRP_JPY = "xrp_jpy"
	XPR_BTC = "xrp_btc"
	ZEC_JPY = "zec_jpy"
	ZEC_BTC = "zec_btc"
)

type IClient interface {
	DoRequest(method, endpoint string, content map[string]string) (io.Reader, error)
}

//Client is a client for the CoinCheck Api
//apiKey and apiSecret are required for non-public endpoints (e.g. Orders, Account, etc)
type Client struct {
	apiKey     string
	apiSecret  string
	httpClient http.Client

	Ticker *Ticker
	Rate   *Rate
}

//New returns a new instance of the CoicCheck client with the specified API key & secret
func New(key, secret string) *Client {
	c := &Client{
		apiKey:     key,
		apiSecret:  secret,
		httpClient: http.Client{},
	}

	c.Ticker = &Ticker{client: c}
	c.Rate = &Rate{client: c}

	return c
}

func createSignature(nonce int64, url, secret string) string {
	message := fmt.Sprintf("%d%s%s", nonce, url, "")
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

//DoRequest create a request for the given endpoint
func (client Client) DoRequest(method, endpoint string, content map[string]string) (io.Reader, error) {
	var body io.Reader
	nonce := time.Now().UnixNano()
	data := url.Values{}
	for key, value := range content {
		data.Add(key, value)
	}

	switch method {
	case "", "GET":
		if content != nil {
			endpoint = endpoint + "?" + data.Encode()
		}
	case "POST", "DELETE":
		body = bytes.NewBufferString(data.Encode())
	default:
		return nil, errors.NewGenericError("Invalid method (" + method + ")")

	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("ACCESS-KEY", client.apiKey)
	req.Header.Add("ACCESS-NONCE", strconv.FormatInt(nonce, 10))
	req.Header.Add("ACCESS-SIGNATURE", createSignature(nonce, endpoint, client.apiSecret))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		tmp := &struct {
			Message string `json:"error"`
		}{}
		json.NewDecoder(resp.Body).Decode(tmp)
		return nil, errors.NewEndPointError(tmp.Message)
	}

	return resp.Body, nil
}
