package rabbitmq

// AllInit 初始化所有队列
func AllInit() {

	// 初始化用户队列
	userQueueInit()

	// 初始化练习队列
	practiceQueueInit()
}
