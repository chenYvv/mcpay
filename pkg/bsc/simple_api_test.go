package bsc

import (
	"testing"
)

// 测试简化API - 初始化和获取客户端
func TestSimplifiedAPI(t *testing.T) {
	t.Log("测试简化BSC客户端API...")

	// 1. 测试测试网初始化
	err := InitBSCClient(true)
	if err != nil {
		t.Fatalf("初始化BSC测试网失败: %v", err)
	}

	client := GetBSCClient()
	if client == nil {
		t.Fatal("获取BSC客户端失败")
	}

	// 验证测试网配置
	if !client.IsTestnet() {
		t.Error("客户端应该配置为测试网")
	}

	if client.GetChainID().Int64() != BSC_TESTNET_CHAIN_ID {
		t.Errorf("链ID错误, 期望: %d, 实际: %d", BSC_TESTNET_CHAIN_ID, client.GetChainID().Int64())
	}

	if client.GetUSDTContract() != BSC_TESTNET_USDT_CONTRACT {
		t.Errorf("USDT合约地址错误, 期望: %s, 实际: %s", BSC_TESTNET_USDT_CONTRACT, client.GetUSDTContract())
	}

	t.Logf("✅ 测试网客户端配置正确")
	t.Logf("链ID: %d", client.GetChainID().Int64())
	t.Logf("RPC地址: %s", client.GetRPCURL())
	t.Logf("USDT合约: %s", client.GetUSDTContract())

	// 2. 测试重置到主网（模拟）
	err = ResetBSCClient(false)
	if err != nil {
		// 主网可能连接失败，这是正常的
		t.Logf("主网连接失败（预期）: %v", err)

		// 重置回测试网
		err = ResetBSCClient(true)
		if err != nil {
			t.Fatalf("重置回测试网失败: %v", err)
		}
	} else {
		// 如果主网连接成功，验证配置
		client = GetBSCClient()
		if client.IsTestnet() {
			t.Error("客户端应该配置为主网")
		}

		if client.GetChainID().Int64() != BSC_MAINNET_CHAIN_ID {
			t.Errorf("链ID错误, 期望: %d, 实际: %d", BSC_MAINNET_CHAIN_ID, client.GetChainID().Int64())
		}

		t.Logf("✅ 主网客户端配置正确")

		// 重置回测试网以确保其他测试正常
		err = ResetBSCClient(true)
		if err != nil {
			t.Fatalf("重置回测试网失败: %v", err)
		}
	}

	t.Logf("✅ 简化API测试完成")
}

// 测试全局客户端单例行为
func TestGlobalClientSingleton(t *testing.T) {
	t.Log("测试全局客户端单例行为...")

	// 初始化测试网客户端
	err := InitBSCClient(true)
	if err != nil {
		t.Fatalf("初始化BSC测试网失败: %v", err)
	}

	// 获取客户端实例
	client1 := GetBSCClient()
	client2 := GetBSCClient()

	// 验证是同一个实例
	if client1 != client2 {
		t.Error("GetBSCClient应该返回同一个全局实例")
	}

	// 验证配置一致
	if client1.IsTestnet() != client2.IsTestnet() {
		t.Error("客户端实例配置不一致")
	}

	if client1.GetChainID().Int64() != client2.GetChainID().Int64() {
		t.Error("客户端实例链ID不一致")
	}

	t.Logf("✅ 全局客户端单例行为正确")
}

// 测试客户端重置功能
func TestClientReset(t *testing.T) {
	t.Log("测试客户端重置功能...")

	// 初始化测试网
	err := InitBSCClient(true)
	if err != nil {
		t.Fatalf("初始化BSC测试网失败: %v", err)
	}

	client1 := GetBSCClient()
	if client1 == nil {
		t.Fatal("获取第一个客户端实例失败")
	}

	originalRPC := client1.GetRPCURL()
	t.Logf("原始RPC地址: %s", originalRPC)

	// 重置客户端（仍然是测试网）
	err = ResetBSCClient(true)
	if err != nil {
		t.Fatalf("重置BSC测试网失败: %v", err)
	}

	client2 := GetBSCClient()
	if client2 == nil {
		t.Fatal("获取重置后的客户端实例失败")
	}

	// 验证是否为新实例
	if client1 == client2 {
		t.Error("重置后应该是新的客户端实例")
	}

	// 验证配置仍然正确
	if !client2.IsTestnet() {
		t.Error("重置后的客户端应该仍然是测试网")
	}

	if client2.GetChainID().Int64() != BSC_TESTNET_CHAIN_ID {
		t.Errorf("重置后的客户端链ID错误, 期望: %d, 实际: %d", BSC_TESTNET_CHAIN_ID, client2.GetChainID().Int64())
	}

	t.Logf("重置后RPC地址: %s", client2.GetRPCURL())
	t.Logf("✅ 客户端重置功能正常")
}
