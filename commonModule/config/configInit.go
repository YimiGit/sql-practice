package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

var (
	DB            *gorm.DB
	RedisClient   *redis.Client
	ServerPort    string
	ServerHost    string
	ServerName    string
	RabbitConnect *amqp.Connection
	EtcdHostPort  string
	v             *viper.Viper
)

// AdminEnvInit 管理员模块环境初始化
func AdminEnvInit() {
	readConfig("admin")
	mysqlInit()
	redisInit(context.Background())
	serverInit()
	rabbitConnectionInit()
	grpcConnectionInit()
	etcdInit()
}

// PracticeEnvInit 练习模块环境初始化
func PracticeEnvInit() {
	readConfig("practice")
	mysqlInit()
	redisInit(context.Background())
	serverInit()
	rabbitConnectionInit()
	grpcConnectionInit()
	etcdInit()
}

// GatewayEnvInit 网关模块环境初始化
func GatewayEnvInit() {
	readConfig("gateway")
	etcdInit()
}

// ScheduleEnvInit 定时任务模块环境初始化
func ScheduleEnvInit() {
	readConfig("schedule")
	mysqlInit()
	redisInit(context.Background())
}

// readConfig 读取配置文件
func readConfig(yamlName string) {
	v1 := viper.New()
	v1.SetConfigName(yamlName)
	v1.SetConfigType("yaml")
	v1.AddConfigPath(".././commonModule/config")
	err := v1.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置文件失败", err)
	}
	v = v1
}

// mysqlInit 初始化mysql
func mysqlInit() {
	client, err := gorm.Open(mysql.Open(
		strings.Join([]string{
			v.GetString("mysql.username"),
			":",
			v.GetString("mysql.password"),
			"@tcp(",
			v.GetString("mysql.host"),
			":",
			v.GetString("mysql.port"),
			")/",
			v.GetString("mysql.database"),
			"?charset=utf8&parseTime=True&loc=Local"}, "")), nil)
	if err != nil {
		log.Fatal(err)
	}
	db, err := client.DB()
	if err != nil {
		log.Fatal("mysql初始化失败", err)
	}
	db.SetMaxIdleConns(v.GetInt("mysql.maxIdleConns"))
	db.SetMaxOpenConns(v.GetInt("mysql.maxOpenConns"))
	db.SetConnMaxLifetime(time.Second * 30)
	DB = client
	log.Println("mysql初始化成功")
}

// redisInit 初始化redis
func redisInit(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     v.GetString("redis.host") + ":" + v.GetString("redis.port"),
		Password: v.GetString("redis.password"),
		DB:       v.GetInt("redis.database"),
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("redis初始化失败", err)
	}
	RedisClient = client
	log.Println("redis初始化成功")
}

// serverInit 初始化gin-web 端口
func serverInit() {
	ServerHost = v.GetString("service.host")
	ServerPort = v.GetString("service.port")
	ServerName = v.GetString("service.name")
}

// rabbitConnectionInit 初始化rabbitmq
func rabbitConnectionInit() {
	connection, err := amqp.Dial(
		strings.Join([]string{
			"amqp://",
			v.GetString("rabbitmq.username"),
			":",
			v.GetString("rabbitmq.password"),
			"@",
			v.GetString("rabbitmq.host"),
			":",
			v.GetString("rabbitmq.port")}, ""))
	if err != nil {
		log.Println("rabbitmq初始化失败", err)
	}
	RabbitConnect = connection
	log.Println("rabbitmq初始化成功")
}

func etcdInit() {
	EtcdHostPort = v.GetString("etcd.host") + ":" + v.GetString("etcd.port")
}

// grpcConnectionInit 初始化grpc(服务端)
func grpcConnectionInit() {

}

func ViperGetString(key string) string {
	return v.GetString(key)
}

func ViperGetInt(key string) int {
	return v.GetInt(key)
}
