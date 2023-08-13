package userGrpc

import (
	"common/config"
	"common/proto/userProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// UserServiceClient 用户服务客户端
var UserServiceClient userProto.UserServiceClient

// GRPCClientInit grpc客户端初始化
func GRPCClientInit() {
	//创建grpc连接
	connection, err := grpc.Dial(config.ViperGetString("grpc.host"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("rpc连接失败", err)
		return
	}

	//创建grpc客户端
	userServiceClient := userProto.NewUserServiceClient(connection)

	UserServiceClient = userServiceClient
}
