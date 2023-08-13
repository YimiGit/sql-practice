package grpcService

import (
	"common/config"
	"common/proto/practiceProto"
	"common/proto/userProto"
	"google.golang.org/grpc"
	"log"
	"net"
)

// GRPCInit 初始化grpc服务端
func GRPCInit() {
	//创建grpc服务端
	server := grpc.NewServer()
	userProto.RegisterUserServiceServer(server, &UserService{})
	practiceProto.RegisterPaperServiceServer(server, &PracticeService{})

	listener, runListenerErr := net.Listen("tcp", config.ViperGetString("grpc.host"))
	if runListenerErr != nil {
		log.Println("grpc监听器-初始化失败", runListenerErr)
		return
	} else {
		log.Println("grpc监听器-初始化成功")
	}
	go func() {
		runServerErr := server.Serve(listener)
		if runServerErr != nil {
			log.Println("grpc服务端-初始化失败", runServerErr)
		}
	}()
}
