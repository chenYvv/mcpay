package bootstrap

import (
	"mcpay/pkg/config"
)

func IniConfig() {
	//config.InitConfig("./cmd/api/config.yml")
	config.InitConfig("api")
}
