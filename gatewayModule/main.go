package main

import (
	"common/config"
	"gateway/grpc/userGrpc"
	"gateway/router"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

// main gateway服务
func main() {

	//初始化所有配置文件
	config.GatewayEnvInit()

	engine := gin.Default()

	//gRPC客户端初始化
	userGrpc.GRPCClientInit()

	//获取etcd 客户端
	client, err2 := clientv3.New(clientv3.Config{Endpoints: []string{config.EtcdHostPort}})
	if err2 != nil {
		log.Println("etcd连接失败", err2)
	}

	router.AllRouterInit(engine, client)

	_ = engine.Run(":19600")
}
