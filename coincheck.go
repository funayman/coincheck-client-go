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
	BaseURL = "https://coincheck.com/api"

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

	Rate    *Rate
	Account *Account
}

//New returns a new instance of the CoicCheck client with the specified API key & secret
func New(key, secret string) *Client {
	c := &Client{
		apiKey:     key,
		apiSecret:  secret,
		httpClient: http.Client{},
	}

	c.Rate = &Rate{client: c}
	c.Account = &Account{client: c}

	return c
}

func createSignature(nonce, url, secret, body string) string {
	message := fmt.Sprintf("%s%s%s", nonce, url, body)
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

//DoRequest create a request for the given endpoint
func (client Client) DoRequest(method, endpoint string, content map[string]string) (body io.Reader, err error) {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
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

	req, err := http.NewRequest(method, BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ACCESS-KEY", client.apiKey)
	req.Header.Add("ACCESS-NONCE", nonce)
	req.Header.Add("ACCESS-SIGNATURE", createSignature(nonce, endpoint, client.apiSecret, data.Encode()))

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
