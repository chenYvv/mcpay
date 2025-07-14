package bsc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// BSC网络配置
const (
	// BSC主网
	//BSC_MAINNET_RPC           = "https://bsc-dataseed1.binance.org/"
	BSC_MAINNET_RPC           = "https://bsc-dataseed.bnbchain.org/"
	BSC_MAINNET_RPC_BACKUP1   = "https://bsc-dataseed.nariox.org/"
	BSC_MAINNET_RPC_BACKUP2   = "https://bsc-dataseed.defibit.io/"
	BSC_MAINNET_RPC_BACKUP3   = "https://bsc-dataseed.ninicoin.io/"
	BSC_MAINNET_RPC_BACKUP4   = "https://bsc.nodereal.io/"
	BSC_MAINNET_RPC_BACKUP5   = "https://bsc-dataseed-public.bnbchain.org/"
	BSC_MAINNET_RPC_BACKUP6   = "https://bnb.rpc.subquery.network/public"
	BSC_MAINNET_CHAIN_ID      = 56
	BSC_MAINNET_USDT_CONTRACT = "0x55d398326f99059ff775485246999027b3197955"

	// BSC测试网 - 多个备用RPC端点
	BSC_TESTNET_RPC           = "https://data-seed-prebsc-1-s1.binance.org:8545/"
	BSC_TESTNET_RPC_BACKUP1   = "https://data-seed-prebsc-2-s1.binance.org:8545/"
	BSC_TESTNET_RPC_BACKUP2   = "https://data-seed-prebsc-1-s2.binance.org:8545/"
	BSC_TESTNET_RPC_BACKUP3   = "https://data-seed-prebsc-2-s2.binance.org:8545/"
	BSC_TESTNET_RPC_BACKUP4   = "https://bsc-testnet.public.blastapi.io/"
	BSC_TESTNET_RPC_BACKUP5   = "https://bsc-testnet.blockpi.network/v1/rpc/public"
	BSC_TESTNET_CHAIN_ID      = 97
	BSC_TESTNET_USDT_CONTRACT = "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd"

	// Gas配置
	DEFAULT_GAS_LIMIT = 21000
	TOKEN_GAS_LIMIT   = 60000
	DEFAULT_GAS_PRICE = 5000000000 // 5 Gwei

	// 代币小数位数配置
	COMMON_DECIMALS = 18
	BNB_DECIMALS    = 18 // BNB使用18位小数
	USDT_DECIMALS   = 18 // BSC上的USDT使用18位小数

	// 显示精度配置
	DISPLAY_PRECISION = 6 // 余额显示时的小数精度
)

// ERC20 ABI - 简化版本，包含必要的方法
const ERC20_ABI = `[
	{
		"constant": true,
		"inputs": [{"name": "_owner", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "balance", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{"name": "_to", "type": "address"},
			{"name": "_value", "type": "uint256"}
		],
		"name": "transfer",
		"outputs": [{"name": "", "type": "bool"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [{"name": "", "type": "uint8"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	}
]`

// BSCClient BSC客户端结构体
type BSCClient struct {
	client  *ethclient.Client
	chainID *big.Int
	rpcURL  string
	isTest  bool
}

// WalletInfo 钱包信息
type WalletInfo struct {
	PrivateKey string `json:"private_key"`
	Address    string `json:"address"`
}

// BalanceInfo 余额信息
type BalanceInfo struct {
	BNBBalance  *big.Int `json:"bnb_balance"`
	USDTBalance *big.Int `json:"usdt_balance"`
	Address     string   `json:"address"`
}

// TransferParams 转账参数
type TransferParams struct {
	PrivateKey string   `json:"private_key"`
	ToAddress  string   `json:"to_address"`
	Amount     float64  `json:"amount"`
	GasPrice   *big.Int `json:"gas_price,omitempty"`
	GasLimit   uint64   `json:"gas_limit,omitempty"`
}

// 全局BSC客户端实例
var (
	Client *BSCClient
	once   sync.Once
)

// InitBSCClient 初始化BSC客户端（系统启动时调用）
func InitBSCClient(isTest bool) error {
	var initErr error
	once.Do(func() {
		client, err := newBSCClient(isTest)
		if err != nil {
			initErr = err
			return
		}
		Client = client
	})
	return initErr
}

// GetClient 获取全局BSC客户端实例
func GetClient() *BSCClient {
	return Client
}

// newBSCClient 创建新的BSC客户端
func newBSCClient(isTest bool) (*BSCClient, error) {
	var rpcEndpoints []string
	var expectedChainID int64

	if isTest {
		rpcEndpoints = []string{
			BSC_TESTNET_RPC_BACKUP4, // BlastAPI
			BSC_TESTNET_RPC_BACKUP5, // BlockPI
			BSC_TESTNET_RPC_BACKUP1,
			BSC_TESTNET_RPC_BACKUP2,
			BSC_TESTNET_RPC_BACKUP3,
			BSC_TESTNET_RPC,
		}
		expectedChainID = BSC_TESTNET_CHAIN_ID
	} else {
		rpcEndpoints = []string{
			BSC_MAINNET_RPC,
			BSC_MAINNET_RPC_BACKUP1,
			BSC_MAINNET_RPC_BACKUP2,
			BSC_MAINNET_RPC_BACKUP3,
			BSC_MAINNET_RPC_BACKUP4,
			BSC_MAINNET_RPC_BACKUP5,
			BSC_MAINNET_RPC_BACKUP6,
		}
		expectedChainID = BSC_MAINNET_CHAIN_ID
	}

	for _, rpcURL := range rpcEndpoints {
		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			continue
		}

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		chainID, err := client.ChainID(ctx)
		cancel()

		if err != nil {
			client.Close()
			continue
		}

		// 验证链ID
		if chainID.Int64() != expectedChainID {
			client.Close()
			continue
		}

		return &BSCClient{
			client:  client,
			chainID: chainID,
			rpcURL:  rpcURL,
			isTest:  isTest,
		}, nil
	}

	return nil, fmt.Errorf("failed to connect to any BSC RPC endpoint")
}

// ResetBSCClient 重置BSC客户端（用于重新连接）
func ResetBSCClient(isTest bool) error {
	if Client != nil {
		Client.Close()
	}

	return InitBSCClient(isTest)
}

// Close 关闭客户端连接
func (c *BSCClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// GetRPCURL 获取当前使用的RPC URL
func (c *BSCClient) GetRPCURL() string {
	return c.rpcURL
}

// GetChainID 获取链ID
func (c *BSCClient) GetChainID() *big.Int {
	return c.chainID
}

// IsTestnet 检查是否是测试网
func (c *BSCClient) IsTestnet() bool {
	return c.isTest
}

// GetUSDTContract 获取USDT合约地址
func (c *BSCClient) GetUSDTContract() string {
	if c.isTest {
		return BSC_TESTNET_USDT_CONTRACT
	} else {
		return BSC_MAINNET_USDT_CONTRACT
	}
}

// ========== 1. 创建BSC钱包 ==========

// CreateWallet 创建新的BSC钱包
func CreateWallet() (*WalletInfo, error) {
	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// 转换为十六进制字符串
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// 获取地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &WalletInfo{
		PrivateKey: privateKeyHex,
		Address:    address,
	}, nil
}

// GetAddressFromPrivateKey 从私钥获取地址
func GetAddressFromPrivateKey(privateKeyHex string) (string, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

// ========== 2. 查询余额 ==========

// GetBNBBalance 获取BNB余额
func (c *BSCClient) GetBNBBalance(address string) (*big.Int, error) {
	if !IsValidAddress(address) {
		return nil, fmt.Errorf("invalid address: %s", address)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	balance, err := c.client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get BNB balance: %v", err)
	}

	return balance, nil
}

// GetUSDTBalance 获取USDT余额
func (c *BSCClient) GetUSDTBalance(address string) (*big.Int, error) {
	return c.GetTokenBalance(address, c.GetUSDTContract())
}

// GetTokenBalance 获取指定代币余额
func (c *BSCClient) GetTokenBalance(address, contractAddress string) (*big.Int, error) {
	if !IsValidAddress(address) || !IsValidAddress(contractAddress) {
		return nil, fmt.Errorf("invalid address")
	}

	// 解析ABI
	parsedABI, err := abi.JSON(strings.NewReader(ERC20_ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// 创建合约实例
	contract := bind.NewBoundContract(common.HexToAddress(contractAddress), parsedABI, c.client, c.client, c.client)

	// 调用balanceOf方法
	var result []interface{}
	err = contract.Call(&bind.CallOpts{}, &result, "balanceOf", common.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("failed to call balanceOf: %v", err)
	}

	if len(result) == 0 {
		return big.NewInt(0), nil
	}

	balance, ok := result[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected balance type")
	}

	return balance, nil
}

// GetBalance 获取钱包完整余额信息
func (c *BSCClient) GetBalance(address string) (*BalanceInfo, error) {
	bnbBalance, err := c.GetBNBBalance(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get BNB balance: %v", err)
	}

	usdtBalance, err := c.GetUSDTBalance(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get USDT balance: %v", err)
	}

	return &BalanceInfo{
		BNBBalance:  bnbBalance,
		USDTBalance: usdtBalance,
		Address:     address,
	}, nil
}

// ========== 3. 转移BNB ==========

// TransferBNB 转账BNB
func (c *BSCClient) TransferBNB(params *TransferParams) (string, error) {
	return c.transferNative(params)
}

// transferNative 转账原生代币（BNB）
func (c *BSCClient) transferNative(params *TransferParams) (string, error) {
	// 验证参数
	if err := c.validateTransferParams(params); err != nil {
		return "", err
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(params.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// 获取发送者地址
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 获取nonce
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 设置gas价格
	gasPrice := params.GasPrice
	if gasPrice == nil {
		gasPrice = big.NewInt(DEFAULT_GAS_PRICE)
	}

	// 设置gas限制
	gasLimit := params.GasLimit
	if gasLimit == 0 {
		gasLimit = DEFAULT_GAS_LIMIT
	}

	// 创建交易
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(params.ToAddress),
		NumToWei(params.Amount),
		gasLimit,
		gasPrice,
		nil,
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// ========== 4. 转移USDT ==========

// TransferUSDT 转账USDT
func (c *BSCClient) TransferUSDT(params *TransferParams) (string, error) {
	return c.TransferToken(params, c.GetUSDTContract())
}

// TransferToken 转账指定代币
func (c *BSCClient) TransferToken(params *TransferParams, contractAddress string) (string, error) {
	// 验证参数
	if err := c.validateTransferParams(params); err != nil {
		return "", err
	}

	if !IsValidAddress(contractAddress) {
		return "", fmt.Errorf("invalid contract address: %s", contractAddress)
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(params.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// 获取发送者地址
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 解析ABI
	parsedABI, err := abi.JSON(strings.NewReader(ERC20_ABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %v", err)
	}

	// 创建合约实例
	contract := bind.NewBoundContract(common.HexToAddress(contractAddress), parsedABI, c.client, c.client, c.client)

	// 准备交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, c.chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %v", err)
	}

	// 设置gas价格和限制
	auth.GasPrice = params.GasPrice
	if auth.GasPrice == nil {
		auth.GasPrice = big.NewInt(DEFAULT_GAS_PRICE)
	}

	auth.GasLimit = params.GasLimit
	if auth.GasLimit == 0 {
		auth.GasLimit = TOKEN_GAS_LIMIT
	}

	// 获取nonce
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// 调用transfer方法
	tx, err := contract.Transact(auth, "transfer", common.HexToAddress(params.ToAddress), params.Amount)
	if err != nil {
		return "", fmt.Errorf("failed to send token transfer: %v", err)
	}

	return tx.Hash().Hex(), nil
}

// ========== 5. 辅助工具函数 ==========

// IsValidAddress 验证地址是否有效
func IsValidAddress(address string) bool {
	// BSC地址必须以0x开头且长度为42
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	if len(address) != 42 {
		return false
	}
	return common.IsHexAddress(address)
}

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

// ...existing code...

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

// validateTransferParams 验证转账参数
func (c *BSCClient) validateTransferParams(params *TransferParams) error {
	amountWei := NumToWei(params.Amount)

	if params == nil {
		return fmt.Errorf("transfer params cannot be nil")
	}

	if params.PrivateKey == "" {
		return fmt.Errorf("private key cannot be empty")
	}

	if !IsValidAddress(params.ToAddress) {
		return fmt.Errorf("invalid to address: %s", params.ToAddress)
	}

	if amountWei == nil || amountWei.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	return nil
}

// GetNetworkInfo 获取网络信息
func (c *BSCClient) GetNetworkInfo() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取最新区块号
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %v", err)
	}

	// 获取gas价格
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	info := map[string]interface{}{
		"chainID":     c.chainID.String(),
		"rpcURL":      c.rpcURL,
		"blockNumber": blockNumber,
		"gasPrice":    gasPrice.String(),
	}

	return info, nil
}

func (c *BSCClient) BlockNumber() int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 获取最新区块号
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0
	}
	return int(blockNumber)
}

// WeiStrToNum 将最小单位转换为金额
func WeiStrToNum(str string) float64 {
	wei := new(big.Int)
	wei.SetString(str, 10)
	// 转换为float64
	floatResult, _ := WeiToNumWithDecimals(wei, COMMON_DECIMALS).Float64()
	return floatResult
}
