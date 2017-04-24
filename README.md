# CoinCheck Client
A simple [CoinCheck](https://www.coincheck.com) client written in Go.
Lacking many features and I will implement more as I need them
But right now, it works

```go
var client = coincheck.New("APIKEY", "APISECRET")

//Get the current rate & current balance at CoinCheck
ccRate, _:= client.Rate(coincheck.BtcJpy)
ccBalance, _ := client.AccountBalance()

//Purchase Bitcoin
ccAmount := ccBalance.JPY / ccRate
order, err := ccClient.NewOrder(ccRate, ccAmount, coincheck.Buy, coincheck.BtcJpy)
if err != nil {
	log.Fatal(err)
}

//Current Orders
orders, _ := ccClient.UnsettledOrders()
for _, order := range orders {
  client.CancleOrder(order.ID)
}

//Send Bitcoin
ccBalance, _ = client.AccountBalance()
transaction, _:= ccClient.SendBTC("1NVESChAyPMaM4jfJSwSnDRQwBEJmD77Mv", ccBalance.BTC - coincheck.SendingFee)
```
