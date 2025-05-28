package models

import (
	"time"

	"gorm.io/gorm"
)

// AddressPool 钱包地址池表
type AddressPool struct {
	Id         uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Network    int       `gorm:"column:network;type:int(11);default:0;comment:网络：1:波场；2:币安;NOT NULL" json:"network"`
	Address    string    `gorm:"column:address;type:varchar(255);comment:地址;NOT NULL" json:"address"`
	Balance    string    `gorm:"column:balance;type:varchar(50);default:0;comment:剩余额度;NOT NULL" json:"balance"`
	PrivateKey string    `gorm:"column:private_key;type:varchar(255);comment:私钥;NOT NULL" json:"private_key"`
	Status     int       `gorm:"column:status;type:tinyint(1);default:1;comment:1:可用；2:使用中；3:禁用;NOT NULL" json:"status"`
	UsedTimes  int       `gorm:"column:used_times;type:int(11);default:0;comment:使用次数;NOT NULL" json:"used_times"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
	LastUsedAt time.Time `gorm:"column:last_used_at;type:datetime;comment:最后使用时间" json:"last_used_at"`
}

// TableName 指定表名
func (m *AddressPool) TableName() string {
	return "wallet_address_pool"
}

// 状态常量
const (
	AddressStatusAvailable = 1 // 可用
	AddressStatusInUse     = 2 // 使用中
	AddressStatusDisabled  = 3 // 禁用
)

// 网络类型常量
const (
	NetworkTRON = 1
	NetworkBSC  = 2
)

// BeforeCreate 创建前的钩子
func (m *AddressPool) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加地址验证逻辑
	return nil
}

// IsAvailable 检查地址是否可用
func (m *AddressPool) IsAvailable() bool {
	return m.Status == AddressStatusAvailable
}

// IsInUse 检查地址是否在使用中
func (m *AddressPool) IsInUse() bool {
	return m.Status == AddressStatusInUse
}

// MarkAsUsed 标记为使用中
func (m *AddressPool) MarkAsUsed(tx *gorm.DB) error {
	now := time.Now()
	return tx.Model(m).Updates(map[string]interface{}{
		"status":       AddressStatusInUse,
		"used_times":   gorm.Expr("used_times + 1"),
		"last_used_at": &now,
	}).Error
}

// MarkAsAvailable 标记为可用
func (m *AddressPool) MarkAsAvailable(tx *gorm.DB) error {
	return tx.Model(m).Update("status", AddressStatusAvailable).Error
}

// MarkAsDisabled 标记为禁用
func (m *AddressPool) MarkAsDisabled(tx *gorm.DB) error {
	return tx.Model(m).Update("status", AddressStatusDisabled).Error
}
