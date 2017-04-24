package coincheck

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

//Ticker holds all information provided by the Ticker endpoint
//Timestamp is converted into a time.Time type rather than a Unix Timestamp
//Raw houses the original JSON request data as []byte
type Ticker struct {
	Last      int `json:"last"`
	Bid       int `json:"bid"`
	Ask       int `json:"ask"`
	High      int `json:"high"`
	Low       int `json:"low"`
	Volume    float64
	Timestamp time.Time
	Raw       []byte
}

//UnmarshalJSON is used to parse timestamp into time.Time and volume into float64
func (t *Ticker) UnmarshalJSON(b []byte) (err error) {
	type Alias Ticker
	tmp := &struct {
		Timestamp int64  `json:"timestamp"`
		Volume    string `json:"volume"`
		*Alias
	}{Alias: (*Alias)(t)}

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

//Trade JSON object in array of the Public trades endpoint
//More information can be found at https://coincheck.com/documents/exchange/api#public-trades
type Trade struct {
	ID        int    `json:"id"`
	Amount    string `json:"amount"`
	Rate      int64  `json:"rate"`
	OrderType string `json:"order_type"`
	CreatedAt time.Time
	Raw       []byte
}

//Trades data retruned from the Public trades endpoint
//More information can be found at https://coincheck.com/documents/exchange/api#public-trades
type Trades []Trade

//OrderBook weird ass data structure returned from the Order book endpoint
type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Raw  []byte
}

//Ticker gets the latest information and returns a Ticker
//More information can be found at https://coincheck.com/documents/exchange/api#ticker
func (client Client) Ticker() (t Ticker, err error) {
	url := "/ticker"
	body, err := client.DoRequest("GET", url, nil)
	if err != nil {
		return
	}

	if err = json.NewDecoder(body).Decode(&t); err != nil {
		return
	}
	return
}

//PublicTrades You can get current order transactions
//More information can be found at https://coincheck.com/documents/exchange/api#public-trades
func (client *Client) PublicTrades(offset string) (trades Trades, err error) {
	endpoint := "/trades"
	body, err := client.DoRequest("GET", endpoint, map[string]string{"offset": offset})
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&trades)
	if err != nil {
		return
	}

	return
}

//OrderBook Fetch order book information.
//More information can be found at https://coincheck.com/documents/exchange/api#order-book
func (client *Client) OrderBook() (orders OrderBook, err error) {
	endpoint := "/order_books"

	body, err := client.DoRequest("GET", endpoint, nil)
	if err != nil {
		return
	}

	orders.Raw, err = ioutil.ReadAll(body)
	if err != nil {
		return
	}

	err = json.Unmarshal(orders.Raw, &orders)
	if err != nil {
		return
	}

	return
}

//UnmarshalJSON custom Unmarshaler for Trade
func (t *Trade) UnmarshalJSON(b []byte) error {
	type Alias Trade
	tmp := &struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{Alias: (*Alias)(t)}

	err := json.Unmarshal(b, tmp)
	if err != nil {
		return err
	}

	ca, err := time.Parse(time.RFC3339, tmp.CreatedAt)
	if err != nil {
		return err
	}

	t.CreatedAt = ca
	t.Raw = b

	return nil
}

//Rate returns the price of the coin selected in the currency selected
//BTC_JPY would return the rate of BitCoin in Japanese Yen
//ETH_BTC would return the rate of Ethereum in BitCoins
//More information can be found at https://coincheck.com/documents/exchange/api#buy-rate
func (client *Client) Rate(pair string) (float64, error) {
	url := "/rate/" + pair

	body, err := client.DoRequest("GET", url, nil)
	if err != nil {
		return 0.0, err
	}

	tmp := &struct {
		Rate string `json:"rate"`
	}{}

	err = json.NewDecoder(body).Decode(tmp)
	if err != nil {
		return 0.0, err
	}

	rate, err := strconv.ParseFloat(tmp.Rate, 64)
	if err != nil {
		return 0.0, err
	}

	return rate, nil
}
