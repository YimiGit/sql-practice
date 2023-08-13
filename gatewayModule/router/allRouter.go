package router

import (
	"common/middleware"
	"context"
	"gateway/controller"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

// ResponseData 结构体用于存储异步请求和响应
type responseData struct {
	Bytes   []byte
	Code    int
	Context *gin.Context
}

// responseChannel 定义全局ResponseChannel 用于存储异步请求和响应 + 请求限流
var responseChannel = make(chan *responseData)

func AllRouterInit(engine *gin.Engine, client *clientv3.Client) {

	go listenAndResponse()

	background := context.Background()

	//跨域
	engine.Use(middleware.Cors())
	{
		engine.POST("/login", controller.Login)
		engine.POST("/register", controller.Register)

		//user-server转发前缀
		userServerRoute := "/user-server/*path"
		//user-server
		userServerName := "user-server"
		engine.Any(userServerRoute, forwardRequest(getHostByServer(userServerName, client, background)))

		//practice-server转发前缀
		practiceServerRoute := "/practice-server/*path"
		//practice-server
		practiceServerName := "practice-server"
		engine.Any(practiceServerRoute, forwardRequest(getHostByServer(practiceServerName, client, background)))
	}
}

// forwardRequest 服务转发中间件
func forwardRequest(target *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(target)

		c.Request.Host = target.Host
		c.Request.URL.Path = c.Param("path")

		// 创建ResponseWriter 用于获取响应数据
		var respWriter = httptest.NewRecorder()
		proxy.ServeHTTP(respWriter, c.Request)

		dataPoint := &responseData{Bytes: respWriter.Body.Bytes(), Code: respWriter.Code, Context: c}

		//写入ResponseChannel
		select {
		case responseChannel <- dataPoint:
		//case <-time.After(time.Second): // 超时处理
		default:
			c.String(500, "服务器繁忙")
			c.Abort()
			return
		}
	}
}

// getHostByServer 从etcd中获取服务地址,etcd做负载均衡
func getHostByServer(serverName string, client *clientv3.Client, c context.Context) *url.URL {
	serverHostByte, err := client.Get(c, serverName)
	if err != nil {
		log.Println("etcd获取服务失败", err)
		panic(err)
	}
	serverUrl, err := url.Parse("http://" + string(serverHostByte.Kvs[0].Value))
	if err != nil {
		log.Println("url解析失败", err)
		panic(err)
	}
	return serverUrl
}

// 监听全局responseChannel并将响应写回客户端
func listenAndResponse() {
	for {
		// 从responseChannel中获取响应数据
		realData := <-responseChannel
		realData.Context.String(realData.Code, string(realData.Bytes))
	}

	//10线程读channel， 10线程响应  存在竞争态
	//j := cap(responseChannel)
	//log.Println("responseChannel容量: ", j)
	//for i := 0; i < j; i++ {
	//	go func() {
	//		for {
	//			// 从responseChannel中获取响应数据
	//			realData := <-responseChannel
	//			realData.Context.String(realData.Code, string(realData.Bytes))
	//		}
	//	}()
	//}

	// 单线程读channel，多线程响应  每次都创建新的goroutine
	//for {
	//	// 从responseChannel中获取响应数据
	//	realData := <-responseChannel
	//	go func(data *responseData) {
	//		// 写回客户端
	//		data.Context.String(data.Code, string(data.Bytes))
	//	}(realData)
	//}
}
