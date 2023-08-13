package practiceGrpc

import (
	"common/config"
	"common/proto/practiceProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// PracticeServiceClient grpc客户端
var PracticeServiceClient practiceProto.PaperServiceClient

// GRPCClientInit grpc客户端初始化
func GRPCClientInit() {
	//创建grpc连接
	connection, err := grpc.Dial(config.ViperGetString("grpc.host"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("rpc连接失败", err)
		return
	}

	//创建grpc客户端
	practiceServiceClient := practiceProto.NewPaperServiceClient(connection)

	PracticeServiceClient = practiceServiceClient
}
