package models

import (
	"math/rand"
	"mcpay/pkg/chain/bsc"
	"mcpay/pkg/chain/tron"
	"mcpay/pkg/helpers"
	"time"
)

type Transaction struct {
	Id            uint    `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	OrderId       int     `gorm:"column:order_id;type:int(11);comment:订单ID" json:"order_id"`
	AddressId     int     `gorm:"column:address_id;type:int(11);comment:地址ID" json:"address_id"`
	Network       int     `gorm:"column:network;type:int(11);default:0;comment:网络：1:波场；2:币安;NOT NULL" json:"network"`
	TxHash        string  `gorm:"column:tx_hash;type:varchar(255);comment:交易hash" json:"tx_hash"`
	FromAddress   string  `gorm:"column:from_address;type:varchar(255);comment:转账地址" json:"from_address"`
	Amount        float64 `gorm:"column:amount;type:decimal(10,2);comment:交易金额" json:"amount"`
	BlockTime     int64   `gorm:"column:block_time;type:bigint(20);default:0;comment:交易时间;NOT NULL" json:"block_time"`
	CreatedAt     int64   `gorm:"column:created_at;type:bigint(20);default:0;comment:创建时间;NOT NULL" json:"created_at"`
	CallbackState int     `gorm:"column:callback_state;type:int(11);default:0;comment:通知状态：0:失败；1:成功；;NOT NULL" json:"callback_state"`
	CallbackErr   string  `gorm:"column:callback_err;type:varchar(500);comment:通知报错信息" json:"callback_err"`
	CallbackDate  int64   `gorm:"column:callback_date;type:bigint(20);default:0;comment:最后通知时间;NOT NULL" json:"callback_date"`
	CallbackTimes int     `gorm:"column:callback_times;type:int(11);default:0;comment:通知次数;NOT NULL" json:"callback_times"`
}

// TableName 指定表名
func (m *Transaction) TableName() string {
	return "transactions"
}

type TransactionItem struct {
	Hash        string  `json:"hash"`
	FromAddress string  `json:"from_address"`
	Amount      float64 `json:"amount"`
	BlockTime   int64   `json:"block_time"`
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
				BlockTime:   item.BlockTimestamp,
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
				BlockTime:   helpers.StrToInt64(item.TimeStamp),
			})
		}
	}
	return transactions, nil
}
