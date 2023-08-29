package router

import (
	"github.com/gin-gonic/gin"
	"stock/controller"
)

func stockInit(engine *gin.Engine) {

	engine.POST("/stock/priceOrder", controller.PriceOrder)

	engine.POST("/stock/cancelOrder", controller.CancelOrder)
}
