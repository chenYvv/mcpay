package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	//"mcpay/model"
	"mcpay/pkg/database"
)

func InitDb() {
	var err error
	database.DB, err = database.MysqlDB(
		viper.GetString("database.mysql.DSN"),
		viper.GetInt("database.mysql.max_open_connections"),
		viper.GetInt("database.mysql.max_idle_connections"),
		viper.GetInt("database.mysql.max_life_seconds"))
	if err != nil {
		fmt.Println(err.Error())
	}

	//uri := fmt.Sprintf("%s:%s@%s", viper.GetString("database.mongo.user"), viper.GetString("database.mongo.pwd"), viper.GetString("database.mongo.host"))

	//database.NewConnectionV1(uri)

	//models.IniGlobalConfigs()

}
