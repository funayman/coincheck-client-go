//Package coincheck provides minimal interaction with the CoinCheck API (https://coincheck.com/documents/exchange/api)
package coincheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	//BaseURL Coincheck API URL
	BaseURL = "https://coincheck.com/api"

	//SendingFee to send BTC out of CoinCheck
	SendingFee = 0.0005

	BtcJpy = "btc_jpy"
	EthJpy = "eth_jpy"
	EthBtc = "eth_btc"
	EtcJpy = "etc_jpy"
	EtcBtc = "etc_btc"
	DaoJpy = "dao_jpy"
	LskJpy = "lsk_jpy"
	LskBtc = "lsk_btc"
	FctJpy = "fct_jpy"
	FtcBtc = "fct_btc"
	XmrJpy = "xmr_jpy"
	XmrBtc = "xmr_btc"
	RepJpy = "rep_jpy"
	RepBtc = "rep_btc"
	XrpJpy = "xrp_jpy"
	XprBtc = "xrp_btc"
	ZecJpy = "zec_jpy"
	ZecBtc = "zec_btc"
)

//IClient a client interface
type IClient interface {
	DoRequest(method, endpoint string, content map[string]string) (io.Reader, error)
}

//Client is a client for the CoinCheck Api
//apiKey and apiSecret are required for non-public endpoints (e.g. Orders, Account, etc)
type Client struct {
	apiKey     string
	apiSecret  string
	httpClient http.Client
}

//New returns a new instance of the CoicCheck client with the specified API key & secret
func New(key, secret string) *Client {
	return &Client{
		apiKey:     key,
		apiSecret:  secret,
		httpClient: http.Client{},
	}
}

func createSignature(nonce, url, secret, body string) string {
	message := fmt.Sprintf("%s%s%s", nonce, url, body)
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

//DoRequest create a request for the given endpoint
func (client *Client) DoRequest(method, endpoint string, content map[string]string) (io.Reader, error) {
	var data string
	var body io.Reader

	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)

	switch method {
	case "", "GET":
		m := url.Values{}
		for key, value := range content {
			m.Add(key, value)
		}
		data = m.Encode()
		if content != nil {
			endpoint = endpoint + "?" + data
		}
	case "POST", "DELETE":
		b, _ := json.Marshal(content)
		data = string(b)
		body = bytes.NewReader(b)
	default:
		return nil, errors.New("Invalid method (" + method + ")")

	}

	req, err := http.NewRequest(method, BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ACCESS-KEY", client.apiKey)
	req.Header.Add("ACCESS-NONCE", nonce)
	req.Header.Add("ACCESS-SIGNATURE", createSignature(nonce, BaseURL+endpoint, client.apiSecret, data))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		tmp := &struct {
			Message string `json:"error"`
		}{}
		json.NewDecoder(resp.Body).Decode(tmp)
		msg := strings.Replace(tmp.Message, "\n", " | ", -1)
		endPointError := fmt.Sprintf("StatusCode[%d], Error[%s]", resp.StatusCode, msg)
		return nil, errors.New(endPointError)
	}

	return resp.Body, nil
}
