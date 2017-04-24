package coincheck

import (
	"encoding/json"
	"strconv"
)

type CalculatedRate struct {
	Rate   float64
	Amount string `json:"amount"`
	Price  string `json:"price"`
	Raw    []byte
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
	cr.Raw = b

	return nil
}

//CalculateByAmount calculates the current rate given the amount
func (client *Client) CalculateByAmount(orderType, pair string, amount float64) (cr CalculatedRate, err error) {
	return client.calculateRate(map[string]string{
		"order_type": orderType,
		"pair":       pair,
		"amount":     strconv.FormatFloat(amount, 'G', -1, 64),
	})
}

//CalculateByPrice calculates the current rate given the price
func (client *Client) CalculateByPrice(orderType, pair string, price float64) (cr CalculatedRate, err error) {
	return client.calculateRate(map[string]string{
		"order_type": orderType,
		"pair":       pair,
		"price":      strconv.FormatFloat(price, 'G', -1, 64),
	})
}

func (client *Client) calculateRate(data map[string]string) (cr CalculatedRate, err error) {
	url := "/exchange/orders/rate"

	body, err := client.DoRequest("GET", url, data)
	if err != nil {
		return
	}

	if err = json.NewDecoder(body).Decode(&cr); err != nil {
		return
	}

	return

}
