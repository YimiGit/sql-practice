package main

import (
	"common/config"
	"context"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"practice/grpc/practiceGrpc"
	"practice/rabbitmq"
	"practice/routers"
)

func main() {

	//初始化配置
	config.PracticeEnvInit()

	engine := gin.Default()

	//初始化rabbitmq
	rabbitmq.AllInit()

	//初始化路由
	routers.AllRoutersInit(engine)

	//初始化grpc客户端
	practiceGrpc.GRPCClientInit()

	//注册到etcd
	client, err2 := clientv3.New(clientv3.Config{Endpoints: []string{config.EtcdHostPort}})
	if err2 != nil {
		log.Println("etcd连接失败", err2)
	}
	_, err2 = client.Put(context.Background(), config.ServerName, config.ServerHost+config.ServerPort)
	if err2 != nil {
		log.Println("etcd注册失败", err2)
	}

	_ = engine.Run(config.ServerPort)
}
