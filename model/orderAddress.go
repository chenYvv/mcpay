package models

// 订单表

type OrderAddress struct {
	OrderId   int `gorm:"column:order_id;type:int(11);primary_key" json:"order_id"`
	AddressId int `gorm:"column:address_id;type:int(11);comment:地址ID;NOT NULL" json:"address_id"`
	BlockNum  int `gorm:"column:block_num;type:int(11);default:0;comment:区块高度;NOT NULL" json:"block_num"`
}

// TableName 指定表名
func (m *OrderAddress) TableName() string {
	return "order_address"
}
