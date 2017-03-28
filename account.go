package coincheck

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/funayman/coincheck-client-go/errors"
)

type Transaction struct {
	ID          string    `json:"id"`
	Success     bool      `json:"success"`
	Address     string    `json:"address"`
	Amount      string    `json:"amount"`
	Fee         string    `json:"fee"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}

type Transactions []Transaction

type BankAccount struct {
	ID     int64  `json:"id"`
	Name   string `json:"bank_name"`
	Branch string `json:"branch_name"`
	Type   string `json:"bank_account_type"`
	Number string `json:"number"`
	Owner  string `json:"name"`
}

type BankAccounts []BankAccount

type Account struct {
	ID              int64  `json:"id"`
	Email           string `json:"email"`
	IdentityStatus  string `json:"identity_status"`
	Address         string `json:"bitcoin_address"`
	LendingLeverate string `json:"lending_leverage"`
	TakerFee        string `json:"taker_fee"`
	MakerFee        string `json:"maker_fee"`

	Accounts       BankAccounts
	DepositHistory []Transactions
	SendHistory    []Transactions

	client *Client
}

func (a *Account) UnmarshalJSON(b []byte) error {
	type Alias Account
	tmp := &struct {
		Error string `json:"error"`
		*Alias
	}{Alias: (*Alias)(a)}

	if err := json.Unmarshal(b, tmp); err != nil {
		return err
	}
	if tmp.Error != "" {
		return errors.NewEndPointError(tmp.Error)
	}

	return nil
}

func (a *Account) Update() (err error) {

	return
}

func (a Account) SendBTC(address string, amount float64) (t Transaction, err error) {
	url := "https://coincheck.com/api/send_money"
	content := map[string]string{
		"address": address,
		"amount":  strconv.FormatFloat(amount, 'G', -1, 64),
	}

	body, err := a.client.DoRequest("POST", url, content)
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&t)
	return
}
