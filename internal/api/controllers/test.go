package controllers

import (
	models "mcpay/model"
)

type TestController struct {
	BaseController
}

func (ctr *TestController) CreateAddress() {
	models.CreateAddress(models.NetworkBsc, 10)
	//models.CreateAddress(models.NetworkTron, 10)
}
