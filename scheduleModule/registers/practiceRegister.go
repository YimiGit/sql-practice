package registers

import (
	"github.com/xxl-job/xxl-job-executor-go"
	"schedule/controller"
)

// allPracticeRegister 注册practice所有任务
func practiceRegister(exec xxl.Executor) {
	exec.RegTask("SqlPassTotal", controller.SqlPassTotal)
	exec.RegTask("CleanJogLog", controller.CleanJogLog)
	exec.RegTask("UserLoginMonth", controller.UserLoginMonth)
}
