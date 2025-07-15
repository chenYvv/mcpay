package models

type Currency struct {
	Id        int     `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Name      string  `gorm:"column:name;type:varchar(255);comment:法币;NOT NULL" json:"name"`
	RateUSDT  float64 `gorm:"column:rate_USDT;type:decimal(10,6);comment:汇率：1USDT=0.999USD;NOT NULL" json:"rate_usdt"`
	UpdatedAt int64   `gorm:"column:updated_at;type:bigint(20);comment:更新时间戳;NOT NULL" json:"updated_at"`
}

// TableName 指定表名
func (m *Currency) TableName() string {
	return "currency"
}
