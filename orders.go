package coincheck

import (
	"encoding/json"
	"strconv"
	"time"
)

//OrderType used to limit the input values for NewOrder
type OrderType string

const (
	Buy          OrderType = "buy"
	Sell         OrderType = "sell"
	MarketBuy    OrderType = "market_buy"
	MarketSell   OrderType = "market_sell"
	LeverageBuy  OrderType = "leverage_buy"
	LeverageSell OrderType = "leverage_sell"
	CloseLong    OrderType = "close_long"
	CloseShort   OrderType = "close_short"
)

//Order Response from NewOrder endpoint on the coincheck api
type Order struct {
	ID           int64 `json:"id"`
	Rate         float64
	Amount       float64
	Type         OrderType `json:"order_type"`
	StopLossRate string    `json:"stop_loss_rate"`
	Pair         string    `json:"pair"`
	CreatedAt    time.Time

	Error string `json:"error"`
	Raw   []byte
}

func (o *Order) UnmarshalJSON(b []byte) (err error) {
	type alias Order
	tmp := &struct {
		Rate   string `json:"rate"`
		Amount string `json:"amount"`
		*alias
	}{alias: (*alias)(o)}

	err = json.Unmarshal(b, tmp)
	if err != nil {
		return
	}

	oRate, err := strconv.ParseFloat(tmp.Rate, 64)
	if err != nil {
		return
	}

	oAmount, err := strconv.ParseFloat(tmp.Amount, 64)
	if err != nil {
		return
	}

	o.Rate = oRate
	o.Amount = oAmount

	return
}

//UnsettledOrder similar to Order but with a few extra params
type UnsettledOrder struct {
	ID   int64     `json:"id"`
	Type OrderType `json:"order_type"`
	Rate string    `json:"rate"`
	//Rate                   float64   `json:"rate"`
	Pair string `json:"pair"`
	//PendingAmount          float64   `json:"pending_amount"`
	PendingAmount          string `json:"pending_amount"`
	PendingMarketBuyAmount string `json:"pending_market_buy_amount"`
	CreatedAt              time.Time
	StopLossRate           string `json:"stop_loss_rate"`

	Error string `json:"error"`
	Raw   []byte
}

//OrderTransaction Type of Transaction returned from OrdeTransactions API endpoint
type OrderTransaction struct {
	ID        int64 `json:"id"`
	OrderID   int64 `json:"order_id"`
	CreatedAt time.Time
	Funds     struct {
		Btc string `json:"btc"`
		Jpy string `json:"jpn"`
	}
	Pair        string `json:"pair"`
	Rate        string `json:"rate"`
	FeeCurrency string `json:"fee_currency"`
	Fee         string `json:"fee"`
	Liquidity   string `json:"liquidity"`
	Side        string `json:"side"`
}

//NewOrder Publish new order to exchange.
//For example if you'd like buy 10 BTC as 30000 JPY/BTC, you need to specify following parameters.
//rate: 30000, amount: 10, order_type: "buy", pair: "btc_jpy"
//In case of buying order, if selling order which is lower than specify price is already exist, we settle in low price order.
//In case of to be short selling order, remain as unsettled order. It is also true for selling order.
func (client *Client) NewOrder(rate, amount float64, orderType OrderType, pair string) (order Order, err error) {
	endpoint := "/exchange/orders"

	cRate := strconv.FormatFloat(rate, 'f', -1, 64)
	if err != nil {
		return
	}

	cAmount := strconv.FormatFloat(amount, 'f', -1, 64)
	if err != nil {
		return
	}

	content := map[string]string{
		"pair":       pair,
		"order_type": string(orderType),
	}

	switch orderType {
	case Buy, Sell:
		content["rate"] = cRate
		content["amount"] = cAmount
	case MarketBuy:
		content["market_buy"] = cRate
	case MarketSell:
		content["amount"] = cAmount
	case LeverageBuy, LeverageSell:
		content["amount"] = cAmount
		if rate != -1.0 {
			content["rate"] = cRate
		}
	}

	body, err := client.DoRequest("POST", endpoint, content)
	if err != nil {
		return
	}

	err = json.NewDecoder(body).Decode(&order)

	return
}

//UnsettledOrders You can get a unsettled order list
func (client *Client) UnsettledOrders() (orders []UnsettledOrder, err error) {
	endpoint := "/exchange/orders/opens"

	body, err := client.DoRequest("GET", endpoint, nil)
	if err != nil {
		return
	}

	tmp := &struct {
		Data []UnsettledOrder `json:"orders"`
	}{}
	err = json.NewDecoder(body).Decode(tmp)

	orders = tmp.Data

	return
}

//CancelOrder Canceles a NewOrder or, you can cancel by specifying an UnsettledOrder list's ID.
func (client *Client) CancelOrder(id string) (err error) {
	endpoint := "/exchange/orders/" + id
	body, err := client.DoRequest("DELETE", endpoint, nil)
	if err != nil {
		return
	}

	tmp := &struct {
		ID    int64  `json:"id"`
		Error string `json:"error"`
	}{}

	err = json.NewDecoder(body).Decode(tmp)

	return
}

//OrderTransactions Display your transaction history
func (client Client) OrderTransactions() (t []OrderTransaction, err error) {
	endpoint := "/exchange/orders/transactions"
	body, err := client.DoRequest("GET", endpoint, nil)
	if err != nil {
		return
	}

	tmp := &struct {
		Data []OrderTransaction `json:"transactions"`
	}{}
	err = json.NewDecoder(body).Decode(tmp)

	t = tmp.Data

	return
}
