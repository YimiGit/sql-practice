package controller

import (
	"common/results"
	"github.com/gin-gonic/gin"
	"stock/service"
)

// PriceOrder 市价下单/限价下单
func PriceOrder(ctx *gin.Context) {
	results.ResultHandle(ctx, service.PriceOrder(ctx))
}

// CancelOrder 撤单
func CancelOrder(ctx *gin.Context) {
	results.ResultHandle(ctx, service.CancelOrder(ctx))
}

// KLineData 获取K线数据
func KLineData(ctx *gin.Context) {
	results.ResultHandle(ctx, service.KLineData(ctx))
}
