package grpcService

import (
	"common/config"
	model "common/model/adminModel"
	"common/proto/userProto"

	"common/utils"
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"log"
	"strconv"
	"time"
	"user/rabbitmq"
)

// UserService 用户rpc服务
type UserService struct {
	userProto.UnimplementedUserServiceServer
}

func (*UserService) Login(c context.Context, req *userProto.LoginRequest) (*userProto.LoginResponse, error) {

	//客户端额外携带的元数据
	//incomingContext, _ := metadata.FromIncomingContext(c)
	//value := incomingContext.Get("token")
	//log.Println(value)

	//登录逻辑
	//1.校验数据
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("数据类型不合法")
	}

	//2.校验用户名密码
	user := model.User{Username: req.Username}
	if result := config.DB.Model(&user).Where("username = ?", user.Username).First(&user); result.Error != nil {
		return nil, errors.New("用户名不存在")
	}

	if req.Password != user.Password {
		return nil, errors.New("密码错误")
	}

	//3.生成token
	token, err := utils.CreateToken(utils.StringToInt(user.ID))
	if err != nil {
		return nil, errors.New("token生成失败")
	}
	err = config.RedisClient.Set(c, user.ID, token, time.Minute*30).Err()
	if err != nil {
		return nil, errors.New("token存入redis失败")
	}

	//4.gRPC返回数据
	res := userProto.LoginResponse{ID: user.ID, Username: user.Username, Password: user.Password, Token: token}

	//登录统计
	marshal, _ := sonic.Marshal(user.ID)
	channel, err2 := rabbitmq.UserLoginTotalChannelPool.GetChannel()
	if err2 != nil {
		log.Println("userLoginConsumer获取channel失败", err2)
		return nil, errors.New("登录失败")
	}

	utils.PushMessage(marshal, channel, "userExchange", "userLoginKey")
	rabbitmq.UserLoginTotalChannelPool.ReleaseChannel(channel)
	return &res, nil
}

// Register 注册 gRPC服务 供gateway调用
func (*UserService) Register(c context.Context, req *userProto.RegisterRequest) (*userProto.RegisterResponse, error) {

	//1.校验数据
	if req.Username == "" || req.Password == "" || req.PasswordConfirm == "" {
		return nil, errors.New("参数错误")
	}
	if req.Password != req.PasswordConfirm {
		return nil, errors.New("两次密码不一致")
	}

	//2.校验用户名是否存在
	user := model.User{Username: req.Username}
	err := config.DB.Model(&user).Where("username = ?", user.Username).First(&user)
	if err.Error == nil {
		return nil, errors.New("用户名已存在")
	}
	user.Password = req.Password

	//3.注册
	id, _ := utils.DistributedID(1)
	user.ID = strconv.Itoa(int(id))
	err = config.DB.Create(&user)
	if err.Error != nil {
		return nil, errors.New("注册失败")
	}

	//4.gRPC返回数据
	res := userProto.RegisterResponse{ID: user.ID, Username: user.Username, Password: user.Password}
	return &res, nil
}
