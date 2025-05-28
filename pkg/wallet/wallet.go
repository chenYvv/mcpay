package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"image/color"
	"log"
	"math/big"
	"mcpay/pkg/helpers"
	"mcpay/pkg/logger"
	"regexp"
	"strings"
)

// 币安主网 USDT 合约地址
//const BSC_USDT_ContractAddress = "0x55d398326f99059ff775485246999027b3197955"
//const BSC_USDT_ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// 币安测试 USDT 合约地址
//const BSC_USDT_ContractAddress = "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd"

// const BSC_USDT_ContractAddress = "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd" // 测试
const BSC_USDT_ABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"_decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// const ethNodeURL = "http://127.0.0.1:7545"
// const ethNodeURL = "https://bsc-dataseed1.binance.org" // BSC主网
// const ethNodeURL = "https://data-seed-prebsc-1-s1.binance.org:8545" // BSC测试网
// const ethNodeURL = "https://data-seed-prebsc-2-s1.binance.org:8545" // BSC测试网
// const ethNodeURL = "https://data-seed-prebsc-1-s2.binance.org:8545" // BSC测试网
// const ethNodeURL = "https://data-seed-prebsc-2-s2.binance.org:8545" // BSC测试网

type WalletClient struct {
	client *ethclient.Client
}

var WalletClientIns *WalletClient

func GetWalletClient(url string) *WalletClient {
	if WalletClientIns != nil {
		return WalletClientIns
	}
	WalletClientIns = new(WalletClient)
	client, _ := ConnectRPC(url)
	WalletClientIns.client = client
	return WalletClientIns
}

// 创建钱包地址
func CreateWallet() (privateKey, address string, err error) {
	privateKeyEcdsa, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	// 私钥
	privateKey = hex.EncodeToString(privateKeyEcdsa.D.Bytes())
	address = GetAddressByPrivateKey(privateKey)
	return
}

// GetAddressByPrivateKey 通过私钥获取钱包地址
func GetAddressByPrivateKey(privateKey string) string {
	//privateKey = privateKey[2:]
	privateKeyCDSA, err := crypto.HexToECDSA(privateKey)

	if err != nil {
		return ""
	}

	//privateKeyBytes := crypto.FromECDSA(privateKeyCDSA)
	//fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xa0a40acd844e21fa97d74dcea66b61bae5acdbb46f2848e117611480580fbe81

	publicKey := privateKeyCDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	//fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
	// c6ffe41c2261e9fad4b1357ffd5e95dd03707aa982c18c994e5734a37868bd5f094f15d2aadac035cdc0c268f269d258916bdf60179e31f1bb490bfaf2f653b5

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}

// Binance Smart Chain 公共测试链的网络参数如下：
//
// 网络名称：Binance Smart Chain Testnet
// RPC 地址：https://data-seed-prebsc-1-s1.binance.org:8545
// ChainID：97
// 浏览器：https://testnet.bscscan.com
// 您可以使用上述网络参数在测试环境中与 Binance Smart Chain 进行交互。请注意，这是一个测试网络，不涉及真实的资金交易。

func ConnectRPC(url string) (*ethclient.Client, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	conn := ethclient.NewClient(rpcClient)
	return conn, nil
}

// 获取主币剩余额度
func (m *WalletClient) GetBalance(address string) (*big.Int, error) {
	balanceAt, err := m.client.BalanceAt(context.TODO(), common.HexToAddress(address), nil)
	if err != nil {
		return big.NewInt(0), err
	}
	//res1 := Wei2ETH(balanceAt)
	//fmt.Println(res1)
	fmt.Println(balanceAt)
	//fmt.Println(balanceAt.Float64())
	//fmt.Println(float64(balanceAt.Int64()))
	//balanceV := float64(balanceAt.Int64()) * math.Pow(10, -18)
	//fmt.Println(balanceV)
	//res := balanceAt.Div(balanceAt, big.NewInt(1e18))
	//fmt.Println(res)
	return balanceAt, nil
}

// 转移基础币
func (m *WalletClient) Transfer(privateKey string, toAddress string, amount *big.Int) (err error, hash string) {
	fromAddress := GetAddressByPrivateKey(privateKey)
	nonce, err := m.client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		return
	}

	//value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000) // in units
	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}

	var data []byte
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, gasLimit, gasPrice, data)

	//chainID, err := client.NetworkID(context.Background())
	chainID, err := m.client.ChainID(context.Background())
	if err != nil {
		return
	}

	privateKeyEcdsa, _ := crypto.HexToECDSA(privateKey)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyEcdsa)
	if err != nil {
		return
	}

	err = m.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return
	}

	fmt.Printf("tx sent: %s", hash)
	hash = signedTx.Hash().Hex()
	return
}

// https://api.bscscan.com/api?module=account&action=tokenbalance&contractaddress=0x55d398326f99059ff775485246999027b3197955&address=0xf38a409da2823585b66d1164d2520b78d4fab049&tag=latest&apikey=5BZ4BHK758HAFGVRRVDPSEC3H5CZUZ647R

func (m *WalletClient) Get_BSC_USDT_Balance(address string, contractAddress string) (*big.Int, error) {
	// 解析钱包地址
	walletAddress := common.HexToAddress(address)

	balanceBNB, err := m.client.BalanceAt(context.Background(), walletAddress, nil)
	if err != nil {
		return nil, err
	}
	//fmt.Println("bnb:", balanceBNB)

	// 解析 USDT 代币合约地址
	usdtContractAddress := common.HexToAddress(contractAddress)

	// 解析 USDT 代币 ABI
	//const usdtABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"}]`

	contractABI, err := abi.JSON(strings.NewReader(BSC_USDT_ABI))
	if err != nil {
		return nil, err
	}

	// 创建一个新的查询消息，调用 balanceOf 方法
	data, err := contractABI.Pack("balanceOf", walletAddress)
	if err != nil {
		return nil, err
	}

	// 发送调用交易到以太坊网络
	callMsg := ethereum.CallMsg{To: &usdtContractAddress, Data: data}
	result, err := m.client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	// 解码余额
	var balance [32]byte
	copy(balance[:], result)
	balanceInt := new(big.Int).SetBytes(balance[:])

	// 创建一个新的查询消息，调用 decimals 方法
	data, err = contractABI.Pack("decimals")
	if err != nil {
		return nil, err
	}

	// 发送调用交易到以太坊网络
	callMsg = ethereum.CallMsg{To: &usdtContractAddress, Data: data}
	result, err = m.client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	// 解码小数位数
	var decimals [32]byte
	copy(decimals[:], result)
	decimalsInt := new(big.Int).SetBytes(decimals[:])

	// 将余额从大整数转换为可读的格式
	usdtBalanceFloat := new(big.Float).SetInt(balanceInt)
	decimalsExp := new(big.Int).Exp(big.NewInt(10), decimalsInt, nil)
	usdtBalanceFloat.Quo(usdtBalanceFloat, new(big.Float).SetInt(decimalsExp))
	usdtBalanceString := usdtBalanceFloat.Text('f', 2)

	//fmt.Printf("Wallet : %v\n", address)
	//fmt.Printf("Wallet USDT Balance: %v WEI\n", balanceInt)
	//fmt.Printf("Wallet USDT Balance: %s USDT\n", usdtBalanceString)

	logger.Info(
		"钱包余额", zap.String("address", address),
		zap.String("BNB", balanceBNB.String()),
		zap.String("USDT WEI", balanceInt.String()),
		zap.String("USDT", usdtBalanceString),
	)

	return balanceInt, nil
}

// 交易
// const ethNodeURL = "https://bsc-dataseed1.binance.org"
// BSC USDT 交易
// 如果你不指定合约地址，交易将被认为是发送以太币（ETH）而不是任何代币。因为以太坊和BSC网络上的每个地址都可以同时持有以太币和各种代币，因此如果你想发送以太币以外的代币（如USDT），你需要指定相应代币的合约地址。
// 换句话说，如果你想发送USDT或其他代币，你必须将代币合约地址作为交易的To字段，并在交易的Data字段中传递代币转账函数的数据编码，就像上面的示例代码一样。这是因为代币的转账不是通过简单的交易进行的，而是通过调用代币合约中的特定函数来完成的。
func (m *WalletClient) BSC_USDT_Transfer(privateKeyStr string, toAddress string, coinNum *big.Int, tokenAddress string) (hex string, err error) {
	// Load your private key
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		logger.Info("转移代币：", zap.String("错误信息1", err.Error()))
		return
	}

	// Create an authenticated binding instance of the contract
	auth := bind.NewKeyedTransactor(privateKey)

	// Set gas price and gas limit
	auth.GasLimit = uint64(21000) // Replace with appropriate gas limit

	// 获取gas价格
	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Info("转移代币：", zap.String("错误信息2", err.Error()))
		return
	}
	auth.GasPrice = gasPrice

	// Initialize the contract address
	contractAddress := common.HexToAddress(tokenAddress)

	// Create a new contract instance
	token, err := NewUSDT(contractAddress, m.client)

	if err != nil {
		logger.Info("转移代币：", zap.String("错误信息3", err.Error()))
		return
	}

	// Specify the recipient address and the amount to transfer
	recipient := common.HexToAddress(toAddress)

	// Call the transfer function to send USDT tokens
	tx, err := token.Transfer(auth, recipient, coinNum)
	if err != nil {
		logger.Info("转移代币：", zap.String("错误信息4", err.Error()))
		//log.Fatalf("Failed to send USDT tokens: %v", err)
		return
	}

	//fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())
	hex = tx.Hash().Hex()
	logger.Info("转移代币：", zap.String("hash", hex))
	return
}

// GetSuggestGasPrice 获取
func (m *WalletClient) GetSuggestGasPrice() (*big.Int, error) {
	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return gasPrice, err
}

func Wei2ETH(wei *big.Int) float64 {
	//float, _ := wei.Float64()
	//a := decimal.NewFromFloat(float)
	//b := decimal.NewFromFloat(1000000000080000000)
	//res, _ := a.Div(b).Round(20).Float64()
	//balance := new(big.Float)
	//balance.SetString(wei.String())
	//ethValue := new(big.Float).Quo(balance, big.NewFloat(math.Pow10(18)))
	//return ethValue
	// 将Wei转换为USDT
	usdt := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt64(int64(1e18)))
	usdtFloat, _ := usdt.Float64()
	return usdtFloat
}

// EthToWei eth单位安全转wei
// https://stackoverrun.com/cn/q/13021596
func Eth2Wei(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	// Set precision if required.
	// bigval.SetPrec(64)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func (m *WalletClient) GetGasCost() *big.Int {
	gasPrice, _ := m.GetSuggestGasPrice()
	fmt.Println(gasPrice)
	gasLimit := big.NewInt(2100000) // BSC的标准gas限制
	// 计算 gas 费用
	gasCost := new(big.Int).Mul(gasPrice, gasLimit)
	return gasCost
}

// GetRemain 还剩下多少
//func (m *WalletClient) GetRemain(address string) {
//	weiNum, _ := m.GetBalance(address)
//	gasPrice, _ := m.GetSuggestGasPrice()
//	fmt.Println(weiNum)
//	fmt.Println(gasPrice)
//	gasLimit := big.NewInt(21000) // BSC的标准gas限制
//	// 计算 gas 费用
//	gasCost := new(big.Int).Mul(gasPrice, gasLimit)
//
//	// 计算可转出金额
//	remainingAmount := new(big.Int).Sub(weiNum, gasCost)
//
//	// 打印结果
//	fmt.Println("Remaining amount after gas cost:", remainingAmount.String())
//	fmt.Println("Remaining amount after gas cost:", Wei2ETH(remainingAmount))
//}

// Get_BSC_USDT_Remain 还剩下多少可以转出
//func (m *WalletClient) Get_BSC_USDT_Remain(address string, contractAddress string) (canOut *big.Int) {
//	remainWei, err := m.Get_BSC_USDT_Balance(address, contractAddress)
//	if err != nil {
//		return
//	}
//	if remainWei.Cmp(big.NewInt(0)) != 1 {
//		return
//	}
//	gasPrice, _ := m.GetSuggestGasPrice()
//	gasLimit := big.NewInt(21000) // BSC的标准gas限制
//	// 计算 gas 费用
//	gasCost := new(big.Int).Mul(gasPrice, gasLimit)
//
//	// 计算可转出金额
//	remainingAmount := new(big.Int).Sub(remainWei, gasCost)
//
//	// 打印结果
//	fmt.Println("Remaining amount after gas cost:", remainingAmount.String())
//	fmt.Println("Remaining amount after gas cost:", Wei2ETH(remainingAmount))
//
//	canOut = remainingAmount
//	return
//}

// 检测BSC地址
func IsValidBSCAddress(address string) bool {
	// 正则表达式模式，用于匹配以太坊地址格式
	pattern := `^0x[0-9a-fA-F]{40}$`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 检查地址是否匹配模式
	return re.MatchString(address)
}

// 生成分享图
func CreateShareImg(network, address string) (err error, path string) {
	qrCodeFilePath := fmt.Sprintf("./upload/share/%s_qr.png", address)
	shareFilePath := fmt.Sprintf("./upload/share/%s_share.png", address)
	path = fmt.Sprintf("/share/%s_share.png", address)
	var exit bool
	exit, err = helpers.FileExists(shareFilePath)
	if err != nil {
		return
	}

	if exit {
		return
	}

	// 生成二维码
	width := 360
	height := 500
	radius := 50
	err = qrcode.WriteFile(address, qrcode.Medium, 256, qrCodeFilePath)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	// 打开背景图
	//backgroundFilePath := "./static/sharebackground.jpg"
	//backgroundImage, err := gg.LoadImage(backgroundFilePath)
	//if err != nil {
	//	fmt.Println("Error loading background image:", err)
	//	return
	//}

	// 创建一个新的图像
	dc := gg.NewContext(width, height)

	// 绘制背景图
	//dc.DrawImage(backgroundImage, 0, 0)

	// 设置背景色为白色
	dc.SetColor(color.White)
	dc.Clear()

	// 绘制圆角矩形
	dc.SetRGB(0, 0, 0) // 设置描边颜色为黑色
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), float64(radius))

	// 加载二维码图像
	qrCodeImage, err := gg.LoadImage(qrCodeFilePath)
	if err != nil {
		fmt.Println("Error loading QR code image:", err)
		return
	}

	// 绘制二维码
	dc.DrawImage(qrCodeImage, (width-256)/2, 80)

	// 设置文字颜色
	//dc.SetColor(color.Black)

	// 设置字体和字体大小
	if err = dc.LoadFontFace("./static/font.ttf", 20); err != nil {
		fmt.Println("Error loading font:", err)
		return
	}

	// 网络
	textWidth, _ := dc.MeasureString(network)
	textX := (float64(width) - textWidth) / 2
	dc.DrawString(network, textX, 70)

	// address
	addressWidth, _ := dc.MeasureString("Address")
	addressTx := (float64(width) - addressWidth) / 2
	dc.DrawString("Address", addressTx, 360)

	// 地址
	address1 := address[:28]
	address1Width, _ := dc.MeasureString(address1)
	address2 := address[28:]
	address2Width, _ := dc.MeasureString(address2)
	dc.DrawString(address1, (float64(width)-address1Width)/2, 400)
	dc.DrawString(address2, (float64(width)-address2Width)/2, 420)

	// 保存图片
	if err = dc.SavePNG(shareFilePath); err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Image generated successfully!")
	return
}

func (m *WalletClient) Transfer1(privateKey, toAddress, contract string, num *big.Int) (string, error) {

	//从私钥推导出 公钥
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("crypto.HexToECDSA error ,", err)
		return "", err
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("publicKeyECDSA error ,", err)
		return "", err
	}
	//从公钥推导出钱包地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("钱包地址：", fromAddress.Hex())
	//构造请求参数
	var data []byte
	methodName := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]
	paddedToAddress := common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32)
	//amount, _ := new(big.Int).SetString("100000000000000000", 10)
	paddedAmount := common.LeftPadBytes(num.Bytes(), 32)
	data = append(data, methodName...)
	data = append(data, paddedToAddress...)
	data = append(data, paddedAmount...)

	//获取nonce
	nonce, err := m.client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", err
	}
	//获取小费
	gasTipCap, _ := m.client.SuggestGasTipCap(context.Background())
	//transfer 默认是 使用 21000 gas
	gas := uint64(100000)
	//最大gas fee
	gasFeeCap := big.NewInt(38694000460)

	contractAddress := common.HexToAddress(contract)
	//创建交易
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gas,
		To:        &contractAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})
	// 获取当前区块链的ChainID
	chainID, err := m.client.ChainID(context.Background())
	if err != nil {
		fmt.Println("获取ChainID失败:", err)
		return "", err
	}

	fmt.Println("当前区块链的ChainID:", chainID)
	//创建签名者
	signer := types.NewLondonSigner(chainID)
	//对交易进行签名
	signTx, err := types.SignTx(tx, signer, privateKeyECDSA)
	if err != nil {
		return "", err
	}
	//发送交易
	err = m.client.SendTransaction(context.Background(), signTx)
	if err != nil {
		return "", err
	}
	//返回交易哈希
	return signTx.Hash().Hex(), err

}
