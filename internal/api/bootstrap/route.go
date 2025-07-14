// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	"mcpay/internal/api/controllers"
	"mcpay/internal/api/middlewares"
	"reflect"

	"github.com/gin-gonic/gin"
)

// InitRoutes 注册路由
func InitRoutes(r *gin.Engine) {
	r.Use(middlewares.Recovery(), middlewares.Cors())

	r.POST("/test/createaddress", routeHandle(&controllers.TestController{}, "CreateAddress"))

	r.POST("/order/create", routeHandle(&controllers.OrderController{}, "Create"))

	//r.GET("/test", routeHandle(&controllers.PayController{}, "Test"))
	//
	//r.GET("/pay", routeHandle(&controllers.PayController{}, "Pay"))
	//r.GET("/callback", routeHandle(&controllers.PayController{}, "PayRes"))
	//r.POST("/bsc", routeHandle(&controllers.PayController{}, "Wallet"))
	//r.POST("/uploadImg", routeHandle(&controllers.UploadController{}, "UploadImg"))
	//r.POST("/withdraw", routeHandle(&controllers.WalletController{}, "Withdraw"))
	//r.POST("/api", routeHandle(&controllers.PayController{}, "Api"))
	//
	//r.POST("/XwinPayRechargeNotifyUrl", routeHandle(&controllers.PayController{}, "XwinPayRechargeNotifyUrl"))
	//r.POST("/XwinPayConvertNotifyUrl", routeHandle(&controllers.PayController{}, "XwinPayConvertNotifyUrl"))
	//r.POST("/XwinPHPPayRechargeNotifyUrl", routeHandle(&controllers.PayController{}, "XwinPHPPayRechargeNotifyUrl"))
	//r.POST("/XwinPHPPayConvertNotifyUrl", routeHandle(&controllers.PayController{}, "XwinPHPPayConvertNotifyUrl"))

}

func routeHandle(ctr controllers.ControllerInterface, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctr.Init(c)
		value := reflect.ValueOf(ctr)
		methodValue := value.MethodByName(action)
		methodValue.Call(nil)
	}
}
