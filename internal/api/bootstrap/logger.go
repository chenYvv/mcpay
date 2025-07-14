package bootstrap

import (
	"mcpay/pkg/logger"

	"github.com/spf13/viper"
)

// InitLogger 初始化 Logger
func InitLogger() {
	// logger.InitLogger(
	// 	fmt.Sprintf("logs/%s/", viper.GetString("server.name")),
	// 	5,
	// 	5,
	// 	30,
	// 	false,
	// 	"daily",
	// 	viper.GetString("log.level"),
	// )
	// 使用新的 slog 日志系统
	logger.InitSlog(
		viper.GetString("APP_NAME"),
		viper.GetBool("APP_DEBUG"),
	)
}
