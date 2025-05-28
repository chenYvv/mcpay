package models

import (
	"time"

	"gorm.io/gorm"
)

// 商户
type Merchant struct {
	Id        uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(50);default:0;comment:名称;NOT NULL" json:"name"`
	AppId     string    `gorm:"column:app_id;type:varchar(255);comment:地址;NOT NULL" json:"app_id"`
	AppSecret string    `gorm:"column:app_secret;type:varchar(50);comment:剩余额度;NOT NULL" json:"app_secret"`
	Status    int       `gorm:"column:status;type:tinyint(1);default:1;comment:1:正常；2:禁用;NOT NULL" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (m *Merchant) TableName() string {
	return "merchant"
}

// 状态常量
const (
	MerchantStatusAvailable = 1 // 可用
	MerchantStatusDisabled  = 2 // 禁用
)

// IsAvailable 检查是否可用
func (m *Merchant) IsAvailable() bool {
	return m.Status == MerchantStatusAvailable
}

// MarkAsAvailable 标记为可用
func (m *Merchant) MarkAsAvailable(tx *gorm.DB) error {
	return tx.Model(m).Update("status", MerchantStatusAvailable).Error
}

// MarkAsDisabled 标记为禁用
func (m *Merchant) MarkAsDisabled(tx *gorm.DB) error {
	return tx.Model(m).Update("status", MerchantStatusDisabled).Error
}
