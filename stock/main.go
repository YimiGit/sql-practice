package stock

import (
	"github.com/gin-gonic/gin"
	"stock/config"
	"stock/rabbitMq"
	"stock/router"
)

func main() {

	//初始化配置
	config.StockEnvInit()

	engine := gin.Default()

	//初始化路由
	router.AllInit(engine)

	//初始化rabbitMq
	rabbitMq.AllInit()

	_ = engine.Run(config.ServerPort)
}
