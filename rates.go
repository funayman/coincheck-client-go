package coincheck

import (
	"encoding/json"
	"strconv"
)

type Rate struct {
	client *Client
}

type CalculatedRate struct {
	Rate   float64
	Amount string `json:"amount"`
	Price  string `json:"price"`
}

func (cr *CalculatedRate) UnmarshalJSON(b []byte) error {
	type Alias CalculatedRate
	tmp := &struct {
		Rate string `json:"rate"`
		*Alias
	}{Alias: (*Alias)(cr)}

	err := json.Unmarshal(b, tmp)
	if err != nil {
		return nil
	}

	rate, err := strconv.ParseFloat(tmp.Rate, 64)
	if err != nil {
		return err
	}

	cr.Rate = rate

	return nil
}

//Coin returns the price of the coin selected in the currency selected
//BTC_JPY would return the rate of BitCoin in Japanese Yen
//ETH_BTC would return the rate of Ethereum in BitCoins
//More information can be found at https://coincheck.com/documents/exchange/api#buy-rate
func (r Rate) Coin(coin string) (float64, error) {
	url := "/rate" + coin

	body, err := r.client.DoRequest("GET", url, nil)
	if err != nil {
		return 0.0, err
	}

	tmp := &struct {
		Rate string `json:"rate"`
	}{}
	json.NewDecoder(body).Decode(tmp)

	rate, err := strconv.ParseFloat(tmp.Rate, 64)
	if err != nil {
		return 0.0, err
	}

	return rate, nil
}

//CalculateByAmount calculates the current rate given the amount
func (r *Rate) CalculateByAmount(orderType, pair string, amount float64) (cr CalculatedRate, err error) {
	return r.calculateRate(map[string]string{
		"order_type": orderType,
		"pair":       pair,
		"amount":     strconv.FormatFloat(amount, 'G', -1, 64),
	})
}

//CalculateByPrice calculates the current rate given the price
func (r *Rate) CalculateByPrice(orderType, pair string, price float64) (cr CalculatedRate, err error) {
	return r.calculateRate(map[string]string{
		"order_type": orderType,
		"pair":       pair,
		"price":      strconv.FormatFloat(price, 'G', -1, 64),
	})
}

func (r Rate) calculateRate(data map[string]string) (cr CalculatedRate, err error) {
	url := "/exchange/orders/rate"

	body, err := r.client.DoRequest("GET", url, data)
	if err != nil {
		return
	}

	if err = json.NewDecoder(body).Decode(&cr); err != nil {
		return
	}

	return

}
