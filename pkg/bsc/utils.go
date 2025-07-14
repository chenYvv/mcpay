package bsc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BscScanBEP20Transfer
// Get a list of 'BEP-20 Token Transfer Events' by Address
type BscScanBEP20Transfer struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`
}

// 获取充值情况
// 文档 https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-a-list-of-erc20-token-transfer-events-by-address
func GetUSDTTransactionsByEtherscan(apiKey string, address string, startBlock int64) ([]BscScanBEP20Transfer, error) {

	chainId := GetClient().GetChainID().Int64()
	contractAddress := GetClient().GetUSDTContract()

	// API URL
	url := fmt.Sprintf(
		"https://api.etherscan.io/v2/api?chainid=%d&module=account&action=tokentx&contractaddress=%s&address=%s&page=%d&offset=%d&startblock=%d&endblock=999999999&sort=desc&apikey=%s",
		chainId,
		contractAddress,
		address,
		1,
		50,
		startBlock,
		apiKey,
	)
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

	type Response struct {
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Result  []BscScanBEP20Transfer `json:"result"`
	}

	var response Response

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return response.Result, nil
}

// 获取充值情况
// 文档
func GetBlockByEtherscan(apiKey string, timestamp int64) (string, error) {

	chainId := GetClient().GetChainID().Int64()

	// API URL
	url := fmt.Sprintf(
		"https://api.etherscan.io/v2/api?chainid=%d&module=block&timestamp=%d&action=getblocknobytime&apikey=%s",
		chainId,
		timestamp,
		apiKey,
	)
	fmt.Println(url)
	// 发起 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	var response Response

	if err = json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return response.Result, nil
}
