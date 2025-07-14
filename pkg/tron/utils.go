package tron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/shopspring/decimal"
)

// 通用版本 - 支持任意小数位数
func WeiToNumWithDecimals(wei *big.Int, decimals int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}

	// 使用 decimal 库处理，提高精度
	d := decimal.NewFromBigInt(wei, 0)

	// 计算除数 10^decimals
	divisor := decimal.New(1, int32(decimals))

	// 执行除法运算
	result := d.Div(divisor)

	// 转换回 *big.Float
	// 使用字符串作为中介，保持最高精度
	resultStr := result.String()
	resultBigFloat := new(big.Float)
	resultBigFloat.SetPrec(256) // 设置高精度
	resultBigFloat.SetString(resultStr)

	return resultBigFloat
}

func NumToWeiWithDecimals(num float64, decimals int) *big.Int {
	if num < 0 {
		return big.NewInt(0)
	}

	// 使用 decimal 库处理精度问题
	d := decimal.NewFromFloat(num)

	// 计算10^decimals
	multiplier := decimal.New(1, int32(decimals)) // 等于 10^decimals

	// 相乘
	result := d.Mul(multiplier)

	// 转换为 big.Int
	return result.BigInt()
}

// NumToWei 将金额转换为最小单位
func NumToWei(num float64) *big.Int {
	return NumToWeiWithDecimals(num, COMMON_DECIMALS)
}

// WeiToNum 将最小单位转换为金额
func WeiToNum(wei *big.Int) float64 {
	// 转换为float64
	floatResult, _ := WeiToNumWithDecimals(wei, COMMON_DECIMALS).Float64()
	return floatResult
}

// WeiStrToNum 将最小单位转换为金额
func WeiStrToNum(str string) float64 {
	wei := new(big.Int)
	wei.SetString(str, 10)
	// 转换为float64
	floatResult, _ := WeiToNumWithDecimals(wei, COMMON_DECIMALS).Float64()
	return floatResult
}

type Transaction struct {
	TransactionId string `json:"transaction_id"`
	TokenInfo     struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		Name     string `json:"name"`
	} `json:"token_info"`
	BlockTimestamp int64  `json:"block_timestamp"`
	From           string `json:"from"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Value          string `json:"value"`
}

// 获取充值情况
// 文档 https://developers.tron.network/reference/get-trc20-transaction-info-by-account-address
func GetUSDTTransactions(address string, startTimestamp int64) ([]Transaction, error) {

	TEST_API := "https://nile.trongrid.io"
	MAIN_API := "https://api.trongrid.io"

	// 根据 GetClient().IsTest() 返回判断使用哪个 TEST_API 或者 MAIN_API
	apiUrl := MAIN_API
	contract_address := GetClient().config.GetUSDTContract()
	if GetClient().IsTest() {
		apiUrl = TEST_API
	}

	// API URL
	url := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20?contract_address=%s&only_confirmed=true&only_to=true&min_timestamp=%d&limit=100", apiUrl, address, contract_address, startTimestamp)
	fmt.Println(url)
	// 发起 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应
	var response struct {
		Data []Transaction `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return response.Data, nil
}
