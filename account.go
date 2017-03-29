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

type Balance struct {
	JPY          string `json:"jpy"`
	BTC          string `json:"btc"`
	JPYReserved  string `json:"jpy_reserved"`
	BTCReserved  string `json:"btc_reserved"`
	JPYLendInUse string `json:"jpy_lend_in_use"`
	BTCLendInUse string `json:"btc_lend_in_use"`
	JPYLent      string `json:"jpy_lent"`
	BTCLent      string `json:"btc_lent"`
	JPYDebt      string `json:"jpy_debt"`
	BTCDebt      string `json:"btc_debt"`
}

type Account struct {
	ID              int64  `json:"id"`
	Email           string `json:"email"`
	IdentityStatus  string `json:"identity_status"`
	Address         string `json:"bitcoin_address"`
	LendingLeverate string `json:"lending_leverage"`
	TakerFee        string `json:"taker_fee"`
	MakerFee        string `json:"maker_fee"`

	Accounts       BankAccounts
	DepositHistory Transactions
	SentHistory    Transactions
	Balance        Balance

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
	if err = a.UpdateAccountInfo(); err != nil {
		return err
	}
	if err = a.UpdateBalance(); err != nil {
		return err
	}
	if err = a.UpdateBankInfo(); err != nil {
		return err
	}
	/* Currently getting invalid auth errors
	if err = a.UpdateHistory(); err != nil {
		return err
	}
	*/
	return
}

func (a *Account) UpdateAccountInfo() (err error) {
	url := "https://coincheck.com/api/accounts"
	body, err := a.client.DoRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return json.NewDecoder(body).Decode(a)
}

func (a *Account) UpdateBalance() (err error) {
	url := "https://coincheck.com/api/accounts/balance"
	body, err := a.client.DoRequest("GET", url, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&a.Balance)
	return
}

func (a *Account) UpdateBankInfo() (err error) {
	url := "https://coincheck.com/api/bank_accounts"
	body, err := a.client.DoRequest("GET", url, nil)
	if err != nil {
		return
	}

	tmp := &struct {
		Data BankAccounts `json:"data"`
	}{}
	err = json.NewDecoder(body).Decode(tmp)
	if err != nil {
		return
	}

	a.Accounts = tmp.Data
	return
}

/***************
 **  HISTORY  **
 ***************/

//UpdateHistory updates both Sent and Deposit History
func (a *Account) UpdateHistory() (err error) {
	if err = a.updateDepositHistory(); err != nil {
		return err
	}
	if err = a.updateSentHistory(); err != nil {
		return err
	}
	return
}

func (a *Account) updateSentHistory() (err error) {
	url := "http://www.coincheck.com/api/send_money"
	body, err := a.client.DoRequest("GET", url, map[string]string{"currency": "btc"})
	if err != nil {
		return
	}

	tmp := &struct {
		Data Transactions `json:"sends"`
	}{}

	if err = json.NewDecoder(body).Decode(tmp); err != nil {
		return
	}

	a.SentHistory = tmp.Data

	return
}

func (a *Account) updateDepositHistory() (err error) {
	url := "http://www.coincheck.com/api/deposit_money"
	body, err := a.client.DoRequest("GET", url, map[string]string{"currency": "BTC"})
	if err != nil {
		return
	}

	tmp := &struct {
		Data Transactions `json:"deposits"`
	}{}
	if err = json.NewDecoder(body).Decode(tmp); err != nil {
		return
	}
	a.DepositHistory = tmp.Data

	return
}

//SendBTC sends bitcoin to the specified address
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
