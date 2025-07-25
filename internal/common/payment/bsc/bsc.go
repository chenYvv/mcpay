package tron

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mcpay/internal/common/payment"
	models "mcpay/model"
	"net/http"
)

// 自动注册
func init() {
	// 这里可以从配置文件读取配置
	payment.Register(
		payment.ChannelBsc,
		&BscPayment{
			apiKey: viper.GetString("server.port"),
		},
	)
}

// BscPayment TRON支付实现
type BscPayment struct {
	apiKey string
}

// Collect 实现代收功能
func (t *BscPayment) Collect(tx *gorm.DB, order models.Order, metaData map[string]string) (payment.CollectResult, error) {
	err := order.BindAvailableAddressByNetwork(tx, models.NetworkBsc)
	if err != nil {
		return payment.CollectResult{}, err
	}

	// TODO: 实现TRON代收逻辑
	return payment.CollectResult{}, nil
}

// CollectCallback 实现代收回调验证
func (t *BscPayment) CollectCallback(r *http.Request) (payment.CollectCallbackResult, error) {
	// TODO: 实现TRON回调验证逻辑
	return payment.CollectCallbackResult{}, nil
}

// Payout 实现代付功能
func (t *BscPayment) Payout() (payment.PayoutResult, error) {
	// TODO: 实现TRON代付逻辑
	return payment.PayoutResult{}, nil
}

// PayoutCallback 实现代付回调验证
func (t *BscPayment) PayoutCallback(r *http.Request) (payment.PayoutCallbackResult, error) {
	// TODO: 实现TRON代付回调验证逻辑
	return payment.PayoutCallbackResult{}, nil
}
