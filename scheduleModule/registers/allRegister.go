package registers

import "github.com/xxl-job/xxl-job-executor-go"

// AllRegister 注册所有任务
func AllRegister(exec xxl.Executor) {
	practiceRegister(exec)
}
