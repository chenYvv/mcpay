package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mcpay/internal/api/bootstrap"
	_ "mcpay/internal/common/payment/bsc"
	_ "mcpay/internal/common/payment/tron"
	models "mcpay/model"
	"mcpay/pkg/chain/bsc"
	"mcpay/pkg/chain/tron"
	"mcpay/pkg/config"
	"mcpay/pkg/idcode"
	"mcpay/pkg/logger"
)

func main() {
	// Initialize the configuration
	// config.InitConfig("config.yml")
	// Initialize the database
	// db.InitDB()
	// Initialize the router
	r := gin.Default()
	bootstrap.IniConfig()
	bootstrap.InitLogger()
	bootstrap.InitDb()
	bootstrap.InitRoutes(r)

	err := idcode.Init(models.Global().OrderCodeSalt, 8)
	if err != nil {
		panic(err)
	}

	// 波场
	err = tron.InitTronClient(viper.GetBool("APP_DEBUG"))
	if err != nil {
		panic(err)
	} else {
		logger.Info("InitTronClient SUCCESS")
	}

	// 币安
	err = bsc.InitBSCClient(viper.GetBool("APP_DEBUG"))
	if err != nil {
		panic(err)
	} else {
		logger.Info("InitBSCClient SUCCESS")
	}

	// 检测超时订单
	models.UpdateTimeoutOrders()

	// 定时任务
	bootstrap.InitCrontab()

	addr := fmt.Sprintf(":%s", config.GetString("APP_PORT", 8000))
	// Start the server
	r.Run(addr)

}
