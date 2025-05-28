package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"mcpay/pkg/helpers"
	"mcpay/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type OrderController struct {
	BaseController
}

func (ctr *OrderController) Test() {

	params := map[string]interface{}{
		"network": "network.NetworkId",
		"tx_hash": "record.Hash",
		"amount":  "record.Num",
		"uid":     "record.UserId",
	}

	paramsJSON, _ := json.Marshal(params)

	//logger.Info("测试 1111111111111111111")

	gameCallbackUrl := viper.GetString("G_MConfig.game_callback_url")
	res, errPost := helpers.HttpPost(gameCallbackUrl, nil, paramsJSON)

	if errPost != nil {
		logger.Info("测试 2222222222222222222", zap.String("err", errPost.Error()))
	}

	logger.Info("测试 3333333333333333333", zap.String("err", string(res)))

	return
}

// 创建订单
func (ctr *OrderController) Create() {
	reqBody, _ := ctr.Gin.GetRawData()
	logger.Info("创建订单", slog.String("请求参数", string(reqBody)))
	ctr.Gin.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	request := new(struct {
		AppId       string `json:"app_id"`
		OrderId     string `json:"order_id"`
		Network     string `json:"network"`
		Amount      string `json:"amount"`
		CallbackUrl string `json:"callback_url"`
		Sign        string `json:"sign"`
	})

	if err := ctr.Gin.ShouldBind(request); err != nil {
		ctr.JSON(-1, "params err", nil)
		return
	}

	logger.Info("创建订单", slog.Any("请求参数", helpers.Josn(request)))

	ctr.ResponseSuc("", nil)
	return

}

//func (ctr *OrderController) Wallet() {
//	reqBody, _ := ctr.Gin.GetRawData()
//	logger.Info("生成钱包", zap.String("数据", string(reqBody)))
//	ctr.Gin.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
//
//	request := new(struct {
//		UserId int    `json:"user_id"`
//		Sign   string `json:"sign"`
//	})
//
//	if err := ctr.Gin.ShouldBind(request); err != nil {
//		ctr.JSON(-1, "params err", nil)
//		return
//	}
//
//	mapReq, _ := helpers.JsonToMap(string(reqBody))
//
//	if ctr.Sign(mapReq) != request.Sign {
//		ctr.JSON(-1, "sign err", nil)
//		return
//	}
//
//	if request.UserId == 0 {
//		ctr.JSON(-1, "user err", nil)
//		return
//	}
//
//	walletBsc := models.GetUserWalletByNetworkId(request.UserId, models.CoinNetworkBSC)
//	walletTrc := models.GetUserWalletByNetworkId(request.UserId, models.CoinNetworkTRON)
//
//	networkBsc := models.GetNetworkById(models.CoinNetworkBSC)
//	networkTrc := models.GetNetworkById(models.CoinNetworkTRON)
//
//	ctr.JSON(0, "ok", gin.H{
//		"user_id":      request.UserId,
//		"bsc_address":  walletBsc.Address,
//		"bsc_name":     networkBsc.Chain,
//		"tron_address": walletTrc.Address,
//		"tron_name":    networkTrc.Chain,
//	})
//
//	return
//}
