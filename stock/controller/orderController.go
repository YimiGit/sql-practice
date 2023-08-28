package controller

import (
	"common/results"
	"github.com/gin-gonic/gin"
	"stock/service"
)

// StockController 股票模块控制器

// PriceOrder 市价下单/限价下单

func PriceOrder(ctx *gin.Context) {
	results.ResultHandle(ctx, service.PriceOrder(ctx))
}

// CancelOrder 撤单
func CancelOrder(ctx *gin.Context) {
	results.ResultHandle(ctx, service.CancelOrder(ctx))
}
