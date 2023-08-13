package controller

import (
	"context"
	"github.com/xxl-job/xxl-job-executor-go"
	"log"
	"schedule/service"
)

//trigger_code 调度状态
//handle_code  执行状态 200成功 500失败 0调度失败
//handle_msg   执行结果 return msg 或 panic msg

// SqlPassTotal sql作答正确 数据统计 入口
func SqlPassTotal(c context.Context, param *xxl.RunReq) (msg string) {
	err := service.SqlPassTotal(c, param)
	if err != nil {
		log.Println("sql作答正确 数据统计 失败", err.Error())
		panic(err.Error())
	}
	return "200"
}

// CleanJogLog 清理调度日志(保留三天)
func CleanJogLog(c context.Context, param *xxl.RunReq) (msg string) {
	err := service.CleanJogLog(c, param)
	if err != nil {
		log.Println("清理调度日志失败", err)
		panic(err)
	}
	return "200"
}

// UserLoginMonth 用户登录月统计
func UserLoginMonth(c context.Context, param *xxl.RunReq) (msg string) {
	err := service.UserLoginMonth(c, param)
	if err != nil {
		log.Println("用户登录月统计失败", err)
		panic(err)
	}
	return "200"
}
