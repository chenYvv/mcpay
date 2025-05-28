package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

/**
主网络：grpc.trongrid.io:50051
主网络 USDT 合约地址：TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t

nile测试网络：grpc.nile.trongrid.io:50051
nile测试网络 USDT 合约地址：TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj
*/

type TronClient struct {
	node string
	GRPC *client.GrpcClient
}

var TronClientIns *TronClient

func NewTronClient(node string) (*TronClient, error) {
	if TronClientIns != nil {
		return TronClientIns, nil
	}
	TronClientIns = new(TronClient)
	TronClientIns.node = node
	TronClientIns.GRPC = client.NewGrpcClient(node)
	err := TronClientIns.GRPC.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc client start error: %v", err)
	}
	return TronClientIns, nil
}

// 生成地址
//
//	func CreateTRC20Address() (privateKey, address string) {
//		privateKey, address = genkeys.GenerateKey()
//		return
//	}
//
// GenerateAddressFromPrivateKey 从私钥字符串生成波场地址
func GenerateAddressFromPrivateKey(privateKeyHex string) (string, error) { // 将私钥字符串解码为字节格式
	privateKeyBytes, err := hexutil.Decode(privateKeyHex)
	if err != nil {
		return "", err
	}

	// 通过私钥字节生成私钥对象
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", err
	}

	// 从私钥获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}

	// 从公钥获取地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

func GetTronAddressByPrivateKey(hexPrivateKey string) string {
	// Decode hex string to bytes
	privateKeyBytes, _ := hex.DecodeString(hexPrivateKey)

	// Generate private key from bytes
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	address = "41" + address[2:]
	addb, _ := hex.DecodeString(address)
	firstHash := sha256.Sum256(addb)
	secondHash := sha256.Sum256(firstHash[:])
	secret := secondHash[:4]
	addb = append(addb, secret...)
	return base58.Encode(addb)
}

func CreateTronAddress() (pk, b5 string) {
	privateKey, _ := crypto.GenerateKey()
	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	address = "41" + address[2:]
	addb, _ := hex.DecodeString(address)
	firstHash := sha256.Sum256(addb)
	secondHash := sha256.Sum256(firstHash[:])
	secret := secondHash[:4]
	addb = append(addb, secret...)
	return hexutil.Encode(privateKeyBytes)[2:], base58.Encode(addb)
}

func (c *TronClient) SetTimeout(timeout time.Duration) error {
	if c == nil {
		return errors.New("client is nil ptr")
	}
	c.GRPC = client.NewGrpcClientWithTimeout(c.node, timeout)
	err := c.GRPC.Start()
	if err != nil {
		return fmt.Errorf("grpc start error: %v", err)
	}
	return nil
}

/*
保持连接，如果中途连接失败，就重连
*/
func (c *TronClient) keepConnect() error {
	_, err := c.GRPC.GetNodeInfo()
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return c.GRPC.Reconnect(c.node)
		}
		return fmt.Errorf("node connect error: %v", err)
	}
	return nil
}

func (c *TronClient) Transfer(from, to string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.Transfer(from, to, amount)
}

func (c *TronClient) TransferBroadcast(fromPrivateKey, from, to string, amount int64) (*api.TransactionExtention, error) {
	tx, err := c.Transfer(from, to, amount)
	if err != nil {
		return nil, err
	}

	signTx, err := c.SignTransaction(tx.Transaction, fromPrivateKey)
	if err != nil {
		return nil, err
	}

	err = c.BroadcastTransaction(signTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *TronClient) GetTrc10Balance(addr, assetId string) (int64, error) {
	err := c.keepConnect()
	if err != nil {
		return 0, err
	}
	acc, err := c.GRPC.GetAccount(addr)
	if err != nil || acc == nil {
		return 0, fmt.Errorf("get %s account error: %v", addr, err)
	}
	for key, value := range acc.AssetV2 {
		if key == assetId {
			return value, nil
		}
	}
	return 0, fmt.Errorf("%s do not find this assetID=%s amount", addr, assetId)
}

func (c *TronClient) GetTrxBalance(addr string) (*core.Account, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.GetAccount(addr)
}

func (c *TronClient) GetTrc20Balance(addr, contractAddress string) (*big.Int, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.TRC20ContractBalance(addr, contractAddress)
}

func (c *TronClient) TransferTrc10(from, to, assetId string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	fromAddr, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, fmt.Errorf("from address is not equal")
	}
	toAddr, err := address.Base58ToAddress(to)
	if err != nil {
		return nil, fmt.Errorf("to address is not equal")
	}
	return c.GRPC.TransferAsset(fromAddr.String(), toAddr.String(), assetId, amount)
}

func (c *TronClient) TransferTrc20(from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	tx, err := c.GRPC.TRC20Send(from, to, contract, amount, feeLimit)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// TransferTrc20Broadcast 转账并广播
func (c *TronClient) TransferTrc20Broadcast(fromPrivateKey, from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	tx, err := c.TransferTrc20(from, to, contract, amount, feeLimit)
	if err != nil {
		return nil, err
	}

	signTx, err := c.SignTransaction(tx.Transaction, fromPrivateKey)
	if err != nil {
		return nil, err
	}

	err = c.BroadcastTransaction(signTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *TronClient) BroadcastTransaction(transaction *core.Transaction) error {
	err := c.keepConnect()
	if err != nil {
		return err
	}
	result, err := c.GRPC.Broadcast(transaction)
	if err != nil {
		return fmt.Errorf("broadcast transaction error: %v", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("bad transaction: %v", string(result.GetMessage()))
	}
	if result.Result == true {
		return nil
	}
	d, _ := json.Marshal(result)
	return fmt.Errorf("tx send fail: %s", string(d))
}

// TronUSDT2BigInt 将TRON的USDT数量转换为bigint
func TronUSDT2BigInt(amount float64) *big.Int {
	// 将USDT数量乘以10^6，以将小数点精度移动到整数部分
	usdtBigInt := big.NewFloat(amount * 1e6)
	// 将结果转换为bigint类型
	bigInt, _ := usdtBigInt.Int(nil)
	return bigInt
}

func TronUSDT2Num(wei *big.Int) float64 {
	if wei.Cmp(big.NewInt(0)) == 0 {
		return 0
	}
	usdt := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt64(int64(1e6)))
	usdtFloat, _ := usdt.Float64()
	return usdtFloat
}

func TronTRX2Num(wei *big.Int) float64 {
	if wei.Cmp(big.NewInt(0)) == 0 {
		return 0
	}
	float := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt64(int64(1e6)))
	num, _ := float.Float64()
	return num
}

func Num2Wei(amount, y float64) *big.Int {
	// 将USDT数量乘以10^6，以将小数点精度移动到整数部分
	usdtBigInt := big.NewFloat(amount * math.Pow(10, y))
	// 将结果转换为bigint类型
	bigInt, _ := usdtBigInt.Int(nil)
	return bigInt
}

func (c *TronClient) SignTransaction(transaction *core.Transaction, privateKey string) (*core.Transaction, error) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("hex decode private key error: %v", err)
	}
	priv := crypto.ToECDSAUnsafe(privateBytes)
	defer c.zeroKey(priv)
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("proto marshal tx raw data error: %v", err)
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, priv)
	if err != nil {
		return nil, fmt.Errorf("sign error: %v", err)
	}
	transaction.Signature = append(transaction.Signature, signature)
	return transaction, nil
}

// zeroKey zeroes a private key in memory.
func (c *TronClient) zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func IsValidTronAddress(address string) bool {
	if _, err := common.DecodeCheck(address); err != nil {
		return false
	}
	return true
}
