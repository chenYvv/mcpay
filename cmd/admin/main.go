package main

import (
	"github.com/gin-gonic/gin"
	"mcpay/internal/admin/bootstrap"
)

func main() {
	// Initialize the configuration
	// config.InitConfig("config.yml")
	// Initialize the database
	// db.InitDB()
	// Initialize the router
	router := gin.Default()
	bootstrap.IniConfig()

	// Start the server
	router.Run(":8080")
}
