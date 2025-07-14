package models

import (
	"go.uber.org/zap"
	"log/slog"
	"mcpay/pkg/bsc"
	"mcpay/pkg/database"
	"mcpay/pkg/logger"
	"time"

	"gorm.io/gorm"
)

// 订单表
type Order struct {
	Id              int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Uid             int       `gorm:"column:uid;type:int(11)" json:"uid"`
	AppId           int       `gorm:"column:app_id;type:int(11);NOT NULL" json:"app_id"`
	OrderId         string    `gorm:"column:order_id;type:varchar(255);comment:订单号;NOT NULL" json:"order_id"`
	Amount          float64   `gorm:"column:amount;type:decimal(10,2) unsigned;default:0.00;comment:订单金额;NOT NULL" json:"amount"`
	AmountTrue      float64   `gorm:"column:amount_true;type:decimal(10,2) unsigned;default:0.00;comment:到账金额;NOT NULL" json:"amount_true"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;NOT NULL" json:"created_at"`
	CallbackUrl     string    `gorm:"column:callback_url;type:varchar(255);comment:通知地址" json:"callback_url"`
	CallbackState   int       `gorm:"column:callback_state;type:int(11);default:0;comment:通知状态：0:失败；1:成功；;NOT NULL" json:"callback_state"`
	CallbackErr     string    `gorm:"column:callback_err;type:varchar(500);comment:通知报错信息" json:"callback_err"`
	CallbackDate    time.Time `gorm:"column:callback_date;type:datetime;comment:最后通知时间" json:"callback_date"`
	CallbackTimes   int       `gorm:"column:callback_times;type:int(11);default:0;comment:通知次数;NOT NULL" json:"callback_times"`
	MerchantOrderId string    `gorm:"column:merchant_order_id;type:varchar(255);comment:商户订单号" json:"merchant_order_id"`
	ThirdOrderId    string    `gorm:"column:third_order_id;type:varchar(255);comment:第三方订单号" json:"third_order_id"`
	Status          int       `gorm:"column:status;type:int(11);default:1;comment:1:待支付；2:成功；3:失败" json:"status"`
	RedirectUrl     string    `gorm:"column:redirect_url;type:varchar(255);comment:支付成功跳转地址" json:"redirect_url"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`

	// 非数据库
	Address []Address `gorm:"-" json:"address"`
}

// TableName 指定表名
func (m *Order) TableName() string {
	return "orders"
}

// 状态常量
const (
	OrderStatusPending = 1 // 待支付 可用
	OrderStatusSuccess = 2 // 支付完成
	OrderStatusTimeout = 3 // 超时
	OrderStatusFail    = 4 // 失败

	OrderTimeLimit = 10 * 60 // 订单超时时间
)

// 通过 bsc tron 字符串 获取 NetworkTRON NetworkBSC 常量
func GetNetworkByString(network string) int {
	switch network {
	case "TRON":
		return NetworkTron
	case "BSC":
		return NetworkBsc
	default:
		return 0 // 未知网络
	}
}

// MarkStatusSuccess 标记为支付完成
func (m *Order) MarkStatusSuccess() error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(m).UpdateColumn("status", OrderStatusSuccess).Error
		if err != nil {
			return err
		}

		// 对应地址的状态修改
		err = tx.Exec(
			`UPDATE address AS a LEFT JOIN order_address oa ON oa.address_id = a.id SET a.status = ? WHERE oa.order_id = ?`,
			AddressStatusAvailable,
			m.Id,
		).Error

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("更新支付完成订单状态失败:", err.Error())
		return err
	}

	return nil
}

// 创建一个 OrderId 唯一的订单号
func CreateOrderId() string {
	return time.Now().Format("20060102150405") + "-" + RandomString(8)
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[RandomInt(len(charset))]
	}
	return string(result)
}

// RandomInt 生成指定范围内的随机整数
func RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return int(time.Now().UnixNano() % int64(max))
}

// BindAvailableAddressByNetwork 给订单绑定一个可以支付的地址
func (m *Order) BindAvailableAddressByNetwork(tx *gorm.DB, network int) error {
	timeStart := time.Now().Add(-time.Duration(OrderTimeLimit+300) * time.Second)
	var address *Address
	err := tx.Where("network = ? AND status = ? AND (last_used_at <= ? OR last_used_at = '00:00:00 00:00:00')",
		network, AddressStatusAvailable, timeStart).
		Order("used_times DESC").First(&address).Error
	if err != nil {
		return err
	}

	// 设置为使用中
	err = address.MarkAsUsed(tx, m.CreatedAt)
	if err != nil {
		return err
	}

	blockNum := 0
	if network == NetworkTron {
		//tron.GetClient().
	} else if network == NetworkBsc {
		blockNum = bsc.GetClient().BlockNumber()
	}

	// 绑定
	data := OrderAddress{
		OrderId:   m.Id,
		AddressId: address.Id,
		BlockNum:  blockNum,
	}

	err = tx.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

// 订单回调
func (m *Order) Callback() {
	// 支付中就看看金额改成已经支付完成
	if m.Status == OrderStatusPending {
		if m.AmountTrue >= m.Amount {
			err := m.MarkStatusSuccess()
			if err != nil {
			}
			logger.Info("订单回调:更新支付状态：已完成", slog.String("订单号", m.OrderId))
		}
	}

	if m.Status == OrderStatusSuccess {
		// 更新金额
		err := database.DB.Model(&Order{}).
			Where("id = ?", m.Id).
			Updates(map[string]interface{}{
				"callback_times": gorm.Expr("callback_times + 1"),
				"callback_date":  time.Now(),
			})
		if err != nil {

		}
		logger.Info("订单回调:完成", slog.String("订单号", m.OrderId))
	}

}

// 获取通过链上交易的且没有超时的订单
func GetChainPendingOrdersWithAddress() []Order {
	var orders []Order

	database.DB.Table("  as oa").
		Joins("LEFT JOIN orders as o on oa.order_id = o.id").
		Where("o.status = ?", OrderStatusPending).
		Select("o.*").
		Find(&orders)

	var address []Address
	database.DB.Table("address as a").
		Joins("LEFT JOIN order_address as oa on a.id = oa.address_id").
		Joins("LEFT JOIN order as o on o.id = oa.order_id").
		Select("a.*").
		Where("o.status = ?", OrderStatusPending).
		Find(&address)

	//if len(orders) > 0 {
	//	for i, order := range orders {
	//		for _, addr := range address {
	//			if addr.
	//		}
	//	}
	//}

	return orders
}

// UpdateTimeoutOrders 订单超时状态修改
func UpdateTimeoutOrders() {

	thresholdTime := time.Now().Add(-time.Duration(OrderTimeLimit) * time.Second)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 订单超时 状态修改
		err := tx.Model(&Order{}).
			Where("status = ? AND created_at < ?", OrderStatusPending, thresholdTime).
			Updates(map[string]interface{}{"status": OrderStatusTimeout}).Error
		if err != nil {
			return err
		}

		// 对应地址的状态修改
		err = tx.Model(&Address{}).
			Where("status = ? AND last_used_at < ?", AddressStatusInUse, thresholdTime).
			Updates(map[string]interface{}{"status": AddressStatusAvailable}).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("超时订单状态更新失败", zap.String("错误", err.Error()))
		return
	} else {
		logger.Info("超时订单状态更新完成")
	}

}
