package models

import (
	"mcpay/pkg/database"
	"mcpay/pkg/helpers"
	"strings"
	"time"
)

type Apps struct {
	Id         uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(50);default:0;comment:名称;NOT NULL" json:"name"`
	AppId      string    `gorm:"column:app_id;type:varchar(255);NOT NULL" json:"app_id"`
	AppSecret  string    `gorm:"column:app_secret;type:varchar(50);NOT NULL" json:"app_secret"`
	Status     int       `gorm:"column:status;type:tinyint(1);default:1;comment:1:正常；2:禁用;NOT NULL" json:"status"`
	PayChannel string    `gorm:"column:pay_channel;type:varchar(255);comment:支付渠道" json:"pay_channel"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (m *Apps) TableName() string {
	return "apps"
}

func GetAppById(appId string) (*Apps, error) {
	var data *Apps
	err := database.DB.Where("app_id = ?", appId).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *Apps) GetPayChannel() []int {
	var payChannels []int
	if m.PayChannel != "" {
		slicedStr := strings.Split(m.PayChannel, ",")
		for _, v := range slicedStr {
			payChannels = append(payChannels, helpers.Str2Int(v))
		}
	}
	return payChannels

}
