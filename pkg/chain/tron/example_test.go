package tron_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"mcpay/pkg/tron"
)

// 创建一个钱包
func TestClient(t *testing.T) {
	addr, err := tron.GenerateAddress()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated address: %s\n", addr.Address)
	fmt.Printf("Private key: %s\n", addr.PrivateKey)
}

// 查看钱包余额
func TestWalletBalance(t *testing.T) {
	// 使用默认配置创建客户端
	tron.InitTronClient(true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	address := "TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s"

	// 查询 TRX 余额
	balance, err := tron.GetClient().GetTRXBalance(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TRX Balance: %f\n", tron.WeiToNum(balance)) // 修复：使用 %.6f 格式化 float64

	// 查询 USDT 余额
	usdtBalance, err := tron.GetClient().GetUSDTBalance(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDT Balance: %f\n", tron.WeiToNum(usdtBalance)) // 修复：使用 %.6f
}

// 测试 TRX 转账
func TestTransferTRX(t *testing.T) {
	// 使用测试网配置
	// client, err := tron.NewClient(tron.TestNetConfig())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Close()

	ctx := context.Background()

	// 转账 TRX
	req := &tron.TransferRequest{
		From:       "TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s",
		To:         "TGESHsRbgoo72QoMN1nfmzbCMGm362eeSn",
		Amount:     10.5, // 10.5 TRX
		PrivateKey: "b97a2695f42882a11926b44edf45f50d92efb59c8a0f98b33dc4d9d78b8f975d",
	}

	res, err := tron.GetClient().TransferTRX(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transfer result: %+v\n", res)
}

// 测试 USDT 转账
func TestTransferUSDT(t *testing.T) {
	// 使用测试网配置
	// client, err := tron.NewClient(tron.TestNetConfig())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Close()

	// 使用默认配置创建客户端
	tron.InitTronClient(true)

	ctx := context.Background()

	// 转账 USDT
	req := &tron.TransferRequest{
		From:       "TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s",
		To:         "TGESHsRbgoo72QoMN1nfmzbCMGm362eeSn",
		Amount:     10, // 10 USDT
		PrivateKey: "b97a2695f42882a11926b44edf45f50d92efb59c8a0f98b33dc4d9d78b8f975d",
	}

	res, err := tron.GetClient().TransferUSDT(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDT transfer result: %+v\n", res)
	fmt.Printf("USDT transfer txid: %s\n", res.TxID)

}

// 测试 USDT 转账
func TestOrderUSDT(t *testing.T) {
	// 使用测试网配置
	// client, err := tron.NewClient(tron.TestNetConfig())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Close()

	// 使用默认配置创建客户端
	tron.InitTronClient(true)

	transactions, err := tron.GetUSDTTransactions("TGESHsRbgoo72QoMN1nfmzbCMGm362eeSn", 0)
	if err != nil {
		log.Fatal(err)
	}

	// 计算总充值量
	totalRecharge := big.NewFloat(0)
	for _, tx := range transactions {
		wei := new(big.Int)
		wei.SetString(tx.Value, 10)                                         // 转换为 big.Int
		totalRecharge.Add(totalRecharge, tron.WeiToNumWithDecimals(wei, 6)) // 转换为 USDT 单位
	}

	fmt.Printf("Total transactions: %d\n", len(transactions))
	fmt.Printf("Total totalRecharge: %f\n", totalRecharge)

}
