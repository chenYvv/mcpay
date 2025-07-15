package chain

import (
	"encoding/json"
	"net/http"
)

type TetherPrice struct {
	USD float64 `json:"usd"`
	INR float64 `json:"inr"`
}

// 获取USD和法币的汇率
func GetUSDTPrice() (TetherPrice, error) {
	var data TetherPrice
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=tether&vs_currencies=usd,inr")
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}
