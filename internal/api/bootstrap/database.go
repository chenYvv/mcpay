package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	models "mcpay/model"
	"mcpay/pkg/database"
)

func InitDb() {
	var err error
	database.DB, err = database.MysqlDB(
		viper.GetString("MYSQL_DSN"),
		viper.GetInt("MYSQL_MAX_OPEN_CONNECTIONS"),
		viper.GetInt("MYSQL_MAX_IDLE_CONNECTIONS"),
		viper.GetInt("MYSQL_MAX_LIFE_SECONDS"),
		viper.GetBool("APP_DEBUG"),
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	models.InitGlobalConfig()
}
