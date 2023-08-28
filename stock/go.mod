module stock

go 1.20

require (
	common v1.0.0
	github.com/bytedance/sonic v1.9.2
	github.com/gin-gonic/gin v1.9.1
	github.com/redis/go-redis/v9 v9.0.5
	github.com/spf13/viper v1.16.0
	github.com/streadway/amqp v1.1.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.2
	github.com/go-redsync/redsync/v4 v4.8.1
)

require (
	github.com/go-redsync/redsync/v4 v4.8.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
)

replace common => ../commonModule
