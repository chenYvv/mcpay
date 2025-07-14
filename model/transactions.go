package models

import (
	"math/rand"
	"mcpay/pkg/bsc"
	timePkg "mcpay/pkg/time"
	"mcpay/pkg/tron"
	"time"
)

type Transaction struct {
	Id            uint      `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	OrderId       int       `gorm:"column:order_id;type:int(11);comment:订单ID" json:"order_id"`
	AddressId     int       `gorm:"column:address_id;type:int(11);comment:地址ID" json:"address_id"`
	Network       int       `gorm:"column:network;type:int(11);default:0;comment:网络：1:波场；2:币安;NOT NULL" json:"network"`
	TxHash        string    `gorm:"column:tx_hash;type:varchar(255);comment:交易hash" json:"tx_hash"`
	FromAddress   string    `gorm:"column:from_address;type:varchar(255);comment:转账地址" json:"from_address"`
	Amount        float64   `gorm:"column:amount;type:decimal(10,2);comment:交易金额" json:"amount"`
	BlockTime     time.Time `gorm:"column:block_time;type:datetime;comment:交易时间" json:"block_time"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	CallbackState int       `gorm:"column:callback_state;type:int(11);default:0;comment:通知状态：0:失败；1:成功；;NOT NULL" json:"callback_state"`
	CallbackErr   string    `gorm:"column:callback_err;type:varchar(500);comment:通知报错信息" json:"callback_err"`
	CallbackDate  time.Time `gorm:"column:callback_date;type:datetime;comment:最后通知时间" json:"callback_date"`
	CallbackTimes int       `gorm:"column:callback_times;type:int(11);default:0;comment:通知次数;NOT NULL" json:"callback_times"`
}

// TableName 指定表名
func (m *Transaction) TableName() string {
	return "transactions"
}

type TransactionItem struct {
	Hash        string    `json:"hash"`
	FromAddress string    `json:"from_address"`
	Amount      float64   `json:"amount"`
	BlockTime   time.Time `json:"block_time"`
}

type tronTransactionApiFunc func(addr string, startTimestamp int64) ([]TransactionItem, error)

// 波场API
var tronTransactionFuncs = []tronTransactionApiFunc{
	trongridApi,
}

// 随机获取波场API
func GetTronTransaction(address string, startTimestamp int64) ([]TransactionItem, error) {
	rand.Seed(time.Now().UnixNano())
	apiFunc := tronTransactionFuncs[rand.Intn(len(tronTransactionFuncs))]
	list, err := apiFunc(address, startTimestamp)
	return list, err
}

// 波场API - trongrid
func trongridApi(addr string, startTimestamp int64) ([]TransactionItem, error) {
	transactions := []TransactionItem{}
	list, err := tron.GetUSDTTransactions(addr, startTimestamp)
	if err != nil {
		return transactions, err
	}
	if len(list) > 0 {
		for _, item := range list {
			transactions = append(transactions, TransactionItem{
				Hash:        item.TransactionId,
				FromAddress: item.From,
				Amount:      tron.WeiStrToNum(item.Value),
				BlockTime:   time.UnixMilli(item.BlockTimestamp),
			})
		}
	}
	return transactions, nil
}

type bscTransactionApiFunc func(addr string, blockNum int64) ([]TransactionItem, error)

// 币安API
var bscTransactionFuncs = []bscTransactionApiFunc{
	etherscanApi,
}

// 随机获取波场API
func GetBscTransaction(address string, blockNum int64) ([]TransactionItem, error) {
	rand.Seed(time.Now().UnixNano())
	apiFunc := bscTransactionFuncs[rand.Intn(len(bscTransactionFuncs))]
	list, err := apiFunc(address, blockNum)
	return list, err
}

// 币安API - etherscan
func etherscanApi(addr string, blockNum int64) ([]TransactionItem, error) {
	transactions := []TransactionItem{}
	list, err := bsc.GetUSDTTransactionsByEtherscan(Global().EtherscanApikey, addr, blockNum)
	if err != nil {
		return transactions, err
	}
	if len(list) > 0 {
		for _, item := range list {
			transactions = append(transactions, TransactionItem{
				Hash:        item.Hash,
				FromAddress: item.From,
				Amount:      bsc.WeiStrToNum(item.Value),
				BlockTime:   timePkg.StrToTime(item.TimeStamp),
			})
		}
	}
	return transactions, nil
}
