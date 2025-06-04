package main

import (
	"fmt"
	"mcpay/internal/api/bootstrap"
	"mcpay/pkg/bsc"
	"mcpay/pkg/tron"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	tron.InitTronClient(true)
	bsc.InitBSCClient(true)

	addr := fmt.Sprintf(":%s", viper.GetString("server.port"))
	// Start the server
	r.Run(addr)
}
