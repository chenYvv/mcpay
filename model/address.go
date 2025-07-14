package models

import (
	"mcpay/pkg/database"
	"mcpay/pkg/hdwallet"
	"time"

	"gorm.io/gorm"
)

// Address 钱包地址池表
type Address struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Network    int       `gorm:"column:network;type:int(11);default:0;comment:网络：1:波场；2:币安;NOT NULL" json:"network"`
	BlockNum   int       `gorm:"column:blockNum;type:int(11);default:0;comment:最新区块高度;NOT NULL" json:"blockNum"`
	PathIndex  int       `gorm:"column:path_index;type:int(11);default:1;NOT NULL" json:"path_index"`
	Address    string    `gorm:"column:address;type:varchar(255);comment:地址;NOT NULL" json:"address"`
	Balance    float64   `gorm:"column:balance;type:decimal(36,18);default:0.000000000000000000;comment:剩余额度;NOT NULL" json:"balance"`
	PrivateKey string    `gorm:"column:private_key;type:varchar(255);comment:私钥;NOT NULL" json:"private_key"`
	Status     int       `gorm:"column:status;type:tinyint(1);default:1;comment:1:可用；2:使用中；3:禁用;NOT NULL" json:"status"`
	UsedTimes  int       `gorm:"column:used_times;type:int(11);default:0;comment:使用次数;NOT NULL" json:"used_times"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
	LastUsedAt time.Time `gorm:"column:last_used_at;type:datetime;comment:最后使用时间" json:"last_used_at"`
}

// TableName 指定表名
func (m *Address) TableName() string {
	return "address"
}

// 状态常量
const (
	AddressStatusAvailable = 1 // 可用
	AddressStatusInUse     = 2 // 使用中
	AddressStatusDisabled  = 3 // 禁用
)

// 网络类型常量
const (
	NetworkTron = 1
	NetworkBsc  = 2
)

// BeforeCreate 创建前的钩子
func (m *Address) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加地址验证逻辑
	return nil
}

// IsAvailable 检查地址是否可用
func (m *Address) IsAvailable() bool {
	return m.Status == AddressStatusAvailable
}

// IsInUse 检查地址是否在使用中
func (m *Address) IsInUse() bool {
	return m.Status == AddressStatusInUse
}

// MarkAsUsed 标记为使用中
func (m *Address) MarkAsUsed(tx *gorm.DB, lastUsedAt time.Time) error {
	//now := time.Now()
	return tx.Model(m).Updates(map[string]interface{}{
		"status":       AddressStatusInUse,
		"used_times":   gorm.Expr("used_times + 1"),
		"last_used_at": lastUsedAt,
	}).Error
}

// MarkAsAvailable 标记为可用
func (m *Address) MarkAsAvailable(tx *gorm.DB) error {
	return tx.Model(m).UpdateColumn("status", AddressStatusAvailable).Error
}

// MarkAsDisabled 标记为禁用
func (m *Address) MarkAsDisabled(tx *gorm.DB) error {
	return tx.Model(m).UpdateColumn("status", AddressStatusDisabled).Error
}

// IsTron 是波场
func (m *Address) IsTron() bool {
	return m.Network == NetworkTron
}

// IsBsc 是BSC
func (m *Address) IsBsc() bool {
	return m.Network == NetworkBsc
}

// 获取波场使用中的钱包地址列表
func GetTronInUseAddresses() ([]Address, error) {
	var addresses []Address
	err := database.DB.Where("network = ? AND status = ?", NetworkTron, AddressStatusInUse).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// 获取币安使用中的钱包地址列表
func GetBSCInUseAddresses() ([]Address, error) {
	var addresses []Address
	err := database.DB.Where("network = ? AND status = ?", NetworkBsc, AddressStatusInUse).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// 获取当前支付中的订单
func (m *Address) GetPendingOrder() *Order {
	order := &Order{}
	err := database.DB.Table("orders as o").
		Joins("LEFT JOIN order_address as oa on oa.order_id = o.id").
		Where("o.status = ? AND oa.address_id = ?", OrderStatusPending, m.Id).
		Select("o.*").
		Find(&order).Error

	if err != nil {
		return nil
	}
	return order
}

// 获取当前支付中的绑定关系数据
func (m *Address) GetPendingOrderAddress() *OrderAddress {
	res := &OrderAddress{}
	err := database.DB.Table("order_address as oa").
		Joins("LEFT JOIN address as a on oa.address_id = a.id").
		Where("a.status = ? AND a.id = ?", AddressStatusInUse, m.Id).
		Select("oa.*").
		Find(&res).Error

	if err != nil {
		return nil
	}
	return res
}

// 获取一个可用的钱包地址
func GetAvailableAddressByNetwork(network int) (*Address, error) {
	var address *Address

	//err := database.DB.Transaction(func(tx *gorm.DB) error {
	//	err := tx.Where("network = ? AND status = ?", network, AddressStatusAvailable).
	//		Order("used_times DESC").First(&address).Error
	//	if err != nil {
	//		return err
	//	}
	//	err = address.MarkAsUsed(tx)
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//
	//if err != nil {
	//	return nil, err
	//}

	return address, nil

}

// 批量写入地址
func CreateAddress(network int, count int) {
	// 获取最大 path_index
	var startIndex int
	database.DB.Model(&Address{}).
		Where("network = ?", network).
		Select("IFNULL(MAX(path_index), 0)").
		Scan(&startIndex)

	startIndex++

	mnemonic := GetMnemonic()

	wallet, err := hdwallet.LoadWallet(mnemonic)

	if err != nil {
		panic(err)
	}

	data := []Address{}
	now := time.Now()
	for i := 0; i < count; i++ {
		address := ""
		switch network {
		case NetworkTron:
			address, _, _ = wallet.DeriveTRONAddress(uint32(startIndex))
		case NetworkBsc:
			address, _, _ = wallet.DeriveETHAddress(uint32(startIndex))
		}

		data = append(data, Address{
			Network:    network,
			PathIndex:  startIndex,
			Address:    address,
			Balance:    0,
			PrivateKey: "",
			Status:     AddressStatusAvailable,
			UsedTimes:  0,
			CreatedAt:  now,
		})

		startIndex++

	}

	// 批量写入数据库
	if len(data) > 0 {
		err = database.DB.Create(&data).Error
		if err != nil {
			panic(err)
		}
	}

}
