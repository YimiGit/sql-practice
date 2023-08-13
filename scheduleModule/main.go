package main

import (
	"common/config"
	"github.com/gin-gonic/gin"
	xxl_job_executor_gin "github.com/gin-middleware/xxl-job-executor"
	"github.com/xxl-job/xxl-job-executor-go"
	"schedule/registers"
)

func main() {

	config.ScheduleEnvInit()
	engine := gin.Default()
	//初始化执行器
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1:19889/xxl-job-admin"),
		xxl.AccessToken("default_token"),  //请求令牌(默认为空)
		xxl.ExecutorIp("127.0.0.1"),       //可自动获取
		xxl.ExecutorPort("19888"),         //默认9999（此处要与gin服务启动port必需一至）
		xxl.RegistryKey("sql-pass-total"), //执行器名称
	)
	exec.Init()
	//defer exec.Stop()

	//添加到gin路由
	xxl_job_executor_gin.XxlJobMux(engine, exec)

	//注册任务handler
	//exec.RegTask("task.test", task.Test)
	//exec.RegTask("task.test2", task.Test2)
	//exec.RegTask("task.panic", task.Panic)
	registers.AllRegister(exec)

	_ = engine.Run(config.ViperGetString("service.port"))
}
