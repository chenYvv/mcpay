package models

import (
	"fmt"
	"log"
	"mcpay/pkg/database"
	"sync"
)

type Configs struct {
	Id     uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	K      string `gorm:"column:k;type:varchar(255)" json:"k"`
	V      string `gorm:"column:v;type:varchar(255)" json:"v"`
	Beizhu string `gorm:"column:beizhu;type:varchar(255)" json:"beizhu"`
}

// TableName 指定表名
func (m *Configs) TableName() string {
	return "configs"
}

func GetMnemonic() string {
	var data *Configs
	err := database.DB.Where("k = ?", "mnemonic").Find(&data).Error
	if err != nil {
		return ""
	}
	return data.V
}

// GlobalConfig 用于封装系统全局配置项
type GlobalConfig struct {
	Mnemonic        string // 助记词
	OrderCodeSalt   string // 订单号盐
	EtherscanApikey string // EtherscanApikey
}

var (
	global     *GlobalConfig
	initOnce   sync.Once
	configLock sync.RWMutex
)

// Global 返回当前全局配置对象（只读）
func Global() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return global
}

// InitConfig 初始化配置（只执行一次）
func InitGlobalConfig() {
	initOnce.Do(func() {
		// 数据库 config 配置
		var configs []Configs
		err := database.DB.Find(&configs).Error
		if err != nil {
			panic(fmt.Sprintf("[config] Failed to load config from DB: %v", err))
		}

		c := &GlobalConfig{}

		for _, cfg := range configs {
			switch cfg.K {
			case "mnemonic":
				c.Mnemonic = cfg.V
			case "order_code_salt":
				c.OrderCodeSalt = cfg.V
			case "etherscan_apikey":
				c.EtherscanApikey = cfg.V
			}
		}

		configLock.Lock()
		global = c
		configLock.Unlock()
		log.Println("[config] Global config initialized")
	})
}
