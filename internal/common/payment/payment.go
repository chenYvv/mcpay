package payment

import (
	"fmt"
	"gorm.io/gorm"
	models "mcpay/model"
	"net/http"
	"sync"
)

// 支付渠道常量
const (
	ChannelTron = 1
	ChannelBsc  = 2
)

// Payment 支付接口定义
type Payment interface {
	// Collect 代收功能（从用户收款）
	Collect(tx *gorm.DB, order models.Order, metaData map[string]string) (CollectResult, error)
	// CollectCallback 代收回调验证
	CollectCallback(r *http.Request) (CollectCallbackResult, error)

	// Payout 代付功能（向用户付款）
	Payout() (PayoutResult, error)
	// PayoutCallback 代付回调验证
	PayoutCallback(r *http.Request) (PayoutCallbackResult, error)
}

// CollectResult 统一支付结果结构
type CollectResult struct {
	Url           string                 `json:"url"`            // 跳转URL
	ThirdOrderId  string                 `json:"third_order_id"` // 第三方订单号
	PaymentParams map[string]interface{} `json:"payment_params"` // 支付参数（如APPID、订单号等）
}

type CollectCallbackResult struct {
	Success       bool                   // 验证结果
	TransactionID string                 // 交易ID
	OrderID       string                 // 订单ID
	Amount        float64                // 金额
	Status        string                 // 状态
	RawData       map[string]interface{} // 原始回调数据
}

// 代付相关数据结构
type RecipientInfo struct {
	AccountType   string // 账户类型（银行卡/微信/支付宝等）
	AccountNumber string // 账号
	Name          string // 姓名
	BankCode      string // 银行代码（可选）
}

type PayoutResult struct {
	TransferID    string                 // 转账ID
	TransactionID string                 // 交易ID
	Status        string                 // 状态
	ExtraData     map[string]interface{} // 额外数据
}

type PayoutCallbackResult struct {
	Success    bool                   // 验证结果
	TransferID string                 // 转账ID
	Amount     float64                // 金额
	Status     string                 // 状态
	RawData    map[string]interface{} // 原始回调数据
}

// PaymentRegistry 支付注册器
type PaymentRegistry struct {
	providers map[int]Payment
	mu        sync.RWMutex
}

// 全局注册器实例
var globalRegistry = &PaymentRegistry{
	providers: make(map[int]Payment),
}

// Register 注册支付提供商
func Register(channel int, provider Payment) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	globalRegistry.providers[channel] = provider
	fmt.Printf("Payment provider registered: %s\n", channel)
}

// GetProvider 获取支付提供商
func GetProvider(channel int) (Payment, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	provider, exists := globalRegistry.providers[channel]
	if !exists {
		return nil, fmt.Errorf("payment provider not found: %s", channel)
	}

	return provider, nil
}
