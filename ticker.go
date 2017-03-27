package coincheck

import (
	"encoding/json"
	"strconv"
	"time"
)

//Ticker holds all information provided by the Ticker endpoint
//Timestamp is converted into a time.Time type rather than a Unix Timestamp
//Raw houses the original JSON request data as []byte
type Ticker struct {
	Last      int       `json:"last"`
	Bid       int       `json:"bid"`
	Ask       int       `json:"ask"`
	High      int       `json:"high"`
	Low       int       `json:"low"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
	Raw       []byte

	client *Client
}

//UnmarshalJSON is used to parse timestamp into time.Time and volume into float64
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

//Update gets the latest information and updates the struct
//More information can be found at https://coincheck.com/documents/exchange/api#ticker
func (t *Ticker) Update() (err error) {
	url := "https://coincheck.com/api/ticker"
	body, err := t.client.DoRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(body).Decode(&t); err != nil {
		return err
	}
	return nil
}
