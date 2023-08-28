package config

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"stock/util"
	"strings"
	"time"
)

var (
	DB              *gorm.DB
	RedisClient     *redis.Client
	RedisLockClient *util.RedisLockClient
	ServerPort      string
	ServerHost      string
	ServerName      string
	RabbitConnect   *amqp.Connection
	v               *viper.Viper
)

// StockEnvInit 管理员模块环境初始化
func StockEnvInit() {
	readConfig("stock")
	mysqlInit()
	redisInit(context.Background())
	serverInit()
	rabbitConnectionInit()
}

// readConfig 读取配置文件
func readConfig(yamlName string) {
	v1 := viper.New()
	v1.SetConfigName(yamlName)
	v1.SetConfigType("yaml")
	v1.AddConfigPath("config")
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
	RedisLockClient = &util.RedisLockClient{
		RedisClient: client,
		Redsync:     redsync.New(goredis.NewPool(client)),
	}
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
