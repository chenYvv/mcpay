package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mcpay/pkg/database"
	"mcpay/pkg/helpers"
	"mcpay/pkg/logger"
	"net/http"
	"strconv"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

const ResponseCodeSuccess = 1
const ResponseCodeFail = 2

type ControllerInterface interface {
	Init(c *gin.Context)
}

type BaseController struct {
	Gin *gin.Context
}

func (ctr *BaseController) Init(c *gin.Context) {
	ctr.Gin = c
}

func (ctr *BaseController) IsGet() bool {
	return ctr.Gin.Request.Method == "GET"
}

func (ctr *BaseController) IsPost() bool {
	return ctr.Gin.Request.Method == "POST"
}

func (ctr *BaseController) IsAjax() bool {
	return ctr.Gin.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

// JSON 返回 JSON 数据
func (ctr *BaseController) JSON(status int, message string, data interface{}) {
	ctr.Gin.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
	return
}

// ResponseSuc 返回 JSON 数据
func (ctr *BaseController) ResponseSuc(message string, data interface{}) {
	response := Response{
		Code: ResponseCodeSuccess,
		Data: "",
		Msg:  message,
	}
	if data != nil {
		response.Data = string(helpers.Josn(data))
	}
	logger.Info("http 返回", zap.String("原始数据", string(helpers.Josn(response))))
	ctr.Gin.JSON(http.StatusOK, response)
	return
}

// ResponseErr 返回 JSON 数据
func (ctr *BaseController) ResponseErr(message string) {
	response := Response{
		Code: ResponseCodeFail,
		Data: "",
		Msg:  message,
	}
	logger.Info("http 返回", zap.String("原始数据", string(helpers.Josn(response))))
	ctr.Gin.JSON(http.StatusOK, response)
	return
}

func (ctr *BaseController) GetQueryInt(key string, defaultValue int) int {
	valueS, have := ctr.Gin.GetQuery(key)
	if !have {
		return defaultValue
	}
	valueI, err := strconv.Atoi(valueS)
	if err != nil {
		return defaultValue
	}
	return valueI
}

func (ctr *BaseController) GetQueryPages() (offset int, pageSize int) {
	page := ctr.GetQueryInt("page", 1)
	pageSize = ctr.GetQueryInt("page_size", 100)
	offset = (page - 1) * pageSize
	return offset, pageSize
}

func (ctr *BaseController) Page(DB *gorm.DB, out interface{}) {
	offset, limit := ctr.GetQueryPages()
	var count int64
	var data []byte
	if err := DB.Offset(offset).Limit(limit).Find(out).Error; err == nil {
		data, _ = json.Marshal(out)
		database.DB.Table("(?) as tb", DB.Find(out).Offset(-1).Limit(-1)).Count(&count)
	}
	_ = json.Unmarshal(data, out)
	ctr.JSON(200, "success", gin.H{
		"data":        out,
		"total_count": count,
	})
}

func (ctr *BaseController) Sign(params map[string]interface{}) (signStr string) {
	key := "dew78943uw2s667s"
	paramsStr := helpers.QueryAsciiSortNoEmptyNoSign(params) + "&key=" + key
	signStr = helpers.Md5(paramsStr)
	logger.Info("请求", zap.String("字符串", paramsStr), zap.String("MD5", signStr))
	return
}

func (ctr *BaseController) GetClientIP() string {
	reqIP := ctr.Gin.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}
