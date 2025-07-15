package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mcpay/internal/common/payment"
	models "mcpay/model"
	"mcpay/pkg/constants"
	"mcpay/pkg/database"
	"mcpay/pkg/helpers"
	"mcpay/pkg/logger"
	"time"

	"gorm.io/gorm"

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
	//reqBody, _ := ctr.Gin.GetRawData()
	//logger.Info("创建订单", string(reqBody))
	//ctr.Gin.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	request := new(struct {
		AppId       string `json:"app_id"`
		Uid         int    `json:"uid"`
		OrderId     string `json:"order_id"`
		Currency    string `json:"currency"`
		Amount      string `json:"amount"`
		CallbackUrl string `json:"callback_url"`
		RedirectUrl string `json:"redirect_url"`
		Sign        string `json:"sign"`
	})

	if err := ctr.Gin.ShouldBind(request); err != nil {
		ctr.ResponseCode(constants.ParamsError, nil)
		return
	}

	logger.Info("创建订单", slog.String("params", helpers.Json2Struct(request)))

	// 获取 app
	app, err := models.GetAppById(request.AppId)
	if err != nil {
		ctr.ResponseCode(constants.ParamsError, nil)
		return
	}

	payChannel := app.GetPayChannel()
	if len(payChannel) == 0 {
		ctr.ResponseCode(constants.ChannelError, nil)
		return
	}

	// 转换网络字符串为常量
	// network := models.GetNetworkByString(request.Network)

	// 获取空闲地址钱包
	// models.GetAvailableAddressByNetwork(network)

	timeNow := time.Now().Unix()

	//response := payment.CollectResult{}

	err = database.DB.Transaction(func(tx *gorm.DB) error {

		// 创建订单
		order := models.Order{
			OrderId:         helpers.GenerateOrderNo(""),
			Amount:          helpers.StringToFloat64(request.Amount),
			Amount:          helpers.StringToFloat64(request.Amount),
			Uid:             request.Uid,
			AppId:           helpers.Str2Int(request.AppId),
			CallbackUrl:     request.CallbackUrl,
			RedirectUrl:     request.RedirectUrl,
			MerchantOrderId: request.OrderId,
			CreatedAt:       timeNow,
			UpdatedAt:       timeNow,
		}

		err = tx.Create(&order).Error
		if err != nil {
			return err
		}

		// 唤起支付
		collectResult := payment.CollectResult{}
		var payRes error
		var pay payment.Payment
		for _, channel := range payChannel {
			pay, payRes = payment.GetProvider(channel)
			if payRes != nil {
				continue
			}
			collectResult, payRes = pay.Collect(tx, order, nil)
			if payRes != nil {
				continue
			}
			break
		}

		if payRes != nil {
			return payRes
		}

		fmt.Println(payChannel)

		// 更新订单信息
		err = tx.Model(models.Order{}).Where("order_id = ?", order.OrderId).Updates(
			map[string]interface{}{
				"third_order_id": collectResult.ThirdOrderId,
			},
		).Error

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		ctr.ResponseCode(constants.SystemlError, nil)
		return
	}

	ctr.ResponseCode(constants.Success, nil)
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
