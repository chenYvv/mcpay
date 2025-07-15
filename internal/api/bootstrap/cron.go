package bootstrap

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	models "mcpay/model"
	"mcpay/pkg/database"
	"sync"
	"time"
)

// 钱包类型
type WalletType string

// 全局cron调度器
var cronScheduler *cron.Cron

// 初始化定时任务
func InitCrontab() {
	cronScheduler = cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	// 更新费率
	_, err := cronScheduler.AddFunc("@every 1m", models.UpdateCurrencyRateUSDT)
	if err != nil {
		log.Printf("Failed to add UpdateCurrencyRateUSDT cron job: %v", err)
	} else {
		log.Println("UpdateCurrencyRateUSDT monitoring job added")
	}

	// 更新超时订单状态
	_, err = cronScheduler.AddFunc("@every 1m", models.UpdateTimeoutOrders)
	if err != nil {
		log.Printf("Failed to add TRON cron job: %v", err)
	} else {
		log.Println("TRON wallet monitoring job added")
	}

	// 添加TRON监控任务
	_, err = cronScheduler.AddFunc("@every 15s", monitorTronTransactions)
	if err != nil {
		log.Printf("Failed to add TRON cron job: %v", err)
	} else {
		log.Println("TRON wallet monitoring job added")
	}

	// 添加BSC监控任务
	_, err = cronScheduler.AddFunc("@every 15s", monitorBscTransactions)
	if err != nil {
		log.Printf("Failed to add BSC cron job: %v", err)
	} else {
		log.Println("BSC wallet monitoring job added")
	}

	// 启动定时任务
	cronScheduler.Start()
	log.Println("Cron scheduler initialized successfully")
}

// 监控TRON钱包地址交易（批次处理）
func monitorTronTransactions() {
	addressList, _ := models.GetTronInUseAddresses()
	monitorWalletTransactions(addressList, 15, "TRON")
}

// 监控BSC钱包地址交易（批次处理）
func monitorBscTransactions() {
	addressList, _ := models.GetBSCInUseAddresses()
	monitorWalletTransactions(addressList, 10, "BSC")
}

func monitorWalletTransactions(addressList []models.Address, batchSize int, networkName string) {
	if len(addressList) == 0 {
		return
	}

	log.Printf("Starting %s monitoring for %d wallets", networkName, len(addressList))

	for i := 0; i < len(addressList); i += batchSize {
		end := i + batchSize
		if end > len(addressList) {
			end = len(addressList)
		}

		batch := addressList[i:end]
		log.Printf("Processing %s batch %d-%d (%d wallets)", networkName, i+1, end, len(batch))

		var wg sync.WaitGroup
		for _, wallet := range batch {
			wg.Add(1)
			go func(address models.Address) {
				defer wg.Done()
				getTransactions(address)
			}(wallet)
		}
		wg.Wait()
		log.Printf("Completed %s batch %d-%d", networkName, i+1, end)
	}

	log.Printf("Completed %s monitoring for %d wallets", networkName, len(addressList))
}

func getTransactions(address models.Address) {
	var transactions []models.TransactionItem
	if address.IsTron() {
		transactions, _ = models.GetTronTransaction(address.Address, address.LastUsedAt*1000)
	} else if address.IsBsc() {
		orderAddress := address.GetPendingOrderAddress()
		transactions, _ = models.GetBscTransaction(address.Address, int64(orderAddress.BlockNum))
	}

	if len(transactions) > 0 {
		timeNow := time.Now().Unix()
		for _, transaction := range transactions {
			// 检查是否写入过
			history := models.Transaction{}
			database.DB.Where("tx_hash = ? and network = ?", transaction.Hash, address.Network).Find(&history)
			if history.Id == 0 {

				order := address.GetPendingOrder()
				err := database.DB.Transaction(func(tx *gorm.DB) error {
					// 创建订单
					data := models.Transaction{
						OrderId:       order.Id,
						AddressId:     address.Id,
						Network:       address.Network,
						TxHash:        transaction.Hash,
						FromAddress:   transaction.FromAddress,
						Amount:        transaction.Amount,
						BlockTime:     transaction.BlockTime,
						CreatedAt:     timeNow,
						CallbackTimes: 0,
					}

					err := tx.Create(&data).Error
					if err != nil {
						return err
					}

					// 更新金额
					tx.Model(&models.Order{}).
						Where("id = ?", order.Id).
						Updates(map[string]interface{}{
							"amount_true": gorm.Expr("amount_true + ?", transaction.Amount),
						})

					return nil
				})

				// 支付回调
				if err != nil {
				}

				order = address.GetPendingOrder()
				order.Callback()

			}
		}
	}

}

// 停止定时任务
func StopCrontab() {
	if cronScheduler != nil {
		cronScheduler.Stop()
		log.Println("Cron scheduler stopped")
	}
}
