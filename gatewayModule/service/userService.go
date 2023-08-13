package service

import (
	model "common/model/adminModel"
	"common/proto/userProto"
	"common/results"
	"context"
	"gateway/grpc/userGrpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func Login(c *gin.Context) *results.JsonResult {
	var body model.User

	//参数绑定
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return results.Error("参数错误", nil)
	}

	//构建grpc请求参数
	req := &userProto.LoginRequest{
		Username: body.Username,
		Password: body.Password,
	}

	//grpc携带参数(仅在当前服务)
	//reqCtx := context.WithValue(context.Background(), "token", ctx.GetHeader("token"))

	//grpc携带参数(跨服务（客户端 -> 服务端）)
	//构建context
	reqCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", c.GetHeader("token")))
	//发起grpc请求
	res, err := userGrpc.UserServiceClient.Login(reqCtx, req)

	if err != nil {
		return results.Fail("rpc请求失败", err)
	}

	return results.Success("登录成功", res)
}

func Register(c *gin.Context) *results.JsonResult {
	type Body struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
	}

	//参数绑定
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return results.Error("参数错误", nil)
	}

	//构建grpc请求参数
	req := &userProto.RegisterRequest{
		Username:        body.Username,
		Password:        body.Password,
		PasswordConfirm: body.PasswordConfirm,
	}
	reqCtx := context.WithValue(context.Background(), "token", c.GetHeader("token"))

	//发起grpc请求
	res, err := userGrpc.UserServiceClient.Register(reqCtx, req)
	if err != nil {
		return results.Fail("rpc请求失败", err)
	}
	return results.Success("注册成功", res)
}
