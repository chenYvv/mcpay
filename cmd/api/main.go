package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mcpay/internal/api/bootstrap"
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

	addr := fmt.Sprintf(":%s", viper.GetString("server.port"))
	// Start the server
	r.Run(addr)
}
