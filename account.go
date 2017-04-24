package coincheck

import (
	"encoding/json"
	"strconv"
	"time"
)

//Transaction response from CoinCheck API
type Transaction struct {
	ID          int64     `json:"id"`
	Success     bool      `json:"success"`
	Address     string    `json:"address"`
	Amount      string    `json:"amount"`
	Fee         string    `json:"fee"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}

//Account response from CoinCheck API
type Account struct {
	ID              int64  `json:"id"`
	Email           string `json:"email"`
	IdentityStatus  string `json:"identity_status"`
	Address         string `json:"bitcoin_address"`
	LendingLeverate string `json:"lending_leverage"`
	TakerFee        string `json:"taker_fee"`
	MakerFee        string `json:"maker_fee"`
}

//Balance response from CoinCheck API
type Balance struct {
	JPY          float64
	BTC          float64
	JPYReserved  string `json:"jpy_reserved"`
	BTCReserved  string `json:"btc_reserved"`
	JPYLendInUse string `json:"jpy_lend_in_use"`
	BTCLendInUse string `json:"btc_lend_in_use"`
	JPYLent      string `json:"jpy_lent"`
	BTCLent      string `json:"btc_lent"`
	JPYDebt      string `json:"jpy_debt"`
	BTCDebt      string `json:"btc_debt"`
}

//UnmarshalJSON custom unmarshaling to convert JPY and BTC to float64
func (bal *Balance) UnmarshalJSON(b []byte) (err error) {
	type alias Balance
	tmp := &struct {
		JPY string `json:"jpy"`
		BTC string `json:"btc"`
		*alias
	}{alias: (*alias)(bal)}

	err = json.Unmarshal(b, tmp)
	if err != nil {
		return
	}

	fJpy, err := strconv.ParseFloat(tmp.JPY, 64)
	if err != nil {
		return
	}
	fBtc, err := strconv.ParseFloat(tmp.BTC, 64)
	if err != nil {
		return
	}

	bal.JPY = fJpy
	bal.BTC = fBtc
	return
}

//AccountInfo of the user
func (client *Client) AccountInfo() (a Account, err error) {
	url := "/accounts"
	body, err := client.DoRequest("GET", url, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&a)
	return
}

//AccountBalance of the user (BTC, JPY, etc)
func (client *Client) AccountBalance() (b Balance, err error) {
	url := "/accounts/balance"
	body, err := client.DoRequest("GET", url, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&b)
	return
}

//SendBTC to the specified address
func (client *Client) SendBTC(address string, amount float64) (t Transaction, err error) {
	url := "/send_money"
	content := map[string]string{
		"address": address,
		"amount":  strconv.FormatFloat(amount, 'G', -1, 64),
	}

	body, err := client.DoRequest("POST", url, content)
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&t)
	return
}
