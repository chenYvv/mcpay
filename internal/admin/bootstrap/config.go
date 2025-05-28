package bootstrap

import "mcpay/pkg/config"

func IniConfig() {
	config.InitConfig("./config.yml")
}
