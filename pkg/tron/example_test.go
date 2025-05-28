package tron_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"mcpay/pkg/tron"
)

func TestClient(t *testing.T) {
	// 使用默认配置创建客户端
	client, err := tron.NewClient(tron.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 生成新地址
	addr, err := tron.GenerateAddress()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated address: %s\n", addr.Address)
	fmt.Printf("Private key: %s\n", addr.PrivateKey)

	// 查询余额
	balance, err := client.GetTRXBalance(ctx, addr.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TRX Balance: %.6f\n", tron.WeiToTRX(balance)) // 修复：使用 %.6f 格式化 float64

	// 查询 USDT 余额
	usdtBalance, err := client.GetUSDTBalance(ctx, addr.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDT Balance: %.6f\n", tron.WeiToUSDT(usdtBalance)) // 修复：使用 %.6f
}

func ExampleClient_transfer() {
	// 使用测试网配置
	client, err := tron.NewClient(tron.TestNetConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// 转账 TRX
	req := &tron.TransferRequest{
		From:       "TFromAddress...",
		To:         "TToAddress...",
		Amount:     tron.TRXToWei(10.5), // 10.5 TRX
		PrivateKey: "your_private_key",
	}

	result, err := client.TransferTRX(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transfer result: %+v\n", result)

	// 转账 USDT
	usdtReq := &tron.TransferRequest{
		From:       "TFromAddress...",
		To:         "TToAddress...",
		Amount:     tron.USDTToWei(100), // 100 USDT
		PrivateKey: "your_private_key",
		FeeLimit:   50000000, // 50 TRX
	}

	usdtResult, err := client.TransferUSDT(ctx, usdtReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDT transfer result: %+v\n", usdtResult)
}

// 或者如果你想要一个简单的测试函数而不是 Example，可以这样写：
func TestBasicUsage(t *testing.T) {
	// 跳过测试，因为需要真实的网络连接
	t.Skip("Skipping integration test")

	client, err := tron.NewClient(tron.TestNetConfig())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 生成新地址
	addr, err := tron.GenerateAddress()
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	t.Logf("Generated address: %s", addr.Address)
	t.Logf("Private key: %s", addr.PrivateKey)

	// 查询余额
	balance, err := client.GetTRXBalance(ctx, addr.Address)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}

	t.Logf("TRX Balance: %.6f", tron.WeiToTRX(balance))
}
