package main

import (
	"common/config"
	"context"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"user/grpcService"
	"user/rabbitmq"
	"user/routers"
)

func main() {

	//初始化配置
	config.AdminEnvInit()

	engine := gin.Default()

	//初始化路由
	routers.AllRoutersInit(engine)

	//初始化rabbitmq
	rabbitmq.AllInit()

	//初始化gRPC
	grpcService.GRPCInit()

	//注册到etcd
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{config.EtcdHostPort}})
	if err != nil {
		log.Println("etcd连接失败", err)
	}
	_, err = client.Put(context.Background(), config.ServerName, config.ServerHost+config.ServerPort)
	if err != nil {
		log.Println("etcd注册失败", err)
	}

	_ = engine.Run(config.ServerPort)
}
