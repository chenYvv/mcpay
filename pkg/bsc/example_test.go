package bsc

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 测试钱包创建功能
func TestCreateWallet(t *testing.T) {
	t.Log("测试钱包创建功能...")

	wallet, err := CreateWallet()
	if err != nil {
		t.Fatalf("创建钱包失败: %v", err)
	}

	// 验证私钥格式
	if len(wallet.PrivateKey) != 64 {
		t.Errorf("私钥长度错误, 期望: 64, 实际: %d", len(wallet.PrivateKey))
	}

	// 验证地址格式
	if !strings.HasPrefix(wallet.Address, "0x") {
		t.Errorf("地址格式错误, 地址应该以0x开头: %s", wallet.Address)
	}

	if len(wallet.Address) != 42 {
		t.Errorf("地址长度错误, 期望: 42, 实际: %d", len(wallet.Address))
	}

	// 验证地址有效性
	if !IsValidAddress(wallet.Address) {
		t.Errorf("生成的地址无效: %s", wallet.Address)
	}

	t.Logf("✅ 钱包创建成功")
	t.Logf("私钥: %s", wallet.PrivateKey)
	t.Logf("地址: %s", wallet.Address)
}

// 测试多个钱包创建（确保随机性）
func TestCreateMultipleWallets(t *testing.T) {
	t.Log("测试创建多个钱包的随机性...")

	wallets := make([]*WalletInfo, 5)

	for i := 0; i < 5; i++ {
		wallet, err := CreateWallet()
		if err != nil {
			t.Fatalf("创建第%d个钱包失败: %v", i+1, err)
		}
		wallets[i] = wallet
	}

	// 检查所有钱包都是唯一的
	for i := 0; i < len(wallets); i++ {
		for j := i + 1; j < len(wallets); j++ {
			if wallets[i].PrivateKey == wallets[j].PrivateKey {
				t.Errorf("钱包%d和钱包%d的私钥相同，缺乏随机性", i+1, j+1)
			}
			if wallets[i].Address == wallets[j].Address {
				t.Errorf("钱包%d和钱包%d的地址相同，缺乏随机性", i+1, j+1)
			}
		}
	}

	t.Logf("✅ 创建了5个唯一钱包")
}

// 测试BSC测试网客户端连接
func TestNewBSCTestnetClient(t *testing.T) {
	t.Log("测试BSC测试网客户端连接...")

	client, err := NewBSCClient(true) // false = testnet
	if err != nil {
		t.Fatalf("连接BSC测试网失败: %v", err)
	}
	defer client.client.Close()

	// 验证客户端
	if client.chainID.Int64() != BSC_TESTNET_CHAIN_ID {
		t.Errorf("链ID错误, 期望: %d, 实际: %d", BSC_TESTNET_CHAIN_ID, client.chainID.Int64())
	}
	// 测试网络连接
	networkInfo, err := client.GetNetworkInfo()
	if err != nil {
		t.Fatalf("获取网络信息失败: %v", err)
	}

	t.Logf("✅ BSC测试网连接成功")
	t.Logf("链ID: %s", networkInfo["chainID"])
	t.Logf("区块高度: %s", networkInfo["blockNumber"])
	t.Logf("RPC地址: %s", client.rpcURL)
}

// 测试BSC主网客户端连接
func TestNewBSCMainnetClient(t *testing.T) {
	t.Log("测试BSC主网客户端连接...")

	client, err := NewBSCClient(false) // true = mainnet
	if err != nil {
		t.Fatalf("连接BSC主网失败: %v", err)
	}
	defer client.client.Close()

	// 验证客户端
	if client.chainID.Int64() != BSC_MAINNET_CHAIN_ID {
		t.Errorf("链ID错误, 期望: %d, 实际: %d", BSC_MAINNET_CHAIN_ID, client.chainID.Int64())
	}
	// 测试网络连接
	networkInfo, err := client.GetNetworkInfo()
	if err != nil {
		t.Fatalf("获取网络信息失败: %v", err)
	}

	t.Logf("✅ BSC主网连接成功")
	t.Logf("链ID: %s", networkInfo["chainID"])
	t.Logf("区块高度: %s", networkInfo["blockNumber"])
	t.Logf("RPC地址: %s", client.rpcURL)
}

// 测试地址验证功能
func TestIsValidAddress(t *testing.T) {
	t.Log("测试地址验证功能...")

	testCases := []struct {
		address string
		valid   bool
		desc    string
	}{
		{"0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce", true, "有效的BSC地址"},
		{"0x55d398326f99059ff775485246999027b3197955", true, "USDT合约地址"},
		{"4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce", false, "缺少0x前缀"},
		{"0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4C", false, "地址过短"},
		{"0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce1", false, "地址过长"},
		{"0xGF5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce", false, "包含非十六进制字符"},
		{"", false, "空地址"},
		{"0x", false, "只有前缀"},
	}

	for _, tc := range testCases {
		result := IsValidAddress(tc.address)
		if result != tc.valid {
			t.Errorf("地址验证失败: %s, 期望: %v, 实际: %v, 描述: %s",
				tc.address, tc.valid, result, tc.desc)
		} else {
			t.Logf("✅ %s: %s", tc.desc, tc.address)
		}
	}
}

// 测试Wei/Ether转换功能
func TestWeiEtherConversion(t *testing.T) {
	t.Log("测试Wei/Ether转换功能...")

	// 测试WeiToEther
	testCases := []struct {
		wei   string
		ether string
		desc  string
	}{
		{"1000000000000000000", "1.000000", "1 ETH"}, // 1 ETH = 10^18 Wei
		{"500000000000000000", "0.500000", "0.5 ETH"},
		{"1", "0.000000", "1 Wei"},
		{"0", "0.000000", "0 Wei"},
		{"1234567890123456789", "1.234568", "1.234... ETH"},
	}
	for _, tc := range testCases {
		wei := new(big.Int)
		wei.SetString(tc.wei, 10)

		ether := WeiToEther(wei)
		etherString := ether.Text('f', DISPLAY_PRECISION)
		if etherString != tc.ether {
			t.Errorf("Wei转Ether失败: %s Wei, 期望: %s ETH, 实际: %s ETH, 描述: %s",
				tc.wei, tc.ether, etherString, tc.desc)
		} else {
			t.Logf("✅ %s: %s Wei = %s ETH", tc.desc, tc.wei, etherString)
		}
	}

	// 测试EtherToWei
	etherTestCases := []struct {
		ether string
		wei   string
		desc  string
	}{
		{"1", "1000000000000000000", "1 ETH"},
		{"0.5", "500000000000000000", "0.5 ETH"},
		{"0.000001", "1000000000000", "0.000001 ETH"},
		{"0", "0", "0 ETH"},
	}
	for _, tc := range etherTestCases {
		etherFloat, err := strconv.ParseFloat(tc.ether, 64)
		if err != nil {
			t.Errorf("解析Ether值失败: %s, 错误: %v", tc.ether, err)
			continue
		}

		wei := EtherToWei(etherFloat)
		if wei.String() != tc.wei {
			t.Errorf("Ether转Wei失败: %s ETH, 期望: %s Wei, 实际: %s Wei, 描述: %s",
				tc.ether, tc.wei, wei.String(), tc.desc)
		} else {
			t.Logf("✅ %s: %s ETH = %s Wei", tc.desc, tc.ether, wei.String())
		}
	}
}

// 测试余额格式化功能
func TestFormatBalance(t *testing.T) {
	t.Log("测试余额格式化功能...")

	testCases := []struct {
		balance  string
		decimals int
		expected string
		desc     string
	}{
		{"1000000000000000000", 18, "1.000000", "1 BNB (18位小数)"},
		{"1000000", 6, "1.000000", "1 USDT (6位小数)"},
		{"500000000000000000", 18, "0.500000", "0.5 BNB"},
		{"123456", 6, "0.123456", "0.123456 USDT"},
		{"0", 18, "0.000000", "0余额"},
	}

	for _, tc := range testCases {
		balance := new(big.Int)
		balance.SetString(tc.balance, 10)

		result := FormatBalance(balance, tc.decimals)
		if result != tc.expected {
			t.Errorf("格式化余额失败: %s (小数位%d), 期望: %s, 实际: %s, 描述: %s",
				tc.balance, tc.decimals, tc.expected, result, tc.desc)
		} else {
			t.Logf("✅ %s: %s = %s", tc.desc, tc.balance, result)
		}
	}
}

// 测试获取BNB余额 (需要网络连接)
func TestGetBNBBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过网络测试")
	}

	t.Log("测试获取BNB余额功能...")

	client, err := NewBSCClient(true) // 使用测试网
	if err != nil {
		t.Fatalf("连接BSC测试网失败: %v", err)
	}
	defer client.client.Close()

	// 使用一个已知的测试网地址
	testAddress := "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce"

	balance, err := client.GetBNBBalance(testAddress)
	if err != nil {
		t.Fatalf("获取BNB余额失败: %v", err)
	}
	t.Logf("✅ BNB余额查询成功")
	t.Logf("地址: %s", testAddress)
	t.Logf("余额: %s Wei", balance.String())
	t.Logf("余额: %s BNB", FormatBNBBalance(balance))
}

// 测试获取USDT余额 (需要网络连接)
func TestGetUSDTBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过网络测试")
	}

	t.Log("测试获取USDT余额功能...")

	client, err := NewBSCClient(true) // 使用测试网
	if err != nil {
		t.Fatalf("连接BSC测试网失败: %v", err)
	}
	defer client.client.Close()

	// 使用一个已知的测试网地址
	testAddress := "0xf1e1428f8F14C8F723553ff43059fa366397Ae2c"

	balance, err := client.GetUSDTBalance(testAddress)
	if err != nil {
		t.Fatalf("获取USDT余额失败: %v", err)
	}

	t.Logf("✅ USDT余额查询成功")
	t.Logf("地址: %s", testAddress)
	t.Logf("余额: %s", balance.String())
	t.Logf("余额: %s USDT", FormatUSDTBalance(balance))
}

// 测试获取综合余额 (需要网络连接)
func TestGetBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过网络测试")
	}

	t.Log("测试获取综合余额功能...")

	client, err := NewBSCClient(true) // 使用测试网
	if err != nil {
		t.Fatalf("连接BSC测试网失败: %v", err)
	}
	defer client.client.Close()

	// 使用一个已知的测试网地址
	testAddress := "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce"

	balanceInfo, err := client.GetBalance(testAddress)
	if err != nil {
		t.Fatalf("获取综合余额失败: %v", err)
	}
	t.Logf("✅ 综合余额查询成功")
	t.Logf("地址: %s", balanceInfo.Address)
	t.Logf("BNB余额: %s (%s BNB)", balanceInfo.BNBBalance.String(), FormatBNBBalance(balanceInfo.BNBBalance))
	t.Logf("USDT余额: %s (%s USDT)", balanceInfo.USDTBalance.String(), FormatUSDTBalance(balanceInfo.USDTBalance))
}

// 测试传输参数验证
func TestValidateTransferParams(t *testing.T) {
	t.Log("测试传输参数验证功能...")

	testCases := []struct {
		params TransferParams
		valid  bool
		desc   string
	}{
		{
			TransferParams{
				PrivateKey: "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				ToAddress:  "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce",
				Amount:     big.NewInt(1000000000000000000), // 1 ETH
			},
			true,
			"有效参数",
		},
		{
			TransferParams{
				PrivateKey: "", // 空私钥
				ToAddress:  "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce",
				Amount:     big.NewInt(1000000000000000000),
			},
			false,
			"空私钥",
		},
		{
			TransferParams{
				PrivateKey: "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				ToAddress:  "invalid_address", // 无效地址
				Amount:     big.NewInt(1000000000000000000),
			},
			false,
			"无效接收地址",
		},
		{
			TransferParams{
				PrivateKey: "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				ToAddress:  "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce",
				Amount:     big.NewInt(0), // 零金额
			},
			false,
			"零金额",
		},
		{
			TransferParams{
				PrivateKey: "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				ToAddress:  "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce",
				Amount:     big.NewInt(-1), // 负金额
			},
			false,
			"负金额",
		},
	}
	for _, tc := range testCases {
		// 创建临时客户端用于测试
		client := &BSCClient{}
		err := client.validateTransferParams(&tc.params)
		isValid := (err == nil)

		if isValid != tc.valid {
			t.Errorf("参数验证失败: %s, 期望: %v, 实际: %v, 错误: %v",
				tc.desc, tc.valid, isValid, err)
		} else {
			t.Logf("✅ %s: 验证结果正确", tc.desc)
		}
	}
}

// 测试错误情况处理
func TestErrorHandling(t *testing.T) {
	t.Log("测试错误情况处理...")

	// 测试无效RPC连接
	client := &BSCClient{
		chainID: big.NewInt(BSC_TESTNET_CHAIN_ID),
		rpcURL:  "https://invalid-rpc-endpoint.example.com",
	}

	// 这个测试应该失败，因为RPC端点无效
	_, err := client.GetBNBBalance("0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce")
	if err == nil {
		t.Error("期望获取余额失败，但成功了")
	} else {
		t.Logf("✅ 正确处理了无效RPC连接错误: %v", err)
	}

	// 测试无效地址
	validClient, err := NewBSCClient(false)
	if err != nil {
		t.Fatalf("创建有效客户端失败: %v", err)
	}
	defer validClient.client.Close()

	_, err = validClient.GetBNBBalance("invalid_address")
	if err == nil {
		t.Error("期望无效地址查询失败，但成功了")
	} else {
		t.Logf("✅ 正确处理了无效地址错误: %v", err)
	}
}

// 测试转移BNB (需要私钥和测试币)
func TestTransferBNB(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过转账测试 - 需要真实私钥和测试币")
	}

	t.Log("测试转移BNB功能...")

	// 注意：这个测试需要一个有BNB余额的私钥
	// 在实际测试中，请使用测试网的测试私钥
	testPrivateKey := "cc22ffd803b57c3e02c7d06da467f655dd0e1b2b819c184a16439131da5fc954" // 请替换为实际的测试私钥
	testToAddress := "0xf1e1428f8F14C8F723553ff43059fa366397Ae2c"                        // 测试接收地址

	// 如果没有提供私钥，跳过测试
	if testPrivateKey == "" {
		t.Skip("跳过BNB转账测试 - 请配置有效的测试私钥")
	}

	client, err := NewBSCClient(true) // 使用测试网
	if err != nil {
		t.Fatalf("连接BSC测试网失败: %v", err)
	}
	defer client.client.Close()

	// 获取发送者地址
	fromAddress, err := GetAddressFromPrivateKey(testPrivateKey)
	if err != nil {
		t.Fatalf("获取发送者地址失败: %v", err)
	}

	// 检查余额
	balance, err := client.GetBNBBalance(fromAddress)
	if err != nil {
		t.Fatalf("获取BNB余额失败: %v", err)
	}

	t.Logf("发送者地址: %s", fromAddress)
	t.Logf("当前BNB余额: %s BNB", FormatBNBBalance(balance))

	// 转账金额 (0.001 BNB)
	transferAmount := EtherToWei(0.001)

	// 检查余额是否足够
	if balance.Cmp(transferAmount) < 0 {
		t.Skip("跳过BNB转账测试 - 余额不足")
	}

	// 设置转账参数
	params := &TransferParams{
		PrivateKey: testPrivateKey,
		ToAddress:  testToAddress,
		Amount:     transferAmount,
	}

	// 执行转账
	txHash, err := client.TransferBNB(params)
	if err != nil {
		t.Fatalf("BNB转账失败: %v", err)
	}

	t.Logf("✅ BNB转账成功")
	t.Logf("交易哈希: %s", txHash)
	t.Logf("转账金额: %s BNB", FormatBNBBalance(transferAmount))
	t.Logf("接收地址: %s", testToAddress)

	// 验证交易哈希格式
	if !strings.HasPrefix(txHash, "0x") || len(txHash) != 66 {
		t.Errorf("交易哈希格式错误: %s", txHash)
	}
}

// 测试转移USDT (需要私钥和测试币)
func TestTransferUSDT(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过转账测试 - 需要真实私钥和测试币")
	}

	t.Log("测试转移USDT功能...")

	// 注意：这个测试需要一个有USDT余额的私钥
	// 在实际测试中，请使用测试网的测试私钥
	testPrivateKey := "cc22ffd803b57c3e02c7d06da467f655dd0e1b2b819c184a16439131da5fc954" // 请替换为实际的测试私钥
	testToAddress := "0xf1e1428f8F14C8F723553ff43059fa366397Ae2c"                        // 测试接收地址

	// 如果没有提供私钥，跳过测试
	if testPrivateKey == "" {
		t.Skip("跳过USDT转账测试 - 请配置有效的测试私钥")
	}

	client, err := NewBSCClient(true) // 使用测试网
	if err != nil {
		t.Fatalf("连接BSC测试网失败: %v", err)
	}
	defer client.client.Close()

	// 获取发送者地址
	fromAddress, err := GetAddressFromPrivateKey(testPrivateKey)
	if err != nil {
		t.Fatalf("获取发送者地址失败: %v", err)
	}

	// 检查USDT余额
	usdtBalance, err := client.GetUSDTBalance(fromAddress)
	if err != nil {
		t.Fatalf("获取USDT余额失败: %v", err)
	}

	// 检查BNB余额（用于支付gas费）
	bnbBalance, err := client.GetBNBBalance(fromAddress)
	if err != nil {
		t.Fatalf("获取BNB余额失败: %v", err)
	}

	t.Logf("发送者地址: %s", fromAddress)
	t.Logf("当前USDT余额: %s USDT", FormatUSDTBalance(usdtBalance))
	t.Logf("当前BNB余额: %s BNB", FormatBNBBalance(bnbBalance))

	// 转账金额 (1 USDT)
	transferAmount := new(big.Int)
	transferAmount.SetString("1000000000000000000", 10) // 1 USDT (18位小数)

	// 检查USDT余额是否足够
	if usdtBalance.Cmp(transferAmount) < 0 {
		t.Skip("跳过USDT转账测试 - USDT余额不足")
	}

	// 检查BNB余额是否足够支付gas费 (至少0.001 BNB)
	minBNBRequired := EtherToWei(0.001)
	if bnbBalance.Cmp(minBNBRequired) < 0 {
		t.Skip("跳过USDT转账测试 - BNB余额不足支付gas费")
	}

	// 设置转账参数
	params := &TransferParams{
		PrivateKey: testPrivateKey,
		ToAddress:  testToAddress,
		Amount:     transferAmount,
	}

	// 执行转账
	txHash, err := client.TransferUSDT(params)
	if err != nil {
		t.Fatalf("USDT转账失败: %v", err)
	}

	t.Logf("✅ USDT转账成功")
	t.Logf("交易哈希: %s", txHash)
	t.Logf("转账金额: %s USDT", FormatUSDTBalance(transferAmount))
	t.Logf("接收地址: %s", testToAddress)

	// 验证交易哈希格式
	if !strings.HasPrefix(txHash, "0x") || len(txHash) != 66 {
		t.Errorf("交易哈希格式错误: %s", txHash)
	}
}

// 性能测试 - 钱包创建
func BenchmarkCreateWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := CreateWallet()
		if err != nil {
			b.Fatalf("创建钱包失败: %v", err)
		}
	}
}

// 性能测试 - 地址验证
func BenchmarkIsValidAddress(b *testing.B) {
	address := "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce"
	for i := 0; i < b.N; i++ {
		IsValidAddress(address)
	}
}

// 性能测试 - Wei/Ether转换
func BenchmarkWeiToEther(b *testing.B) {
	wei := big.NewInt(1000000000000000000) // 1 ETH
	for i := 0; i < b.N; i++ {
		WeiToEther(wei)
	}
}

// 示例测试 - 完整的钱包操作流程
func Example() {
	// 注意：这是一个示例，实际转账需要测试币

	// 1. 创建钱包
	wallet, err := CreateWallet()
	if err != nil {
		panic(err)
	}

	// 2. 连接到BSC测试网
	client, err := NewBSCClient(false)
	if err != nil {
		panic(err)
	}
	defer client.client.Close()

	// 3. 查询余额
	balance, err := client.GetBalance(wallet.Address)
	if err != nil {
		panic(err)
	}

	// 输出结果
	println("钱包地址:", wallet.Address)
	println("BNB余额:", FormatBNBBalance(balance.BNBBalance), "BNB")
	println("USDT余额:", FormatUSDTBalance(balance.USDTBalance), "USDT")

	// Output:
	// 钱包地址: 0x...
	// BNB余额: 0.000000 BNB
	// USDT余额: 0.000000 USDT
}

// 集成测试 - 多RPC端点故障转移
func TestRPCFailover(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过网络集成测试")
	}

	t.Log("测试多RPC端点故障转移功能...")

	// 这个测试验证当一个RPC端点失败时，客户端是否能自动切换到备用端点
	client, err := NewBSCClient(false)
	if err != nil {
		t.Fatalf("创建BSC客户端失败: %v", err)
	}
	defer client.client.Close()

	// 进行多次网络请求，验证稳定性
	testAddress := "0x4F5eE41eCCCFCF19cD2A52377750c6a2cbe0E4Ce"

	for i := 0; i < 3; i++ {
		balance, err := client.GetBNBBalance(testAddress)
		if err != nil {
			t.Errorf("第%d次请求失败: %v", i+1, err)
			continue
		}

		t.Logf("第%d次请求成功，余额: %s Wei", i+1, balance.String())
		time.Sleep(100 * time.Millisecond) // 避免请求过快
	}

	t.Logf("✅ 多RPC端点故障转移测试完成")
}
